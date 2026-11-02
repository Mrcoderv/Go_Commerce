package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv    string
	AppPort   string
	LogLevel  string
	TestDBURL string

	DB_Host     string
	DB_Port     int
	DB_User     string
	DB_Password string
	DB_Name     string
	DB_SSLMode  string
}

// DSN builds a PostgreSQL connection string from the DB config.
func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.DB_User, c.DB_Password, c.DB_Host, c.DB_Port, c.DB_Name, c.DB_SSLMode)
}

func Load() (*Config, error) {
	appPort, err := getInt("PORT", 8080)
	if err != nil {
		return nil, err
	}
	dbPort, err := getInt("DB_PORT", 5432) //empty check
	if err != nil {
		return nil, err
	}
	return &Config{
		AppEnv:      getString("APP_ENV", "development"),
		AppPort:     appPort,
		LogLevel:    getString("LOG_LEVEL", "info"),
		TestDBURL:   getString("TEST_DB_URL", ""),
		DB_Host:     getString("DB_HOST", "localhost"),
		DB_Port:     dbPort,
		DB_User:     getString("DB_USER", "postgres"),
		DB_Password: getString("DB_PASSWORD", ""),
		DB_Name:     getString("DB_NAME", "postgres"),
		DB_SSLMode:  getString("SSL_MODE", "disable"),
	}, nil
}
func LoadConfig() (*Config, error) {
	return Load()
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

func getDuration(key string, def time.Duration) (time.Duration, error) {
	v := getString(key, "")
	if v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be a valid duration: %w", key, err)
	}
	return d, nil
}
