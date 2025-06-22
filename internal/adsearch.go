package internal

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// initializeBrowser sets up the browser and returns the browser instance and the search results page
func initializeBrowser(query, searchEngineURL string) (*rod.Browser, *rod.Page, error) {
	chromePath, _ := launcher.LookPath()

	launcherURL := launcher.New().Bin(chromePath).Set("disable-features", "Translate").MustLaunch()
	browser := rod.New().ControlURL(launcherURL).MustConnect().MustIncognito()

	page := stealth.MustPage(browser).MustEmulate(Laptop)
	page.MustNavigate(searchEngineURL + query)
	wait := page.MustWaitNavigation()
	wait()
	return browser, page, nil
}

// saveHTML saves the HTML content of the page to a file
func saveHTML(page *rod.Page, outputFilePrefix string, query string) {

	// Get the HTML content of the page
	htmlContent, err := page.HTML()
	if err != nil {
		log.Fatalf("failed to get HTML content: %v\n", err)
	}
	if Logger {
		log.Printf("Save search engine result is on\n")
	}
	fileHtmlPath := fmt.Sprintf("%s-%s-%d.html", outputFilePrefix, query, time.Now().UnixNano())

	// Write the HTML content to a file
	err = os.WriteFile(filepath.Join(HtmlPath, fileHtmlPath), []byte(htmlContent), 0644)
	if err != nil {
		log.Fatalf("failed to save HTML to file: %v\n", err)
	} else {
		if Logger {
			log.Printf("Visited page saved to %s", fileHtmlPath)
		}
	}
}

// takeScreenshot saves a screenshot of the page to a file
func takeScreenshot(page *rod.Page, outputFilePrefix string, query string) {
	if Logger {
		log.Printf("Save screenshot is on\n")
	}
	filename := fmt.Sprintf("%s-%s-%d.png", outputFilePrefix, query, time.Now().UnixNano())
	if Logger {
		log.Printf("Taking screenshot... ")
	}
	page.MustScreenshotFullPage(filepath.Join(ScreenshotPath, filename))
	if Logger {
		log.Printf("Screenshot saved at %s", filename)
	}
}

// getAdInfo retrieves the advertiser name and location, works only in google, syndicated and adsenseads
func getAdInfo(browser *rod.Browser, adDetail rod.Element) ([]string, error) {

	adInfoResult := []string{}

	// get ad info URL
	adInfoURL, err := adDetail.Attribute("href")
	if err != nil || adInfoURL == nil {
		return adInfoResult, fmt.Errorf("unable to find ad info URL: %v", err)
	}

	// navigate to ad info page
	adInfoPage := browser.MustPage()
	defer adInfoPage.Close()
	if Logger {
		log.Printf("Navigating to advertisement info page: %s\n", *adInfoURL)
	}

	err = adInfoPage.Navigate(*adInfoURL)
	if err != nil {
		return adInfoResult, fmt.Errorf("\nerror trying to open: %s -> %s\n", *adInfoURL, err)
	}
	adInfoPage.MustWaitLoad()

	// save advertiser name and location in advertisersInfo: [0] is name, [1] is location
	advertisersInfo, err := adInfoPage.ElementsX("//div[div[text()=\"Location\"]]/div[2]/text() | //div[div[text()=\"Location\"]]/preceding-sibling::div[1]/div[2]/text()")
	if err != nil || len(advertisersInfo) < 2 {
		return adInfoResult, fmt.Errorf("error trying to find ad details: %s -> %v\n", *adInfoURL, err)
	}

	// save advertiser name and location in advertisersInfo: [0] is name, [1] is location
	adInfoResult = append(adInfoResult, advertisersInfo[0].MustText(), advertisersInfo[1].MustText())

	return adInfoResult, nil
}

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
				ad.Advertiser, ad.Location = getAdInfoSafe(browser, adDetails[i])
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

// getAdInfoSafe safely extracts advertiser info and location
func getAdInfoSafe(browser *rod.Browser, adDetail *rod.Element) (advertiser, location string) {
	adInfo, err := getAdInfo(browser, *adDetail)
	if err == nil && len(adInfo) >= 2 {
		return adInfo[0], adInfo[1]
	}
	return "", ""
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

// ResolveAdUrl resolves the ad URL to its final destination and updates the AdResult
func ResolveAdUrl(adURL string, currentAd *AdResult) {
	// Skip resolution if already resolved
	if currentAd.FinalRedirectURL != "" {
		return
	}

	// First resolution attempt
	redirectURL, finalDomain := resolveAdURLByDomain(adURL)

	// Handle DoubleClick nested redirects
	if finalDomain == doubleclickdomain {
		redirectURL, finalDomain = resolveAdURLByDomain(redirectURL)
	}

	// Handle d.adx.io nested redirects
	if finalDomain == dadxio {
		redirectURL, finalDomain = resolveAdURLByDomain(redirectURL)
	}

	// Update the AdResult with final values
	currentAd.FinalRedirectURL = redirectURL
	currentAd.FinalDomainURL = finalDomain
}

// resolveAdURLByDomain handles URL resolution based on domain type
func resolveAdURLByDomain(adURL string) (string, string) {
	adDomain, err := extractDomain(adURL)
	if err != nil {
		if Logger {
			log.Printf("Error extracting domain from URL: %s", adURL)
		}
		return adURL, ""
	}

	// Resolve URL based on domain
	resolvers := map[string]func(string) (string, error){
		googledomain:      ResolveGoogleAdURL,
		adsenseadsdomain:  ResolveGoogleAdURL,
		syndicateddomain:  ResolveGoogleAdURL,
		bingdomain:        ResolveBingAdURL,
		ddgdomain:         ResolveDuckDuckGoAdURL,
		doubleclickdomain: ResolveDoubleClickAdURL,
		googleadsservices: ResolveGoogleAdURL,
		dadxio:            ResolveDadxioAdURL,
	}

	if resolver, exists := resolvers[adDomain]; exists {
		if resolvedURL, err := resolver(adURL); err == nil {
			finalDomain, _ := extractDomain(resolvedURL)
			return resolvedURL, finalDomain
		}
	}

	// Default case: return original URL and its domain
	return adURL, adDomain
}

// searchAdsWithEngine performs concurrent ad searches using a specific search engine
func searchAdsWithEngine(
	engineFunc func(string, string, string, bool) ([]AdResult, error),
	query SearchQuery, engineName string, userAgent string, noRedirection bool) ([]AdResult, error) {
	encodedQuery := url.QueryEscape(query.SearchTerm)

	if Logger {
		log.Printf("Searching ads on %s\n", searchEngineURLs[engineName]+query.SearchTerm)
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
					fmt.Printf("\n\n******\nPanic in runConcurrentSearch for %s\n*******\n\n", engineName)
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

// processSearchResults handles post-search processing of ads and returns unique Ads
func processSearchResults(ads []AdResult, userAgent string, noRedirection bool) ([]AdResult, error) {
	// Remove duplicates
	uniqueAds, err := removeDuplicateAds(ads)
	if err != nil {
		return nil, fmt.Errorf("failed to remove duplicates: %v", err)
	}

	// Follow redirects if enabled
	if !noRedirection {
		for i := range uniqueAds {
			redirectChain, _ := findRedirectionChain(uniqueAds[i].OriginalAdURL, userAgent)
			uniqueAds[i].RedirectChain = redirectChain
		}
	}

	return uniqueAds, nil
}

// RunAdSearch returns the ads found in the search engines for the specified config
func RunAdSearch(config Config) ([]AdResult, []AdResult, []AdResult, error) {
	var notifications []AdResult
	var allAdResults []AdResult
	var submitToURLScan []AdResult

	// Get global domain exclusion list
	globalDomainExclusionList := config.GlobalDomainExclusion.GlobalDomainExclusionList

	for _, searchQuery := range config.Queries {
		// Merge expected/exclusion individual expected domain with global domain lists
		expectedDomainList := mergeLists(globalDomainExclusionList, searchQuery.ExpectedDomains)
		log.Printf("\n* SEARCHING FOR: '%s'\n\n", searchQuery.SearchTerm)

		for _, engine := range searchEnginesFunctions {
			log.Printf("> Search Engine lookup using '%s' for keyword '%s'\n", engine.EngineName, searchQuery.SearchTerm)
			adResults, err := searchAdsWithEngine(engine.SearchFunction, searchQuery, engine.EngineName, UserAgentString, NoRedirection)
			if err != nil {
				return nil, nil, nil, err
			}
			if len(adResults) == 0 {
				italic.Printf("  no ads found\n\n")
			} else {
				processAdResults(adResults, expectedDomainList, &allAdResults, &notifications, &submitToURLScan)
			}
		}
	}
	return allAdResults, notifications, submitToURLScan, nil
}

// processAdResults processes the ad results and updates the respective lists
func processAdResults(adResults []AdResult, expectedDomainList []string, allAdResults *[]AdResult, notifications *[]AdResult, submitToURLScan *[]AdResult) error {
	// Iterate over each ad result
	for _, adResult := range adResults {
		// Append the ad result to the allAdResults list
		*allAdResults = append(*allAdResults, adResult)

		if !IsExpectedDomain(adResult.FinalDomainURL, expectedDomainList) {
			if Logger {
				log.Printf("\nURL's domain not on expectedDomain: %s not in '%s'\n", adResult.FinalDomainURL, expectedDomainList)
			}
			printDomainInfo(adResult, false)

			// Append the ad result to submitToURLScan list if enabled
			if EnableURLScan {
				*submitToURLScan = append(*submitToURLScan, adResult)
			}

			// Append the ad result to the notifications list if notifications are enabled
			if EnableNotifications {
				*notifications = append(*notifications, adResult)
			}

			// Print the redirection chain if enabled
			if PrintRedirectChain {
				if err := printRedirectionChain(adResult.RedirectChain); err != nil {
					return fmt.Errorf("failed to print redirection chain: %w", err)
				}
			}
		} else {
			// add is in the expected domain list
			printDomainInfo(adResult, true)
		}

	}
	return nil
}
