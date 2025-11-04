package telegram

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"ai-telegram-bot/internal/configs"
	"ai-telegram-bot/internal/logger"
	"ai-telegram-bot/internal/provider"
	"ai-telegram-bot/internal/render"
	"ai-telegram-bot/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func FireBot(provider provider.AI, repo repository.DialogStorage, config configs.Bot) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithMiddlewares(verifyAccessMiddleware(config.Whitelist)),
		bot.WithDefaultHandler(handler(provider, repo)),
	}
	if logger.IsDebugEnabled() {
		opts = append(opts, bot.WithDebug())
	}
	b, err := bot.New(config.Key, opts...)
	if err != nil {
		log.Fatalf("unable to init bot: %v", err)
	}

	configureCommands(ctx, b, repo)

	b.Start(ctx)
}

func handler(provider provider.AI, repo repository.DialogStorage) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		replier := newReplier(b, update)
		message := update.Message
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic in handler", slog.Any("error", r), slog.Int64("chat_id", message.Chat.ID))
				replier.reply(ctx, fmt.Sprintf("panic handled: %v", r))
			}
		}()

		replier.sendChatTyping(ctx)

		providerResponse, err := provider.Send(ctx, message.Text)
		if err != nil {
			slog.Error("provider error", slog.Any("error", err), slog.Int64("chat_id", message.Chat.ID))
			replier.reply(ctx, fmt.Sprintf("provider respond with error: %v", err))
			return
		}

		content := providerResponse.Content[0]
		adjusted, err := render.AdjustMdToTelegramFormat(content.Text)
		if err != nil {
			slog.Error("markdown adjustment error", slog.Any("error", err), slog.Int64("chat_id", message.Chat.ID))
			return
		}

		_, err = replier.reply(ctx, adjusted)
		if err != nil {
			slog.Error("unable to send in response to telegram", slog.Any("error", err), slog.Int64("chat_id", message.Chat.ID))
		}
	}
}
