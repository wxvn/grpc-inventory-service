package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env         string
	PostgresURL string
	GRPC        GRPCConfig
}

type GRPCConfig struct {
	Port    int
	Timeout time.Duration
}

func New() *Config {
	return &Config{
		Env:         getEnv("ENV", "local"),
		PostgresURL: getEnv("POSTGRES_URL", ""),

		GRPC: GRPCConfig{
			Port:    getEnvInt("GRPC_PORT", 0),
			Timeout: getEnvDuration("GRPC_TIMEOUT", 5*time.Second),
		},
	}
}

func getEnv(key string, defValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defValue
}

func getEnvInt(key string, defValue int) int {
	valueStr := getEnv(key, "")

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defValue
	}

	return value
}

func getEnvDuration(key string, defValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defValue
	}

	return value
}
