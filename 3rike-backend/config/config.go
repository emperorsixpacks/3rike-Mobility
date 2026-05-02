package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	CantonURL   string
	CantonToken string
	Env         string

	// Keycloak OIDC — when set, the backend auto-fetches Canton bearer tokens.
	OIDCTokenURL string
	OIDCClientID string
	OIDCUsername string
	OIDCPassword string

	// Canton operator party ID (your platform's party on the ledger).
	CantonOperatorParty string
	CantonValidatorURL  string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me"),
		CantonURL:   getEnv("CANTON_URL", "http://localhost:7575"),
		CantonToken: getEnv("CANTON_TOKEN", ""),
		Env:         getEnv("ENV", "development"),

		OIDCTokenURL: getEnv("OIDC_TOKEN_URL", ""),
		OIDCClientID: getEnv("OIDC_CLIENT_ID", ""),
		OIDCUsername: getEnv("OIDC_USERNAME", ""),
		OIDCPassword: getEnv("OIDC_PASSWORD", ""),

		CantonOperatorParty: getEnv("CANTON_OPERATOR_PARTY", ""),
		CantonValidatorURL:  getEnv("CANTON_VALIDATOR_URL", "https://wallet.validator.hackcanton-01.devnet.naas.noders.services"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
