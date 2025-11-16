package main

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application/service"
	"ev-warranty-go/internal/infrastructure/client/dotnet"
	"ev-warranty-go/internal/infrastructure/cloudinary"
	"ev-warranty-go/internal/infrastructure/config"
	"ev-warranty-go/internal/infrastructure/database"
	"ev-warranty-go/internal/infrastructure/oauth"
	"ev-warranty-go/internal/infrastructure/oauth/providers"
	"ev-warranty-go/internal/infrastructure/persistence"
	"ev-warranty-go/internal/interface/api"
	"ev-warranty-go/internal/interface/api/handler"
	"ev-warranty-go/pkg/logger"
	"ev-warranty-go/pkg/security"
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

	dotnetClient := dotnet.NewClient(cfg.ExternalService.DotnetBackendURL)

	officeService := service.NewOfficeService(officeRepo)
	tokenService := service.NewTokenService(tokenRepo,
		cfg.AccessTokenTTL, cfg.RefreshTokenTTL, security.PrivateKey(), security.PublicKey())
	authService := service.NewAuthService(userRepo, tokenService)
	userService := service.NewUserService(userRepo, officeRepo, claimRepo)
	oauthService := oauth.NewOAuthService(googleProvider, userRepo)
	claimService := service.NewClaimService(log, claimRepo, userRepo, claimItemRepo, claimAttachmentRepo,
		claimHistoryRepo, cloudinaryService)
	claimItemService := service.NewClaimItemService(claimRepo, claimItemRepo, userRepo, dotnetClient)
	claimAttachmentService := service.NewClaimAttachmentService(log, claimRepo, claimAttachmentRepo,
		cloudinaryService)

	officeHandler := handler.NewOfficeHandler(log, officeService)
	authHandler := handler.NewAuthHandler(log, authService, tokenService, userService)
	oauthHandler := handler.NewOAuthHandler(log, cfg.OAuth.FrontendBaseURL, oauthService, authService)
	userHandler := handler.NewUserHandler(log, userService)
	claimHandler := handler.NewClaimHandler(log, txManager, claimService)
	claimItemHandler := handler.NewClaimItemHandler(log, txManager, claimItemService)
	claimAttachmentHandler := handler.NewClaimAttachmentHandler(log, txManager, claimAttachmentService)

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
