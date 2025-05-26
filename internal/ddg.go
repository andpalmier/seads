package internal

import (
	"fmt"
	"net/url"
)

// searchDuckDuckGoAds searches for ads on DuckDuckGo for a given encoded string
func searchDuckDuckGoAds(query, userAgent, engine string, noRedirectionFlag bool) ([]AdResult, error) {

	// hardcoded user agent string below seems to trigger more ads
	browser, page, err := initializeBrowser(query, searchEngineURLs[engine])

	if err != nil {
		return nil, err
	}
	page.MustWaitLoad()
	defer browser.MustClose()

	if len(ScreenshotPath) > 0 {
		takeScreenshot(page, engine, query)
	}
	if len(HtmlPath) > 0 {
		saveHTML(page, engine, query)
	}

	adLinks, err := extractAds(browser, page, userAgent, ddgSelector, "href", query, engine, noRedirectionFlag)
	if err != nil {
		return nil, err
	}

	return adLinks, nil
}

// ResolveDuckDuckGoAdURL parses a DuckDuckGo URL and extracts the final redirect URL
func ResolveDuckDuckGoAdURL(ddgURL string) (string, error) {

	// Parse the unescaped URL
	parsedURL, err := url.Parse(ddgURL)
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("Skipping invalid DDG URL: %s, Error: %v\n", ddgURL, err)
	}

	// Extract query parameters from the parsed URL
	queryParams := parsedURL.Query()
	unescapedDuckDuckGoURL, err := ResolveBingAdURL(queryParams.Get("u3"))
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("Skipping invalid Bing URL: %s, Error: %v\n", ddgURL, err)
	}

	return unescapedDuckDuckGoURL, nil
}
