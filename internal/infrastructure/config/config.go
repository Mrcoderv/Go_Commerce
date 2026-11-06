package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv    string
	AppPort   int
	LogLevel  string
	JWTSecret string

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func (c Config) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   fmt.Sprintf("%s:%d", c.DBHost, c.DBPort),
		Path:   c.DBName,
	}
	q := u.Query()
	q.Set("sslmode", c.DBSSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

func (c Config) Validate() error {
	if c.AppPort <= 0 || c.AppPort > 65535 {
		return fmt.Errorf("config: invalid APP_PORT %d", c.AppPort)
	}
	if c.DBPort <= 0 || c.DBPort > 65535 {
		return fmt.Errorf("config: invalid DB_PORT %d", c.DBPort)
	}
	if c.DBHost == "" {
		return errors.New("config: DB_HOST is required")
	}
	if c.DBUser == "" {
		return errors.New("config: DB_USER is required")
	}
	if c.DBName == "" {
		return errors.New("config: DB_NAME is required")
	}
	if c.AppEnv == "production" && c.DBPassword == "" {
		return errors.New("config: DB_PASSWORD is required in production")
	}
	if c.AppEnv == "production" && c.JWTSecret == "" {
		return errors.New("config: JWT_SECRET is required in production")
	}
	return nil
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	appPort, err := getInt("APP_PORT", 8080)
	if err != nil {
		return nil, err
	}
	dbPort, err := getInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppEnv:     getString("APP_ENV", "development"),
		AppPort:    appPort,
		LogLevel:   getString("LOG_LEVEL", "info"),
		JWTSecret:  getString("JWT_SECRET", "dev-secret-change-me"),
		DBHost:     getString("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getString("DB_USER", "postgres"),
		DBPassword: getString("DB_PASSWORD", ""),
		DBName:     getString("DB_NAME", "postgres"),
		DBSSLMode:  getString("DB_SSLMODE", "disable"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
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
