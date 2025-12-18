package utils

import (
	"context"
	"go-archetype/internal/ports"
)

func WithTransaction(
	ctx context.Context,
	uow ports.UnitOfWork,
	fn func(tx ports.UnitOfWorkTx) error,
) error {
	tx, err := uow.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
