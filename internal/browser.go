package internal

/*
import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// BrowserConfig holds all configuration needed for browser operations
type BrowserConfig struct {
	UserAgent      string
	NoRedirection  bool
	ScreenshotPath string
	HtmlPath       string
	Engine         string
	Query          string
	LinkSelector   string
	AttributeName  string
}

// Browser represents a wrapped rod.Browser with its configuration
type Browser struct {
	rod    *rod.Browser
	config BrowserConfig
}

// NewBrowser creates and initializes a new Browser instance
func NewBrowser(config BrowserConfig) (*Browser, error) {
	chromePath, _ := launcher.LookPath()
	launcherURL := launcher.New().
		Bin(chromePath).
		Set("disable-features", "Translate").
		MustLaunch()

	browser := rod.New().
		ControlURL(launcherURL).
		MustConnect().
		MustIncognito()

	return &Browser{
		rod:    browser,
		config: config,
	}, nil
}

// Close properly closes the browser
func (b *Browser) Close() {
	b.rod.MustClose()
}


// NavigateToSearch opens a new page and navigates to the search URL
func (b *Browser) NavigateToSearch(searchURL string) (*rod.Page, error) {
	page := stealth.MustPage(b.rod).MustEmulate(Laptop)
	page.MustNavigate(searchURL)
	wait := page.MustWaitNavigation()
	wait()

	return page, nil
}

// ExtractAds processes a page and extracts all advertisements
func (b *Browser) ExtractAds(page *rod.Page) ([]AdResult, error) {
	adElements, err := page.Elements(b.config.LinkSelector)
	if err != nil {
		return nil, fmt.Errorf("unable to find ad elements: %v", err)
	}

	var adDetails rod.Elements
	if isGoogleLikeEngine(b.config.Engine) {
		adDetails, _ = page.Elements(adinfoSelector)
	}

	return b.processAdElements(adElements, adDetails)
}

// processAdElements handles concurrent processing of ad elements
func (b *Browser) processAdElements(adElements, adDetails rod.Elements) ([]AdResult, error) {
	var (
		adsFound []AdResult
		mu       sync.Mutex
		wg       sync.WaitGroup
	)

	for i, element := range adElements {
		adURL, err := element.Attribute(b.config.AttributeName)
		if err != nil || adURL == nil {
			continue
		}

		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			ad := b.createAdResult(i, url, adDetails)
			mu.Lock()
			adsFound = append(adsFound, ad)
			mu.Unlock()
		}(i, *adURL)
	}

	wg.Wait()
	return adsFound, nil
}

// createAdResult creates a single ad result with all necessary information
func (b *Browser) createAdResult(index int, adURL string, adDetails rod.Elements) AdResult {
	ad := AdResult{
		OriginalAdURL: adURL,
		Query:         b.config.Query,
		Time:          time.Now(),
		Engine:        b.config.Engine,
	}

	if isGoogleLikeEngine(b.config.Engine) && index < len(adDetails) {
		adInfo := processAdInfo(adDetails[index])
		ad.Advertiser = adInfo.Advertiser
		ad.Location = adInfo.Location
	}

	if !b.config.NoRedirection {
		ad.FinalRedirectURL, ad.FinalDomainURL = b.followRedirect(adURL)
	}

	processAdURL(adURL, &ad)
	return ad
}



// processAdInfo retrieves advertiser information from an ad detail element
func processAdInfo(adDetail *rod.Element) struct {
	Advertiser string
	Location   string
} {
	text, err := adDetail.Text()
	if err != nil {
		return struct {
			Advertiser string
			Location   string
		}{"", ""}
	}

	parts := splitAdInfo(text)
	if len(parts) >= 2 {
		return struct {
			Advertiser string
			Location   string
		}{parts[0], parts[1]}
	}

	return struct {
		Advertiser string
		Location   string
	}{"", ""}
}

// splitAdInfo splits the ad info text into its components
func splitAdInfo(text string) []string {
	// Remove leading/trailing spaces and split by dot separator
	parts := strings.Split(text, "·")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// followRedirect follows an ad URL to its final destination
func (b *Browser) followRedirect(adURL string) (finalURL, finalDomain string) {
	page := b.rod.MustPage()
	defer page.Close()

	if b.config.UserAgent != "" {
		_ = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: b.config.UserAgent,
		})
	}

	wait := page.MustWaitNavigation()
	if err := page.Navigate(adURL); err == nil {
		wait()
		finalURL = page.MustInfo().URL
		finalDomain, _ = extractDomain(finalURL)
	}
	return
}

// SavePageContent saves both screenshot and HTML content if configured
func (b *Browser) SavePageContent(page *rod.Page) error {
	if b.config.ScreenshotPath != "" {
		if err := b.saveScreenshot(page); err != nil {
			return err
		}
	}

	if b.config.HtmlPath != "" {
		if err := b.saveHTML(page); err != nil {
			return err
		}
	}
	return nil
}

// saveScreenshot saves a screenshot of the current page
func (b *Browser) saveScreenshot(page *rod.Page) error {
	filename := fmt.Sprintf("%s-%s-%d.png",
		b.config.Engine,
		b.config.Query,
		time.Now().UnixNano(),
	)
	fullPath := filepath.Join(b.config.ScreenshotPath, filename)

	data := page.MustScreenshotFullPage(fullPath)
	return os.WriteFile(fullPath, data, 0644)
}

// saveHTML saves the HTML content of the current page
func (b *Browser) saveHTML(page *rod.Page) error {
	content, err := page.HTML()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("search-page--%s-%s-%d.html",
		b.config.Engine,
		b.config.Query,
		time.Now().UnixNano(),
	)

	return os.WriteFile(filepath.Join(b.config.HtmlPath, filename), []byte(content), 0644)
}



// processAdURL processes the ad URL and updates the AdResult
func processAdURL(adURL string, ad *AdResult) {
	parsedURL, err := url.Parse(adURL)
	if err != nil {
		return
	}

	// Extract redirect URL from Google-style ads
	if redirectURL := parsedURL.Query().Get("adurl"); redirectURL != "" {
		ad.RedirectChain = append(ad.RedirectChain, redirectURL)
	}
}

*/
