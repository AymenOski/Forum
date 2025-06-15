package config

import (
	"os"
)
type Config struct {
	DatabasePath  string
	ServerPort    string
}

func Load() *Config {
	return &Config{
		DatabasePath: getEnv("DATABASE_PATH", "./forum.db"),
		ServerPort:   getEnv("SERVER_PORT", ":8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
