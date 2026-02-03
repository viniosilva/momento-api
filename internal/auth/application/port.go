package application

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"pinnado/internal/auth/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	HasByEmail(ctx context.Context, email domain.Email) (bool, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (domain.User, error)
}
