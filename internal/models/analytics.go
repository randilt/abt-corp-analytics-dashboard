package models

import (
	"errors"
	"time"
)

var (
	ErrInvalidCSVRow = errors.New("invalid CSV row format")
)

// CountryRevenue represents revenue data by country and product
type CountryRevenue struct {
	Country          string  `json:"country"`
	ProductName      string  `json:"product_name"`
	TotalRevenue     float64 `json:"total_revenue"`
	TransactionCount int     `json:"transaction_count"`
}

// ProductFrequency represents frequently purchased products
type ProductFrequency struct {
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	PurchaseCount int    `json:"purchase_count"`
	StockQuantity int    `json:"current_stock"`
}

// MonthlySales represents sales volume by month
type MonthlySales struct {
	Month       string  `json:"month"`
	SalesVolume float64 `json:"sales_volume"`
	ItemCount   int     `json:"item_count"`
}

// RegionRevenue represents revenue data by region
type RegionRevenue struct {
	Region       string  `json:"region"`
	TotalRevenue float64 `json:"total_revenue"`
	ItemsSold    int     `json:"items_sold"`
}

// AnalyticsResponse wraps all dashboard data
type AnalyticsResponse struct {
	CountryRevenue   []CountryRevenue   `json:"country_revenue"`
	TopProducts      []ProductFrequency `json:"top_products"`
	MonthlySales     []MonthlySales     `json:"monthly_sales"`
	TopRegions       []RegionRevenue    `json:"top_regions"`
	ProcessingTimeMs int64              `json:"processing_time_ms"`
	TotalRecords     int                `json:"total_records"`
	CacheHit         bool               `json:"cache_hit"`
}

// ProcessingStats holds statistics about data processing
type ProcessingStats struct {
	TotalRecords     int           `json:"total_records"`
	ProcessedRecords int           `json:"processed_records"`
	ErrorCount       int           `json:"error_count"`
	ProcessingTime   time.Duration `json:"processing_time"`
	MemoryUsageMB    float64       `json:"memory_usage_mb"`
}
