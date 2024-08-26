// internal/common/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server
	GRPCAddr string
	HTTPAddr string

	// PostgreSQL
	PostgresURI string

	// MongoDB
	MongoURI string
	MongoDB  string

	// App
	Debug           bool
	ShutdownTimeout time.Duration

	// Auth (for future use)
	JWTSecret string
}

func Load() (*Config, error) {
	config := &Config{
		GRPCAddr:        getEnv("GRPC_ADDR", ":50051"),
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		PostgresURI:     getEnv("POSTGRES_URI", "postgresql://user:password@localhost:5432/devicedb?sslmode=disable"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:         getEnv("MONGO_DB", "devicedb"),
		Debug:           getEnvAsBool("DEBUG", false),
		ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.PostgresURI == "" {
		return fmt.Errorf("POSTGRES_URI is required")
	}
	if c.MongoURI == "" {
		return fmt.Errorf("MONGO_URI is required")
	}
	if c.MongoDB == "" {
		return fmt.Errorf("MONGO_DB is required")
	}
	return nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return defaultValue
}
