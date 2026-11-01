package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

type BotServer struct {
	bot *bot.Bot
}

func NewBotServer(token string, handler *BotHandler, clientOpts []bot.Option) (*BotServer, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler.Default),
	}

	opts = append(opts, clientOpts...)

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start bot: %w", err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handler.Start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handler.Help)

	return &BotServer{
		bot: b,
	}, nil
}

func (b *BotServer) Start(ctx context.Context) {
	b.bot.Start(ctx)
}
