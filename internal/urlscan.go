package internal

import (
	"bytes"
	"encoding/json"
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
	Visibility string `yaml:"visibility"`
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

// SubmitURLScan submit single URL to the URLScan service for scanning
func SubmitURLScan(config Config, urlToScan string) (URLScanSubmissionResponse, error) {

	token := config.URLScanSubmitter.Token
	tags := config.URLScanSubmitter.Tags
	visibility := config.URLScanSubmitter.Visibility

	taglist := strings.Split(tags, ",")
	log.Printf("\n*** URLScan Enabled ***\n")
	log.Printf("Endpoint URL: %v\n", config.URLScanSubmitter.ScanURL)
	log.Printf("Visibility: %v\n", visibility)
	log.Printf("Tags: %v\n\n", tags)
	// UNCOMMENT
	//log.Printf("Total URLs to submit: %d\n", len(uniqueAdLinks))

	log.Printf("URL for submission: %s", urlToScan)
	data := map[string]interface{}{
		"url":        urlToScan,
		"visibility": visibility,
		"tags":       taglist,
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON: %v\n", err)
		return URLScanSubmissionResponse{}, err
	}

	// Create a POST request object and set appropriate headers
	req, err := http.NewRequest("POST", config.URLScanSubmitter.ScanURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return URLScanSubmissionResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", token)

	// Send the POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v\n", err)
		return URLScanSubmissionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("server returned non-200 status: %d\n", resp.StatusCode)
		return URLScanSubmissionResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return URLScanSubmissionResponse{}, err
	}

	var response URLScanSubmissionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		return URLScanSubmissionResponse{}, err
	}

	log.Printf("*************************\n")
	log.Printf("Message: %s", response.Message)
	log.Printf("UUID: %s", response.UUID)
	log.Printf("Result URL: %s", response.Result)
	log.Printf("API URL: %s", response.API)
	log.Printf("Visibility: %s", response.Visibility)
	log.Printf("User Agent: %s", response.Options.UserAgent)
	log.Printf("Original URL: %s", response.URL)
	log.Printf("Country: %s", response.Country)
	log.Printf("\n*************************\n")

	return response, err
}
