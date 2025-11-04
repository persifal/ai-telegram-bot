package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type replier struct {
	b      *bot.Bot
	update *models.Update
}

func newReplier(b *bot.Bot, update *models.Update) *replier {
	return &replier{
		b:      b,
		update: update,
	}
}

func (r *replier) sendChatTyping(ctx context.Context) {
	_, err := r.b.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: r.update.Message.Chat.ID,
		Action: models.ChatActionTyping,
	})
	if err != nil {
		slog.Error("unable to send chat typing action", slog.Any("error", err), slog.Int64("chat_id", r.update.Message.Chat.ID))
	}
}

func (r *replier) reply(ctx context.Context, text string) (*models.Message, error) {
	sendMessageParams := &bot.SendMessageParams{
		ChatID:    r.update.Message.Chat.ID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
		ReplyParameters: &models.ReplyParameters{
			MessageID: r.update.Message.ID,
		},
	}

	retryMax := 64
	retryStart := 1
	var msg *models.Message
	var err error
	for retryStart <= retryMax {
		msg, err = r.b.SendMessage(ctx, sendMessageParams)
		if err != nil {
			l := fmt.Sprintf("unable to send message to telegram. Retry with %d.\nerr: %v\nparams: %v", retryStart, err, sendMessageParams)
			slog.Error(l)
			retryStart <<= 1
			err = nil
		} else {
			break
		}
	}

	return msg, err
}

func (r *replier) getChatId() int64 {
	return r.update.Message.Chat.ID
}
