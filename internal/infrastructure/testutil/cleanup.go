package testutil

import (
	"fmt"

	"gorm.io/gorm"
)

func Truncate(db *gorm.DB, tables ...string) error {
	for _, table := range tables {
		if table == "" {
			return fmt.Errorf("empty table name")
		}

		sql := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}
