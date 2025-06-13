package services_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"analytics-dashboard-api/internal/config"
	"analytics-dashboard-api/internal/models"
	"analytics-dashboard-api/internal/services"
)

func createTestAnalyticsResponse() *models.AnalyticsResponse {
	return &models.AnalyticsResponse{
		CountryRevenue: []models.CountryRevenue{
			{
				Country:          "USA",
				ProductName:      "Product A",
				TotalRevenue:     1000.0,
				TransactionCount: 10,
			},
		},
		TopProducts: []models.ProductFrequency{
			{
				ProductID:     "P1",
				ProductName:   "Product A",
				PurchaseCount: 100,
				StockQuantity: 50,
			},
		},
		MonthlySales: []models.MonthlySales{
			{
				Month:       "2023-01",
				SalesVolume: 5000.0,
				ItemCount:   200,
			},
		},
		TopRegions: []models.RegionRevenue{
			{
				Region:       "California",
				TotalRevenue: 2000.0,
				ItemsSold:    150,
			},
		},
		ProcessingTimeMs: 1000,
		TotalRecords:     100,
		CacheHit:         false,
	}
}

func createTestCacheConfig() *config.CacheConfig {
	return &config.CacheConfig{
		FilePath: "./test_cache.json",
		TTL:      24 * time.Hour,
	}
}

func TestCacheService_SaveToMemory_LoadFromCache(t *testing.T) {
	logger := &mockLogger{}
	cacheConfig := createTestCacheConfig()
	cacheService := services.NewCacheService(logger, cacheConfig)
	testData := createTestAnalyticsResponse()

	_, hit := cacheService.LoadFromCache()
	if hit {
		t.Error("Cache should be empty initially")
	}

	cacheService.SaveToMemory(testData)

	cached, hit := cacheService.LoadFromCache()
	if !hit {
		t.Error("Cache hit should be true after saving")
	}

	if cached == nil {
		t.Fatal("Cached data should not be nil")
	}

	if !cached.CacheHit {
		t.Error("CacheHit flag should be true when loading from cache")
	}

	if len(cached.CountryRevenue) != len(testData.CountryRevenue) {
		t.Errorf("CountryRevenue length mismatch: got %d, want %d",
			len(cached.CountryRevenue), len(testData.CountryRevenue))
	}

	if cached.TotalRecords != testData.TotalRecords {
		t.Errorf("TotalRecords mismatch: got %d, want %d",
			cached.TotalRecords, testData.TotalRecords)
	}
}

func TestCacheService_SaveToFile_LoadFromFile(t *testing.T) {
	logger := &mockLogger{}
	testData := createTestAnalyticsResponse()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_cache.json")

	cacheConfig := &config.CacheConfig{
		FilePath: filePath,
		TTL:      24 * time.Hour,
	}
	cacheService := services.NewCacheService(logger, cacheConfig)

	err := cacheService.SaveToFile(filePath, testData)
	if err != nil {
		t.Fatalf("SaveToFile() error = %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Cache file should exist after saving")
	}

	loaded, err := cacheService.LoadFromFile(filePath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	if loaded == nil {
		t.Fatal("Loaded data should not be nil")
	}

	if len(loaded.CountryRevenue) != len(testData.CountryRevenue) {
		t.Errorf("CountryRevenue length mismatch: got %d, want %d",
			len(loaded.CountryRevenue), len(testData.CountryRevenue))
	}

	if loaded.CountryRevenue[0].Country != testData.CountryRevenue[0].Country {
		t.Errorf("Country mismatch: got %s, want %s",
			loaded.CountryRevenue[0].Country, testData.CountryRevenue[0].Country)
	}

	if loaded.TotalRecords != testData.TotalRecords {
		t.Errorf("TotalRecords mismatch: got %d, want %d",
			loaded.TotalRecords, testData.TotalRecords)
	}
}

func TestCacheService_LoadFromFile_NonexistentFile(t *testing.T) {
	logger := &mockLogger{}
	cacheConfig := createTestCacheConfig()
	cacheService := services.NewCacheService(logger, cacheConfig)

	_, err := cacheService.LoadFromFile("/nonexistent/path/cache.json")
	if err == nil {
		t.Error("LoadFromFile() should return error for non-existent file")
	}
}

func TestCacheService_LoadFromFile_InvalidJSON(t *testing.T) {
	logger := &mockLogger{}
	cacheConfig := createTestCacheConfig()
	cacheService := services.NewCacheService(logger, cacheConfig)

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "invalid_cache.json")

	err := os.WriteFile(filePath, []byte("invalid json content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = cacheService.LoadFromFile(filePath)
	if err == nil {
		t.Error("LoadFromFile() should return error for invalid JSON")
	}
}

func TestCacheService_SaveToFile_InvalidPath(t *testing.T) {
	logger := &mockLogger{}
	cacheConfig := createTestCacheConfig()
	cacheService := services.NewCacheService(logger, cacheConfig)
	testData := createTestAnalyticsResponse()

	err := cacheService.SaveToFile("/invalid/path/cache.json", testData)
	if err == nil {
		t.Error("SaveToFile() should return error for invalid path")
	}
}

func TestCacheService_CacheTTL(t *testing.T) {
	logger := &mockLogger{}
	cacheConfig := createTestCacheConfig()
	cacheService := services.NewCacheService(logger, cacheConfig)
	testData := createTestAnalyticsResponse()

	cacheService.SaveToMemory(testData)

	_, hit := cacheService.LoadFromCache()
	if !hit {
		t.Error("Cache should hit immediately after saving")
	}
}

func TestCacheService_LoadFromFile_AutoSaveToMemory(t *testing.T) {
	logger := &mockLogger{}
	testData := createTestAnalyticsResponse()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_cache.json")

	cacheConfig := &config.CacheConfig{
		FilePath: filePath,
		TTL:      24 * time.Hour,
	}
	cacheService := services.NewCacheService(logger, cacheConfig)

	jsonData, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	loaded, err := cacheService.LoadFromFile(filePath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	if loaded == nil {
		t.Fatal("Loaded data should not be nil")
	}

	cached, hit := cacheService.LoadFromCache()
	if !hit {
		t.Error("Memory cache should have data after LoadFromFile")
	}

	if cached == nil {
		t.Fatal("Memory cached data should not be nil")
	}

	if cached.TotalRecords != loaded.TotalRecords {
		t.Errorf("Memory cache TotalRecords mismatch: got %d, want %d",
			cached.TotalRecords, loaded.TotalRecords)
	}
}
