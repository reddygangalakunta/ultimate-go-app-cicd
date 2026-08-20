package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all server and application configurations.
type Config struct {
	AppName         string
	Environment     string
	Port            int
	LogLevel        string
	Version         string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

// LoadConfig reads configuration from environment variables with sensible defaults.
func LoadConfig(version string) (*Config, error) {
	port, err := getEnvAsInt("PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT environment variable: %w", err)
	}

	shutdownTimeoutSec, err := getEnvAsInt("SHUTDOWN_TIMEOUT_SECONDS", 15)
	if err != nil {
		return nil, fmt.Errorf("invalid SHUTDOWN_TIMEOUT_SECONDS: %w", err)
	}

	cfg := &Config{
		AppName:         getEnv("APP_NAME", "order-service"),
		Environment:     getEnv("ENVIRONMENT", "production"),
		Port:            port,
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		Version:         version,
		ShutdownTimeout: time.Duration(shutdownTimeoutSec) * time.Second,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     60 * time.Second,
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) (int, error) {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback, nil
	}
	val, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, err
	}
	return val, nil
}
