package internal

import (
	"os"

	"gopkg.in/yaml.v2"
)

// SearchQuery represents a search query
type SearchQuery struct {
	SearchTerm      string   `yaml:"query"`
	ExpectedDomains []string `yaml:"expected-domains"`
}

type GlobalDomainExclusion struct {
	GlobalDomainExclusionList []string `yaml:"exclusion-list"`
}

// Config holds the overall configuration
type Config struct {
	TelegramNotifier      *TelegramNotifier      `yaml:"telegram"`
	SlackNotifier         *SlackNotifier         `yaml:"slack"`
	MailNotifier          *MailNotifier          `yaml:"mail"`
	DiscordNotifier       *DiscordNotifier       `yaml:"discord"`
	URLScanSubmitter      *URLScanSubmitter      `yaml:"urlscan"`
	GlobalDomainExclusion *GlobalDomainExclusion `yaml:"global-domain-exclusion"`
	Queries               []SearchQuery          `yaml:"queries"`
}

// parseConfig parses the specified config file
func ParseConfig(configPath string) (Config, error) {
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
