package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv    string
	AppPort   int
	LogLevel  string
	TestDBURL string

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// DSN builds a PostgreSQL connection string from the DB config.
func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	appPort, err := getInt("APP_PORT", 8080)
	if err != nil {
		return nil, err
	}
	dbPort, err := getInt("DB_PORT", 5432) //empty check
	if err != nil {
		return nil, err
	}
	return &Config{
		AppEnv:     getString("APP_ENV", "development"),
		AppPort:    appPort,
		LogLevel:   getString("LOG_LEVEL", "info"),
		TestDBURL:  getString("TEST_DATABASE_URL", ""),
		DBHost:     getString("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getString("DB_USER", "postgres"),
		DBPassword: getString("DB_PASSWORD", ""),
		DBName:     getString("DB_NAME", "postgres"),
		DBSSLMode:  getString("DB_SSLMODE", "disable"),
	}, nil
}

func getString(key, def string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return def
}

func getInt(key string, def int) (int, error) {
	v := getString(key, "")
	if v == "" {
		return def, nil
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be an integer: %w", key, err)
	}
	return i, nil
}
