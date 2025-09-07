package main

import (
	"encoding/json"
	"log"

	"stockbot/config"
	botamqp "stockbot/internal/amqp"
	"stockbot/internal/echo"
)

func main() {
	cfg := config.Load()

	client, err := botamqp.New(cfg.RabbitMQURI)
	if err != nil {
		log.Fatalf("amqp connect: %v", err)
	}
	defer client.Close()

	if err := client.Declare(cfg.ExchangeName, cfg.ExchangeType, cfg.QueueName, cfg.RequestedRoutingKey); err != nil {
		log.Fatalf("amqp declare: %v", err)
	}

	msgs, err := client.Consume(cfg.QueueName)
	if err != nil {
		log.Fatalf("amqp consume: %v", err)
	}

	log.Printf("stockbot echo listening on %s (%s)", cfg.QueueName, cfg.RequestedRoutingKey)

	for d := range msgs {
		var req echo.BotRequested
		if err := json.Unmarshal(d.Body, &req); err != nil {
			log.Printf("skip invalid payload: %v", err)
			continue
		}
		if resp, ok := echo.Handle(req); ok {
			if err := client.PublishJSON(cfg.ExchangeName, cfg.ResponseRoutingKey, resp); err != nil {
				log.Printf("publish response error: %v", err)
			} else {
				log.Printf("echoed to room %s", resp.RoomID)
			}
		}
	}
}
