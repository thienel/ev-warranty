package config

import (
	"os"
	"time"
)

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	FrontendBaseURL    string
}

type CloudinaryConfig struct {
	URL          string
	UploadFolder string
}

type ExternalServiceConfig struct {
	DotnetBackendURL string
}

type Config struct {
	Port            string
	DatabaseURL     string
	LogLevel        string
	PublicKeyPath   string
	PrivateKeyPath  string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	OAuth           OAuthConfig
	Cloudinary      CloudinaryConfig
	ExternalService ExternalServiceConfig
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
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		panic("Cloudinary URL are not set in environment variables")
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
			GoogleClientID:     ggClientID,
			GoogleClientSecret: ggClientSecret,
			GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost/api/v1/oauth/google/callback"),
			FrontendBaseURL:    getEnv("FRONTEND_BASE_URL", "http://localhost:3000"),
		},
		Cloudinary: CloudinaryConfig{
			URL:          cloudinaryURL,
			UploadFolder: getEnv("CLOUDINARY_UPLOAD_FOLDER", "ev-warranty-claim-attachment"),
		},
		ExternalService: ExternalServiceConfig{
			DotnetBackendURL: getEnv("DOTNET_BACKEND_URL", "http://localhost"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultVal
}
