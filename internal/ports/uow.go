package ports

import "context"

type UnitOfWork interface {
	Begin(ctx context.Context) (UnitOfWorkTx, error)
}

type UnitOfWorkTx interface {
	Commit() error
	Rollback() error

	// expose repository getter here`
}
