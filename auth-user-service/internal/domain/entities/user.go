package entities

import (
	"auth-service/internal/errors/apperrors"
	"database/sql/driver"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusPending   UserStatus = "pending"
)

type User struct {
	ID            uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username      string         `gorm:"uniqueIndex;size:24"`
	Email         string         `gorm:"uniqueIndex;size:100"`
	PasswordHash  *string        `gorm:"size:255"`
	Status        UserStatus     `gorm:"default:pending"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Tokens        []RefreshToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	OAuthProvider *string        `gorm:"size:32;column:oauth_provider"`
	OAuthID       *string        `gorm:"size:64;column:oauth_id"`
}

func NewUser(username, email, passwordHash string) *User {
	return &User{
		Username:     strings.TrimSpace(username),
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: &passwordHash,
		Status:       UserStatusPending,
	}
}

func NewOAuthUser(username, email, provider, oauthID string) *User {
	return &User{
		Username:      username,
		Email:         strings.ToLower(strings.TrimSpace(email)),
		Status:        UserStatusActive,
		OAuthProvider: &provider,
		OAuthID:       &oauthID,
	}
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidUserName(username string) bool {
	pattern := `^[A-Za-z0-9._]{2,24}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(username)
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

func (user *User) IsActive() bool {
	return user.Status == UserStatusActive && user.DeletedAt.Time.IsZero()
}

func (user *User) Activate() {
	user.Status = UserStatusActive
}

func (user *User) Suspend() {
	user.Status = UserStatusSuspended
}

func (user *User) Deactivate() {
	user.Status = UserStatusInactive
}

func (user *User) IsOAuthUser() bool {
	return user.OAuthProvider != nil && user.OAuthID != nil
}

func (user *User) LinkToOAuth(oauthProvider, oauthID string) {
	user.OAuthProvider = &oauthProvider
	user.OAuthID = &oauthID
}

func (us UserStatus) Value() (driver.Value, error) {
	return string(us), nil
}

func (us *UserStatus) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch s := value.(type) {
	case string:
		*us = UserStatus(s)
		return nil
	case []byte:
		*us = UserStatus(s)
		return nil
	default:
		return apperrors.ErrScanValue
	}
}

func (us UserStatus) IsValid() bool {
	switch us {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended, UserStatusPending:
		return true
	default:
		return false
	}
}
