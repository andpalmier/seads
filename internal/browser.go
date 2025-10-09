package internal

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// initializeBrowser sets up the browser and returns the browser instance and the search results page
func initializeBrowser(query, searchEngineURL string) (*rod.Browser, *rod.Page, error) {
	chromePath, _ := launcher.LookPath()

	launcherURL := launcher.New().Bin(chromePath).Set("disable-features", "Translate").MustLaunch()
	browser := rod.New().ControlURL(launcherURL).MustConnect().MustIncognito()

	// randomly select a user agent from the list if concurrency is enabled
	if ConcurrencyLevel > 1 {
		Laptop.UserAgent = userAgents[rand.Intn(len(userAgents))]
	}

	page := stealth.MustPage(browser).MustEmulate(Laptop)
	page.MustNavigate(searchEngineURL + query)
	wait := page.MustWaitNavigation()
	wait()
	return browser, page, nil
}

// saveHTML saves the HTML content of the page to a file
func saveHTML(page *rod.Page, outputFilePrefix string, query string) {

	// Get the HTML content of the page
	htmlContent, err := page.HTML()
	if err != nil {
		log.Fatalf("failed to get HTML content: %v\n", err)
	}
	if Logger {
		safePrintf(nil, "Save search engine result is on\n")
	}
	fileHtmlPath := fmt.Sprintf("%s-%s-%s.html",
		outputFilePrefix,
		query,
		time.Now().Format("20060102-150405"))

	// Write the HTML content to a file
	err = os.WriteFile(filepath.Join(HtmlPath, fileHtmlPath), []byte(htmlContent), 0644)
	if err != nil {
		log.Fatalf("failed to save HTML to file: %v\n", err)
	} else {
		if Logger {
			safePrintf(nil, "Visited page saved to %s", fileHtmlPath)
		}
	}
}

// takeScreenshot saves a screenshot of the page to a file
func takeScreenshot(page *rod.Page, outputFilePrefix string, query string) {
	if Logger {
		safePrintf(nil, "Save screenshot is on\n")
	}
	filename := fmt.Sprintf("%s-%s-%s.png",
		outputFilePrefix,
		query,
		time.Now().Format("20060102-150405"))
	if Logger {
		safePrintf(nil, "Taking screenshot... ")
	}
	page.MustScreenshotFullPage(filepath.Join(ScreenshotPath, filename))
	if Logger {
		safePrintf(nil, "Screenshot saved at %s", filename)
	}
}
