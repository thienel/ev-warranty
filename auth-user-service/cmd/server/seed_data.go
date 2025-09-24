package main

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/security"
	"os"
)

func (app *App) seedDbData() {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = "Admin@123"
	}
	hashed, err := security.HashPassword(password)
	if err != nil {
		app.Log.Error("failed to hash password", "error", err)
		os.Exit(1)
		return
	}

	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		email = "admin@localhost"
	}
	admin := entities.NewUser("System Admin", email, entities.UserRoleAdmin, hashed, true, nil)
	if err := app.DB.DB.FirstOrCreate(&admin, entities.User{Email: email}).Error; err != nil {
		app.Log.Error("failed to seed admin user", "error", err)
	}
}
