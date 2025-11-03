package database

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/logger"
	"log"

	"gorm.io/gorm"
)

type Tx struct {
	tx  *gorm.DB
	ctx context.Context
}

func (t *Tx) GetTx() any {
	return t.tx
}

func (t *Tx) GetCtx() context.Context {
	return t.ctx
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback().Error
}

func (t *Tx) Commit() error {
	return t.tx.Commit().Error
}

type txManager struct {
	log logger.Logger
	db  *gorm.DB
}

func NewTxManager(log logger.Logger, db *gorm.DB) application.TxManager {
	return &txManager{
		log: log,
		db:  db,
	}
}

func (m *txManager) Do(ctx context.Context, fn func(tx application.Tx) error) error {
	tx := m.db.Begin()
	t := &Tx{
		tx:  tx,
		ctx: ctx,
	}

	if err := fn(t); err != nil {
		if rbErr := t.Rollback(); rbErr != nil {
			log.Printf("[TX ROLLBACK FAILED] original error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err := t.Commit(); err != nil {
		log.Printf("[TX COMMIT FAILED] commit error: %v", err)
		return apperror.NewInternalServerError(err)
	}

	return nil
}
