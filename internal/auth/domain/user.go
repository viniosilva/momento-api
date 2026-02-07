package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const UsersCollectionName = "users"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     Email              `bson:"email"`
	Password  Password           `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func NewUser(email Email, password Password) User {
	now := time.Now()

	return User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Update(email Email, password Password) {
	u.Email = email
	u.Password = password
	u.UpdatedAt = time.Now()
}
