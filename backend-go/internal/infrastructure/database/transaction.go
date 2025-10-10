package database

import (
	"context"
	"ev-warranty-go/internal/application"

	"gorm.io/gorm"
)

type transaction struct {
	tx  *gorm.DB
	ctx context.Context
}

func NewTransaction(ctx context.Context, db *gorm.DB) application.Transaction {
	return &transaction{
		tx:  db.WithContext(ctx).Begin(),
		ctx: ctx,
	}
}

func (t *transaction) GetTx() any {
	return t.tx
}

func (t *transaction) GetCtx() context.Context {
	return t.ctx
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback().Error
}

func (t *transaction) Commit() error {
	return t.tx.Commit().Error
}
