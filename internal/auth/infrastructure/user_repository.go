package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"pinnado/internal/auth/domain"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: %s", domain.ErrUserAlreadyExists, user.Email)
		}

		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	filter := bson.M{"email": string(email)}

	var user domain.User
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, fmt.Errorf("%w: email %s", domain.ErrUserNotFound, email)
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	filter := bson.M{"_id": id}

	var user domain.User
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, fmt.Errorf("%w: id %s", domain.ErrUserNotFound, id.Hex())
		}

		return domain.User{}, err
	}

	return user, nil
}
