package main

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/infrastructure/cloudinary"
	"ev-warranty-go/internal/infrastructure/config"
	"ev-warranty-go/internal/infrastructure/database"
	"ev-warranty-go/internal/infrastructure/oauth"
	"ev-warranty-go/internal/infrastructure/oauth/providers"
	"ev-warranty-go/internal/infrastructure/persistence"
	"ev-warranty-go/internal/interfaces/api"
	"ev-warranty-go/internal/interfaces/api/handlers"
	"ev-warranty-go/internal/security"
	"ev-warranty-go/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "ev-warranty-go/docs"

	"github.com/joho/godotenv"
)

// @title EV Warranty API
// @version 1.0
// @description API for EV Warranty Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@ev-warranty.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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

	txManager := database.NewTxManager(log, db.DB)

	app := &App{
		Cfg: cfg,
		DB:  db,
		Log: log,
	}

	app.seedDbData()

	officeRepo := persistence.NewOfficeRepository(db.DB)
	tokenRepo := persistence.NewTokenRepository(db.DB)
	userRepo := persistence.NewUserRepository(db.DB)
	claimRepo := persistence.NewClaimRepository(db.DB)
	claimItemRepo := persistence.NewClaimItemRepository(db.DB)
	claimAttachmentRepo := persistence.NewClaimAttachmentRepository(db.DB)
	claimHistoryRepo := persistence.NewClaimHistoryRepository(db.DB)

	googleProvider := providers.NewGoogleProvider(
		cfg.OAuth.GoogleClientID, cfg.OAuth.GoogleClientSecret, cfg.OAuth.GoogleRedirectURL)
	cloudinaryService, err := cloudinary.NewCloudinaryService(&cfg.Cloudinary)
	if err != nil {
		log.Error("Failed to initialize cloudinary", "error", err)
		os.Exit(1)
	}

	officeService := services.NewOfficeService(officeRepo)
	tokenService := services.NewTokenService(tokenRepo,
		cfg.AccessTokenTTL, cfg.RefreshTokenTTL, security.PrivateKey(), security.PublicKey())
	authService := services.NewAuthService(userRepo, tokenService)
	userService := services.NewUserService(userRepo, officeRepo)
	oauthService := oauth.NewOAuthService(googleProvider, userRepo)
	claimService := services.NewClaimService(log, claimRepo, claimItemRepo, claimAttachmentRepo, claimHistoryRepo, cloudinaryService)
	claimItemService := services.NewClaimItemService(claimRepo, claimItemRepo)
	claimAttachmentService := services.NewClaimAttachmentService(log, claimRepo, claimAttachmentRepo, cloudinaryService)

	officeHandler := handlers.NewOfficeHandler(log, officeService)
	authHandler := handlers.NewAuthHandler(log, authService, tokenService, userService)
	oauthHandler := handlers.NewOAuthHandler(log, cfg.OAuth.FrontendBaseURL, oauthService, authService)
	userHandler := handlers.NewUserHandler(log, userService)
	claimHandler := handlers.NewClaimHandler(log, txManager, claimService)
	claimItemHandler := handlers.NewClaimItemHandler(log, txManager, claimItemService)
	claimAttachmentHandler := handlers.NewClaimAttachmentHandler(log, txManager, claimAttachmentService)

	r := api.NewRouter(app.DB, authHandler, oauthHandler, officeHandler,
		userHandler, claimHandler, claimItemHandler, claimAttachmentHandler)
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
