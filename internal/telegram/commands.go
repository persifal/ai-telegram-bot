package telegram

import (
	"context"
	"log/slog"

	"ai-telegram-bot/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func configureCommands(ctx context.Context, client *bot.Bot, repo repository.DialogStorage) {
	client.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeExact, startHandler(repo))
	client.RegisterHandler(bot.HandlerTypeMessageText, "end", bot.MatchTypeExact, endHandler(repo))
	client.RegisterHandler(bot.HandlerTypeMessageText, "reset", bot.MatchTypeExact, resetHandler(repo))

	client.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "start a new dialog"},
			{Command: "end", Description: "end current dialog"},
			{Command: "reset", Description: "reset current dialog if any"},
		},
	})

	slog.Info("bot commands configured")
}

func startHandler(repo repository.DialogStorage) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		replier := newReplier(b, update)
		chatId := replier.getChatId()
		if repo.Exists(chatId) {
			replier.reply(ctx, "you already in dialog")
		} else {
			repo.Create(chatId)
			replier.reply(ctx, "dialog has started")
		}
	}
}

func endHandler(repo repository.DialogStorage) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		replier := newReplier(b, update)
		repo.Remove(update.Message.Chat.ID)
		replier.reply(ctx, "conversation closed and context cleared")
	}
}

func resetHandler(repo repository.DialogStorage) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		replier := newReplier(b, update)
		chatId := replier.getChatId()
		ok := repo.Exists(chatId)
		if ok {
			repo.Clear(chatId)
			replier.reply(ctx, "context cleared")
		} else {
			replier.reply(ctx, "no conversation found")
		}
	}
}
