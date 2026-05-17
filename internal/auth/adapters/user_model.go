package adapters

import (
	"time"

	"momento/internal/auth/domain"
)

type userRow struct {
	ID              string     `db:"id"`
	Email           string     `db:"email"`
	Password        string     `db:"password"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	EmailVerifiedAt *time.Time `db:"email_verified_at"`
}

func toUserRow(u domain.User) userRow {
	return userRow{
		ID:              u.ID,
		Email:           string(u.Email),
		Password:        string(u.Password),
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		EmailVerifiedAt: u.EmailVerifiedAt,
	}
}

func toUserDomain(row userRow) domain.User {
	return domain.User{
		ID:              row.ID,
		Email:           domain.Email(row.Email),
		Password:        domain.Password(row.Password),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		EmailVerifiedAt: row.EmailVerifiedAt,
	}
}
