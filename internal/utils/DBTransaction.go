package utils

import (
	"context"
	"go-archetype/internal/ports/outbound"
)

func WithTransaction(
	ctx context.Context,
	uow outbound.UnitOfWork,
	fn func(tx outbound.UnitOfWorkTx) error,
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
