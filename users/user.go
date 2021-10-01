package users

import (
	"context"

	"github.com/Bwise1/zuri_bot/mongo"
	"github.com/dghubble/go-twitter/twitter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	AccessToken  string             `json:"access_token,omitempty" bson:"access_token,omitempty"`
	AccessSecret string             `json:"access_secret" bson:"access_secret"`
	*twitter.User
}

type UserService interface {
	CreateUser(context.Context, *User) error
}

type userMongo struct {
	db *mongo.DB
}

func NewUserService(db *mongo.DB) UserService {
	return &userMongo{db}
}

func (um userMongo) CreateUser(ctx context.Context, u *User) error {
	coll := um.db.GetCollection("users")
	_, err := coll.InsertOne(ctx, u)
	return err
}
