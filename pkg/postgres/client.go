package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(ctx context.Context, dsn string, maxRetries int, retryDelay, connectTimeout time.Duration) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	var db *sqlx.DB
	var err error

	driverName := "pgx"

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = sqlx.ConnectContext(ctx, driverName, dsn)
		if err == nil {
			if err = db.PingContext(ctx); err == nil {
				return db, nil
			}
			db.Close()
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %w", maxRetries, err)
}
