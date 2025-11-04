package provider

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/url"

	"ai-telegram-bot/internal/configs"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AI interface {
	Send(context.Context, string) (*anthropic.Message, error)
}

// TODO update anthropic deps
// TODO manage tokens sizing
func (a *AnthropicProvider) Send(ctx context.Context, text string) (*anthropic.Message, error) {
	messageParams := anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 4096,
		System:    []anthropic.TextBlockParam{{Text: a.systemPrompt}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(text)),
		},
	}

	response, err := a.client.Messages.New(ctx, messageParams)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type AnthropicProvider struct {
	client       *anthropic.Client
	systemPrompt string
}

func NewAnthropic(configs configs.Anthropic) AI {
	var client anthropic.Client
	if configs.Proxy.Enabled {
		proxyUrl, err := url.Parse(configs.Proxy.Url)
		if err != nil {
			log.Fatal("Invalid proxy URL:", err)
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}

		client = anthropic.NewClient(
			option.WithAPIKey(configs.Key),
			option.WithHTTPClient(httpClient),
		)
	} else {
		client = anthropic.NewClient(option.WithAPIKey(configs.Key))
	}

	slog.Info("current AI provider: Claude")

	return &AnthropicProvider{
		&client,
		configs.System,
	}
}
