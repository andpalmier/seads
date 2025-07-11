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

/*
// ResolveGoogleAdURL parses a Google URL and extracts the final redirect URL
func ResolveGoogleAdURL(googleURL string) (string, error) {

	// Parse the unescaped URL
	parsedURL, err := url.Parse(googleURL)
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("Skipping invalid Google URL: %s, Error: %v\n", googleURL, err)
	}

	// Extract query parameters from the parsed URL
	queryParams := parsedURL.Query()

	unescapedGoogleAdURL := queryParams.Get("adurl")

	test, err := url.Parse(unescapedGoogleAdURL)
	if err != nil || test.Host == "" {
		return "", fmt.Errorf("Skipping invalid Google URL: %s, Error: %v\n", unescapedGoogleAdURL, err)
	}

	return unescapedGoogleAdURL, nil
}

*/
