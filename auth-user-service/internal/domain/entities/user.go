package entities

import (
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email         string         `gorm:"uniqueIndex;size:100"`
	PasswordHash  *string        `gorm:"size:255"`
	IsActive      bool           `gorm:"default:true"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Tokens        []RefreshToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	OAuthProvider *string        `gorm:"size:32;column:oauth_provider"`
	OAuthID       *string        `gorm:"size:64;column:oauth_id"`
}

func NewUser(email, passwordHash string) *User {
	return &User{
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: &passwordHash,
		IsActive:     true,
	}
}

func NewOAuthUser(email, provider, oauthID string) *User {
	return &User{
		Email:         strings.ToLower(strings.TrimSpace(email)),
		IsActive:      true,
		OAuthProvider: &provider,
		OAuthID:       &oauthID,
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

func (user *User) Activate() {
	user.IsActive = true
}

func (user *User) Deactivate() {
	user.IsActive = false
}

func (user *User) IsOAuthUser() bool {
	return user.OAuthProvider != nil && user.OAuthID != nil
}

func (user *User) LinkToOAuth(oauthProvider, oauthID string) {
	user.OAuthProvider = &oauthProvider
	user.OAuthID = &oauthID
}
