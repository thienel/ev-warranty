package main

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/security"
	"os"
)

func (app *App) seedDbData() {
	office := entities.Office{}
	if err := app.DB.DB.
		Where(entities.Office{
			OfficeName: "Head Office",
			OfficeType: entities.OfficeTypeEVM,
		}).
		FirstOrCreate(&office, entities.Office{
			OfficeName: "Head Office",
			OfficeType: entities.OfficeTypeEVM,
			Address:    "Main Street",
			IsActive:   true,
		}).Error; err != nil {
		app.Log.Error("failed to seed office", "error", err)
		os.Exit(1)
	}

	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		app.Log.Error("missing ADMIN_PASSWORD env")
		os.Exit(1)
	}
	hashed, err := security.HashPassword(password)
	if err != nil {
		app.Log.Error("failed to hash password", "error", err)
		os.Exit(1)
	}

	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		app.Log.Error("missing ADMIN_EMAIL env")
	}

	admin := entities.User{}
	if err := app.DB.DB.
		Where(&entities.User{Email: email}).
		Attrs(entities.User{
			Name:         "System Admin",
			Role:         entities.UserRoleAdmin,
			PasswordHash: hashed,
			IsActive:     true,
			OfficeID:     office.ID,
		}).
		FirstOrCreate(&admin).Error; err != nil {
		app.Log.Error("failed to seed admin user", "error", err)
		os.Exit(1)
	}
}
