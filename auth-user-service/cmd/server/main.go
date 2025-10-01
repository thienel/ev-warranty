package main

import (
	"auth-service/internal/application/services"
	"auth-service/internal/infrastructure/config"
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/infrastructure/oauth"
	"auth-service/internal/infrastructure/oauth/providers"
	"auth-service/internal/infrastructure/persistence"
	"auth-service/internal/interfaces/api"
	"auth-service/internal/security"
	"auth-service/pkg/logger"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

type App struct {
	Cfg *config.Config
	DB  *database.Database
	Log logger.Logger
}

func main() {
	_ = godotenv.Load(".env")
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	err := security.InitRSAKeys(cfg.PublicKeyPath, cfg.PrivateKeyPath)
	if err != nil {
		log.Error("Failed to initialize RSA keys", "error", err)
		os.Exit(1)
	}

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			log.Error("Failed to close database connection", "error", err)
		}
	}(db)

	app := &App{
		Cfg: cfg,
		DB:  db,
		Log: log,
	}

	app.seedDbData()

	officeRepo := persistence.NewOfficeRepository(db.DB)
	tokenRepo := persistence.NewTokenRepository(db.DB)
	userRepo := persistence.NewUserRepository(db.DB)

	officeService := services.NewOfficeService(officeRepo)
	tokenService := services.NewTokenService(tokenRepo,
		cfg.AccessTokenTTL, cfg.RefreshTokenTTL, security.PrivateKey(), security.PublicKey())
	authService := services.NewAuthService(userRepo, tokenService)
	oauthService := oauth.NewOAuthService(userRepo)
	userService := services.NewUserService(userRepo, officeService)

	googleProvider := providers.NewGoogleProvider(
		cfg.OAuth.GoogleClientID, cfg.OAuth.GoogleClientSecret, cfg.OAuth.GoogleRedirectURL)
	oauthService.RegisterProvider(googleProvider)

	officeHandler := api.NewOfficeHandler(log, officeService)
	authHandler := api.NewAuthHandler(log, authService, tokenService, userService)
	oauthHandler := api.NewOAuthHandler(log, cfg.OAuth.FrontendBaseURL, oauthService, authService)
	userHandler := api.NewUserHandler(log, userService)

	r := api.NewRouter(app.DB, authHandler, oauthHandler, officeHandler, userHandler)
	log.Info("Server starting on port " + cfg.Port)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}
}
