package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	Enabled  bool   `yaml:"enabled"`
	FilePath string `yaml:"file_path"`
	Level    string `yaml:"level"`
}

var globalLogger *slog.Logger

func Init(cfg Config) error {
	var writers []io.Writer

	writers = append(writers, os.Stderr)

	if cfg.Enabled && cfg.FilePath != "" {
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}

		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		fullname := cfg.FilePath + "-" + timestamp

		file, err := os.OpenFile(fullname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return err
		}

		writers = append(writers, file)
	}

	multiWriter := io.MultiWriter(writers...)
	level := parseLevel(cfg.Level)
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(multiWriter, opts)
	globalLogger = slog.New(handler)

	slog.SetDefault(globalLogger)

	return nil
}

func IsDebugEnabled() bool {
	return slog.Default().Enabled(context.Background(), slog.LevelDebug)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Get() *slog.Logger {
	if globalLogger == nil {
		return slog.Default()
	}
	return globalLogger
}
