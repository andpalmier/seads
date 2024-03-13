package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"path/filepath"
	"time"
)

// getGoogleAds searches for ads on Google for a given encoded string
func getGoogleAds(encoded string) ([]string, error) {
	var ads []string

	// Create a new Rod browser instance
	browser := rod.New().MustConnect().MustIncognito()
	page := browser.MustPage()
	defer browser.MustClose()

	wait := page.MustWaitNavigation()
	page.MustNavigate(seURLs["Google"] + encoded)
	wait()

	cookieButton := page.MustElements(`button#W0wltc`)
	if len(cookieButton) > 0 {
		cookieButton[0].MustClick()
	}

	// Open Google search page and search for encoded string
	adList, err := page.Elements(`div#tads a.sVXRqc`)
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
		filename := fmt.Sprintf("google-%s-%d.png", encoded, time.Now().UnixNano())
		page.MustWaitStable().MustScreenshotFullPage(filepath.Join(*screenshotPath, filename))
	}

	return ads, nil
}
