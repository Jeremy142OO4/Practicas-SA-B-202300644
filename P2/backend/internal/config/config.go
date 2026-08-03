package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Cookie   CookieConfig
	Security SecurityConfig
}

type AppConfig struct {
	Environment string
	Port        int
	FrontendURL string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret              string
	ExpirationMinutes   int
	RenewalGraceMinutes int
	Issuer              string
}

type CookieConfig struct {
	Secure   bool
	SameSite string
}

type SecurityConfig struct {
	AESKey        []byte
	EmailIndexKey []byte
}

func Load() (Config, error) {
	_ = godotenv.Load()

	port, err := positiveInt("APP_PORT")
	if err != nil {
		return Config{}, err
	}

	expiration, err := positiveInt("JWT_EXPIRATION_MINUTES")
	if err != nil {
		return Config{}, err
	}

	grace, err := positiveInt("JWT_RENEWAL_GRACE_MINUTES")
	if err != nil {
		return Config{}, err
	}

	secure, err := boolValue("COOKIE_SECURE")
	if err != nil {
		return Config{}, err
	}

	aesKey, err := base64Key("AES_KEY", 32)
	if err != nil {
		return Config{}, err
	}

	indexKey, err := base64Key("EMAIL_INDEX_KEY", 32)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		App: AppConfig{
			Environment: valueOrDefault("APP_ENV", "development"),
			Port:        port,
			FrontendURL: strings.TrimSpace(os.Getenv("FRONTEND_URL")),
		},
		Database: DatabaseConfig{URL: strings.TrimSpace(os.Getenv("DATABASE_URL"))},
		JWT: JWTConfig{
			Secret:              strings.TrimSpace(os.Getenv("JWT_SECRET")),
			ExpirationMinutes:   expiration,
			RenewalGraceMinutes: grace,
			Issuer:              strings.TrimSpace(os.Getenv("JWT_ISSUER")),
		},
		Cookie: CookieConfig{
			Secure:   secure,
			SameSite: strings.ToLower(strings.TrimSpace(os.Getenv("COOKIE_SAME_SITE"))),
		},
		Security: SecurityConfig{AESKey: aesKey, EmailIndexKey: indexKey},
	}

	if err := config.validate(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (config Config) validate() error {
	required := map[string]string{
		"FRONTEND_URL": config.App.FrontendURL,
		"DATABASE_URL": config.Database.URL,
		"JWT_SECRET":   config.JWT.Secret,
		"JWT_ISSUER":   config.JWT.Issuer,
	}
	for name, value := range required {
		if value == "" {
			return fmt.Errorf("environment variable %s is required", name)
		}
	}

	if len(config.JWT.Secret) < 32 {
		return errors.New("JWT_SECRET must contain at least 32 characters")
	}

	validSameSite := map[string]bool{"lax": true, "strict": true, "none": true}
	if !validSameSite[config.Cookie.SameSite] {
		return errors.New("COOKIE_SAME_SITE must be Lax, Strict, or None")
	}

	if config.Cookie.SameSite == "none" && !config.Cookie.Secure {
		return errors.New("COOKIE_SECURE must be true when COOKIE_SAME_SITE is None")
	}

	return nil
}

func positiveInt(name string) (int, error) {
	value := strings.TrimSpace(os.Getenv(name))
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("environment variable %s must be a positive integer", name)
	}
	return parsed, nil
}

func boolValue(name string) (bool, error) {
	value := strings.TrimSpace(os.Getenv(name))
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("environment variable %s must be true or false", name)
	}
	return parsed, nil
}

func base64Key(name string, expectedLength int) ([]byte, error) {
	encoded := strings.TrimSpace(os.Getenv(name))
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(decoded) != expectedLength {
		return nil, fmt.Errorf("environment variable %s must be a Base64-encoded %d-byte key", name, expectedLength)
	}
	return decoded, nil
}

func valueOrDefault(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}
