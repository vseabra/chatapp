package config

import (
	"log"
	"os"
)

// AppConfig holds application configuration values.
type AppConfig struct {
	RabbitMQURI         string
	ExchangeName        string
	ExchangeType        string
	RequestedRoutingKey string
	ResponseRoutingKey  string
	QueueName           string
}

// Load returns configuration populated from environment variables with sane defaults.
func Load() AppConfig {

	rabbitMQURI := os.Getenv("RABBITMQ_URI")
	if rabbitMQURI == "" {
		log.Fatalf("RABBITMQ_URI is required")
		return AppConfig{}
	}

	return AppConfig{
		RabbitMQURI:         rabbitMQURI,
		ExchangeName:        getEnv("RABBITMQ_EXCHANGE", "chat.events"),
		ExchangeType:        getEnv("RABBITMQ_EXCHANGE_TYPE", "topic"),
		RequestedRoutingKey: getEnv("RABBITMQ_KEY_REQUESTED", "bot.requested"),
		ResponseRoutingKey:  getEnv("RABBITMQ_KEY_RESPONSE", "bot.response.submit"),
		QueueName:           getEnv("RABBITMQ_QUEUE", "bot.stockbot"),
	}
}

// getEnv returns env var value or the provided default when empty.
func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
