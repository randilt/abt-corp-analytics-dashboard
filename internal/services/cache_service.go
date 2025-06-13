package services

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"analytics-dashboard-api/internal/models"
	"analytics-dashboard-api/pkg/logger"
)

type CacheService struct {
	logger    logger.Logger
	cacheData *models.AnalyticsResponse
	mu        sync.RWMutex
	cacheTime time.Time
	cacheTTL  time.Duration
}

func NewCacheService(logger logger.Logger) *CacheService {
	return &CacheService{
		logger:   logger,
		cacheTTL: 24 * time.Hour, // Cache for 24 hours
	}
}

// LoadFromCache loads analytics data from cache if valid
func (c *CacheService) LoadFromCache() (*models.AnalyticsResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cacheData == nil || time.Since(c.cacheTime) > c.cacheTTL {
		return nil, false
	}

	// Mark as cache hit
	result := *c.cacheData
	result.CacheHit = true
	return &result, true
}

// SaveToMemory saves analytics data to memory cache
func (c *CacheService) SaveToMemory(data *models.AnalyticsResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheData = data
	c.cacheTime = time.Now()
}

// SaveToFile saves analytics data to file cache
func (c *CacheService) SaveToFile(filePath string, data *models.AnalyticsResponse) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	c.logger.Info("Cache saved to file", "path", filePath, "size_kb", len(jsonData)/1024)
	return nil
}

// LoadFromFile loads analytics data from file cache
func (c *CacheService) LoadFromFile(filePath string) (*models.AnalyticsResponse, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cache file does not exist: %s", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var analytics models.AnalyticsResponse
	if err := json.Unmarshal(data, &analytics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	// Save to memory cache
	c.SaveToMemory(&analytics)

	c.logger.Info("Cache loaded from file", "path", filePath, "records", analytics.TotalRecords)
	return &analytics, nil
}
