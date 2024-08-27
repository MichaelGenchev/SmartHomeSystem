package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      AppConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	GatewayPort    string
	UserGRPCAddr   string
	UserHTTPAddr   string
	DeviceGRPCAddr string
	DeviceHTTPAddr string
}

type DatabaseConfig struct {
	UserPostgresURI   string
	DevicePostgresURI string
	DeviceMongoURI    string
	DeviceMongoDB     string
}

type AppConfig struct {
	Debug           bool
	ShutdownTimeout time.Duration
}

type AuthConfig struct {
	JWTSecret string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			UserGRPCAddr:   getEnv("USER_GRPC_ADDR", ":50051"),
			UserHTTPAddr:   getEnv("USER_HTTP_ADDR", ":8080"),
			DeviceGRPCAddr: getEnv("DEVICE_GRPC_ADDR", ":50052"),
			DeviceHTTPAddr: getEnv("DEVICE_HTTP_ADDR", ":8081"),
		},
		Database: DatabaseConfig{
			UserPostgresURI:   getEnv("USER_POSTGRES_URI", "postgresql://user:password@localhost:5432/userdb?sslmode=disable"),
			DevicePostgresURI: getEnv("DEVICE_POSTGRES_URI", "postgresql://user:password@localhost:5432/devicedb?sslmode=disable"),
			DeviceMongoURI:    getEnv("DEVICE_MONGO_URI", "mongodb://localhost:27017"),
			DeviceMongoDB:     getEnv("DEVICE_MONGO_DB", "devicedb"),
		},
		App: AppConfig{
			Debug:           getEnvAsBool("DEBUG", false),
			ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		},
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.Database.UserPostgresURI == "" {
		return fmt.Errorf("USER_POSTGRES_URI is required")
	}
	if c.Database.DevicePostgresURI == "" {
		return fmt.Errorf("DEVICE_POSTGRES_URI is required")
	}
	if c.Database.DeviceMongoURI == "" {
		return fmt.Errorf("DEVICE_MONGO_URI is required")
	}
	if c.Database.DeviceMongoDB == "" {
		return fmt.Errorf("DEVICE_MONGO_DB is required")
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
