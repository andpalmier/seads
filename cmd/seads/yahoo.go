package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"path/filepath"
	"time"
)

// getYahooAds searches for ads on Yahoo for a given encoded string
func getYahooAds(encoded string) ([]string, error) {
	var ads []string

	// Create a new Rod browser instance
	browser := rod.New().MustConnect().MustIncognito()
	page := browser.MustPage()
	defer browser.MustClose()

	wait := page.MustWaitNavigation()
	// Open Yahoo search page and scroll to click "reject cookie" button if present
	page.MustNavigate(seURLs["Yahoo"] + encoded)
	wait()

	scrollButtons := page.MustElements(`button#scroll-down-btn`)
	if len(scrollButtons) > 0 {
		scrollButtons[0].MustClick()
		page.MustElement(`button[value="reject"`).MustClick()
	}

	// Get ad links from the search results
	adList, err := page.Elements(`ol.searchCenterTopAds a[data-matarget="ad"]`)
	if err != nil {
		return nil, err
	}

	// Open ads link in a new page to get URL
	for _, ad := range adList {
		href, err := ad.Attribute("href")
		if err != nil {
			return nil, err
		}
		adPage := browser.MustPage(*href)
		defer adPage.Close()
		wait := adPage.MustWaitNavigation()
		wait()
		ads = append(ads, adPage.MustInfo().URL)
	}

	// Capture a screenshot if ads are found and screenshot path is provided
	if len(ads) > 0 && len(*screenshotPath) > 0 {
		filename := fmt.Sprintf("yahoo-%s-%d.png", encoded, time.Now().UnixNano())
		page.MustWaitStable().MustScreenshotFullPage(filepath.Join(*screenshotPath, filename))
	}

	return ads, nil
}
