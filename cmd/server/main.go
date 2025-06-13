package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"analytics-dashboard-api/internal/config"
	"analytics-dashboard-api/internal/handlers"
	"analytics-dashboard-api/internal/middleware"
	"analytics-dashboard-api/internal/services"
	"analytics-dashboard-api/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.NewLogger(cfg.Logger.Level)
	log.Info("Starting analytics dashboard server", "version", "1.0.0")
	// Initialize services
	csvProcessor := services.NewCSVProcessor(log, &cfg.CSV, &cfg.Cache)
	analyticsService := services.NewAnalyticsService(log)
	cacheService := services.NewCacheService(log, &cfg.Cache)
	// Initialize handlers
	analyticsHandler := handlers.NewAnalyticsHandler(
		analyticsService,
		cacheService,
		csvProcessor,
		log,
		cfg.CSV.FilePath,
		cfg.Cache.FilePath,
	)
	healthHandler := handlers.NewHealthHandler(log)

	// Setup router
	router := setupRouter(analyticsHandler, healthHandler, log)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("Server starting", "address", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal or server error

	select {
	case err := <-serverErrors:
		log.Error("Server failed to start", "error", err)
		os.Exit(1)

	case sig := <-interrupt:
		log.Info("Server shutdown initiated", "signal", sig.String())

		// Give outstanding requests 30 seconds to complete

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error("Server shutdown failed", "error", err)
			if err := server.Close(); err != nil {
				log.Error("Server force close failed", "error", err)
			}
			os.Exit(1)
		}
	}

	log.Info("Server shutdown completed")
}

func setupRouter(
	analyticsHandler *handlers.AnalyticsHandler,
	healthHandler *handlers.HealthHandler,
	log logger.Logger,
) *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.Recovery(log))
	router.Use(middleware.Logging(log))
	router.Use(middleware.CORS)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Analytics endpoints
	api.HandleFunc("/analytics", analyticsHandler.GetAnalytics).Methods("GET")
	api.HandleFunc("/analytics/stats", analyticsHandler.GetAnalyticsStats).Methods("GET")
	api.HandleFunc("/analytics/country-revenue", analyticsHandler.GetCountryRevenue).Methods("GET")
	api.HandleFunc("/analytics/top-products", analyticsHandler.GetTopProducts).Methods("GET")
	api.HandleFunc("/analytics/monthly-sales", analyticsHandler.GetMonthlySales).Methods("GET")
	api.HandleFunc("/analytics/top-regions", analyticsHandler.GetTopRegions).Methods("GET")
	api.HandleFunc("/analytics/refresh", analyticsHandler.RefreshCache).Methods("POST")

	// Health endpoints
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")
	router.HandleFunc("/ready", healthHandler.Ready).Methods("GET")

	return router
}
