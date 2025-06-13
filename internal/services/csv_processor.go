package services

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"analytics-dashboard-api/internal/config"
	"analytics-dashboard-api/internal/models"
	"analytics-dashboard-api/pkg/logger"
)

var (
	// Global mutex to prevent concurrent CSV processing
	// The application crashed when multiple getanalytics requests were made (with cache reset)
	// simultaneously, so we use a global mutex to ensure only one CSV processing
	globalProcessingMu sync.Mutex 
)

type CSVProcessor struct {
	logger     logger.Logger
	batchSize  int
	workerPool int
	bufferSize int
	cacheConfig *config.CacheConfig
}

type ProcessingResult struct {
	Transactions []models.Transaction
	Stats        models.ProcessingStats
	Error        error
}

type BatchResult struct {
	Transactions []models.Transaction
	BatchIndex   int
	Error        error
	ParseErrors  int
}

// IndexedBatch contains the batch data with its correct index
type IndexedBatch struct {
	Records [][]string
	Index   int
}

func NewCSVProcessor(logger logger.Logger, csvConfig *config.CSVConfig, cacheConfig *config.CacheConfig) *CSVProcessor {
	return &CSVProcessor{
		logger:      logger,
		batchSize:   csvConfig.BatchSize,
		workerPool:  min(csvConfig.WorkerPool, runtime.NumCPU()),
		bufferSize:  csvConfig.BufferSize,
		cacheConfig: cacheConfig,
	}
}

// ProcessLargeCSV processes a large CSV file in batches using multiple goroutines
func (p *CSVProcessor) ProcessLargeCSV(ctx context.Context, filePath string) (*ProcessingResult, error) {
	// Acquire global lock to prevent concurrent processing
	globalProcessingMu.Lock()
	defer globalProcessingMu.Unlock()

	startTime := time.Now()
	p.logger.Info("Starting CSV processing", "file", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	bufferedReader := bufio.NewReaderSize(file, p.bufferSize)
	csvReader := csv.NewReader(bufferedReader)
	csvReader.ReuseRecord = true

	// Skip header row
	if _, err := csvReader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Setup concurrent processing pipeline with indexed batches
	batchChan := make(chan IndexedBatch, 5)
	resultChan := make(chan BatchResult, 5)

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < p.workerPool; i++ {
		wg.Add(1)
		go p.processBatchWorker(ctx, batchChan, resultChan, &wg, i)
	}

	// Start batch reader goroutine
	batchCount := make(chan int, 1)
	go p.readBatches(ctx, csvReader, batchChan, batchCount)

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results with proper ordering
	var allTransactions []models.Transaction
	var totalRecords, errorCount, totalParseErrors int
	batchResults := make(map[int]BatchResult)
	maxBatchIndex := -1

	// Collect all batch results
	for result := range resultChan {
		if result.Error != nil {
			p.logger.Error("Batch processing error",
				"batch", result.BatchIndex,
				"error", result.Error)
			errorCount++
			continue
		}

		batchResults[result.BatchIndex] = result
		totalParseErrors += result.ParseErrors
		if result.BatchIndex > maxBatchIndex {
			maxBatchIndex = result.BatchIndex
		}

		p.logger.Debug("Batch completed",
			"batch", result.BatchIndex,
			"transactions", len(result.Transactions),
			"parse_errors", result.ParseErrors)
	}

	// Wait for batch count
	totalBatches := <-batchCount
	p.logger.Info("Batch processing summary",
		"total_batches", totalBatches,
		"completed_batches", len(batchResults),
		"max_batch_index", maxBatchIndex)

	// Log missing batches for debugging
	missingBatches := 0
	for i := 0; i < totalBatches; i++ {
		if _, exists := batchResults[i]; !exists {
			p.logger.Error("Missing batch in results", "batch_index", i)
			missingBatches++
		}
	}

	if missingBatches > 0 {
		p.logger.Error("CRITICAL: Missing batches detected",
			"missing_count", missingBatches,
			"total_batches", totalBatches)
	}

	// Reassemble transactions in correct order
	for i := 0; i < totalBatches; i++ {
		if batch, exists := batchResults[i]; exists {
			allTransactions = append(allTransactions, batch.Transactions...)
			totalRecords += len(batch.Transactions)
		} else {
			errorCount++
		}
	}

	// Calculate memory usage
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryUsageMB := float64(memStats.Alloc) / 1024 / 1024

	stats := models.ProcessingStats{
		TotalRecords:     totalRecords,
		ProcessedRecords: len(allTransactions),
		ErrorCount:       errorCount + totalParseErrors,
		ProcessingTime:   time.Since(startTime),
		MemoryUsageMB:    memoryUsageMB,
	}

	p.logger.Info("CSV processing completed",
		"total_records", totalRecords,
		"processed_records", len(allTransactions),
		"parse_errors", totalParseErrors,
		"batch_errors", errorCount,
		"missing_batches", missingBatches,
		"duration", stats.ProcessingTime,
		"memory_mb", memoryUsageMB)

	expectedRecords := totalBatches * p.batchSize
	if float64(len(allTransactions)) < float64(expectedRecords)*0.95 {
		p.logger.Error("CRITICAL: Significant data loss detected",
			"expected_approx", expectedRecords,
			"actual", len(allTransactions),
			"loss_percentage", float64(expectedRecords-len(allTransactions))/float64(expectedRecords)*100)
	}

	return &ProcessingResult{
		Transactions: allTransactions,
		Stats:        stats,
	}, nil
}

// PreprocessAndCache processes CSV and caches results for faster subsequent loads
func (p *CSVProcessor) PreprocessAndCache(ctx context.Context, csvPath, cachePath string) (*models.ProcessingStats, error) {
	// Process CSV
	result, err := p.ProcessLargeCSV(ctx, csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to process CSV: %w", err)
	}

	// Create analytics data
	analyticsService := NewAnalyticsService(p.logger)
	analytics := analyticsService.GenerateAnalytics(result.Transactions)

	cacheService := NewCacheService(p.logger, p.cacheConfig)
	if err := cacheService.SaveToFile(cachePath, analytics); err != nil {
		p.logger.Error("Failed to save cache", "error", err)
		// Don't fail the entire process if caching fails
	}

	return &result.Stats, nil
}

// readBatches reads CSV records in batches and sends them to the batch channel
func (p *CSVProcessor) readBatches(ctx context.Context, reader *csv.Reader, batchChan chan<- IndexedBatch, batchCount chan<- int) {
	defer close(batchChan)

	batchIndex := 0
	var batch [][]string
	totalRowsRead := 0

	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				select {
				case batchChan <- IndexedBatch{Records: batch, Index: batchIndex}:
					batchIndex++
				case <-ctx.Done():
				}
			}
			batchCount <- batchIndex
			return
		default:
		}

		record, err := reader.Read()
		if err == io.EOF {
			// Send final batch if it has records
			if len(batch) > 0 {
				select {
				case batchChan <- IndexedBatch{Records: batch, Index: batchIndex}:
					batchIndex++
				case <-ctx.Done():
				}
			}
			batchCount <- batchIndex
			p.logger.Info("Finished reading CSV",
				"total_rows_read", totalRowsRead,
				"total_batches", batchIndex)
			return
		}

		if err != nil {
			p.logger.Error("CSV read error", "error", err, "row", totalRowsRead)
			continue
		}

		totalRowsRead++

		// Create a copy of the record since csv.Reader reuses the slice
		recordCopy := make([]string, len(record))
		copy(recordCopy, record)
		batch = append(batch, recordCopy)

		if len(batch) >= p.batchSize {
			select {
			case batchChan <- IndexedBatch{Records: batch, Index: batchIndex}:
				p.logger.Debug("Batch sent", "batch_index", batchIndex, "size", len(batch))
				batch = make([][]string, 0, p.batchSize)
				batchIndex++
			case <-ctx.Done():
				batchCount <- batchIndex
				return
			}
		}
	}
}

// 	processBatchWorker processes batches of CSV records concurrently
func (p *CSVProcessor) processBatchWorker(ctx context.Context, batchChan <-chan IndexedBatch, resultChan chan<- BatchResult, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	// Process indexed batches as they come from the channel
	for indexedBatch := range batchChan {
		select {
		case <-ctx.Done():
			return
		default:
		}

		batch := indexedBatch.Records
		batchIndex := indexedBatch.Index

		transactions := make([]models.Transaction, 0, len(batch))
		parseErrors := 0

		for rowIndex, record := range batch {
			var transaction models.Transaction
			if err := transaction.ParseCSVRow(record); err != nil {
				parseErrors++
				if parseErrors <= 5 { // Log first 5 errors per batch
					p.logger.Debug("Failed to parse CSV row",
						"worker", workerID,
						"batch", batchIndex,
						"row", rowIndex,
						"error", err,
						"record_length", len(record))
				}
				continue
			}
			transactions = append(transactions, transaction)
		}

		p.logger.Debug("Worker processed batch",
			"worker", workerID,
			"batch", batchIndex,
			"input_rows", len(batch),
			"output_transactions", len(transactions),
			"parse_errors", parseErrors)

		// Send result with correct batch index
		resultChan <- BatchResult{
			Transactions: transactions,
			BatchIndex:   batchIndex,
			ParseErrors:  parseErrors,
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
