package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"analytics-dashboard-api/internal/handlers"
	"analytics-dashboard-api/internal/models"
	"analytics-dashboard-api/internal/services"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...interface{})  {}
func (m *mockLogger) Error(msg string, args ...interface{}) {}
func (m *mockLogger) Debug(msg string, args ...interface{}) {}
func (m *mockLogger) Warn(msg string, args ...interface{})  {}

type mockAnalyticsService struct {
	generateAnalyticsFunc func([]models.Transaction) *models.AnalyticsResponse
}

func (m *mockAnalyticsService) GenerateAnalytics(transactions []models.Transaction) *models.AnalyticsResponse {
	if m.generateAnalyticsFunc != nil {
		return m.generateAnalyticsFunc(transactions)
	}
	return &models.AnalyticsResponse{
		CountryRevenue: []models.CountryRevenue{
			{Country: "USA", ProductName: "Product A", TotalRevenue: 1000.0, TransactionCount: 10},
		},
		TopProducts: []models.ProductFrequency{
			{ProductID: "P1", ProductName: "Product A", PurchaseCount: 100, StockQuantity: 50},
		},
		MonthlySales: []models.MonthlySales{
			{Month: "2023-01", SalesVolume: 5000.0, ItemCount: 200},
		},
		TopRegions: []models.RegionRevenue{
			{Region: "California", TotalRevenue: 2000.0, ItemsSold: 150},
		},
		ProcessingTimeMs: 1000,
		TotalRecords:     100,
		CacheHit:         false,
	}
}

type mockCacheService struct {
	loadFromCacheFunc func() (*models.AnalyticsResponse, bool)
	loadFromFileFunc  func(string) (*models.AnalyticsResponse, error)
	saveToMemoryFunc  func(*models.AnalyticsResponse)
	saveToFileFunc    func(string, *models.AnalyticsResponse) error
}

func (m *mockCacheService) LoadFromCache() (*models.AnalyticsResponse, bool) {
	if m.loadFromCacheFunc != nil {
		return m.loadFromCacheFunc()
	}
	return nil, false
}

func (m *mockCacheService) LoadFromFile(path string) (*models.AnalyticsResponse, error) {
	if m.loadFromFileFunc != nil {
		return m.loadFromFileFunc(path)
	}
	return nil, os.ErrNotExist
}

func (m *mockCacheService) SaveToMemory(data *models.AnalyticsResponse) {
	if m.saveToMemoryFunc != nil {
		m.saveToMemoryFunc(data)
	}
}

func (m *mockCacheService) SaveToFile(path string, data *models.AnalyticsResponse) error {
	if m.saveToFileFunc != nil {
		return m.saveToFileFunc(path, data)
	}
	return nil
}

type mockCSVProcessor struct {
	processLargeCSVFunc    func(context.Context, string) (*services.ProcessingResult, error)
	preprocessAndCacheFunc func(context.Context, string, string) (*models.ProcessingStats, error)
}

func (m *mockCSVProcessor) ProcessLargeCSV(ctx context.Context, filePath string) (*services.ProcessingResult, error) {
	if m.processLargeCSVFunc != nil {
		return m.processLargeCSVFunc(ctx, filePath)
	}
	return &services.ProcessingResult{
		Transactions: []models.Transaction{},
		Stats: models.ProcessingStats{
			TotalRecords:     100,
			ProcessedRecords: 100,
			ErrorCount:       0,
			ProcessingTime:   time.Second,
			MemoryUsageMB:    10.0,
		},
	}, nil
}

func (m *mockCSVProcessor) PreprocessAndCache(ctx context.Context, csvPath, cachePath string) (*models.ProcessingStats, error) {
	if m.preprocessAndCacheFunc != nil {
		return m.preprocessAndCacheFunc(ctx, csvPath, cachePath)
	}
	return &models.ProcessingStats{
		TotalRecords:     100,
		ProcessedRecords: 100,
		ErrorCount:       0,
		ProcessingTime:   time.Second,
		MemoryUsageMB:    10.0,
	}, nil
}

func createMockAnalyticsHandler() *handlers.AnalyticsHandler {
	logger := &mockLogger{}
	analyticsService := &mockAnalyticsService{}
	cacheService := &mockCacheService{}
	csvProcessor := &mockCSVProcessor{}

	return handlers.NewAnalyticsHandler(
		analyticsService,
		cacheService,
		csvProcessor,
		logger,
		"test.csv",
		"test_cache.json",
	)
}

func TestAnalyticsHandler_GetAnalytics_FromMemoryCache(t *testing.T) {
	logger := &mockLogger{}
	analyticsService := &mockAnalyticsService{}
	csvProcessor := &mockCSVProcessor{}

	cacheService := &mockCacheService{
		loadFromCacheFunc: func() (*models.AnalyticsResponse, bool) {
			return &models.AnalyticsResponse{
				CountryRevenue: []models.CountryRevenue{
					{Country: "USA", ProductName: "Product A", TotalRevenue: 1000.0, TransactionCount: 10},
				},
				TotalRecords: 100,
				CacheHit:     true,
			}, true
		},
	}

	handler := handlers.NewAnalyticsHandler(
		analyticsService,
		cacheService,
		csvProcessor,
		logger,
		"test.csv",
		"test_cache.json",
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics", nil)
	recorder := httptest.NewRecorder()

	handler.GetAnalytics(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetMonthlySales() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("GetMonthlySales() response parsing error: %v", err)
	}

	if _, exists := response["data"]; !exists {
		t.Error("GetMonthlySales() missing data field")
	}

	if _, exists := response["count"]; !exists {
		t.Error("GetMonthlySales() missing count field")
	}
}

func TestAnalyticsHandler_GetTopRegions(t *testing.T) {
	handler := createMockAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/top-regions", nil)
	recorder := httptest.NewRecorder()

	handler.GetTopRegions(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetTopRegions() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("GetTopRegions() response parsing error: %v", err)
	}

	if _, exists := response["data"]; !exists {
		t.Error("GetTopRegions() missing data field")
	}

	if _, exists := response["count"]; !exists {
		t.Error("GetTopRegions() missing count field")
	}
}

func TestAnalyticsHandler_RefreshCache(t *testing.T) {
	logger := &mockLogger{}
	analyticsService := &mockAnalyticsService{}
	cacheService := &mockCacheService{}

	csvProcessor := &mockCSVProcessor{
		preprocessAndCacheFunc: func(ctx context.Context, csvPath, cachePath string) (*models.ProcessingStats, error) {
			return &models.ProcessingStats{
				TotalRecords:     200,
				ProcessedRecords: 200,
				ErrorCount:       0,
				ProcessingTime:   2 * time.Second,
				MemoryUsageMB:    15.0,
			}, nil
		},
	}

	handler := handlers.NewAnalyticsHandler(
		analyticsService,
		cacheService,
		csvProcessor,
		logger,
		"test.csv",
		"test_cache.json",
	)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/analytics/refresh", nil)
	recorder := httptest.NewRecorder()

	handler.RefreshCache(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("RefreshCache() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("RefreshCache() response parsing error: %v", err)
	}

	expectedFields := []string{"message", "stats", "duration_ms"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("RefreshCache() missing field: %s", field)
		}
	}

	if message, ok := response["message"].(string); !ok || message != "Cache refreshed successfully" {
		t.Errorf("RefreshCache() message = %v, want 'Cache refreshed successfully'", response["message"])
	}
}

func TestAnalyticsHandler_RefreshCache_MethodNotAllowed(t *testing.T) {
	handler := createMockAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/refresh", nil)
	recorder := httptest.NewRecorder()

	handler.RefreshCache(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("RefreshCache() with GET status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}

func TestAnalyticsHandler_GetIntQueryParam(t *testing.T) {
	handler := createMockAnalyticsHandler()

	tests := []struct {
		name        string
		queryParams string
		wantLimit   float64
		wantOffset  float64
	}{
		{
			name:        "valid parameters",
			queryParams: "limit=50&offset=10",
			wantLimit:   50,
			wantOffset:  10,
		},
		{
			name:        "default values",
			queryParams: "",
			wantLimit:   100,
			wantOffset:  0,
		},
		{
			name:        "invalid limit",
			queryParams: "limit=abc&offset=5",
			wantLimit:   100,
			wantOffset:  5,
		},
		{
			name:        "negative values",
			queryParams: "limit=-10&offset=-5",
			wantLimit:   100,
			wantOffset:  0,
		},
		{
			name:        "limit over 1000",
			queryParams: "limit=1500&offset=0",
			wantLimit:   1000,
			wantOffset:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/analytics/country-revenue"
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			handler.GetCountryRevenue(recorder, req)

			var response map[string]interface{}
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("Response parsing error: %v", err)
			}

			if limit, ok := response["limit"].(float64); !ok || limit != tt.wantLimit {
				t.Errorf("limit = %v, want %v", response["limit"], tt.wantLimit)
			}

			if offset, ok := response["offset"].(float64); !ok || offset != tt.wantOffset {
				t.Errorf("offset = %v, want %v", response["offset"], tt.wantOffset)
			}
		})
	}
}

func TestAnalyticsHandler_GetCountryRevenue(t *testing.T) {
	logger := &mockLogger{}
	analyticsService := &mockAnalyticsService{}
	csvProcessor := &mockCSVProcessor{}

	cacheService := &mockCacheService{
		loadFromCacheFunc: func() (*models.AnalyticsResponse, bool) {
			return &models.AnalyticsResponse{
				CountryRevenue: []models.CountryRevenue{
					{Country: "USA", ProductName: "Product A", TotalRevenue: 1000.0, TransactionCount: 10},
					{Country: "Canada", ProductName: "Product B", TotalRevenue: 800.0, TransactionCount: 8},
					{Country: "Germany", ProductName: "Product C", TotalRevenue: 600.0, TransactionCount: 6},
				},
				TotalRecords: 100,
			}, true
		},
	}

	handler := handlers.NewAnalyticsHandler(
		analyticsService,
		cacheService,
		csvProcessor,
		logger,
		"test.csv",
		"test_cache.json",
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/country-revenue", nil)
	recorder := httptest.NewRecorder()

	handler.GetCountryRevenue(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetCountryRevenue() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("GetCountryRevenue() response parsing error: %v", err)
	}

	expectedFields := []string{"data", "count", "total", "limit", "offset", "has_more"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("GetCountryRevenue() missing field: %s", field)
		}
	}

	if data, ok := response["data"].([]interface{}); ok {
		if len(data) != 3 {
			t.Errorf("GetCountryRevenue() data length = %d, want 3", len(data))
		}
	} else {
		t.Error("GetCountryRevenue() data should be an array")
	}
}

func TestAnalyticsHandler_GetCountryRevenue_WithPagination(t *testing.T) {
	handler := createMockAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/country-revenue?limit=2&offset=1", nil)
	recorder := httptest.NewRecorder()

	handler.GetCountryRevenue(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetCountryRevenue() with pagination status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("GetCountryRevenue() response parsing error: %v", err)
	}

	if limit, ok := response["limit"].(float64); !ok || limit != 2 {
		t.Errorf("GetCountryRevenue() limit = %v, want 2", response["limit"])
	}

	if offset, ok := response["offset"].(float64); !ok || offset != 1 {
		t.Errorf("GetCountryRevenue() offset = %v, want 1", response["offset"])
	}
}

func TestAnalyticsHandler_GetAnalyticsStats(t *testing.T) {
	handler := createMockAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/stats", nil)
	recorder := httptest.NewRecorder()

	handler.GetAnalyticsStats(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetAnalyticsStats() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("GetAnalyticsStats() response parsing error: %v", err)
	}

	expectedFields := []string{
		"total_records", "processing_time_ms", "cache_hit",
		"country_revenue_count", "top_products_count",
		"monthly_sales_count", "top_regions_count", "endpoints",
	}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("GetAnalyticsStats() missing field: %s", field)
		}
	}

	if endpoints, ok := response["endpoints"].(map[string]interface{}); ok {
		endpointFields := []string{"country_revenue", "top_products", "monthly_sales", "top_regions"}
		for _, field := range endpointFields {
			if _, exists := endpoints[field]; !exists {
				t.Errorf("GetAnalyticsStats() endpoints missing field: %s", field)
			}
		}
	} else {
		t.Error("GetAnalyticsStats() endpoints should be an object")
	}
}

func TestAnalyticsHandler_GetTopProducts(t *testing.T) {
	handler := createMockAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/top-products", nil)
	recorder := httptest.NewRecorder()

	handler.GetTopProducts(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetTopProducts() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("GetTopProducts() response parsing error: %v", err)
	}

	if _, exists := response["data"]; !exists {
		t.Error("GetTopProducts() missing data field")
	}

	if _, exists := response["count"]; !exists {
		t.Error("GetTopProducts() missing count field")
	}
}

func TestAnalyticsHandler_GetMonthlySales(t *testing.T) {
	handler := createMockAnalyticsHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/monthly-sales", nil)
	recorder := httptest.NewRecorder()

	handler.GetMonthlySales(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetMonthlySales() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("GetMonthlySales() response parsing error: %v", err)
	}

	if summary, exists := response["summary"]; !exists {
		t.Error("GetMonthlySales() should return summary section")
	} else if summaryMap, ok := summary.(map[string]interface{}); ok {
		if cacheHit, exists := summaryMap["cache_hit"]; !exists || cacheHit != true {
			t.Error("GetMonthlySales() cache_hit should be true")
		}
	}
}
