package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	AIToken  string
	ProxyURL string
}

func Load() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf(".env not found")
	}
	return &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
		//	AIToken:  ai,
		ProxyURL: os.Getenv("PROXY_URL"),
	}
}
