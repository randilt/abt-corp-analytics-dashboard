package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"analytics-dashboard-api/internal/services"
	"analytics-dashboard-api/pkg/logger"
)

func main() {
	var (
		csvPath   = flag.String("csv", "./data/raw/transactions.csv", "Path to CSV file")
		cachePath = flag.String("cache", "./data/processed/analytics_cache.json", "Path to cache file")
		logLevel  = flag.String("log", "info", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	log := logger.NewLogger(*logLevel)
	log.Info("Starting data preprocessing", "csv", *csvPath, "cache", *cachePath)

	// Check if CSV file exists
	if _, err := os.Stat(*csvPath); os.IsNotExist(err) {
		log.Error("CSV file does not exist", "path", *csvPath)
		os.Exit(1)
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(*cachePath), 0755); err != nil {
		log.Error("Failed to create cache directory", "error", err)
		os.Exit(1)
	}

	// Initialize CSV processor
	processor := services.NewCSVProcessor(log)

	// Process CSV with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	start := time.Now()
	stats, err := processor.PreprocessAndCache(ctx, *csvPath, *cachePath)
	if err != nil {
		log.Error("Preprocessing failed", "error", err)
		os.Exit(1)
	}

	log.Info("Preprocessing completed successfully",
		"records", stats.ProcessedRecords,
		"errors", stats.ErrorCount,
		"duration", time.Since(start),
		"memory_mb", stats.MemoryUsageMB,
	)

	fmt.Printf("✅ Preprocessing completed!\n")
	fmt.Printf("📊 Processed %d records in %v\n", stats.ProcessedRecords, time.Since(start))
	fmt.Printf("💾 Cache saved to: %s\n", *cachePath)
}