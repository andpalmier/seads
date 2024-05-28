package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"path/filepath"
	"time"
)

// getYahooAds searches for ads on Yahoo for a given encoded string
func getYahooAds(encoded string, userAgent string) ([]string, error) {
	var ads []string

	// Search for chromium path
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()

	// Create a new Rod browser instance
	browser := rod.New().ControlURL(u).MustConnect().MustIncognito()
	defer browser.MustClose()

	page := browser.MustPage()
	wait := page.MustWaitNavigation()
	// Open Yahoo search page and scroll to click "reject cookie" button if present
	pagerr := page.Navigate(seURLs["Yahoo"] + encoded)
	if pagerr != nil {
		return nil, pagerr
	}
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
		adPage := browser.MustPage()
		defer adPage.Close()
		if userAgent != "" {
			if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent}); err != nil {
				continue
			}
		}
		wait := adPage.MustWaitNavigation()
		pagerr = adPage.Navigate(*href)
		if pagerr != nil {
			fmt.Printf("\nerror trying to open: %s -> %s\n", *href, pagerr)
			continue
		}
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
