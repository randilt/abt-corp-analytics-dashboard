package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	CSV      CSVConfig
	Cache    CacheConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type CSVConfig struct {
	FilePath    string
	BatchSize   int
	WorkerPool  int
	BufferSize  int
}

type CacheConfig struct {
	FilePath string
	TTL      time.Duration
}

type LoggerConfig struct {
	Level string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", "15s"),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", "15s"),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", "60s"),
		},
		CSV: CSVConfig{
			FilePath:   getEnv("CSV_FILE_PATH", "./data/raw/transactions.csv"),
			BatchSize:  getEnvAsInt("CSV_BATCH_SIZE", 10000),
			WorkerPool: getEnvAsInt("CSV_WORKER_POOL", 8), // reduce this if resource usage becomes an issue
			BufferSize: getEnvAsInt("CSV_BUFFER_SIZE", 65536),
		},
		Cache: CacheConfig{
			FilePath: getEnv("CACHE_FILE_PATH", "./data/processed/analytics_cache.json"),
			TTL:      getEnvAsDuration("CACHE_TTL", "24h"),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.CSV.FilePath == "" {
		return fmt.Errorf("CSV file path is required")
	}

	if c.CSV.BatchSize <= 0 {
		return fmt.Errorf("CSV batch size must be positive")
	}

	if c.CSV.WorkerPool <= 0 {
		return fmt.Errorf("CSV worker pool size must be positive")
	}

	return nil
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
