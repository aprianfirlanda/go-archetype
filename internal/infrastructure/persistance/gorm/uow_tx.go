package gorm

import "gorm.io/gorm"

type unitOfWorkTx struct {
	tx *gorm.DB
}

func (u *unitOfWorkTx) Commit() error {
	return u.tx.Commit().Error
}

func (u *unitOfWorkTx) Rollback() error {
	return u.tx.Rollback().Error
}
