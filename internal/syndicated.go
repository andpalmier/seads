package internal

// searchSyndicatedAds searches for ads on Syndicatedsearch for a given encoded string
func searchSyndicatedAds(query, userAgent, engine string, noRedirectionFlag bool) ([]AdResult, error) {

	browser, page, err := initializeBrowser(query, searchEngineURLs[engine])

	if err != nil {
		return nil, err
	}
	defer browser.MustClose()

	page.MustWaitLoad()

	if len(ScreenshotPath) > 0 {
		takeScreenshot(page, engine, query)
	}
	if len(HtmlPath) > 0 {
		saveHTML(page, engine, query)
	}

	ads, err := extractAds(browser, page, userAgent, syndicatedSelector, "href", query, engine, noRedirectionFlag)

	if err != nil {
		return nil, err
	}
	return ads, nil
}
