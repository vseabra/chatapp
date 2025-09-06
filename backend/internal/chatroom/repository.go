package chatroom

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Create(ctx context.Context, c *ChatRoom) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*ChatRoom, error)
	FindAll(ctx context.Context, limit int64, skip int64) ([]ChatRoom, error)
	UpdateTitle(ctx context.Context, id primitive.ObjectID, title string) error
	Delete(ctx context.Context, id primitive.ObjectID, ownerID string) (bool, error)
}

type mongoRepository struct {
	col *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{col: db.Collection("chatrooms")}
}

func (r *mongoRepository) Create(ctx context.Context, c *ChatRoom) error {
	res, err := r.col.InsertOne(ctx, c)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		c.ID = oid
	}
	return nil
}

func (r *mongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*ChatRoom, error) {
	var c ChatRoom
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *mongoRepository) FindAll(ctx context.Context, limit int64, skip int64) ([]ChatRoom, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []ChatRoom
	for cursor.Next(ctx) {
		var c ChatRoom
		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, cursor.Err()
}

func (r *mongoRepository) UpdateTitle(ctx context.Context, id primitive.ObjectID, title string) error {
	_, err := r.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"title": title}})
	return err
}

func (r *mongoRepository) Delete(ctx context.Context, id primitive.ObjectID, ownerID string) (bool, error) {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id, "ownerId": ownerID})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
