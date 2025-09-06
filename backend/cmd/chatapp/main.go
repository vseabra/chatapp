package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"chatapp/internal/chatroom"
	"chatapp/internal/config"
	"chatapp/internal/db"
	server "chatapp/internal/http"
	"chatapp/internal/user"
)

func main() {
	cfg := config.Load()
	timeout := 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := db.ConnectMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	database := client.Database(cfg.MongoDB)
	if err := db.BuildIndexes(ctx, database); err != nil {
		log.Fatalf("failed ensuring indexes: %v", err)
	}

	r := server.NewRouter()

	userRepo := user.NewRepository(database)
	userService := user.NewService(userRepo)
	userController := user.NewHandler(userService)

	roomRepo := chatroom.NewRepository(database)
	roomService := chatroom.NewService(roomRepo)
	roomHandler := chatroom.NewHandler(roomService)

	// Register routes
	userController.RegisterRoutes(r)
	roomHandler.RegisterRoutes(r, cfg.JWTSecret)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	log.Printf("server listening on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
