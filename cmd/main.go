package main

import (
	"context"
	"corp-rag-tg-bot/config"
	"corp-rag-tg-bot/internal/controller"
	"corp-rag-tg-bot/internal/controller/telegram"
	"corp-rag-tg-bot/internal/domain"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cfg := config.Load()

	opts := []bot.Option{
		bot.WithDefaultHandler(telegram.HelloHandler),
		bot.WithSkipGetMe(),
	}

	proxy := controller.NewProxy(cfg.ProxyURL)

	opts = append(opts, bot.WithHTTPClient(30*time.Second, proxy))

	settings := domain.BotSettings{
		Token: cfg.BotToken,
		Opts:  opts,
	}

	bot, err := domain.NewBot(settings)
	if err != nil {
		return
	}
	bot.Start(ctx)
}
