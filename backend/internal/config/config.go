package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents the validated application configuration.
type Config struct {
	AppEnv             string
	HTTPPort           int
	DatabaseURL        string
	CORSAllowedOrigins []string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
	DBMaxOpenConns     int
	DBMinOpenConns     int
	DBMaxIdleTime      time.Duration
	LogLevel           string
	LogFormat          string
}

// Load reads and validates configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		HTTPPort:        getEnvAsInt("HTTP_PORT", 8080),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://app_user:app_password@localhost:5432/app_db?sslmode=disable"),
		ReadTimeout:     time.Duration(getEnvAsInt("READ_TIMEOUT_SECONDS", 10)) * time.Second,
		WriteTimeout:    time.Duration(getEnvAsInt("WRITE_TIMEOUT_SECONDS", 10)) * time.Second,
		IdleTimeout:     time.Duration(getEnvAsInt("IDLE_TIMEOUT_SECONDS", 120)) * time.Second,
		ShutdownTimeout: time.Duration(getEnvAsInt("SHUTDOWN_TIMEOUT_SECONDS", 15)) * time.Second,
		DBMaxOpenConns:  getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMinOpenConns:  getEnvAsInt("DB_MIN_OPEN_CONNS", 5),
		DBMaxIdleTime:   time.Duration(getEnvAsInt("DB_MAX_IDLE_TIME_MINUTES", 15)) * time.Minute,
		LogLevel:        strings.ToLower(getEnv("LOG_LEVEL", "info")),
		LogFormat:       strings.ToLower(getEnv("LOG_FORMAT", "text")),
	}

	// Parse CORS allowed origins
	corsRaw := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000")
	origins := strings.Split(corsRaw, ",")
	cfg.CORSAllowedOrigins = make([]string, 0, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			cfg.CORSAllowedOrigins = append(cfg.CORSAllowedOrigins, trimmed)
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate checks that critical configuration parameters are within bounds.
func (c *Config) Validate() error {
	if c.HTTPPort < 1 || c.HTTPPort > 65535 {
		return fmt.Errorf("invalid HTTP_PORT: %d (must be between 1 and 65535)", c.HTTPPort)
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL must not be empty")
	}

	if _, err := url.ParseRequestURI(c.DatabaseURL); err != nil {
		return fmt.Errorf("invalid DATABASE_URL URI format: %w", err)
	}

	if c.DBMaxOpenConns <= 0 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS must be greater than 0")
	}

	if c.DBMinOpenConns < 0 || c.DBMinOpenConns > c.DBMaxOpenConns {
		return fmt.Errorf("DB_MIN_OPEN_CONNS must be >= 0 and <= DB_MAX_OPEN_CONNS")
	}

	switch c.LogLevel {
	case "debug", "info", "warn", "error":
		// valid
	default:
		return fmt.Errorf("invalid LOG_LEVEL: %q (expected debug, info, warn, or error)", c.LogLevel)
	}

	switch c.LogFormat {
	case "json", "text":
		// valid
	default:
		return fmt.Errorf("invalid LOG_FORMAT: %q (expected json or text)", c.LogFormat)
	}

	return nil
}

// IsProduction returns true if running in production mode.
func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "production")
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
