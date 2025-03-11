package main

import (
	"github.com/go-rod/rod"
)

// searchGoogleAds searches for ads on Google for a given encoded string
func searchGoogleAds(query, userAgent string, noRedirectionFlag bool) ([]AdLinkPair, error) {
	browser, page, err := initializeBrowser(query, searchEngineURLs["Google"], "")
	if err != nil {
		return nil, err
	}
	defer browser.MustClose()
	handleGooglePageInteraction(page)

	adLinks, err := extractAdLinks(browser, page, userAgent,
		`a.sVXRqc`, "data-rw", "google", query, noRedirectionFlag)
	if err != nil {
		return nil, err
	}
	return adLinks, nil
}

// handleGooglePageInteraction handles interactions with Google search results page (closing cookies button)
func handleGooglePageInteraction(page *rod.Page) {
	acceptCookiesButton, err := page.Elements(`button#W0wltc`)
	if err == nil {
		if len(acceptCookiesButton) > 0 {
			acceptCookiesButton[0].MustClick()
		}
	}
}
