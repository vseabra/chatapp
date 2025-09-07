package main

import (
	"encoding/json"
	"log"

	"stockbot/config"
	botamqp "stockbot/internal/amqp"
	"stockbot/internal/commands"
	"stockbot/internal/contracts"
	"stockbot/internal/echo"
	"stockbot/internal/help"
	"stockbot/internal/stock"
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

	log.Printf("stockbot listening on %s (%s)", cfg.QueueName, cfg.RequestedRoutingKey)
	stockHandler := stock.NewHandler(cfg.StockCSVURLTemplate, nil)

	// Build command registry
	reg := commands.NewRegistry()

	reg.Register("echo", echo.Handle)
	reg.Register("help", help.Handle)
	reg.Register("stock", stockHandler.Handle)

	for d := range msgs {
		var req contracts.BotRequest
		if err := json.Unmarshal(d.Body, &req); err != nil {
			log.Printf("skip invalid payload: %v", err)
			continue
		}
		if resp, ok := reg.Dispatch(req); ok {
			if err := client.PublishJSON(cfg.ExchangeName, cfg.ResponseRoutingKey, resp); err != nil {
				log.Printf("publish response error: %v", err)
			} else {
				log.Printf("handled %s for room %s", req.Command, resp.RoomID)
			}
		}
	}
}
