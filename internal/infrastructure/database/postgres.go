package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/Mrcoderv/Go_Commerce/internal/infrastructure/config"
)

func NewDB() (*sql.DB, error) {
	config, err := config.LoadConfig()

	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", dbConfigValue(config.DB, "DSN", "URL", "ConnectionString"))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(dbConfigInt(config.DB, 10, "MaxOpenConns", "MaxConnections"))
	db.SetMaxIdleConns(dbConfigInt(config.DB, 5, "MaxIdleConns", "MaxIdle"))
	db.SetConnMaxLifetime(dbConfigDuration(config.DB, 5*time.Minute, "ConnMaxLifetime", "ConnectionMaxLifetime"))
	err = db.Ping() // testing ping.
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully  database connected!")
	return (db), nil
}

func dbConfigField(value any, names ...string) reflect.Value {
	v := reflect.Indirect(reflect.ValueOf(value))
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return reflect.Value{}
	}
	for _, name := range names {
		if field := v.FieldByName(name); field.IsValid() && field.CanInterface() {
			return field
		}
	}
	return reflect.Value{}
}

func dbConfigValue(value any, names ...string) string {
	if field := dbConfigField(value, names...); field.IsValid() {
		return fmt.Sprint(field.Interface())
	}
	return fmt.Sprintf("%v", value)
}

func dbConfigInt(value any, fallback int, names ...string) int {
	if field := dbConfigField(value, names...); field.IsValid() {
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int(field.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int(field.Uint())
		}
	}
	return fallback
}

func dbConfigDuration(value any, fallback time.Duration, names ...string) time.Duration {
	if field := dbConfigField(value, names...); field.IsValid() {
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			return time.Duration(field.Int())
		}
		if field.Kind() == reflect.String {
			if duration, err := time.ParseDuration(field.String()); err == nil {
				return duration
			}
		}
	}
	return fallback
}
