package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"path/filepath"
	"time"
)

// getDuckDuckGoAds searches for ads on DuckDuckGo for a given encoded string
func getDuckDuckGoAds(encoded string) ([]string, error) {
	var ads []string

	// Search for chromium path
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()

	// Create a new Rod browser instance
	browser := rod.New().ControlURL(u).MustConnect().MustIncognito()
	page := browser.MustPage()
	defer browser.MustClose()

	wait := page.MustWaitNavigation()
	// Open DuckDuckGo search page and search for encoded string
	page.MustNavigate(seURLs["DuckDuckGo"] + encoded)
	wait()

	// Get ad links from the search results
	adList, err := page.Elements("li[data-layout=\"ad\"] a[data-testid=\"result-extras-url-link\"]")
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
		filename := fmt.Sprintf("duckduckgo-%s-%d.png", encoded, time.Now().UnixNano())
		page.MustWaitStable().MustScreenshotFullPage(filepath.Join(*screenshotPath, filename))
	}

	return ads, nil
}
