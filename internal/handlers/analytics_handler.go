package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"analytics-dashboard-api/internal/models"
	"analytics-dashboard-api/internal/services"
	"analytics-dashboard-api/internal/utils"
	"analytics-dashboard-api/pkg/logger"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	cacheService     *services.CacheService
	csvProcessor     *services.CSVProcessor
	logger           logger.Logger
	csvPath          string
	cachePath        string
}

func NewAnalyticsHandler(
	analyticsService *services.AnalyticsService,
	cacheService *services.CacheService,
	csvProcessor *services.CSVProcessor,
	logger logger.Logger,
	csvPath, cachePath string,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		cacheService:     cacheService,
		csvProcessor:     csvProcessor,
		logger:           logger,
		csvPath:          csvPath,
		cachePath:        cachePath,
	}
}

// GetAnalytics returns all dashboard analytics data
func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	ctx := r.Context()

	h.logger.Info("Analytics request received", "method", r.Method, "path", r.URL.Path)

	// Try loading from memory cache first
	if analytics, hit := h.cacheService.LoadFromCache(); hit {
		h.logger.Info("Serving from memory cache", "duration", time.Since(startTime))

		summary := h.createAnalyticsSummary(analytics)
		utils.WriteJSONResponse(w, http.StatusOK, summary)
		return
	}

	// Try loading from file cache
	if analytics, err := h.cacheService.LoadFromFile(h.cachePath); err == nil {
		h.logger.Info("Serving from file cache", "duration", time.Since(startTime))
		analytics.CacheHit = true

		summary := h.createAnalyticsSummary(analytics)
		utils.WriteJSONResponse(w, http.StatusOK, summary)
		return
	}

	// Process CSV if no cache
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result, err := h.csvProcessor.ProcessLargeCSV(ctx, h.csvPath)
	if err != nil {
		h.logger.Error("Failed to process CSV", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process data")
		return
	}

	analytics := h.analyticsService.GenerateAnalytics(result.Transactions)
	
	// Cache the results
	h.cacheService.SaveToMemory(analytics)
	go func() {
		if err := h.cacheService.SaveToFile(h.cachePath, analytics); err != nil {
			h.logger.Error("Failed to save cache to file", "error", err)
		}
	}()

	h.logger.Info("Analytics generated successfully", 
		"records", len(result.Transactions), 
		"duration", time.Since(startTime))

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

	analytics, err := h.getAnalyticsData(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get analytics data")
		return
	}

	// Apply pagination
	total := len(analytics.CountryRevenue)
	start := offset
	end := offset + limit
	
	if start >= total {
		start = total
		end = total
	} else if end > total {
		end = total
	}

	paginatedData := analytics.CountryRevenue[start:end]

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":   paginatedData,
		"count":  len(paginatedData),
		"total":  total,
		"limit":  limit,
		"offset": offset,
		"has_more": end < total,
	})
}

// GetAnalyticsStats returns summary statistics about the analytics data
func (h *AnalyticsHandler) GetAnalyticsStats(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.getAnalyticsData(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get analytics data")
		return
	}

	stats := map[string]interface{}{
		"total_records":         analytics.TotalRecords,
		"processing_time_ms":    analytics.ProcessingTimeMs,
		"cache_hit":            analytics.CacheHit,
		"country_revenue_count": len(analytics.CountryRevenue),
		"top_products_count":    len(analytics.TopProducts),
		"monthly_sales_count":   len(analytics.MonthlySales),
		"top_regions_count":     len(analytics.TopRegions),
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
	analytics, err := h.getAnalyticsData(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get analytics data")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data": analytics.TopProducts,
		"count": len(analytics.TopProducts),
	})
}

// GetMonthlySales returns monthly sales volume data
func (h *AnalyticsHandler) GetMonthlySales(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.getAnalyticsData(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get analytics data")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data": analytics.MonthlySales,
		"count": len(analytics.MonthlySales),
	})
}

// GetTopRegions returns top 30 regions by revenue
func (h *AnalyticsHandler) GetTopRegions(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.getAnalyticsData(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get analytics data")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data": analytics.TopRegions,
		"count": len(analytics.TopRegions),
	})
}

// RefreshCache forces a cache refresh by reprocessing the CSV
func (h *AnalyticsHandler) RefreshCache(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	startTime := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	h.logger.Info("Cache refresh requested")

	// Process CSV and update cache
	stats, err := h.csvProcessor.PreprocessAndCache(ctx, h.csvPath, h.cachePath)
	if err != nil {
		h.logger.Error("Failed to refresh cache", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to refresh cache")
		return
	}

	h.logger.Info("Cache refreshed successfully", "duration", time.Since(startTime))

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Cache refreshed successfully",
		"stats": stats,
		"duration_ms": time.Since(startTime).Milliseconds(),
	})
}

// getAnalyticsData is a helper method to get analytics data with caching
func (h *AnalyticsHandler) getAnalyticsData(ctx context.Context) (*models.AnalyticsResponse, error) {
	// Try cache first
	if analytics, hit := h.cacheService.LoadFromCache(); hit {
		return analytics, nil
	}

	// Try file cache
	if analytics, err := h.cacheService.LoadFromFile(h.cachePath); err == nil {
		return analytics, nil
	}

	// Process CSV as fallback
	result, err := h.csvProcessor.ProcessLargeCSV(ctx, h.csvPath)
	if err != nil {
		return nil, err
	}

	analytics := h.analyticsService.GenerateAnalytics(result.Transactions)
	h.cacheService.SaveToMemory(analytics)

	return analytics, nil
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

	return map[string]interface{}{
		"summary": map[string]interface{}{
			"total_records":              analytics.TotalRecords,
			"processing_time_ms":         analytics.ProcessingTimeMs,
			"cache_hit":                  analytics.CacheHit,
			"country_revenue_count":      len(analytics.CountryRevenue),
			"top_products_count":         len(analytics.TopProducts),
			"monthly_sales_count":        len(analytics.MonthlySales),
			"top_regions_count":          len(analytics.TopRegions),
		},
		"country_revenue":  countryRevenue,
		"top_products":     topProducts,
		"monthly_sales":    analytics.MonthlySales, // Usually not too large
		"top_regions":      topRegions,
		"message": "Use specific endpoints with pagination for complete data: /api/v1/analytics/country-revenue?limit=100&offset=0",
	}
}

func (h *AnalyticsHandler) getIntQueryParam(r *http.Request, key string, defaultValue int) int {
	if value := r.URL.Query().Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil && intValue >= 0 {
			return intValue
		}
	}
	return defaultValue
}
