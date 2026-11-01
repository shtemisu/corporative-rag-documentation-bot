package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotHandler struct {
}

func NewBotHandler() *BotHandler {
	return &BotHandler{}
}

func (h *BotHandler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello! I am Documento!! I am your assistant in the company! Wish you start registration?",
	})
}

// Help обрабатывает команду /help
func (h *BotHandler) Help(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Инструкция:\nПросто напишите ваш вопрос обычным текстом, и я найду ответ в базе знаний.",
	})
}

// Default обрабатывает всё, что не подошло под конкретные команды (основной канал для вопросов RAG)
func (h *BotHandler) Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	// Сейчас это эхо-дублирование, позже здесь будет вызов RAG-пайплайна:
	// answer := h.ragUseCase.SearchAnswers(update.Message.Text)
	userText := update.Message.Text

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Вы спросили: " + userText,
	})
}
