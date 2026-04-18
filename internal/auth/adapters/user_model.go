package adapters

import (
	"time"

	"momento/internal/auth/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const usersCollectionName = "users"

type userDocument struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func toUserDocument(u domain.User) (userDocument, error) {
	id, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return userDocument{}, err
	}

	return userDocument{
		ID:        id,
		Email:     string(u.Email),
		Password:  string(u.Password),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func toUserDomain(d userDocument) domain.User {
	return domain.User{
		ID:        d.ID.Hex(),
		Email:     domain.Email(d.Email),
		Password:  domain.Password(d.Password),
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
