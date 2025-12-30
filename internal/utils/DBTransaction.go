package utils

import (
	"context"
	"go-archetype/internal/ports/output"
)

func WithTransaction(
	ctx context.Context,
	uow output.UnitOfWork,
	fn func(tx output.UnitOfWorkTx) error,
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
