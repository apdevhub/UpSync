package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port              string
	SupabaseURL       string
	SupabaseServiceKey string
	SupabaseBucket    string
	FrontendURL       string
	MaxFileSizeBytes  int64
}

// Load reads configuration from environment variables and validates required fields.
// Returns an error if any required variable is missing.
func Load() (*Config, error) {
	cfg := &Config{
		Port:             getEnv("PORT", "8080"),
		SupabaseURL:      os.Getenv("SUPABASE_URL"),
		SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		SupabaseBucket:   getEnv("SUPABASE_BUCKET", "upsync-files"),
		FrontendURL:      getEnv("FRONTEND_URL", "http://localhost:5173"),
		MaxFileSizeBytes: parseSize(getEnv("MAX_FILE_SIZE_MB", "50")),
	}

	if cfg.SupabaseURL == "" {
		return nil, fmt.Errorf("SUPABASE_URL is required")
	}
	if cfg.SupabaseServiceKey == "" {
		return nil, fmt.Errorf("SUPABASE_SERVICE_ROLE_KEY is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseSize(mb string) int64 {
	n, err := strconv.ParseInt(mb, 10, 64)
	if err != nil || n <= 0 {
		return 50
	}
	return n * 1024 * 1024
}
