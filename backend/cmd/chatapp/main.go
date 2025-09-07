package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"chatapp/internal/chatroom"
	"chatapp/internal/config"
	"chatapp/internal/db"
	"chatapp/internal/events"
	server "chatapp/internal/http"
	"chatapp/internal/message"
	"chatapp/internal/user"
	"chatapp/internal/ws"
)

func main() {
	cfg := config.Load()
	timeout := 10 * time.Second

	startupCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	appCtx := context.Background()

	client, err := db.ConnectMongo(startupCtx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	database := client.Database(cfg.MongoDB)
	if err := db.BuildIndexes(startupCtx, database); err != nil {
		log.Fatalf("failed ensuring indexes: %v", err)
	}

	r := server.NewRouter()

	userRepo := user.NewRepository(database)
	userService := user.NewService(userRepo)
	userController := user.NewHandler(userService)

	roomRepo := chatroom.NewRepository(database)
	roomService := chatroom.NewService(roomRepo)
	roomHandler := chatroom.NewHandler(roomService)

	msgRepo := message.NewRepository(database)
	msgService := message.NewService(msgRepo)
	msgHandler := message.NewHandler(msgService)

	// Register routes
	userController.RegisterRoutes(r)
	roomHandler.RegisterRoutes(r, cfg.JWTSecret)
	msgHandler.RegisterRoutes(r, cfg.JWTSecret)

	// WebSocket hub
	hub := ws.BuildHub()

	amq, err := events.NewAMQP(startupCtx, cfg.RabbitMQURI, "chat.events")
	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %v", err)
	}
	defer amq.Close()

	// Wire publisher and register ws routes
	pub := &ws.Publisher{AMQP: amq}
	hub.WithPublisher(pub)
	hub.RegisterRoutes(r, cfg)

	// Start consumers
	ingress := &events.IngressConsumer{AMQP: amq, Service: msgService, Users: userRepo}
	_ = ingress.Start(appCtx)
	broadcast := &events.BroadcastConsumer{AMQP: amq, Hub: hub}
	_ = broadcast.Start(appCtx)
	botResp := &events.BotResponseConsumer{AMQP: amq, Service: msgService}
	_ = botResp.Start(appCtx)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	log.Printf("server listening on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
