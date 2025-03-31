package main

import (
	"reflect"
	"testing"
	"time"
)

func TestRemoveDuplicateAds(t *testing.T) {
	tests := []struct {
		adLinks         []AdLinkPair
		noRedirection   bool
		expectedAdLinks []AdLinkPair
		expectedError   bool
	}{
		{
			adLinks: []AdLinkPair{
				{OriginalAdURL: "http://example.com", FinalAdURL: "http://example.com"},
				{OriginalAdURL: "http://example.com", FinalAdURL: "http://example.com"},
			},
			noRedirection: true,
			expectedAdLinks: []AdLinkPair{
				{OriginalAdURL: "http://example.com", FinalAdURL: "http://example.com"},
			},
			expectedError: false,
		},
		{
			adLinks: []AdLinkPair{
				{OriginalAdURL: "http://example.com", FinalAdURL: "http://example.com"},
				{OriginalAdURL: "http://example.com", FinalAdURL: "http://example.com"},
			},
			noRedirection: false,
			expectedAdLinks: []AdLinkPair{
				{OriginalAdURL: "http://example.com", FinalAdURL: "http://example.com"},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, err := removeDuplicateAds(tt.adLinks, tt.noRedirection)
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
			if !reflect.DeepEqual(result, tt.expectedAdLinks) {
				t.Errorf("expected %v, got %v", tt.expectedAdLinks, result)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com", "https://example.com"},
		{"https://example.com", "https://example.com"},
		{"example.com", "https://example.com"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := normalizeURL(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDefangURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com", "http://example[.]com"},
		{"https://example.com", "https://example[.]com"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := defangURL(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"http://example.com", "example.com", false},
		{"https://www.example.com", "example.com", false},
		{"", "", true}, // Adding a case with an empty string
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := extractDomain(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v, got: %v", tt.hasError, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenerateAdResults(t *testing.T) {
	adLinks := []AdLinkPair{
		{OriginalAdURL: "http://example.com", FinalAdURL: "http://example.com"},
	}
	expectedResults := []AdResult{
		{Engine: "Google", Query: "test", OriginalAdURL: "http://example.com", FinalDomainURL: "example.com", FinalRedirectURL: "http://example.com", RedirectChain: nil, Time: time.Now()},
	}

	result, err := generateAdResults(adLinks, "test", "Google", time.Now(), true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != len(expectedResults) {
		t.Fatalf("expected %v results, got %v", len(expectedResults), len(result))
	}
}

func TestIsDomainExpected(t *testing.T) {
	tests := []struct {
		domain          string
		expectedDomains []string
		expected        bool
	}{
		{"example.com", []string{"example.com"}, true},
		{"example.com", []string{"test.com"}, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := isDomainExpected(tt.domain, tt.expectedDomains)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExportAdResults(t *testing.T) {
	adResults := []AdResult{
		{Engine: "Google", Query: "test", OriginalAdURL: "http://example.com", FinalDomainURL: "example.com", FinalRedirectURL: "http://example.com", RedirectChain: nil, Time: time.Now()},
	}
	err := exportAdResults("test.json", adResults)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBeginsWithHTTP(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"example.com", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := beginsWithHTTP(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCheckHostnameEndsWithDomain(t *testing.T) {
	tests := []struct {
		hostname string
		domain   string
		expected bool
	}{
		{"example.com", "com", true},
		{"example.com", "example.com", true},
		{"example.com", "test.com", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := checkHostnameEndsWithDomain(tt.hostname, tt.domain)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDecodeBase64(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"aHR0cDovL2V4YW1wbGUuY29t", "http://example.com"},
		{"invalid%base64", ""},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := decodeBase64(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMergeLists(t *testing.T) {
	tests := []struct {
		firstList  []string
		secondList []string
		expected   []string
	}{
		{
			firstList:  []string{"a", "b", "c"},
			secondList: []string{"b", "c", "d"},
			expected:   []string{"a", "b", "c", "d"},
		},
		{
			firstList:  []string{"apple", "banana"},
			secondList: []string{"banana", "cherry"},
			expected:   []string{"apple", "banana", "cherry"},
		},
		{
			firstList:  []string{"1", "2", "3"},
			secondList: []string{"4", "5", "6"},
			expected:   []string{"1", "2", "3", "4", "5", "6"},
		},
		{
			firstList:  []string{},
			secondList: []string{"x", "y", "z"},
			expected:   []string{"x", "y", "z"},
		},
		{
			firstList:  []string{"x", "y", "z"},
			secondList: []string{},
			expected:   []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := mergeLists(tt.firstList, tt.secondList)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
