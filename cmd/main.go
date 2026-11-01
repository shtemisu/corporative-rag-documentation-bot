package main

import (
	"context"
	"corp-rag-tg-bot/config"
	"corp-rag-tg-bot/internal/controller"
	"corp-rag-tg-bot/internal/controller/telegram"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.Load()
	var networkOpts []bot.Option
	networkOpts = append(networkOpts, bot.WithSkipGetMe())
	if cfg.ProxyURL != "" {
		proxy := controller.NewProxy(cfg.ProxyURL)
		networkOpts = append(networkOpts, bot.WithHTTPClient(30*time.Second, proxy))
	}

	handler := telegram.NewBotHandler()
	server, err := telegram.NewBotServer(cfg.BotToken, handler, networkOpts)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Bot started...")
	server.Start(ctx)
}
