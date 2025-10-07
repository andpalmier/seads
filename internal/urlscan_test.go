package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSubmitURLScan_Success(t *testing.T) {
	// Make the handler assert the incoming request and return a valid JSON response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Method and headers
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("expected Content-Type application/json, got %q", ct)
		}
		if apiKey := r.Header.Get("API-Key"); apiKey != "TEST_TOKEN" {
			t.Fatalf("expected API-Key TEST_TOKEN, got %q", apiKey)
		}

		// Body
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if body["url"] != "https://example.com" {
			t.Fatalf("expected url=https://example.com, got %v", body["url"])
		}
		if body["visibility"] != "public" {
			t.Fatalf("expected visibility=public, got %v", body["visibility"])
		}
		tagsAny, ok := body["tags"]
		if !ok {
			t.Fatalf("expected tags field in request body")
		}
		tags, ok := tagsAny.([]any)
		if !ok {
			t.Fatalf("expected tags to be an array, got %T", tagsAny)
		}
		got := make([]string, 0, len(tags))
		for _, v := range tags {
			got = append(got, v.(string))
		}
		wantTags := []string{"tag1", "tag2", "tag3"}
		if strings.Join(got, ",") != strings.Join(wantTags, ",") {
			t.Fatalf("expected tags %v, got %v", wantTags, got)
		}

		// Respond with a realistic success payload
		resp := URLScanSubmissionResponse{
			Message:    "Submission successful",
			UUID:       "123e4567-e89b-12d3-a456-426614174000",
			Result:     "https://urlscan.io/result/123e4567-e89b-12d3-a456-426614174000/",
			API:        "https://urlscan.io/api/v1/result/123e4567-e89b-12d3-a456-426614174000/",
			Visibility: "public",
			Options:    Options{UserAgent: "test-agent"},
			URL:        "https://example.com",
			Country:    "SG",
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Avoid sleeping in tests
	URLScanSleepSeconds = 0

	cfg := Config{
		URLScanSubmitter: &URLScanSubmitter{
			Token:      "TEST_TOKEN",
			Tags:       "tag1,tag2,tag3",
			Visibility: "public",
			ScanURL:    ts.URL, // send mock server ts
		},
	}

	resp, err := SubmitURLScan(cfg, "https://example.com")
	if err != nil {
		t.Fatalf("SubmitURLScan returned error: %v", err)
	}

	if resp.Message != "Submission successful" ||
		resp.UUID == "" ||
		resp.Result == "" ||
		resp.API == "" ||
		resp.Visibility != "public" ||
		resp.Options.UserAgent != "test-agent" ||
		resp.URL != "https://example.com" ||
		resp.Country != "SG" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
