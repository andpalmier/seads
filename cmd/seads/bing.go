package main

// searchBingAds searches for ads on Bing for a given query
func searchBingAds(query, userAgent string, noRedirectionFlag bool) ([]AdLinkPair, error) {

	// hardcoded user agent string below seems to trigger more ads
	browser, page, err := initializeBrowser(query, searchEngineURLs["Bing"],
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/123.0.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	defer browser.MustClose()

	adLinks, err := extractAdLinks(browser, page, userAgent,
		`li.b_adTop [role="link"]`, "href", "bing", query, noRedirectionFlag)
	if err != nil {
		return nil, err
	}

	return adLinks, nil
}
