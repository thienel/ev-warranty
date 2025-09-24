package main

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/security"
	"log"
)

func (app *App) seedDbData() {
	err := app.DB.DB.AutoMigrate(&entities.Office{}, &entities.User{}, &entities.RefreshToken{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	office := entities.
		NewOffice("Head Office", entities.OfficeTypeEVM, "123 Main Street", true)
	app.DB.DB.FirstOrCreate(&office, entities.Office{OfficeName: "Head Office"})

	password := "Admin@123"
	hashed, err := security.HashPassword(password)
	if err != nil {
		log.Fatal("failed to hash password:", err)
	}

	email := "admin@localhost"
	admin := entities.
		NewUser("System Admin", email, entities.UserRoleAdmin, hashed, true, office.Id)

	app.DB.DB.FirstOrCreate(&admin, entities.User{Email: email})
}
