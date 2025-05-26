package internal

import (
	// "github.com/go-rod/rod"
	"github.com/go-rod/rod"
)

// searchAolAds searches for ads on Aol for a given encoded string
func searchAolAds(query, userAgent, engine string, noRedirectionFlag bool) ([]AdResult, error) {
	browser, page, err := initializeBrowser(query, searchEngineURLs[engine])
	if err != nil {
		return nil, err
	}
	page.MustWaitLoad()
	defer browser.MustClose()
	handleAolPageInteraction(page)

	if len(ScreenshotPath) > 0 {
		takeScreenshot(page, engine, query)
	}
	if len(HtmlPath) > 0 {
		saveHTML(page, engine, query)
	}

	adLinks, err := extractAds(browser, page, userAgent, aolSelector, "href", query, engine, noRedirectionFlag)
	if err != nil {
		return nil, err
	}

	return adLinks, nil
}

// handleAolPageInteraction handles interactions with Aol search results page (closing cookies button)
func handleAolPageInteraction(page *rod.Page) {
	scrollButtons, err := page.Elements(aolScrollBtn)
	if err == nil {
		if len(scrollButtons) > 0 {
			scrollButtons[0].MustClick()
			wait := page.MustWaitNavigation()
			page.MustElement(aolCookieBtn).MustClick()
			wait()
		}
	}
}
