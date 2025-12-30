package gorm

import (
	"context"
	"go-archetype/internal/ports/output"

	"gorm.io/gorm"
)

type unitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) output.UnitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) Begin(ctx context.Context) (output.UnitOfWorkTx, error) {
	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &unitOfWorkTx{tx: tx}, nil
}
