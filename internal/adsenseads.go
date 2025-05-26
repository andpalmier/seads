package internal

// searchAdsenseAds searches for ads on adsensecustomsearchads for a given encoded string
func searchAdsenseAds(query, userAgent, engine string, noRedirectionFlag bool) ([]AdResult, error) {

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

	ads, err := extractAds(browser, page, userAgent, adsenseadsSelector, "href", query, engine, noRedirectionFlag)

	if err != nil {
		return nil, err
	}
	return ads, nil
}
