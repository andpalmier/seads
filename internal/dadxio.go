package internal

import (
	"fmt"
	"net/url"
)

// ResolveDadxioAdURL parses a d.adx.io URL and extracts the final redirect URL
func ResolveDadxioAdURL(DadxioAdURL string) (string, error) {

	// Parse the unescaped URL
	parsedURL, err := url.Parse(DadxioAdURL)
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("Skipping invalid d.adx.io URL: %s, Error: %v\n", DadxioAdURL, err)
	}

	// Extract query parameters from the parsed URL
	queryParams := parsedURL.Query()
	destURL := queryParams.Get("xu")

	test, err := url.Parse(destURL)
	if err != nil || test.Host == "" {
		return "", fmt.Errorf("Skipping invalid d.adx.io URL: %s, Error: %v\n", destURL, err)
	}

	return destURL, nil
}
