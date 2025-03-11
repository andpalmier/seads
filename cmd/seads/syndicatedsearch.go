package main

// searchSyndicatedAds searches for ads on Syndicatedsearch for a given encoded string
func searchSyndicatedAds(query, userAgent string, noRedirectionFlag bool) ([]AdLinkPair, error) {
	browser, page, err := initializeBrowser(query, searchEngineURLs["Syndicated"], "")
	if err != nil {
		return nil, err
	}
	defer browser.MustClose()
	adLinks, err := extractAdLinks(browser, page, userAgent,
		`a.si27`, "href", "syndicated", query, noRedirectionFlag)
	if err != nil {
		return nil, err
	}
	return adLinks, nil
}
