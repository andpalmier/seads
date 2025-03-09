package main

// searchDuckDuckGoAds searches for ads on DuckDuckGo for a given encoded string
func searchDuckDuckGoAds(query, userAgent string, noRedirectionFlag bool) ([]AdLinkPair, error) {

	// hardcoded user agent string below seems to trigger more ads
	browser, page, err := initializeBrowser(query, searchEngineURLs["DuckDuckGo"],
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	defer browser.MustClose()

	adLinks, err := extractAdLinks(browser, page, userAgent,
		`li[data-layout="ad"] a[data-testid="result-extras-url-link"]`, "href",
		"duckduckgo", query, noRedirectionFlag)
	if err != nil {
		return nil, err
	}

	return adLinks, nil
}
