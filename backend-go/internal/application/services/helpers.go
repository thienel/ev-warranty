package services

import (
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"log"
)

func rollbackOnErr(tx application.Transaction, originalErr error) error {
	if err := tx.Rollback(); err != nil {
		log.Printf("[TX ROLLBACK FAILED] original error: %v, rollback error: %v", originalErr, err)
	}
	return originalErr
}

func commitOrLog(tx application.Transaction) error {
	if err := tx.Commit(); err != nil {
		log.Printf("[TX COMMIT FAILED] commit error: %v", err)
		return apperrors.NewInternalServerError(err)
	}
	return nil
}
