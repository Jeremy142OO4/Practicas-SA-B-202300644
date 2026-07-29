package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Host        string
	Port        string
}

type DatabaseConfig struct {
	URL            string
	MaxConnections int32
	MinConnections int32
	ConnectTimeout time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf(
			"cargar archivo .env: %w",
			err,
		)
	}

	maxConnections, err := readInt32(
		"DATABASE_MAX_CONNECTIONS",
		10,
	)
	if err != nil {
		return nil, err
	}

	minConnections, err := readInt32(
		"DATABASE_MIN_CONNECTIONS",
		2,
	)
	if err != nil {
		return nil, err
	}

	timeoutSeconds, err := readInt(
		"DATABASE_CONNECT_TIMEOUT_SECONDS",
		10,
	)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Name: readString(
				"APP_NAME",
				"academia-api",
			),
			Environment: readString(
				"APP_ENV",
				"development",
			),
			Host: readString(
				"APP_HOST",
				"0.0.0.0",
			),
			Port: readString(
				"APP_PORT",
				"8080",
			),
		},
		Database: DatabaseConfig{
			URL:            os.Getenv("DATABASE_URL"),
			MaxConnections: maxConnections,
			MinConnections: minConnections,
			ConnectTimeout: time.Duration(
				timeoutSeconds,
			) * time.Second,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.Database.URL == "" {
		return errors.New(
			"DATABASE_URL es obligatoria",
		)
	}

	if c.Database.MinConnections < 0 ||
		c.Database.MaxConnections < 1 {
		return errors.New(
			"la configuración del pool no es válida",
		)
	}

	if c.Database.MinConnections >
		c.Database.MaxConnections {
		return errors.New(
			"DATABASE_MIN_CONNECTIONS no puede " +
				"superar DATABASE_MAX_CONNECTIONS",
		)
	}

	return nil
}

func readString(
	key string,
	defaultValue string,
) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func readInt(
	key string,
	defaultValue int,
) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf(
			"%s debe ser un número entero: %w",
			key,
			err,
		)
	}

	return parsed, nil
}

func readInt32(
	key string,
	defaultValue int32,
) (int32, error) {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseInt(
		value,
		10,
		32,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"%s debe ser un número entero: %w",
			key,
			err,
		)
	}

	return int32(parsed), nil
}
