package entity

import (
	"net/mail"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	UserRoleAdmin        = "ADMIN"
	UserRoleEvmStaff     = "EVM_STAFF"
	UserRoleScStaff      = "SC_STAFF"
	UserRoleScTechnician = "SC_TECHNICIAN"
)

type User struct {
	ID            uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name          string         `gorm:"not null"`
	Email         string         `gorm:"not null;uniqueIndex;"`
	Role          string         `gorm:"not null"`
	PasswordHash  string         `gorm:"not null"`
	IsActive      bool           `gorm:"not null;default:true"`
	OfficeID      uuid.UUID      `gorm:"not null;type:uuid"`
	Office        Office         `gorm:"foreignKey:OfficeID;references:ID"`
	OAuthProvider *string        `gorm:"column:oauth_provider"`
	OAuthID       *string        `gorm:"column:oauth_id"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func IsValidUserRole(userRole string) bool {
	switch userRole {
	case UserRoleAdmin, UserRoleEvmStaff, UserRoleScStaff, UserRoleScTechnician:
		return true
	default:
		return false
	}
}

func NewUser(name, email, role, passwordHash string, isActive bool, officeID uuid.UUID) *User {
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

func IsValidName(name string) bool {
	if len(name) < 2 || len(name) > 50 {
		return false
	}
	regex := regexp.MustCompile(`^[\p{L}][\p{L}\s'-]*[\p{L}]$`)
	return regex.MatchString(name)
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

func (u *User) IsValidOfficeByRole(officeType string) bool {
	switch officeType {
	case OfficeTypeEVM:
		if u.Role == UserRoleAdmin || u.Role == UserRoleEvmStaff {
			return true
		}
	case OfficeTypeSC:
		if u.Role == UserRoleScTechnician || u.Role == UserRoleScStaff {
			return true
		}
	}
	return false
}

func (u *User) IsOAuthUser() bool {
	return u.OAuthProvider != nil && u.OAuthID != nil
}

func (u *User) LinkToOAuth(oauthProvider, oauthID string) {
	u.OAuthProvider = &oauthProvider
	u.OAuthID = &oauthID
}
