package config

import (
	"os"
)
//sssdsdasdsadsadsds
// Config holds application configuration
type Config struct {
	DatabasePath  string
	ServerPort    string
	//SessionSecret string
}

// Load reads configuration from environment variables or sets defaults
func Load() *Config {
	return &Config{
		DatabasePath: getEnv("DATABASE_PATH", "./forum.db"),
		ServerPort:   getEnv("SERVER_PORT", ":8080"),
		// SessionSecret: getEnv("SESSION_SECRET", "your-secret-key"),
	}
}

// getEnv returns environment variable value or default if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
