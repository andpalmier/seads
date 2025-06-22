package internal

import (
	"fmt"
	"net/url"
)

// ResolveDoubleClickAdURL parses a Doubleclick URL and extracts the final redirect URL
func ResolveDoubleClickAdURL(doubleClickURL string) (string, error) {

	// Parse the unescaped URL
	parsedURL, err := url.Parse(doubleClickURL)
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("Skipping invalid DoubleClick URL: %s, Error: %v\n", doubleClickURL, err)
	}

	// Extract query parameters from the parsed URL
	queryParams := parsedURL.Query()
	destURL := queryParams.Get("ds_dest_url")

	// Special case for d.agkn.com
	if destURL == "https://d.agkn.com/pixel/4102/?che=" {
		destURL = queryParams.Get("l0")
	}

	test, err := url.Parse(destURL)
	if err != nil || test.Host == "" {
		return "", fmt.Errorf("Skipping invalid DoubleClick URL: %s, Error: %v\n", destURL, err)
	}

	return destURL, nil
}
