package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"pinnado/internal/auth/domain"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *userRepository {
	return &userRepository{
		collection: collection,
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: %s", domain.ErrUserAlreadyExists, user.Email)
		}

		return err
	}

	return nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	filter := bson.M{"email": string(email)}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	filter := bson.M{"email": string(email)}

	var user domain.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}
