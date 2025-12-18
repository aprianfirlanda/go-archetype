package gorm

import (
	"context"
	"go-archetype/internal/ports/outbound"

	"gorm.io/gorm"
)

type unitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) outbound.UnitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) Begin(ctx context.Context) (outbound.UnitOfWorkTx, error) {
	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &unitOfWorkTx{
		tx: tx,
	}, nil
}
