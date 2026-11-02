package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"github.com/spf13/viper\"
	"time"
)

type Config struct {
	AppEnv    string
	AppPort   int
	LogLevel  string
	DB        DBConfig
	TestDBURL string
	
 }

type DBConfig struct {
	DB_Host                  string
	DB_Port                  int
	DB_User                  string
	DB_Password              string
	DB_Name                  string
	DB_SSLMode               string
	Min_conns                int
	Max_connections          int
	Max_connection_lifetime time.Duration
	Max_idle                 time.Duration
}

func Load() (*Config, error) {
	appPort, err := getInt("PORT", 8080)
	if err != nil {
		return nil, err
	}
	dbPort, err := getInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}
	minConns, err := getInt("DB_MIN_CONNS", 1)
	if err != nil {
		return nil, err
	}
	maxConns, err := getInt("DB_MAX_CONNS", 10)
	if err != nil {
		return nil, err
	}
	lifetime, err := getDuration("DB_MAX_CONN_LIFETIME", time.Hour)
	if err != nil {
		return nil, err
	}
	idle, err := getDuration("DB_MAX_IDLE_TIME", 30*time.Minute)
	if err != nil {
		return nil, err
	}

	return &Config{
		AppEnv:    getString("APP_ENV", "development"),
		AppPort:   appPort,
		LogLevel:  getString("LOG_LEVEL", "info"),
		TestDBURL: getString("TEST_DB_URL", ""),
		DB: DBConfig{
			DB_Host:                  getString("DB_HOST", "localhost"),
			DB_Port:                  dbPort,
			DB_User:                  getString("DB_USER", "postgres"),
			DB_Password:              getString("DB_PASSWORD", ""),
			DB_Name:                  getString("DB_NAME", "postgres"),
			DB_SSLMode:               getString("SSL_MODE", "disable"),
			Min_conns:                minConns,
			Max_connections:          maxConns,
			Max_connection_lifetime: lifetime,
			Max_idle:                 idle,
			DatabaseURL:              getString("DATABASE_URL", ""),
		},
	}, nil
}

// LoadConfig is kept as an explicit alias for callers that prefer that name.
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
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def, nil
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be an integer: %w", key, err)
	}
	return i, nil
}

func getInt32(key string, def int32) (int32, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def, nil
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be an integer: %w", key, err)
	}
	return int32(i), nil
}

func getDuration(key string, def time.Duration) (time.Duration, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be a valid duration: %w", key, err)
	}
	return d, nil
}


}

type DBConfig struct {
	viper.SetConfigFile `json:"viper_._set_config_file,omitempty"`(".env")
	viper.AutomaticEnv()

	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")

	AppEnv    string
	AppPort   int
	LogLevel  string
	DB        DBConfig
	TestDBURL string
	Env:         viper.GetString("APP_ENV"),
	Port:        viper.GetString("PORT"),
	Host:        viper.GetString("DB_HOST"),
	DB_Host                 string
	DB_Port                 int
	DB_User                 string
	DB_Password             string
	DB_Name                 string
	DB_SSLMode              string
	Min_conns               int
	Max_connections         int
	Max_connection_lifetime time.Duration
	Max_idle                 time.Duration
	DatabaseURL             string
}
		return strings.TrimSpace(v)
	}
	appPort, err := getInt("PORT", 8080)
	if err != nil {
		return nil, err
	}
	dbPort, err := getInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}
	minConns, err := getInt("DB_MIN_CONNS", 1)
	if err != nil {
		return nil, err
	}
	maxConns, err := getInt("DB_MAX_CONNS", 10)
	if err != nil {
		return nil, err
	}
	lifetime, err := getDuration("DB_MAX_CONN_LIFETIME", time.Hour)
	if err != nil {
		return nil, err
	}
	idle, err := getDuration("DB_MAX_IDLE_TIME", 30*time.Minute)
	if err != nil {
		return nil, err
	}

	return &Config{
		AppEnv:    getString("APP_ENV", "development"),
		AppPort:   appPort,
		LogLevel:  getString("LOG_LEVEL", "info"),
		TestDBURL: getString("TEST_DB_URL", ""),
		DB: DBConfig{
			DB_Host:                  getString("DB_HOST", "localhost"),
			DB_Port:                  dbPort,
			DB_User:                  getString("DB_USER", "postgres"),
			DB_Password:              getString("DB_PASSWORD", ""),
			DB_Name:                  getString("DB_NAME", "postgres"),
			DB_SSLMode:               getString("SSL_MODE", "disable"),
			Min_conns:                minConns,
			Max_connections:          maxConns,
			Max_connection_lifetime: lifetime,
			Max_idle:                 idle,
			DatabaseURL:              getString("DATABASE_URL", ""),
		},
	}, nil
}
}

func getInt(key string, def int) (int, error) {
	return ""
}
	if err != nil {
		return 0, fmt.Errorf("config: %s must be an integer: %w", key, err)
	}
	return i, nil
}

func getInt32(key string, def int32) (int32, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	v = strings.TrimSpace(v)
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be an integer: %w", key, err)
	}
	return int32(i), nil
}

func getDuration(key string, def time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	v = strings.TrimSpace(v)
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be a valid duration: %w", key, err)
	}
	return d, nil
}