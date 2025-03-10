package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// URLScanSubmitter holds configurations for sending URL with unexpected domain to URLScan
type URLScanSubmitter struct {
	Token      string `yaml:"token"`
	ScanURL    string `yaml:"scanurl"`
	Tags       string `yaml:"tags"`
	Visibility string `yaml:visibility`
}

// URLScanSubmissionResponse represents the response from URLScan
type URLScanSubmissionResponse struct {
	Message    string  `json:"message"`
	UUID       string  `json:"uuid"`
	Result     string  `json:"result"`
	API        string  `json:"api"`
	Visibility string  `json:"visibility"`
	Options    Options `json:"options"`
	URL        string  `json:"url"`
	Country    string  `json:"country"`
}

type Options struct {
	UserAgent string `json:"useragent"`
}

// Deduplicate URLs
func deduplicateURLs(adsToScan []AdResult) []string {
	var uniqueAdLinks []string
	seenURLs := make(map[string]struct{})
	for _, ads := range adsToScan {
		if _, seen := seenURLs[ads.OriginalAdURL]; !seen {
			uniqueAdLinks = append(uniqueAdLinks, ads.OriginalAdURL)
			seenURLs[ads.OriginalAdURL] = struct{}{}
		}
	}

	return uniqueAdLinks
}

// SubmitURLScan submits the URL to URLScan for scanning
func (config *Config) submitURLScan(adsToScan []AdResult) {
	uniqueAdLinks := deduplicateURLs(adsToScan)

	token := config.URLScanSubmitter.Token
	url := config.URLScanSubmitter.ScanURL
	tags := config.URLScanSubmitter.Tags
	visibility := config.URLScanSubmitter.Visibility

	tagslist := strings.Split(tags, ",")
	fmt.Println()
	fmt.Println("*** URLScan Enable ***")
	fmt.Printf("Endpoint URL: %v\n", url)
	fmt.Printf("Visibility: %v\n", visibility)
	fmt.Printf("Tags: %v\n", tags)
	fmt.Println()

	fmt.Println("Total URLs to submit: %d", len(uniqueAdLinks))

	for _, urlToScan := range uniqueAdLinks {
		// Create the data payload as a map
		log.Printf("URL for submission: %s", urlToScan)

		data := map[string]interface{}{
			"url":        urlToScan,
			"visibility": visibility,
			"tags":       tagslist,
		}

		// Convert data to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		// Create a POST request object
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

		// Set necessary headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("API-Key", token)

		// Send the POST request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("server returned non-200 status: %d", resp.StatusCode)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("failed to read response body: %v", err)
		} else {
			// Parse the URLScan response
			var response URLScanSubmissionResponse
			if err := json.Unmarshal(body, &response); err != nil {
				log.Printf("failed to parse JSON response: %v", err)
			}
			// Print the URLScan response
			log.Println("*************************")
			log.Println("Message:", response.Message)
			log.Println("UUID:", response.UUID)
			log.Println("Result URL:", response.Result)
			log.Println("API URL:", response.API)
			log.Println("Visibility:", response.Visibility)
			log.Println("User Agent:", response.Options.UserAgent)
			log.Println("Original URL:", response.URL)
			log.Println("Country:", response.Country)
			log.Println("*************************")
		}
	}
}
