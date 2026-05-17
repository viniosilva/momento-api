package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"momento/internal/auth/domain"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	row := toUserRow(user)

	query := `INSERT INTO users (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query, row.ID, row.Email, row.Password, row.CreatedAt, row.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrUserAlreadyExists
		}

		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, string(email)).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user exists: %w", err)
	}

	return exists, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at, email_verified_at FROM users WHERE email = $1`

	var row userRow
	err := r.db.GetContext(ctx, &row, query, string(email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return toUserDomain(row), nil
}

func (r *userRepository) FindVerifiedByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at, email_verified_at FROM users WHERE email = $1 AND email_verified_at IS NOT NULL`

	var row userRow
	err := r.db.GetContext(ctx, &row, query, string(email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf("find verified user by email: %w", err)
	}

	return toUserDomain(row), nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at, email_verified_at FROM users WHERE id = $1`

	var row userRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf("find user by id: %w", err)
	}

	return toUserDomain(row), nil
}

func (r *userRepository) Update(ctx context.Context, user domain.User) error {
	row := toUserRow(user)

	query := `UPDATE users SET password = $1, updated_at = $2, email_verified_at = $3 WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query, row.Password, row.UpdatedAt, row.EmailVerifiedAt, row.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate key") || contains(err.Error(), "unique constraint") || contains(err.Error(), "23505"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
