package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL            string
	PORT             string
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
	AppEnv           string
}

func (c *Config) IsProd() bool {
	return c.AppEnv == "production"
}

func LoadConfig() (*Config, error) {
	config := &Config{
		DBURL:            getEnv("DB_URL", ""),
		PORT:             getEnv("PORT", "8080"),
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
		JWTAccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "5m"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "2160h"),
		AppEnv:           getEnv("APP_ENV", "development"),
	}

	if !config.IsProd() {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
