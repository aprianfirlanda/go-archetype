package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Ping checks database reachability with a timeout.
// Exported so it can be reused by other packages (e.g., readiness checks).
func Ping(ctx context.Context, sqlDB *sql.DB) error {
	if sqlDB == nil {
		return fmt.Errorf("sqlDB is nil")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	return sqlDB.PingContext(ctx)
}
