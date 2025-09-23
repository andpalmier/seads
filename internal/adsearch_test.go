package internal

import (
	"reflect"
	"testing"
)

func TestProcessAdResults(t *testing.T) {
	// Mock input data
	adResults := []AdResult{
		{OriginalAdURL: "http://example.com", FinalDomainURL: "example.com", ExpectedDomains: true},
		{OriginalAdURL: "http://unexpected.com", FinalDomainURL: "unexpected.com", ExpectedDomains: false},
	}
	expectedDomainList := []string{"example.com"}

	// Mock output slices
	var allAdResults []AdResult
	var notifications []AdResult

	// Enable necessary flags
	EnableNotifications = true
	PrintRedirectChain = false
	EnableURLScan = false

	// Mock config data
	config := Config{
		GlobalDomainExclusion: &GlobalDomainExclusion{
			GlobalDomainExclusionList: []string{},
		},
		Queries:          []SearchQuery{},
		URLScanSubmitter: nil, // Disable URLScan
	}

	// Call the function
	err := processAdResults(adResults, expectedDomainList, &allAdResults, &notifications, config)
	if err != nil {
		t.Fatalf("processAdResults returned an error: %v", err)
	}

	// Expected results
	expectedAllAdResults := adResults
	expectedNotifications := []AdResult{
		{OriginalAdURL: "http://unexpected.com", FinalDomainURL: "unexpected.com"},
	}

	// Assertions
	if !reflect.DeepEqual(allAdResults, expectedAllAdResults) {
		t.Errorf("allAdResults mismatch. Expected: %v, Got: %v", expectedAllAdResults, allAdResults)
	}
	if !reflect.DeepEqual(notifications, expectedNotifications) {
		t.Errorf("notifications mismatch. Expected: %v, Got: %v", expectedNotifications, notifications)
	}
}
