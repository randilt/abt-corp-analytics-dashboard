package services_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

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

func TestCacheService_SaveToMemory_LoadFromCache(t *testing.T) {
	logger := &mockLogger{}
	cacheService := services.NewCacheService(logger)
	testData := createTestAnalyticsResponse()

	// Initially cache should be empty
	_, hit := cacheService.LoadFromCache()
	if hit {
		t.Error("Cache should be empty initially")
	}

	// Save to memory
	cacheService.SaveToMemory(testData)

	// Load from cache
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

	// Verify data integrity
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
	cacheService := services.NewCacheService(logger)
	testData := createTestAnalyticsResponse()

	// Create temporary file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_cache.json")

	// Save to file
	err := cacheService.SaveToFile(filePath, testData)
	if err != nil {
		t.Fatalf("SaveToFile() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Cache file should exist after saving")
	}

	// Load from file
	loaded, err := cacheService.LoadFromFile(filePath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	if loaded == nil {
		t.Fatal("Loaded data should not be nil")
	}

	// Verify data integrity
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
	cacheService := services.NewCacheService(logger)

	// Try to load from non-existent file
	_, err := cacheService.LoadFromFile("/nonexistent/path/cache.json")
	if err == nil {
		t.Error("LoadFromFile() should return error for non-existent file")
	}
}

func TestCacheService_LoadFromFile_InvalidJSON(t *testing.T) {
	logger := &mockLogger{}
	cacheService := services.NewCacheService(logger)

	// Create temporary file with invalid JSON
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "invalid_cache.json")

	err := os.WriteFile(filePath, []byte("invalid json content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Try to load invalid JSON
	_, err = cacheService.LoadFromFile(filePath)
	if err == nil {
		t.Error("LoadFromFile() should return error for invalid JSON")
	}
}

func TestCacheService_SaveToFile_InvalidPath(t *testing.T) {
	logger := &mockLogger{}
	cacheService := services.NewCacheService(logger)
	testData := createTestAnalyticsResponse()

	// Try to save to invalid path
	err := cacheService.SaveToFile("/invalid/path/cache.json", testData)
	if err == nil {
		t.Error("SaveToFile() should return error for invalid path")
	}
}

func TestCacheService_CacheTTL(t *testing.T) {
	// This test would require modifying CacheService to accept TTL or use dependency injection
	// For now, we'll test the basic TTL concept by manipulating time indirectly

	logger := &mockLogger{}
	cacheService := services.NewCacheService(logger)
	testData := createTestAnalyticsResponse()

	// Save to memory
	cacheService.SaveToMemory(testData)

	// Immediately load should hit
	_, hit := cacheService.LoadFromCache()
	if !hit {
		t.Error("Cache should hit immediately after saving")
	}

	// Note: Testing actual TTL expiration would require either:
	// 1. Dependency injection of time interface
	// 2. Modifying the service to accept TTL parameter
	// 3. Waiting for actual TTL (not practical in unit tests)
	// For comprehensive testing, consider implementing time interface injection
}

func TestCacheService_LoadFromFile_AutoSaveToMemory(t *testing.T) {
	logger := &mockLogger{}
	cacheService := services.NewCacheService(logger)
	testData := createTestAnalyticsResponse()

	// Create temporary file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_cache.json")

	// Manually create cache file
	jsonData, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Load from file (should auto-save to memory)
	loaded, err := cacheService.LoadFromFile(filePath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	if loaded == nil {
		t.Fatal("Loaded data should not be nil")
	}

	// Now memory cache should also have the data
	cached, hit := cacheService.LoadFromCache()
	if !hit {
		t.Error("Memory cache should have data after LoadFromFile")
	}

	if cached == nil {
		t.Fatal("Memory cached data should not be nil")
	}

	// Verify both loaded and cached data are equivalent
	if cached.TotalRecords != loaded.TotalRecords {
		t.Errorf("Memory cache TotalRecords mismatch: got %d, want %d",
			cached.TotalRecords, loaded.TotalRecords)
	}
}
