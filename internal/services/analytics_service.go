package services

import (
	"sort"
	"sync"
	"time"

	"analytics-dashboard-api/internal/models"
	"analytics-dashboard-api/pkg/logger"
)

type AnalyticsService struct {
	logger logger.Logger
}

func NewAnalyticsService(logger logger.Logger) *AnalyticsService {
	return &AnalyticsService{
		logger: logger,
	}
}

// GenerateAnalytics processes transactions and generates all required analytics
func (s *AnalyticsService) GenerateAnalytics(transactions []models.Transaction) *models.AnalyticsResponse {
	startTime := time.Now()
	s.logger.Info("Generating analytics", "records", len(transactions))

	// Use concurrent processing for different analytics
	var wg sync.WaitGroup
	wg.Add(4)

	var countryRevenue []models.CountryRevenue
	var topProducts []models.ProductFrequency
	var monthlySales []models.MonthlySales
	var topRegions []models.RegionRevenue

	// Process country revenue concurrently
	go func() {
		defer wg.Done()
		countryRevenue = s.generateCountryRevenue(transactions)
	}()

	// Process top products concurrently
	go func() {
		defer wg.Done()
		topProducts = s.generateTopProducts(transactions)
	}()

	// Process monthly sales concurrently
	go func() {
		defer wg.Done()
		monthlySales = s.generateMonthlySales(transactions)
	}()

	// Process top regions concurrently
	go func() {
		defer wg.Done()
		topRegions = s.generateTopRegions(transactions)
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// Calculate processing time after all goroutines complete
	processingTime := time.Since(startTime)
	s.logger.Info("Analytics generation completed", "duration", processingTime)

	return &models.AnalyticsResponse{
		CountryRevenue:   countryRevenue,
		TopProducts:      topProducts,
		MonthlySales:     monthlySales,
		TopRegions:       topRegions,
		ProcessingTimeMs: processingTime.Milliseconds(),
		TotalRecords:     len(transactions),
		CacheHit:         false,
	}
}

// generateCountryRevenue creates country-level revenue table sorted by revenue
func (s *AnalyticsService) generateCountryRevenue(transactions []models.Transaction) []models.CountryRevenue {
	// Use map for efficient aggregation: "country|product" -> revenue data
	revenueMap := make(map[string]*models.CountryRevenue)

	for _, t := range transactions {
		key := t.Country + "|" + t.ProductName

		if entry, exists := revenueMap[key]; exists {
			entry.TotalRevenue += t.TotalPrice
			entry.TransactionCount++
		} else {
			revenueMap[key] = &models.CountryRevenue{
				Country:          t.Country,
				ProductName:      t.ProductName,
				TotalRevenue:     t.TotalPrice,
				TransactionCount: 1,
			}
		}
	}

	// Convert map to slice
	result := make([]models.CountryRevenue, 0, len(revenueMap))
	for _, entry := range revenueMap {
		result = append(result, *entry)
	}

	// Sort by total revenue descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalRevenue > result[j].TotalRevenue
	})

	return result
}

// generateTopProducts finds top 20 frequently purchased products with stock
func (s *AnalyticsService) generateTopProducts(transactions []models.Transaction) []models.ProductFrequency {
	// Aggregate by product ID
	productMap := make(map[string]*models.ProductFrequency)

	for _, t := range transactions {
		if entry, exists := productMap[t.ProductID]; exists {
			entry.PurchaseCount += t.Quantity
		} else {
			productMap[t.ProductID] = &models.ProductFrequency{
				ProductID:     t.ProductID,
				ProductName:   t.ProductName,
				PurchaseCount: t.Quantity,
				StockQuantity: t.StockQuantity, // Using latest stock quantity
			}
		}
	}

	// Convert to slice and sort by purchase count
	result := make([]models.ProductFrequency, 0, len(productMap))
	for _, entry := range productMap {
		result = append(result, *entry)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].PurchaseCount > result[j].PurchaseCount
	})

	// Return top 20
	if len(result) > 20 {
		result = result[:20]
	}

	return result
}

// generateMonthlySales creates monthly sales volume chart data
func (s *AnalyticsService) generateMonthlySales(transactions []models.Transaction) []models.MonthlySales {
	monthlyMap := make(map[string]*models.MonthlySales)

	for _, t := range transactions {
		month := t.GetMonth()

		if entry, exists := monthlyMap[month]; exists {
			entry.SalesVolume += t.TotalPrice
			entry.ItemCount += t.Quantity
		} else {
			monthlyMap[month] = &models.MonthlySales{
				Month:       month,
				SalesVolume: t.TotalPrice,
				ItemCount:   t.Quantity,
			}
		}
	}

	// Convert to slice and sort by month
	result := make([]models.MonthlySales, 0, len(monthlyMap))
	for _, entry := range monthlyMap {
		result = append(result, *entry)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Month < result[j].Month
	})

	return result
}

// generateTopRegions finds top 30 regions by revenue and items sold
func (s *AnalyticsService) generateTopRegions(transactions []models.Transaction) []models.RegionRevenue {
	regionMap := make(map[string]*models.RegionRevenue)

	for _, t := range transactions {
		if entry, exists := regionMap[t.Region]; exists {
			entry.TotalRevenue += t.TotalPrice
			entry.ItemsSold += t.Quantity
		} else {
			regionMap[t.Region] = &models.RegionRevenue{
				Region:       t.Region,
				TotalRevenue: t.TotalPrice,
				ItemsSold:    t.Quantity,
			}
		}
	}

	// Convert to slice and sort by revenue
	result := make([]models.RegionRevenue, 0, len(regionMap))
	for _, entry := range regionMap {
		result = append(result, *entry)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalRevenue > result[j].TotalRevenue
	})

	// Return top 30
	if len(result) > 30 {
		result = result[:30]
	}

	return result
}
