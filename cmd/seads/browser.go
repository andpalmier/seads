package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

// initializeBrowser sets up the browser and returns the browser instance and the search results page
func initializeBrowser(query, searchEngineURL, userAgent string) (*rod.Browser, *rod.Page, error) {
	chromePath, _ := launcher.LookPath()

	launcherURL := launcher.New().Bin(chromePath).MustLaunch()
	browser := rod.New().ControlURL(launcherURL).MustConnect().MustIncognito()

	page := stealth.MustPage(browser)
	page.MustNavigate(searchEngineURL + query)
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
func extractAdLinks(browser *rod.Browser, page *rod.Page, userAgent, linkSelector, attrName, screenshotPrefix, query string, noRedirectionFlag bool) ([]AdLinkPair, error) {
	var adLinks []AdLinkPair

	// Get the HTML content of the page
	htmlContent, err := page.HTML()
	if err != nil {
		log.Fatalf("failed to get HTML content: %v", err)
	}

	// Specify the file path to save the HTML
	if len(*htmlPath) > 0 {
		log.Printf("Save search engine result is on")
		fileHtmlPath := fmt.Sprintf("search-page--%s-%s-%d.html", screenshotPrefix, query, time.Now().UnixNano())
		// Write the HTML content to a file
		err = os.WriteFile(filepath.Join(*htmlPath, fileHtmlPath), []byte(htmlContent), 0644)
		if err != nil {
			log.Fatalf("failed to save HTML to file: %v", err)
		} else {
			log.Printf("Visited page saved to %s\n", fileHtmlPath)
		}
	}

	if len(*screenshotPath) > 0 {
		filename := fmt.Sprintf("%s-%s-%d.png", screenshotPrefix, query, time.Now().UnixNano())
		log.Printf("Taking screenshot... ")
		page.MustScreenshotFullPage(filepath.Join(*screenshotPath, filename))
		log.Printf("Screenshot saved at %s", filename)
	}

	adElements, err := page.Elements(linkSelector)
	if err != nil {
		return nil, fmt.Errorf("unable to find ad elements: %v", err)
	}
	for _, adElement := range adElements {
		adURL, err := adElement.Attribute(attrName)
		if err != nil || adURL == nil {
			continue
		}
		log.Printf("AdURL for '%s' ad: %v", query, *adURL)

		// If no redirection is required, skip navigating to the advertisement
		if noRedirectionFlag {
			adLinks = append(adLinks, AdLinkPair{
				OriginalAdURL: *adURL,
				FinalAdURL:    "",
			})
			continue
		} else {
			adPage := browser.MustPage()
			defer adPage.Close()
			// use custom user agent string if provided
			if userAgent != "" {
				if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent}); err != nil {
					continue
				}
			}
			// Navitage to advertisement URL using current host
			wait := adPage.MustWaitNavigation()
			log.Printf("Navigating to advertisement: %v", *adURL)
			err = adPage.Navigate(*adURL)
			if err != nil {
				log.Printf("\nerror trying to open: %s -> %s\n", *adURL, err)
				continue
			}
			wait()

			finalURL := adPage.MustInfo().URL
			adLinks = append(adLinks, AdLinkPair{
				OriginalAdURL: *adURL,
				FinalAdURL:    finalURL,
			})
		}
	}

	return adLinks, nil
}

// searchAdsWithEngine searches ads using a specific search engine function
func searchAdsWithEngine(engineFunc func(string, string, bool) ([]AdLinkPair, error), query SearchQuery, engineName string, userAgent string, noRedirectionFlag bool) ([]AdResult, error) {
	encoded := encodeSearchTerm(query.SearchTerm)

	adsChannel := make(chan []AdLinkPair, *concurrencyLevel)
	var wg sync.WaitGroup

	searchAdFunction := func(i int) {
		defer wg.Done()
		// Browse page and scrap advertisement links
		adLinks, err := engineFunc(encoded, userAgent, noRedirectionFlag)
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
	// Use gathered ads and follow redirect here
	resultAdList, err := generateAdResults(allAdLinks, query.SearchTerm, engineName, time.Now(), noRedirectionFlag)
	if err != nil {
		return nil, err
	}

	return resultAdList, nil
}
