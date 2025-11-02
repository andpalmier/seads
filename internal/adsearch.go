package internal

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// extractAds extracts advertisement information from a search results page
func extractAds(browser *rod.Browser, page *rod.Page, userAgent, linkSelector, attrName, query, engine string, noRedirectionFlag bool) ([]AdResult, error) {
	// Find all ad elements on the page
	adElements, err := page.Elements(linkSelector)
	if err != nil {
		return nil, fmt.Errorf("unable to find ad elements: %v", err)
	}

	// Get additional ad details for Google-like engines (Google, Syndicated, AdsenseAds)
	var adDetails rod.Elements
	if isGoogleLikeEngine(engine) {
		adDetails, _ = page.Elements(adinfoSelector)
	}

	// Set up concurrent processing
	var (
		adsFound []AdResult     // Slice to store found ads
		mu       sync.Mutex     // Mutex to protect concurrent slice access
		wg       sync.WaitGroup // WaitGroup for goroutine synchronization
	)

	// Process each ad element concurrently
	for i, adElement := range adElements {
		// Get the ad URL from the element
		adURL, err := adElement.Attribute(attrName)
		if err != nil || adURL == nil {
			continue
		}

		wg.Add(1)
		go func(i int, adURL string) {
			defer wg.Done()

			// Initialize basic ad information
			ad := AdResult{
				OriginalAdURL: adURL,
				Query:         query,
				Time:          time.Now(),
				Engine:        engine,
			}

			// Extract advertiser info and location for Google-like engines
			if isGoogleLikeEngine(engine) && i < len(adDetails) {
				ad.Advertiser, ad.Location = "", ""
				adInfo, err := getAdInfo(browser, adDetails[i])
				if err == nil && len(adInfo) >= 2 {
					ad.Advertiser, ad.Location, ad.AdInfoURL = adInfo[0], adInfo[1], adInfo[2]
				}
			}

			// Follow redirect chain if enabled
			if !noRedirectionFlag {
				ad.FinalRedirectURL, ad.FinalDomainURL = followAdRedirect(browser, adURL, userAgent)
			}

			// Resolve ad URL for additional information
			ResolveAdUrl(adURL, &ad)

			// Safely append the ad to results
			mu.Lock()
			adsFound = append(adsFound, ad)
			mu.Unlock()
		}(i, *adURL)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	return adsFound, nil
}

// followAdRedirect follows the ad URL and returns the final URL and domain
func followAdRedirect(browser *rod.Browser, adURL, userAgent string) (finalURL, finalDomain string) {
	adPage := browser.MustPage()
	defer adPage.Close()

	// Set custom user agent if provided
	if userAgent != "" {
		_ = adPage.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent})
	}

	// Navigate to the ad URL and wait for completion
	wait := adPage.MustWaitNavigation()
	if err := adPage.Navigate(adURL); err == nil {
		wait()
		finalURL = adPage.MustInfo().URL
		finalDomain, _ = extractDomain(finalURL)
	}
	return
}

// searchAdsWithEngine performs concurrent ad searches using a specific search engine
func searchAdsWithEngine(
	engineFunc func(string, string, string, bool) ([]AdResult, error),
	query, engineName string, userAgent string, noRedirection bool) ([]AdResult, error) {
	encodedQuery := url.QueryEscape(query)

	if Logger {
		safePrintf(nil, "Searching ads on %s\n", searchEngineURLs[engineName]+encodedQuery)
	}

	// Collect ads using concurrent workers
	ads, err := runConcurrentSearch(engineFunc, encodedQuery, engineName, userAgent, noRedirection)
	if err != nil {
		return nil, fmt.Errorf("search failed for %s: %v", engineName, err)
	}

	// Process the collected ads
	return processSearchResults(ads, userAgent, noRedirection)
}

// runConcurrentSearch manages concurrent ad collection using worker pool
func runConcurrentSearch(
	engineFunc func(string, string, string, bool) ([]AdResult, error),
	query string, engineName string, userAgent string, noRedirection bool) ([]AdResult, error) {
	results := make(chan []AdResult, ConcurrencyLevel)
	errors := make(chan error, ConcurrencyLevel)
	var wg sync.WaitGroup

	// Launch workers
	for i := 0; i < ConcurrencyLevel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Fix Recover from panics
			defer func() {
				if r := recover(); r != nil {
					safePrintf(nil, "\n\n******\nPanic in runConcurrentSearch for %s\n*******\n\n %s", engineName, r)
				}
			}()
			ads, err := engineFunc(query, userAgent, engineName, noRedirection)
			if err != nil {
				errors <- err
				return
			}
			results <- ads
		}()
	}

	// Wait for completion and close channels
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results and handle errors
	var allAds []AdResult
	for ads := range results {
		allAds = append(allAds, ads...)
	}

	// Check for errors
	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return allAds, nil
}

// RunAdSearch returns the ads found in the search engines for the specified config
func RunAdSearch(config Config) ([]AdResult, []AdResult, error) {
	var notifications []AdResult
	var allAdResults []AdResult

	// Option for independent DirectQuery from keywords provided from command line
	// DirectQuery gets priority over config file
	if DirectQuery != "" && len(DirectQuery) > 0 {
		safePrintf(bold, "\n* DIRECT QUERY SEARCH FOR: '%s'\n\n", DirectQuery)

		// Iterate search engines
		for _, engine := range searchEnginesFunctions {
			// Check if SelectedEngine engine option is enabled and the search engine name inside the list
			if isInSelectedSearchEngineList(engine.EngineName) == false {
				continue
			}
			safePrintf(nil, "> Search Engine lookup using '%s' for keyword '%s'\n\n", engine.EngineName, DirectQuery)

			// Search for ads on every search engine
			adResults, err := searchAdsWithEngine(engine.SearchFunction, DirectQuery, engine.EngineName, UserAgentString, NoRedirection)
			if err != nil {
				safePrintf(nil, "Error searching using %s: %v\n", engine.EngineName, err)
				return nil, nil, err
			}

			// Process ads if found
			if len(adResults) == 0 {
				safePrintf(italic, "  no ads found\n\n")
			} else {
				err := processAdResults(adResults, GlobalDomainExclusionList, &allAdResults, &notifications, config)
				if err != nil {
					return nil, nil, err
				}
			}
		}
	} else {
		// Iterate keyword provided from the config file
		for _, searchQuery := range config.Queries {
			// Merge expected/exclusion individual expected domain with global domain lists
			expectedDomainList := mergeLists(GlobalDomainExclusionList, searchQuery.ExpectedDomains)
			safePrintf(bold, "\n* SEARCHING FOR: '%s'\n\n", searchQuery.SearchTerm)

			for _, engine := range searchEnginesFunctions {
				// Check if SelectedEngine engine option is enabled and the search engine name inside the list
				if isInSelectedSearchEngineList(engine.EngineName) == false {
					continue
				}
				safePrintf(nil, "> Search Engine lookup using '%s' for keyword '%s'\n", engine.EngineName, searchQuery.SearchTerm)
				adResults, err := searchAdsWithEngine(engine.SearchFunction, searchQuery.SearchTerm, engine.EngineName, UserAgentString, NoRedirection)
				if err != nil {
					safePrintf(nil, "Error searching using %s: %v\n", engine.EngineName, err)
					return nil, nil, err
				}
				if len(adResults) == 0 {
					safePrintf(italic, "  no ads found\n\n")
				} else {
					err := processAdResults(adResults, expectedDomainList, &allAdResults, &notifications, config)
					if err != nil {
						return nil, nil, err
					}
				}
			}
		}
	}
	return allAdResults, notifications, nil
}

func isInSelectedSearchEngineList(engineName string) bool {
	if SelectedEngine == "" {
		return true
	}

	selectedEngineList := strings.Split(SelectedEngine, ",")
	for _, searchEngine := range selectedEngineList {
		searchEngine = strings.TrimSpace(searchEngine)
		if engineName == searchEngine {
			return true
		}
	}
	return false
}
