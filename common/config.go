package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv    string
	MongoURI  string
	JWTSecret string
	PORT      string
}

var AppConfig *Config

func LoadConfig() {
	// Load the .env file
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	err := godotenv.Load(env + ".env")
	if err != nil {
		log.Fatalf("Error loading %s.env file: %v", env, err)
	}

	// Initialize configuration
	AppConfig = &Config{
		AppEnv:    os.Getenv("APP_ENV"),
		MongoURI:  os.Getenv("MONGO_URI"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		PORT:      os.Getenv("PORT"),
	}
}
