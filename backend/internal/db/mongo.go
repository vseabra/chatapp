package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectMongo initializes and returns a Mongo client connected to the given URI.
func ConnectMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}

func BuildIndexes(ctx context.Context, db *mongo.Database) error {
	if err := BuildUserIndexes(ctx, db); err != nil {
		return fmt.Errorf("error ensuring user indexes: %v", err)
	}

	return nil
}

// BuildUserIndexes ensures unique indexes for the users collection.
func BuildUserIndexes(ctx context.Context, db *mongo.Database) error {
	users := db.Collection("users")
	model := mongo.IndexModel{
		Keys:    map[string]int{"email": 1},
		Options: options.Index().SetUnique(true).SetName("uq_users_email"),
	}
	_, err := users.Indexes().CreateOne(ctx, model)
	if err != nil {
		log.Printf("error creating users.email unique index: %v", err)
		return err
	}
	return nil
}
