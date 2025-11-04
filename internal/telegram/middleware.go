package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func verifyAccessMiddleware(list []string) bot.Middleware {
	w := newWhiteList(list)

	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			replier := newReplier(b, update)
			from := update.Message.From
			if !w.contains(from.ID) {
				from := update.Message.From
				about := fmt.Sprintf("@%s, %s %s, is bot: %v", from.Username, from.FirstName, from.LastName, from.IsBot)
				slog.Warn("unauthorized access: " + about)
				replier.reply(ctx, "you are not authorized")
				return
			}

			next(ctx, b, update)
		}
	}
}
