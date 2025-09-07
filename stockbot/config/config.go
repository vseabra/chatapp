package config

import (
	"log"
	"os"
)

// AppConfig holds application configuration values.
type AppConfig struct {
	RabbitMQURI string
}

// Load returns configuration populated from environment variables with sane defaults.
func Load() AppConfig {

	rabbitMQURI := os.Getenv("RABBITMQ_URI")
	if rabbitMQURI == "" {
		log.Fatalf("RABBITMQ_URI is required")
		return AppConfig{}
	}


	return AppConfig{
		RabbitMQURI: rabbitMQURI,
	}
}
