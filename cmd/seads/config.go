package main

import (
	"gopkg.in/yaml.v2"
	"os"
)

// SearchQuery represents a search query
type SearchQuery struct {
	SearchTerm      string   `yaml:"query"`
	ExpectedDomains []string `yaml:"expected-domains"`
}

// Config holds the overall configuration
type Config struct {
	TelegramNotifier *TelegramNotifier `yaml:"telegram"`
	SlackNotifier    *SlackNotifier    `yaml:"slack"`
	MailNotifier     *MailNotifier     `yaml:"mail"`
	Queries          []SearchQuery     `yaml:"queries"`
}

// parseConfig parses the specified config file
func parseConfig(configPath string) (Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
