package main

import (
	"log"
	"log/slog"
	"os"

	"ai-telegram-bot/internal/configs"
	"ai-telegram-bot/internal/logger"
	"ai-telegram-bot/internal/provider"
	"ai-telegram-bot/internal/repository"
	"ai-telegram-bot/internal/telegram"
)

func main() {
	config := configs.New(os.Getenv("AI_TELEGRAM_BOT_CONFIG"))
	loggerCfg := logger.Config{
		Enabled:  config.Enabled,
		FilePath: config.FilePath,
		Level:    config.Level,
	}
	if err := logger.Init(loggerCfg); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	provider := provider.NewAnthropic(config.Anthropic)
	repo := repository.NewInMemoryDialogStorage()

	slog.Info("all configurations applied")

	telegram.FireBot(provider, repo, config.Bot)
}
