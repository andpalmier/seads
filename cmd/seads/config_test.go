package main

import (
	"os"
	"reflect"
	"testing"
)

// TestParseConfig tests the parseConfig function with a valid configuration file
func TestParseConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
mail:
  host: smtp.mailtrap.io
  port: 2525
  username: user123
  password: pass123
  from: no-reply@example.com
  recipients: [recipient1@example.com, recipient2@example.com]

slack:
  token: xoxb-1234-5678-91011-abcdef
  channels: [general, random]

telegram:
  token: 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
  chatids: [123456789, 987654321]

urlscan:
  token: "your-api-key"
  scanurl: "https://urlscan.io/api/v1/scan/"
  visibility: "private"
  tags: "private.mytag"

global-domain-exclusion:
  exclusion-list: [ebay.com, amazon.com]

queries:
  - query: "ipad"
    expected-domains: [apple.com]

  - query: "as roma"
    expected-domains: []
`
	tmpFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Parse the config file
	config, err := parseConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	// Expected config
	expectedConfig := Config{
		TelegramNotifier: &TelegramNotifier{Token: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11", ChatIDs: []string{"123456789", "987654321"}},
		SlackNotifier:    &SlackNotifier{Token: "xoxb-1234-5678-91011-abcdef", Channels: []string{"general", "random"}},
		MailNotifier:     &MailNotifier{Host: "smtp.mailtrap.io", Port: "2525", Username: "user123", Password: "pass123", From: "no-reply@example.com", Recipients: []string{"recipient1@example.com", "recipient2@example.com"}},
		URLScanSubmitter: &URLScanSubmitter{Token: "your-api-key", ScanURL: "https://urlscan.io/api/v1/scan/", Visibility: "private", Tags: "private.mytag"},
		GlobalDomainExclusion: &GlobalDomainExclusion{
			GlobalDomainExclusionList: []string{"ebay.com", "amazon.com"},
		},
		Queries: []SearchQuery{
			{SearchTerm: "ipad", ExpectedDomains: []string{"apple.com"}},
			{SearchTerm: "as roma", ExpectedDomains: []string{}},
		},
	}

	// Compare the parsed config with the expected config
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("Expected config %v, got %v", expectedConfig.TelegramNotifier.ChatIDs, config)
	}
}

// TestParseConfigFileNotFound tests the parseConfig function with a non-existent file
func TestParseConfigFileNotFound(t *testing.T) {
	_, err := parseConfig("non_existent_file.yaml")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
