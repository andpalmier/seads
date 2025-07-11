package internal

import (
	"fmt"
	"github.com/go-rod/rod"
)

// getAdInfo retrieves the advertiser name and location, works only in google, syndicated and adsenseads
func getAdInfo(browser *rod.Browser, adDetail *rod.Element) ([]string, error) {

	adInfoResult := []string{}

	// get ad info URL
	adInfoURL, err := adDetail.Attribute("href")
	if err != nil || adInfoURL == nil {
		return adInfoResult, fmt.Errorf("unable to find ad info URL: %v", err)
	}

	// navigate to ad info page
	adInfoPage := browser.MustPage()
	defer adInfoPage.Close()
	if Logger {
		safePrintf(nil, "Navigating to advertisement info page: %s\n", *adInfoURL)
	}

	err = adInfoPage.Navigate(*adInfoURL)
	if err != nil {
		return adInfoResult, fmt.Errorf("\nerror trying to open: %s -> %s\n", *adInfoURL, err)
	}
	adInfoPage.MustWaitLoad()

	// save advertiser name and location in advertisersInfo: [0] is name, [1] is location
	advertisersInfo, err := adInfoPage.ElementsX(adInfoText)
	if err != nil || len(advertisersInfo) < 2 {
		return adInfoResult, fmt.Errorf("error trying to find ad details: %s -> %v\n", *adInfoURL, err)
	}

	// clean up advertiser name -> remove "Paid for by " prefix if it exists
	name := advertisersInfo[0].MustText()
	const prefix = "Paid for by "
	if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
		name = name[len(prefix):]
	}

	// save advertiser name and location in advertisersInfo: [0] is name, [1] is location
	adInfoResult = append(adInfoResult, name, advertisersInfo[1].MustText())

	return adInfoResult, nil
}
