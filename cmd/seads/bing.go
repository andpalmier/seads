package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"path/filepath"
	"time"
)

// getBingAds searches for ads on Bing for a given encoded string
func getBingAds(encoded string) ([]string, error) {
	var ads []string

	// Create a new Rod browser instance
	browser := rod.New().MustConnect().MustIncognito()
	page := browser.MustPage()
	defer browser.MustClose()

	wait := page.MustWaitNavigation()
	// Open Bing search page and search for encoded string
	page.MustNavigate(seURLs["Bing"] + encoded)
	wait()

	// Get ad links from the search results
	adList, err := page.Elements(`li.b_adTop a.b_restorableLink`)
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
		filename := fmt.Sprintf("bing-%s-%d.png", encoded, time.Now().UnixNano())
		page.MustWaitStable().MustScreenshotFullPage(filepath.Join(*screenshotPath, filename))
	}

	return ads, nil
}
