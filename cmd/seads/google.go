package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"path/filepath"
	"time"
)

// getGoogleAds searches for ads on Google for a given encoded string
func getGoogleAds(encoded string, userAgent string) ([]string, error) {
	var ads []string

	// Search for chromium path
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()

	// Create a new Rod browser instance
	browser := rod.New().ControlURL(u).MustConnect().MustIncognito()
	page := browser.MustPage()
	defer browser.MustClose()
	if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"}); err != nil {
		return nil, err
	}

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
		adPage := browser.MustPage()
		if userAgent != "" {
			if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent}); err != nil {
				return nil, err
			}
		}
		wait := adPage.MustWaitNavigation()
		adPage.MustNavigate(*href)
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
