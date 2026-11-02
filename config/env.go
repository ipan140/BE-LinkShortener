package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file if it exists.
func LoadEnv() {
	_ = godotenv.Load()
}

// GetEnv gets an environment variable or returns a fallback value.
func GetEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallback
}

// GetEnvInt gets an environment variable as an integer or returns a fallback value.
func GetEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
			return intVal
		}
	}
	return fallback
}

// GetEnvBool gets an environment variable as a boolean or returns a fallback value.
func GetEnvBool(key string, fallback bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		if boolVal, err := strconv.ParseBool(strings.TrimSpace(val)); err == nil {
			return boolVal
		}
	}
	return fallback
}
