package services

import (
	"database/sql"
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"log"
	"net/url"
)

func rollbackOrLog(tx application.Transaction) {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		log.Printf("[TX ROLLBACK FAILED] rollback error: %v", err)
	}
}

func commitOrLog(tx application.Transaction) error {
	if err := tx.Commit(); err != nil {
		log.Printf("[TX COMMIT FAILED] commit error: %v", err)
		return apperrors.NewInternalServerError(err)
	}
	return nil
}

func IsValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
