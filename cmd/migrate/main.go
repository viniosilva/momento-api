package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"momento/internal/config"
	"momento/pkg/logger"
	"momento/pkg/postgres"
)

const (
	shutdownTimeout  = 10 * time.Second
	migrationsDir    = "migrations"
)

func main() {
	slog.SetDefault(logger.NewLogger("info"))

	log.Println("loading configuration...")
	cfg := config.LoadConfig()

	log.Println("connecting to PostgreSQL...")
	ctx := context.Background()
	db, err := postgres.Connect(ctx,
		cfg.PG.DSN,
		cfg.PG.MaxRetries,
		cfg.PG.RetryDelay,
		cfg.PG.ConnectTimeout)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer func() {
		log.Println("closing database connection...")
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	log.Println("running migrations...")
	if err := runMigrations(ctx, db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migration completed successfully")
}

func runMigrations(ctx context.Context, db *sqlx.DB) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	sqlFiles := make([]string, 0)
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}

	sort.Strings(sqlFiles)

	for _, name := range sqlFiles {
		path := filepath.Join(migrationsDir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		log.Printf("applying migration: %s", name)

		_, err = db.ExecContext(ctx, string(content))
		if err != nil {
			return fmt.Errorf("execute migration %s: %w", name, err)
		}

		log.Printf("migration %s applied", name)
	}

	return nil
}
