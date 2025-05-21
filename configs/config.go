package configs

import (
	"os"
	"path/filepath"
)

type Config struct {
	ServerAddress string

	DBPath string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8088"),
		DBPath:        getEnv("DB_PATH", filepath.Join("data", "rockets.db")),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
