package internal

import (
	"reflect"
	"testing"
)

func TestProcessAdResults(t *testing.T) {
	// Mock input data
	mockAdResults := []AdResult{
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
	err := processAdResults(mockAdResults, expectedDomainList, &allAdResults, &notifications, config)
	if err != nil {
		t.Fatalf("processAdResults returned an error: %v", err)
	}

	// Expected results
	expectedAllAdResults := mockAdResults
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

// TestGlobalDomainExclusionList ensures the global exclusion list is properly initialized and used.
func TestGlobalDomainExclusionList(t *testing.T) {
	// Mock input data
	// Here the domain inside GlobalDomainExclusionList should have ExpectedDomain flagged with true
	mockAdResults := []AdResult{
		{OriginalAdURL: "http://example.com", FinalDomainURL: "example.com", ExpectedDomains: true},
		{OriginalAdURL: "http://unexpected.com", FinalDomainURL: "unexpected.com", ExpectedDomains: false},
	}

	if GlobalDomainExclusionList == nil {
		t.Fatal("GlobalDomainExclusionList should not be nil. It has to be initialized.")
	}

	// Enable necessary flags
	EnableNotifications = true
	PrintRedirectChain = false
	EnableURLScan = false
	GlobalDomainExclusionList = []string{"example.com"}

	if len(GlobalDomainExclusionList) == 0 {
		t.Fatal("GlobalDomainExclusionList is empty")
	}

	expectedGlobalDomainExclusionList := GlobalDomainExclusionList

	// Mock output slices
	var allAdResults []AdResult
	var notifications []AdResult

	err := processAdResults(mockAdResults, expectedGlobalDomainExclusionList, &allAdResults, &notifications, Config{})
	if err != nil {
		t.Errorf("processAdResults returned error: %v", err)
	}

	// Expected results
	expectedAllAdResults := mockAdResults
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

func TestIsInSelectedSearchEngineList(t *testing.T) {
	// Mock user defined list of SelectedEngine separated by comma
	SelectedEngine = "google, aol    ,syndicated, google1, typo-edsearch engine,yah oo, yahoo"

	// From available engine check if the Selected Engine
	testMockSelection := make(map[string]bool) // instantiate to all false
	for _, e := range searchEnginesFunctions {
		testMockSelection[e.EngineName] = false
	}

	for _, engine := range searchEnginesFunctions {
		testMockSelection[engine.EngineName] = isInSelectedSearchEngineList(engine.EngineName)
	}

	if testMockSelection["google"] == false {
		t.Errorf("Search engine 'google' is in test map but not flagged as expected")
	}

	if testMockSelection["aol"] == false {
		t.Errorf("Search engine 'aol' is in test map but not flagged as expected")
	}

	if testMockSelection["syndicated"] == false {
		t.Errorf("Search engine 'syndicated' is in test map but not flagged as expected")
	}

	if testMockSelection["bing"] == true {
		t.Errorf("Search engine 'bing' is NOT test map but flagged as expected")
	}

	if _, exists := testMockSelection["google1"]; exists {
		t.Errorf("Search engine 'google1' is NOT a valid search engine but flagged as expected")
	}

	if _, exists := testMockSelection["typo-edsearch engine"]; exists {
		t.Errorf("Search engine 'typo-edsearch engine' is NOT a valid search engine but flagged as expected")
	}

	if _, exists := testMockSelection["yah oo"]; exists {
		t.Errorf("Search engine 'typo-edsearch engine' is NOT a valid search engine but flagged as expected")
	}

	if testMockSelection["yahoo"] == false {
		t.Errorf("Search engine 'aol' is in test map but not flagged as expected")
	}
}

func TestSearchEngineUnused(t *testing.T) {
	// Mock user defined list of SelectedEngine separated by comma
	SelectedEngine = ""

	// From available engine check if the Selected Engine
	testMockSelection := make(map[string]bool) // instantiate to all false
	for _, engine := range searchEnginesFunctions {
		testMockSelection[engine.EngineName] = false
	}

	for _, engine := range searchEnginesFunctions {
		testMockSelection[engine.EngineName] = isInSelectedSearchEngineList(engine.EngineName)
	}

	for engineName, boolActive := range testMockSelection {
		if boolActive == false {
			t.Errorf("Search engine '%s' is in test map but not flagged as expected", engineName)
		}
	}
}
