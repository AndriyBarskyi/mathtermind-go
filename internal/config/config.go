package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
// We're using struct tags to map environment variables to our config
// The `env` tag tells us which environment variable to use
// The `default` tag provides a default value if the env var is not set
type Config struct {
	Server struct {
		Port string `env:"SERVER_PORT" default:"8080"`
	}
	Database struct {
		URL string `env:"DATABASE_URL"`
	}
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	var cfg Config

	_ = godotenv.Load()

	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Database.URL = getEnv("DATABASE_URL", "")

	if cfg.Server.Port == "" {
		return nil, fmt.Errorf("server port cannot be empty")
	}

	return &cfg, nil
}

// getEnv is a helper function that reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
