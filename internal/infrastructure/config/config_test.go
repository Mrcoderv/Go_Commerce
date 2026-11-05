package config_test

import (
	"Go_Commerce/internal/infrastructure/config"
	"strings"
	"testing"
)

func TestConfig_Defaults(t *testing.T) {
	t.Setenv("APP_PORT", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("APP_ENV", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error loading defaults: %v", err)
	}

	if cfg.AppPort != 8080 {
		t.Errorf("expected AppPort 8080, got %d", cfg.AppPort)
	}
	if cfg.DBPort != 5432 {
		t.Errorf("expected DBPort 5432, got %d", cfg.DBPort)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected DBHost localhost, got %s", cfg.DBHost)
	}
	if cfg.DBUser != "postgres" {
		t.Errorf("expected DBUser postgres, got %s", cfg.DBUser)
	}
}

func TestConfig_BadPort(t *testing.T) {
	t.Setenv("APP_PORT", "not-a-number")
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for non-numeric APP_PORT, got nil")
	}

	t.Setenv("APP_PORT", "999999")
	_, err = config.Load()
	if err == nil {
		t.Fatal("expected error for out of range APP_PORT, got nil")
	}
}

func TestConfig_DSNWithSpecialCharInPassword(t *testing.T) {
	t.Setenv("DB_PASSWORD", "p@ss:word#123")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "myuser")
	t.Setenv("DB_NAME", "mydb")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dsn := cfg.DSN()
	if !strings.Contains(dsn, "p%40ss%3Aword%23123") {
		t.Errorf("expected DSN to URL-encode password, got: %s", dsn)
	}
}

func TestConfig_ProductionEmptyPassword(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DB_PASSWORD", "")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for empty DB_PASSWORD in production, got nil")
	}
}
