package domain

import "github.com/go-telegram/bot"

type BotSettings struct {
	Token string
	Opts  []bot.Option
}

func NewBot(settings BotSettings) (*bot.Bot, error) {
	bot, err := bot.New(settings.Token, settings.Opts...)
	if err != nil {
		panic(err)
	}
	return bot, nil
}
