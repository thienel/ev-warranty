package application

type Transaction interface {
	GetTx() any
	Rollback() error
	Commit() error
}
