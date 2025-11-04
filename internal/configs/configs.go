package configs

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Bot struct {
	Key       string   `yaml:"key"`
	Whitelist []string `yaml:"whitelist"`
}

type Anthropic struct {
	System string `yaml:"system"`
	Key    string `yaml:"key"`
	Proxy  struct {
		Enabled bool   `yaml:"enabled"`
		Url     string `yaml:"url"`
	} `yaml:"proxy"`
}

type Logger struct {
	Enabled  bool   `yaml:"enabled"`
	FilePath string `yaml:"file_path"`
	Level    string `yaml:"level"`
}

type All struct {
	Logger    `yaml:"logger"`
	Bot       `yaml:"bot"`
	Anthropic `yaml:"anthropic"`
}

func New(path string) *All {
	slog.Info(path)
	conf, err := load(path)
	if err != nil {
		log.Fatalf("unable to load '%s': %v", path, err)
	}

	return conf
}

func load(path string) (*All, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error resolving config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &All{}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing yaml: %w", err)
	}

	if err := validate(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func validate(config *All) error {
	if config.Bot.Key == "" {
		return fmt.Errorf("telegram token is required")
	}
	if len(config.Whitelist) == 0 {
		return fmt.Errorf("whitelist is empty")
	}
	if config.Anthropic.Key == "" {
		return fmt.Errorf("anthropic API key is required")
	}
	if config.Logger.Enabled && config.Logger.FilePath == "" {
		return fmt.Errorf("logging enabled but file is not specified")
	}

	return nil
}
