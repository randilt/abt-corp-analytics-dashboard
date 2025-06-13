package handlers

import (
	"net/http"
	"runtime"
	"time"

	"analytics-dashboard-api/internal/utils"
	"analytics-dashboard-api/pkg/logger"
)

type HealthHandler struct {
	logger    logger.Logger
	startTime time.Time
}

func NewHealthHandler(logger logger.Logger) *HealthHandler {
	return &HealthHandler{
		logger:    logger,
		startTime: time.Now(),
	}
}

// Health returns service health status
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	health := map[string]interface{}{
		"status":     "healthy",
		"timestamp":  time.Now().UTC(),
		"uptime":     time.Since(h.startTime).String(),
		"version":    "1.0.0",
		"memory": map[string]interface{}{
			"alloc_mb":      float64(memStats.Alloc) / 1024 / 1024,
			"total_alloc_mb": float64(memStats.TotalAlloc) / 1024 / 1024,
			"sys_mb":        float64(memStats.Sys) / 1024 / 1024,
			"num_gc":        memStats.NumGC,
		},
		"goroutines": runtime.NumGoroutine(),
	}

	utils.WriteJSONResponse(w, http.StatusOK, health)
}

// Ready returns readiness status
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ready := map[string]interface{}{
		"status": "ready",
		"timestamp": time.Now().UTC(),
	}

	utils.WriteJSONResponse(w, http.StatusOK, ready)
}