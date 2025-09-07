package message

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Insert(ctx context.Context, m *Message) error
	ListByRoom(ctx context.Context, roomID string, limit int64, cursor string) ([]Message, string, error)
}

type mongoRepository struct {
	col *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{col: db.Collection("messages")}
}

func (r *mongoRepository) Insert(ctx context.Context, m *Message) error {
	res, err := r.col.InsertOne(ctx, m)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		m.ID = oid
	}
	return nil
}

// ListByRoom returns messages newest-first with a simple opaque cursor based on ObjectID.
func (r *mongoRepository) ListByRoom(ctx context.Context, roomID string, limit int64, cursor string) ([]Message, string, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	filter := bson.M{"roomId": roomID}
	if cursor != "" {
		if oid, err := primitive.ObjectIDFromHex(cursor); err == nil {
			filter["_id"] = bson.M{"$lt": oid}
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(limit)
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, "", err
	}
	defer cur.Close(ctx)

	var items []Message
	var lastID primitive.ObjectID
	for cur.Next(ctx) {
		var m Message
		if err := cur.Decode(&m); err != nil {
			return nil, "", err
		}
		items = append(items, m)
		lastID = m.ID
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}

	next := ""
	if len(items) == int(limit) && !lastID.IsZero() {
		next = lastID.Hex()
	}
	return items, next, nil
}
