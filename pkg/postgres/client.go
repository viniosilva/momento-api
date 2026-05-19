package postgres

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	driverName = "pgx"
)

func Connect(ctx context.Context, host, port, user, pass, dbname, sslmode string, connectTimeout time.Duration) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, dbname, sslmode)

	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, driverName, dsn)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
