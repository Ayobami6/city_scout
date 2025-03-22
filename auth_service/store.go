package authservice

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	db *mongo.Database
}

func NewUserStore(db *mongo.Database) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) CreateUser(ctx context.Context, user *User) error {
	_, err := s.db.Collection("users").InsertOne(ctx, user)
	return err
}

func (s *UserStore) GetUser(ctx context.Context, username string) (*User, error) {
	var user User
	err := s.db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return &user, err
}

func (s *UserStore) GetUserByAPIKey(ctx context.Context, apiKey string) (*User, error) {
	var user User
	err := s.db.Collection("users").FindOne(ctx, bson.M{"apiKey": apiKey}).Decode(&user)
	return &user, err
}
