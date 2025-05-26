package internal

import (
	"github.com/go-rod/rod"
	//"time"
)

// searchYahooAds searches for ads on Yahoo for a given encoded string
func searchYahooAds(query, userAgent, engine string, noRedirectionFlag bool) ([]AdResult, error) {
	browser, page, err := initializeBrowser(query, searchEngineURLs[engine])
	if err != nil {
		return nil, err
	}
	page.MustWaitLoad()
	defer browser.MustClose()
	handleYahooPageInteraction(page)

	if len(ScreenshotPath) > 0 {
		takeScreenshot(page, engine, query)
	}
	if len(HtmlPath) > 0 {
		saveHTML(page, engine, query)
	}

	adLinks, err := extractAds(browser, page, userAgent, yahooSelector, "href", query, engine, noRedirectionFlag)
	if err != nil {
		return nil, err
	}

	return adLinks, nil
}

// handleYahooPageInteraction handles interactions with Yahoo search results page (closing cookies button)
func handleYahooPageInteraction(page *rod.Page) {
	scrollButtons, err := page.Elements(yahooScrollBtn)
	if err == nil {
		if len(scrollButtons) > 0 {
			scrollButtons[0].MustClick()
			wait := page.MustWaitNavigation()
			page.MustElement(yahooCookieBtn).MustClick()
			wait()
		}
	}
}
