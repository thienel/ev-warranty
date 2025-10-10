package application

import "context"

type Transaction interface {
	GetTx() any
	GetCtx() context.Context
	Rollback() error
	Commit() error
}
