package config

import (
	"log"
	"os"
)

// AppConfig holds application configuration values.
type AppConfig struct {
	Port       string
	MongoURI   string
	MongoDB    string
	JWTSecret  string
	JWTExpires string
}

// Load returns configuration populated from environment variables with sane defaults.
func Load() AppConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://root:example@mongodb:27017/chat?authSource=admin"
		log.Printf("MONGODB_URI not set, using default: %s", mongoURI)
	}

	mongoDB := os.Getenv("MONGODB_DB")
	if mongoDB == "" {
		mongoDB = "chat"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
		log.Printf("JWT_SECRET not set, using development default")
	}
	jwtExpires := os.Getenv("JWT_EXPIRES_IN")
	if jwtExpires == "" {
		jwtExpires = "24h"
		log.Printf("JWT_EXPIRES_IN not set, using default: %s", jwtExpires)
	}

	return AppConfig{
		Port:       port,
		MongoURI:   mongoURI,
		MongoDB:    mongoDB,
		JWTSecret:  jwtSecret,
		JWTExpires: jwtExpires,
	}
}
