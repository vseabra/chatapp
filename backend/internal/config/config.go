package config

import (
	"log"
	"os"
)

// AppConfig holds application configuration values.
type AppConfig struct {
	Port        string
	MongoURI    string
	RabbitMQURI string
	MongoDB     string
	JWTSecret   string
	JWTExpires  string
}

// Load returns configuration populated from environment variables with sane defaults.
func Load() AppConfig {
	port := getEnv("PORT", "8080")

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatalf("MONGODB_URI is required")
		return AppConfig{}
	}

	rabbitMQURI := os.Getenv("RABBITMQ_URI")
	if rabbitMQURI == "" {
		log.Fatalf("RABBITMQ_URI is required")
		return AppConfig{}
	}

	mongoDB := getEnv("MONGODB_DB", "chat")

	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-me")
	if os.Getenv("JWT_SECRET") == "" {
		log.Printf("JWT_SECRET not set, using development default")
	}
	jwtExpires := getEnv("JWT_EXPIRES_IN", "24h")
	if os.Getenv("JWT_EXPIRES_IN") == "" {
		log.Printf("JWT_EXPIRES_IN not set, using default: %s", jwtExpires)
	}

	return AppConfig{
		Port:        port,
		MongoURI:    mongoURI,
		MongoDB:     mongoDB,
		JWTSecret:   jwtSecret,
		JWTExpires:  jwtExpires,
		RabbitMQURI: rabbitMQURI,
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
