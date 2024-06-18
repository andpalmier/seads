package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"path/filepath"
	"sync"
	"time"
)

// initializeBrowser sets up the browser and returns the browser instance and the search results page
func initializeBrowser(query, searchEngineURL, userAgent string) (*rod.Browser, *rod.Page, error) {
	chromePath, _ := launcher.LookPath()

	launcherURL := launcher.New().Bin(chromePath).MustLaunch()
	browser := rod.New().ControlURL(launcherURL).MustConnect().MustIncognito()

	page := browser.MustPage(searchEngineURL + query)
	if userAgent != "" {
		err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent})
		if err != nil {
			browser.MustClose()
			return nil, nil, fmt.Errorf("unable to set user agent: %v", err)
		}
	}
	wait := page.MustWaitNavigation()
	wait()
	return browser, page, nil
}

// extractAdLinks extracts ad links from the search result page based on the provided selectors
func extractAdLinks(browser *rod.Browser, page *rod.Page, userAgent, linkSelector, attrName, screenshotPrefix, query string) ([]AdLinkPair, error) {
	var adLinks []AdLinkPair

	adElements, err := page.Elements(linkSelector)
	if err != nil {
		return nil, fmt.Errorf("unable to find ad elements: %v", err)
	}
	for _, adElement := range adElements {
		adURL, err := adElement.Attribute(attrName)
		if err != nil || adURL == nil {
			continue
		}
		adPage := browser.MustPage()
		defer adPage.Close()

		// use custom user agent string if provided
		if userAgent != "" {
			if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent}); err != nil {
				continue
			}
		}

		wait := adPage.MustWaitNavigation()
		err = adPage.Navigate(*adURL)
		if err != nil {
			fmt.Printf("\nerror trying to open: %s -> %s\n", *adURL, err)
			continue
		}
		wait()

		finalURL := adPage.MustInfo().URL
		adLinks = append(adLinks, AdLinkPair{
			OriginalAdURL: *adURL,
			FinalAdURL:    finalURL,
		})
	}

	if len(adLinks) > 0 && len(*screenshotPath) > 0 {
		filename := fmt.Sprintf("%s-%s-%d.png", screenshotPrefix, query, time.Now().UnixNano())
		page.MustWaitStable().MustScreenshotFullPage(filepath.Join(*screenshotPath, filename))
	}

	return adLinks, nil
}

// searchAdsWithEngine searches ads using a specific search engine function
func searchAdsWithEngine(engineFunc func(string, string) ([]AdLinkPair, error), query SearchQuery, engineName string,
	userAgent string) ([]AdResult, error) {
	encoded := encodeSearchTerm(query.SearchTerm)

	adsChannel := make(chan []AdLinkPair, *concurrencyLevel)
	var wg sync.WaitGroup

	searchAdFunction := func(i int) {
		defer wg.Done()
		adLinks, err := engineFunc(encoded, userAgent)
		if err != nil {
			log.Printf("Error searching %s ad: %v", engineName, err)
			return
		}
		adsChannel <- adLinks
	}

	for i := 0; i < *concurrencyLevel; i++ {
		wg.Add(1)
		go searchAdFunction(i)
	}

	wg.Wait()
	close(adsChannel)

	var allAdLinks []AdLinkPair
	for adLinks := range adsChannel {
		allAdLinks = append(allAdLinks, adLinks...)
	}

	resultAdList, err := generateAdResults(allAdLinks, query.SearchTerm, engineName, time.Now())
	if err != nil {
		return nil, err
	}

	return resultAdList, nil
}
