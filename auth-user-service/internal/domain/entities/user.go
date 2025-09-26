package entities

import (
	"net/mail"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	UserRoleAdmin        = "admin"
	UserRoleEvmStaff     = "evm staff"
	UserRoleScStaff      = "sc staff"
	UserRoleScTechnician = "sc technician"
)

type User struct {
	ID            uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name          string         `gorm:"not null;size:30"`
	Email         string         `gorm:"not null;uniqueIndex;size:100"`
	Role          string         `gorm:"not null;size:20"`
	PasswordHash  string         `gorm:"not null;size:255"`
	IsActive      bool           `gorm:"not null;default:true"`
	OfficeID      *uuid.UUID     `gorm:"type:uuid"`
	Office        *Office        `gorm:"foreignKey:OfficeID;references:ID"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	OAuthProvider *string        `gorm:"size:32;column:oauth_provider"`
	OAuthID       *string        `gorm:"size:64;column:oauth_id"`
}

func IsValidUserRole(userRole string) bool {
	switch userRole {
	case UserRoleAdmin, UserRoleEvmStaff, UserRoleScStaff, UserRoleScTechnician:
		return true
	default:
		return false
	}
}

func NewUser(name, email, role, passwordHash string, isActive bool, officeID *uuid.UUID) *User {
	return &User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		Role:         role,
		PasswordHash: passwordHash,
		IsActive:     isActive,
		OfficeID:     officeID,
	}
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)

	return hasLower && hasUpper && hasDigit && hasSpecial
}

func (user *User) IsOAuthUser() bool {
	return user.OAuthProvider != nil && user.OAuthID != nil
}

func (user *User) LinkToOAuth(oauthProvider, oauthID string) {
	user.OAuthProvider = &oauthProvider
	user.OAuthID = &oauthID
}
