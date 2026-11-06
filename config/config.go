package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken   string
	OllamaURL  string
	GenModel   string
	EmbedModel string
	ProxyURL   string
}

func Load() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf(".env not found")
	}
	return &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		OllamaURL:  os.Getenv("OLLAMA_URL"),
		GenModel:   os.Getenv("GEN_MODEL"),
		EmbedModel: os.Getenv("EMBED_MODEL"),
		ProxyURL:   os.Getenv("PROXY_URL"),
	}
}
