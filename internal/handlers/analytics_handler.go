package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"analytics-dashboard-api/internal/models"
	"analytics-dashboard-api/internal/utils"
	"analytics-dashboard-api/pkg/logger"
)

type DuckDBService interface {
	LoadFromCSV(string) error
	GetCountryRevenue(context.Context, int, int) ([]models.CountryRevenue, error)
	GetTopProducts(context.Context) ([]models.ProductFrequency, error)
	GetMonthlySales(context.Context) ([]models.MonthlySales, error)
	GetTopRegions(context.Context) ([]models.RegionRevenue, error)
	GetTotalRecords(context.Context) (int, error)
	GetCountryRevenueCount(context.Context) (int, error)
	Close() error
}

type AnalyticsHandler struct {
	duckdbService DuckDBService
	logger        logger.Logger
	csvPath       string
	initialized   bool
}

func NewAnalyticsHandler(
	duckdbService DuckDBService,
	logger logger.Logger,
	csvPath string,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		duckdbService: duckdbService,
		logger:        logger,
		csvPath:       csvPath,
		initialized:   false,
	}
}

// ensureInitialized loads CSV data into DuckDB if not already done
func (h *AnalyticsHandler) ensureInitialized(ctx context.Context) error {
	if h.initialized {
		return nil
	}

	h.logger.Info("Initializing DuckDB with CSV data", "file", h.csvPath)
	
	if err := h.duckdbService.LoadFromCSV(h.csvPath); err != nil {
		return fmt.Errorf("failed to load CSV into DuckDB: %w", err)
	}

	h.initialized = true
	h.logger.Info("DuckDB initialization completed")
	return nil
}

// GetAnalytics returns all dashboard analytics data
func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	ctx := r.Context()

	h.logger.Info("Analytics request received", "method", r.Method, "path", r.URL.Path)

	// Ensure DuckDB is initialized
	if err := h.ensureInitialized(ctx); err != nil {
		h.logger.Error("Failed to initialize DuckDB", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to initialize database")
		return
	}

	// Get all analytics data concurrently
	var countryRevenue []models.CountryRevenue
	var topProducts []models.ProductFrequency
	var monthlySales []models.MonthlySales
	var topRegions []models.RegionRevenue
	var totalRecords int
	var countryRevenueCount int

	type result struct {
		name string
		err  error
	}

	results := make(chan result, 6)

	// Get country revenue (first 1000 records)
	go func() {
		data, err := h.duckdbService.GetCountryRevenue(ctx, 1000, 0)
		countryRevenue = data
		results <- result{"country_revenue", err}
	}()

	// Get top products
	go func() {
		data, err := h.duckdbService.GetTopProducts(ctx)
		topProducts = data
		results <- result{"top_products", err}
	}()

	// Get monthly sales
	go func() {
		data, err := h.duckdbService.GetMonthlySales(ctx)
		monthlySales = data
		results <- result{"monthly_sales", err}
	}()

	// Get top regions
	go func() {
		data, err := h.duckdbService.GetTopRegions(ctx)
		topRegions = data
		results <- result{"top_regions", err}
	}()

	// Get total records
	go func() {
		count, err := h.duckdbService.GetTotalRecords(ctx)
		totalRecords = count
		results <- result{"total_records", err}
	}()

	// Get country revenue count
	go func() {
		count, err := h.duckdbService.GetCountryRevenueCount(ctx)
		countryRevenueCount = count
		results <- result{"country_revenue_count", err}
	}()

	// Wait for all goroutines to complete
	var errors []string
	for i := 0; i < 6; i++ {
		res := <-results
		if res.err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", res.name, res.err))
		}
	}

	if len(errors) > 0 {
		h.logger.Error("Failed to get analytics data", "errors", errors)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get analytics data")
		return
	}

	processingTime := time.Since(startTime)
	analytics := &models.AnalyticsResponse{
		CountryRevenue:   countryRevenue,
		TopProducts:      topProducts,
		MonthlySales:     monthlySales,
		TopRegions:       topRegions,
		ProcessingTimeMs: processingTime.Milliseconds(),
		TotalRecords:     totalRecords,
		CacheHit:         false, // DuckDB queries are always fresh
	}

	h.logger.Info("Analytics generated successfully",
		"records", totalRecords,
		"country_revenue_count", countryRevenueCount,
		"processing_time", processingTime)

	// Return summary version
	summary := h.createAnalyticsSummary(analytics)
	utils.WriteJSONResponse(w, http.StatusOK, summary)
}

// GetCountryRevenue returns country-level revenue data
func (h *AnalyticsHandler) GetCountryRevenue(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit := h.getIntQueryParam(r, "limit", 100) // Default 100, max 1000
	offset := h.getIntQueryParam(r, "offset", 0)

	if limit > 1000 {
		limit = 1000 // Cap at 1000 records
	}

	// Ensure DuckDB is initialized
	if err := h.ensureInitialized(r.Context()); err != nil {
		h.logger.Error("Failed to initialize DuckDB", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to initialize database")
		return
	}

	// Get data from DuckDB
	data, err := h.duckdbService.GetCountryRevenue(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to get country revenue", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get country revenue data")
		return
	}

	// Get total count for pagination
	total, err := h.duckdbService.GetCountryRevenueCount(r.Context())
	if err != nil {
		h.logger.Error("Failed to get country revenue count", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get total count")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":     data,
		"count":    len(data),
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"has_more": offset+limit < total,
	})
}

// GetAnalyticsStats returns summary statistics about the analytics data
func (h *AnalyticsHandler) GetAnalyticsStats(w http.ResponseWriter, r *http.Request) {
	// Ensure DuckDB is initialized
	if err := h.ensureInitialized(r.Context()); err != nil {
		h.logger.Error("Failed to initialize DuckDB", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to initialize database")
		return
	}

	// Get counts from DuckDB
	totalRecords, err := h.duckdbService.GetTotalRecords(r.Context())
	if err != nil {
		h.logger.Error("Failed to get total records", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get total records")
		return
	}

	countryRevenueCount, err := h.duckdbService.GetCountryRevenueCount(r.Context())
	if err != nil {
		h.logger.Error("Failed to get country revenue count", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get country revenue count")
		return
	}

	stats := map[string]interface{}{
		"total_records":         totalRecords,
		"processing_time_ms":    0, // DuckDB queries are fast
		"cache_hit":             false, // Always fresh data
		"country_revenue_count": countryRevenueCount,
		"top_products_count":    20, // Fixed limit
		"monthly_sales_count":   "varies", // Depends on data
		"top_regions_count":     30, // Fixed limit
		"endpoints": map[string]string{
			"country_revenue": "/api/v1/analytics/country-revenue?limit=100&offset=0",
			"top_products":    "/api/v1/analytics/top-products",
			"monthly_sales":   "/api/v1/analytics/monthly-sales",
			"top_regions":     "/api/v1/analytics/top-regions",
		},
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetTopProducts returns top 20 frequently purchased products
func (h *AnalyticsHandler) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	// Ensure DuckDB is initialized
	if err := h.ensureInitialized(r.Context()); err != nil {
		h.logger.Error("Failed to initialize DuckDB", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to initialize database")
		return
	}

	// Get data from DuckDB
	data, err := h.duckdbService.GetTopProducts(r.Context())
	if err != nil {
		h.logger.Error("Failed to get top products", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get top products data")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":  data,
		"count": len(data),
	})
}

// GetMonthlySales returns monthly sales volume data
func (h *AnalyticsHandler) GetMonthlySales(w http.ResponseWriter, r *http.Request) {
	// Ensure DuckDB is initialized
	if err := h.ensureInitialized(r.Context()); err != nil {
		h.logger.Error("Failed to initialize DuckDB", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to initialize database")
		return
	}

	// Get data from DuckDB
	data, err := h.duckdbService.GetMonthlySales(r.Context())
	if err != nil {
		h.logger.Error("Failed to get monthly sales", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get monthly sales data")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":  data,
		"count": len(data),
	})
}

// GetTopRegions returns top 30 regions by revenue
func (h *AnalyticsHandler) GetTopRegions(w http.ResponseWriter, r *http.Request) {
	// Ensure DuckDB is initialized
	if err := h.ensureInitialized(r.Context()); err != nil {
		h.logger.Error("Failed to initialize DuckDB", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to initialize database")
		return
	}

	// Get data from DuckDB
	data, err := h.duckdbService.GetTopRegions(r.Context())
	if err != nil {
		h.logger.Error("Failed to get top regions", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get top regions data")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":  data,
		"count": len(data),
	})
}

// RefreshCache forces a cache refresh by reloading the CSV into DuckDB
func (h *AnalyticsHandler) RefreshCache(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	startTime := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	h.logger.Info("DuckDB refresh requested")

	// Reset initialization flag to force reload
	h.initialized = false

	// Reload CSV into DuckDB
	if err := h.duckdbService.LoadFromCSV(h.csvPath); err != nil {
		h.logger.Error("Failed to refresh DuckDB", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to refresh database")
		return
	}

	h.initialized = true

	// Get record count for stats
	totalRecords, err := h.duckdbService.GetTotalRecords(ctx)
	if err != nil {
		h.logger.Error("Failed to get total records", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get record count")
		return
	}

	h.logger.Info("DuckDB refreshed successfully", "duration", time.Since(startTime))

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message":       "Database refreshed successfully",
		"total_records": totalRecords,
		"duration_ms":   time.Since(startTime).Milliseconds(),
	})
}


func (h *AnalyticsHandler) createAnalyticsSummary(analytics *models.AnalyticsResponse) map[string]interface{} {
	// Limit each section to prevent huge responses
	countryRevenue := analytics.CountryRevenue
	if len(countryRevenue) > 50 {
		countryRevenue = countryRevenue[:50]
	}

	topProducts := analytics.TopProducts
	if len(topProducts) > 20 {
		topProducts = topProducts[:20]
	}

	topRegions := analytics.TopRegions
	if len(topRegions) > 30 {
		topRegions = topRegions[:30]
	}

	// Calculate total revenue from monthly sales
	var totalRevenue float64
	for _, sale := range analytics.MonthlySales {
		totalRevenue += sale.SalesVolume
	}

	return map[string]interface{}{
		"summary": map[string]interface{}{
			"total_records":         analytics.TotalRecords,
			"processing_time_ms":    analytics.ProcessingTimeMs,
			"cache_hit":             analytics.CacheHit,
			"country_revenue_count": len(analytics.CountryRevenue),
			"top_products_count":    len(analytics.TopProducts),
			"monthly_sales_count":   len(analytics.MonthlySales),
			"top_regions_count":     len(analytics.TopRegions),
			"total_revenue":         totalRevenue,
		},
		"country_revenue": countryRevenue,
		"top_products":    topProducts,
		"monthly_sales":   analytics.MonthlySales,
		"top_regions":     topRegions,
		"message":         "Use specific endpoints with pagination for complete data: /api/v1/analytics/country-revenue?limit=100&offset=0",
	}
}

// Helper function to get integer query parameter with default value
func (h *AnalyticsHandler) getIntQueryParam(r *http.Request, key string, defaultValue int) int {
	if value := r.URL.Query().Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil && intValue >= 0 {
			return intValue
		}
	}
	return defaultValue
}
