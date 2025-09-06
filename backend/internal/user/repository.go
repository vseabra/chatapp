package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type mongoRepository struct {
	col *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{col: db.Collection("users")}
}

func (r *mongoRepository) Create(ctx context.Context, u *User) error {
	res, err := r.col.InsertOne(ctx, u)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid
	}
	return nil
}

func (r *mongoRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	if err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
