package application

import (
	"context"
)

type Tx interface {
	GetTx() any
	GetCtx() context.Context
	Rollback() error
	Commit() error
}

type TxManager interface {
	Do(ctx context.Context, fn func(tx Tx) error) error
}
