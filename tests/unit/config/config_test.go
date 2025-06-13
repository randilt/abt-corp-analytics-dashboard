package config_test

import (
	"os"
	"testing"
	"time"

	"analytics-dashboard-api/internal/config"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear any environment variables
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT",
		"CSV_FILE_PATH", "CSV_BATCH_SIZE", "CSV_WORKER_POOL", "CSV_BUFFER_SIZE",
		"CACHE_FILE_PATH", "CACHE_TTL", "LOG_LEVEL",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Test server defaults
	if cfg.Server.Host != "localhost" {
		t.Errorf("Server.Host = %s, want localhost", cfg.Server.Host)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want 8080", cfg.Server.Port)
	}

	if cfg.Server.ReadTimeout != 15*time.Second {
		t.Errorf("Server.ReadTimeout = %v, want 15s", cfg.Server.ReadTimeout)
	}

	if cfg.Server.WriteTimeout != 15*time.Second {
		t.Errorf("Server.WriteTimeout = %v, want 15s", cfg.Server.WriteTimeout)
	}

	if cfg.Server.IdleTimeout != 60*time.Second {
		t.Errorf("Server.IdleTimeout = %v, want 60s", cfg.Server.IdleTimeout)
	}

	// Test CSV defaults
	if cfg.CSV.FilePath != "./data/raw/transactions.csv" {
		t.Errorf("CSV.FilePath = %s, want ./data/raw/transactions.csv", cfg.CSV.FilePath)
	}

	if cfg.CSV.BatchSize != 10000 {
		t.Errorf("CSV.BatchSize = %d, want 10000", cfg.CSV.BatchSize)
	}

	if cfg.CSV.WorkerPool != 8 {
		t.Errorf("CSV.WorkerPool = %d, want 8", cfg.CSV.WorkerPool)
	}

	if cfg.CSV.BufferSize != 65536 {
		t.Errorf("CSV.BufferSize = %d, want 65536", cfg.CSV.BufferSize)
	}

	// Test Cache defaults
	if cfg.Cache.FilePath != "./data/processed/analytics_cache.json" {
		t.Errorf("Cache.FilePath = %s, want ./data/processed/analytics_cache.json", cfg.Cache.FilePath)
	}

	if cfg.Cache.TTL != 24*time.Hour {
		t.Errorf("Cache.TTL = %v, want 24h", cfg.Cache.TTL)
	}

	// Test Logger defaults
	if cfg.Logger.Level != "info" {
		t.Errorf("Logger.Level = %s, want info", cfg.Logger.Level)
	}
}

func TestLoadConfig_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	testEnvVars := map[string]string{
		"SERVER_HOST":         "0.0.0.0",
		"SERVER_PORT":         "9090",
		"SERVER_READ_TIMEOUT": "30s",
		"CSV_FILE_PATH":       "/custom/path/data.csv",
		"CSV_BATCH_SIZE":      "5000",
		"CSV_WORKER_POOL":     "4",
		"CACHE_TTL":           "12h",
		"LOG_LEVEL":           "debug",
	}

	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Clean up environment variables after test
	defer func() {
		for key := range testEnvVars {
			os.Unsetenv(key)
		}
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Test overridden values
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %s, want 0.0.0.0", cfg.Server.Host)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %d, want 9090", cfg.Server.Port)
	}

	if cfg.Server.ReadTimeout != 30*time.Second {
		t.Errorf("Server.ReadTimeout = %v, want 30s", cfg.Server.ReadTimeout)
	}

	if cfg.CSV.FilePath != "/custom/path/data.csv" {
		t.Errorf("CSV.FilePath = %s, want /custom/path/data.csv", cfg.CSV.FilePath)
	}

	if cfg.CSV.BatchSize != 5000 {
		t.Errorf("CSV.BatchSize = %d, want 5000", cfg.CSV.BatchSize)
	}

	if cfg.CSV.WorkerPool != 4 {
		t.Errorf("CSV.WorkerPool = %d, want 4", cfg.CSV.WorkerPool)
	}

	if cfg.Cache.TTL != 12*time.Hour {
		t.Errorf("Cache.TTL = %v, want 12h", cfg.Cache.TTL)
	}

	if cfg.Logger.Level != "debug" {
		t.Errorf("Logger.Level = %s, want debug", cfg.Logger.Level)
	}
}

func TestLoadConfig_InvalidValues(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
	}{
		{
			name: "invalid port negative",
			envVars: map[string]string{
				"SERVER_PORT": "-1",
			},
			expectError: true,
		},
		{
			name: "invalid port zero",
			envVars: map[string]string{
				"SERVER_PORT": "0",
			},
			expectError: true,
		},
		{
			name: "invalid port too high",
			envVars: map[string]string{
				"SERVER_PORT": "65536",
			},
			expectError: true,
		},
		{
			name: "empty csv file path",
			envVars: map[string]string{
				"CSV_FILE_PATH": "",
			},
			expectError: true,
		},
		{
			name: "invalid batch size",
			envVars: map[string]string{
				"CSV_BATCH_SIZE": "0",
			},
			expectError: true,
		},
		{
			name: "invalid worker pool",
			envVars: map[string]string{
				"CSV_WORKER_POOL": "0",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Clearenv()

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			config, err := config.LoadConfig()

			if tt.expectError {
				if err == nil {
					t.Errorf("LoadConfig() should return error for %s = %s", tt.name, tt.envVars)
				}
			} else {
				if err != nil {
					t.Errorf("LoadConfig() returned unexpected error: %v", err)
				}
				if config == nil {
					t.Error("LoadConfig() returned nil config")
				}
			}
		})
	}
}

func TestLoadConfig_InvalidEnvironmentValues_FallbackToDefaults(t *testing.T) {
	// Test that invalid environment values fall back to defaults for non-critical configs
	testCases := []struct {
		name     string
		envVar   string
		value    string
		checkFn  func(*config.Config) bool
		errorMsg string
	}{
		{
			name:   "invalid duration falls back to default",
			envVar: "SERVER_READ_TIMEOUT",
			value:  "invalid_duration",
			checkFn: func(cfg *config.Config) bool {
				return cfg.Server.ReadTimeout == 15*time.Second
			},
			errorMsg: "Should fallback to default 15s for invalid duration",
		},
		{
			name:   "invalid int falls back to default",
			envVar: "CSV_BUFFER_SIZE",
			value:  "not_a_number",
			checkFn: func(cfg *config.Config) bool {
				return cfg.CSV.BufferSize == 65536
			},
			errorMsg: "Should fallback to default 65536 for invalid int",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear environment first
			os.Unsetenv(tc.envVar)

			// Set invalid value
			os.Setenv(tc.envVar, tc.value)
			defer os.Unsetenv(tc.envVar)

			cfg, err := config.LoadConfig()
			if err != nil {
				t.Fatalf("LoadConfig() should not error for non-critical invalid values: %v", err)
			}

			if !tc.checkFn(cfg) {
				t.Error(tc.errorMsg)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				Server: config.ServerConfig{
					Port: 8080,
				},
				CSV: config.CSVConfig{
					FilePath:   "data.csv",
					BatchSize:  1000,
					WorkerPool: 4,
				},
				Cache: config.CacheConfig{
					TTL: 5 * time.Minute,
				},
			},
			expectError: false,
		},
		{
			name: "invalid port negative",
			config: &config.Config{
				Server: config.ServerConfig{
					Port: -1,
				},
				CSV: config.CSVConfig{
					FilePath:   "data.csv",
					BatchSize:  1000,
					WorkerPool: 4,
				},
				Cache: config.CacheConfig{
					TTL: 5 * time.Minute,
				},
			},
			expectError: true,
		},
		{
			name: "invalid port too high",
			config: &config.Config{
				Server: config.ServerConfig{
					Port: 65536,
				},
				CSV: config.CSVConfig{
					FilePath:   "data.csv",
					BatchSize:  1000,
					WorkerPool: 4,
				},
				Cache: config.CacheConfig{
					TTL: 5 * time.Minute,
				},
			},
			expectError: true,
		},
		{
			name: "empty csv file path",
			config: &config.Config{
				Server: config.ServerConfig{
					Port: 8080,
				},
				CSV: config.CSVConfig{
					FilePath:   "",
					BatchSize:  1000,
					WorkerPool: 4,
				},
				Cache: config.CacheConfig{
					TTL: 5 * time.Minute,
				},
			},
			expectError: true,
		},
		{
			name: "invalid batch size",
			config: &config.Config{
				Server: config.ServerConfig{
					Port: 8080,
				},
				CSV: config.CSVConfig{
					FilePath:   "data.csv",
					BatchSize:  0,
					WorkerPool: 4,
				},
				Cache: config.CacheConfig{
					TTL: 5 * time.Minute,
				},
			},
			expectError: true,
		},
		{
			name: "invalid worker pool",
			config: &config.Config{
				Server: config.ServerConfig{
					Port: 8080,
				},
				CSV: config.CSVConfig{
					FilePath:   "data.csv",
					BatchSize:  1000,
					WorkerPool: 0,
				},
				Cache: config.CacheConfig{
					TTL: 5 * time.Minute,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("Config validation expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Config validation returned unexpected error: %v", err)
				}
			}
		})
	}
}
