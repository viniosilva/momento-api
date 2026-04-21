package adapters

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"momento/internal/auth/domain"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{
		collection: db.Collection(usersCollectionName),
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	doc, err := toUserDocument(user)
	if err != nil {
		return fmt.Errorf("toUserDocument: %w", err)
	}

	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserAlreadyExists
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

	var doc userDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return toUserDomain(doc), nil
}

func (r *userRepository) Update(ctx context.Context, user domain.User) error {
	doc, err := toUserDocument(user)
	if err != nil {
		return fmt.Errorf("toUserDocument: %w", err)
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{
		"$set": bson.M{
			"password":          doc.Password,
			"updated_at":        doc.UpdatedAt,
			"email_verified_at": doc.EmailVerifiedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
