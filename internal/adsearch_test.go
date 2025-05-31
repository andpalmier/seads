package internal

import (
	"reflect"
	"testing"
)

func TestProcessAdResults(t *testing.T) {
	// Mock input data
	adResults := []AdResult{
		{OriginalAdURL: "http://example.com", FinalDomainURL: "example.com"},
		{OriginalAdURL: "http://unexpected.com", FinalDomainURL: "unexpected.com"},
	}
	expectedDomainList := []string{"example.com"}

	// Mock output slices
	var allAdResults []AdResult
	var notifications []AdResult
	var submitToURLScan []AdResult

	// Enable necessary flags
	EnableURLScan = true
	EnableNotifications = true
	PrintRedirectChain = false

	// Call the function
	err := processAdResults(adResults, expectedDomainList, &allAdResults, &notifications, &submitToURLScan)
	if err != nil {
		t.Fatalf("processAdResults returned an error: %v", err)
	}

	// Expected results
	expectedAllAdResults := adResults
	expectedNotifications := []AdResult{
		{OriginalAdURL: "http://unexpected.com", FinalDomainURL: "unexpected.com"},
	}
	expectedSubmitToURLScan := []AdResult{
		{OriginalAdURL: "http://unexpected.com", FinalDomainURL: "unexpected.com"},
	}

	// Assertions
	if !reflect.DeepEqual(allAdResults, expectedAllAdResults) {
		t.Errorf("allAdResults mismatch. Expected: %v, Got: %v", expectedAllAdResults, allAdResults)
	}
	if !reflect.DeepEqual(notifications, expectedNotifications) {
		t.Errorf("notifications mismatch. Expected: %v, Got: %v", expectedNotifications, notifications)
	}
	if !reflect.DeepEqual(submitToURLScan, expectedSubmitToURLScan) {
		t.Errorf("submitToURLScan mismatch. Expected: %v, Got: %v", expectedSubmitToURLScan, submitToURLScan)
	}
}
