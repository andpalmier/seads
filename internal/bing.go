package internal

import (
	"fmt"
	"net/url"
)

// searchBingAds searches for ads on Bing for a given query
func searchBingAds(query, userAgent, engine string, noRedirectionFlag bool) ([]AdResult, error) {

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

	adLinks, err := extractAds(browser, page, userAgent, bingSelector, "href", query, engine, noRedirectionFlag)
	if err != nil {
		return nil, err
	}

	return adLinks, nil
}

// ResolveBingAdURL parses a Bing URL and extracts the redirect URL
func ResolveBingAdURL(bingURL string) (string, error) {

	// Parse the unescaped URL
	parsedURL, err := url.Parse(bingURL)
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("Skipping invalid Bing URL: %s, Error: %v\n", bingURL, err)
	}

	// Extract query parameters from the parsed URL
	queryParams := parsedURL.Query()
	uBingURL, err := decodeBase64(queryParams.Get("u"))
	if err != nil {
		return "", fmt.Errorf("Skipping invalid Bing URL: %s, Error: %v\n", parsedURL.RawQuery, err)
	}

	unescapedBingURL, err := url.QueryUnescape(uBingURL)
	if err != nil {
		return "", fmt.Errorf("Unescaping not possible in Bing URL: %s, Error: %v\n", unescapedBingURL, err)
	}

	test, err := url.Parse(unescapedBingURL)
	if err != nil || test.Host == "" {
		return "", fmt.Errorf("Skipping invalid Bing URL: %s, Error: %v\n", unescapedBingURL, err)
	}

	return unescapedBingURL, nil
}
