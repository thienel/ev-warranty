package main

import (
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/security"
	"os"
)

func (app *App) seedDbData() {
	office := entity.Office{}
	if err := app.DB.DB.
		Where(entity.Office{
			OfficeName: "Head Office",
			OfficeType: entity.OfficeTypeEVM,
		}).
		FirstOrCreate(&office, entity.Office{
			OfficeName: "Head Office",
			OfficeType: entity.OfficeTypeEVM,
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
		os.Exit(1)
	}

	admin := entity.User{}
	if err := app.DB.DB.
		Where(&entity.User{Email: email}).
		Attrs(entity.User{
			Name:         "System Admin",
			Role:         entity.UserRoleAdmin,
			PasswordHash: hashed,
			IsActive:     true,
			OfficeID:     office.ID,
		}).
		FirstOrCreate(&admin).Error; err != nil {
		app.Log.Error("failed to seed admin user", "error", err)
		os.Exit(1)
	}
}
