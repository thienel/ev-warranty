package config

import (
	"os"
	"time"
)

type OAuthConfig struct {
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
}

type Config struct {
	Port            string
	DatabaseURL     string
	LogLevel        string
	PublicKeyPath   string
	PrivateKeyPath  string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	OAuth           OAuthConfig `mapstructure:",squash"`
}

func Load() *Config {
	accessTokenTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		accessTokenTTL = 15 * time.Minute
	}
	refreshTokenTTL, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		refreshTokenTTL = 168 * time.Hour
	}
	ggClientID := os.Getenv("GOOGLE_CLIENT_ID")
	ggClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if ggClientID == "" || ggClientSecret == "" {
		panic("Google OAuth credentials are not set in environment variables")
	}
	return &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://auth_service:password@localhost:5432/auth_service?sslmode=disable"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		PublicKeyPath:   getEnv("PUBLIC_KEY_PATH", "./keys/public.pem"),
		PrivateKeyPath:  getEnv("PRIVATE_KEY_PATH", "./keys/private.pem"),
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		OAuth: OAuthConfig{
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/oauth/google/callback"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultVal
}
