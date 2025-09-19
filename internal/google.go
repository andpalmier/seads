package internal

import (
	//"fmt"
	"github.com/go-rod/rod"
	//"net/url"
)

// searchGoogleAds searches for ads on Google for a given encoded string
func searchGoogleAds(query, userAgent, engine string, noRedirectionFlag bool) ([]AdResult, error) {

	browser, page, err := initializeBrowser(query, searchEngineURLs[engine])

	if err != nil {
		return nil, err
	}
	defer browser.MustClose()
	page.MustWaitLoad()
	handleGooglePageInteraction(page)

	if len(ScreenshotPath) > 0 {
		takeScreenshot(page, engine, query)
	}
	if len(HtmlPath) > 0 {
		saveHTML(page, engine, query)
	}

	adLinks, err := extractAds(browser, page, userAgent, googleSelector, "data-rw", query, engine, noRedirectionFlag)
	if err != nil {
		return nil, err
	}
	return adLinks, nil
}

// handleGooglePageInteraction handles interactions with Google search results page (closing cookies button)
func handleGooglePageInteraction(page *rod.Page) {
	acceptCookiesButton, err := page.Elements(googleCookieBtn)
	if err == nil {
		if len(acceptCookiesButton) > 0 {
			acceptCookiesButton[0].MustClick()
		}
	}
}

// ResolveGoogleAdURL uses the generic extractor
func ResolveGoogleAdURL(googleURL string) (string, error) {
	return extractDestURL(googleURL, "adurl")
}
