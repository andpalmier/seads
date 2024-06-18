package main

// searchDuckDuckGoAds searches for ads on DuckDuckGo for a given encoded string
func searchDuckDuckGoAds(query, userAgent string) ([]AdLinkPair, error) {

	// hardcoded user agent string below seems to trigger more ads
	browser, page, err := initializeBrowser(query, searchEngineURLs["DuckDuckGo"],
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/123.0.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	defer browser.MustClose()

	adLinks, err := extractAdLinks(browser, page, userAgent,
		`li[data-layout="ad"] a[data-testid="result-extras-url-link"]`, "href",
		"duckduckgo", query)
	if err != nil {
		return nil, err
	}

	return adLinks, nil
}
