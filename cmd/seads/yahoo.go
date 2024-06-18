package main

import (
	"github.com/go-rod/rod"
)

// searchYahooAds searches for ads on Yahoo for a given encoded string
func searchYahooAds(query, userAgent string) ([]AdLinkPair, error) {
	browser, page, err := initializeBrowser(query, searchEngineURLs["Yahoo"], "")
	if err != nil {
		return nil, err
	}
	defer browser.MustClose()

	handleYahooPageInteraction(page)

	adLinks, err := extractAdLinks(browser, page, userAgent,
		`ol.searchCenterTopAds a[data-matarget="ad"]`, "href", "yahoo", query)
	if err != nil {
		return nil, err
	}

	return adLinks, nil
}

// handleYahooPageInteraction handles interactions with Yahoo search results page (closing cookies button)
func handleYahooPageInteraction(page *rod.Page) {
	scrollButtons, err := page.Elements(`button#scroll-down-btn`)
	if err == nil {
		if len(scrollButtons) > 0 {
			scrollButtons[0].MustClick()
			wait := page.MustWaitNavigation()
			page.MustElement(`button[value="reject"`).MustClick()
			wait()
		}
	}
}
