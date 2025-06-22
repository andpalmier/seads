package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

// removeDuplicateAds filters out ads with duplicate domains from the given list
func removeDuplicateAds(ads []AdResult) ([]AdResult, error) {
	var uniqueAds []AdResult
	seenDomains := make(map[string]struct{})

	for _, ad := range ads {
		// Normalize the URL based on the noRedirectionFlag
		normalizedAdURL := normalizeURL(ad.FinalRedirectURL)

		// Parse the normalized URL to extract the domain
		parsedURL, err := url.Parse(normalizedAdURL)
		if err != nil {
			return nil, err
		}

		adDomain := parsedURL.Host
		// Check if the domain has already been seen
		if _, seen := seenDomains[adDomain]; !seen {
			uniqueAds = append(uniqueAds, ad)
			seenDomains[adDomain] = struct{}{}
		}
	}
	return uniqueAds, nil
}

// normalizeURL ensures a URL starts with "https://"
func normalizeURL(adURL string) string {
	if strings.HasPrefix(adURL, "https://") {
		return adURL
	}
	if strings.HasPrefix(adURL, "http://") {
		return strings.ReplaceAll(adURL, "http://", "https://")
	}
	return "https://" + adURL
}

// defangURL modifies a URL to make it non-clickable by replacing "." with "[.]"
func defangURL(url string) string {
	replace := strings.ReplaceAll(url, ".", "[.]")
	replace = strings.Replace(replace, "http", "hxxp", 1)
	return replace
}

// extractDomain extracts the domain name from a given URL or returns an error if invalid
func extractDomain(inputURL string) (string, error) {
	// Ensure the URL has a scheme
	normalizeURL(inputURL)

	// Parse the URL
	parsedURL, err := url.ParseRequestURI(inputURL)
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: %s", inputURL)
	}

	// Extract and return the domain
	adHost := parsedURL.Host
	return strings.TrimPrefix(adHost, "www."), nil
}

// IsExpectedDomain checks if a domain belongs to one of the expected domains
func IsExpectedDomain(domain string, expectedDomains []string) bool {
	for _, expectedDomain := range expectedDomains {
		if domain == expectedDomain || checkHostnameEndsWithDomain(domain, expectedDomain) {
			if Logger {
				log.Printf("URL excluded by expected domain: %s\n", domain)
			}
			return true
		}
	}
	return false
}

// isGoogleLikeEngine checks if the engine is Google, Syndicated, or AdsenseAds
func isGoogleLikeEngine(engine string) bool {
	return engine == "syndicated" || engine == "adsenseads" || engine == "google"
}

// ExportAdResults writes AdResult objects to a beautified JSON file
func ExportAdResults(filepath string, allAds []AdResult) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(allAds)
}

// checkHostnameEndsWithDomain verifies if a hostname ends with a specific domain
func checkHostnameEndsWithDomain(hostname, domain string) bool {
	return strings.HasSuffix(hostname, domain)
}

// decodeBase64 decodes a Base64-encoded string or returns an error if decoding fails
func decodeBase64(encoded string) (string, error) {
	decodedBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}
	return string(decodedBytes), nil
}

// mergeLists combines two lists of strings and returns unique values
func mergeLists(firstList, secondList []string) []string {
	uniqueMap := make(map[string]struct{})
	for _, item := range append(firstList, secondList...) {
		uniqueMap[item] = struct{}{}
	}
	result := make([]string, 0, len(uniqueMap))
	for key := range uniqueMap {
		result = append(result, key)
	}
	return result
}

// printDomainInfo logs domain information based on whether it is expected or unexpected
func printDomainInfo(resultAd AdResult, expected bool) {
	if expected {
		green.Printf("  [+] expected domain: ")
	} else {
		red.Printf("  [!] unexpected domain: ")
	}

	domainToPrint := resultAd.FinalDomainURL
	urlToPrint := resultAd.FinalRedirectURL
	originalURL := resultAd.OriginalAdURL

	if PrintCleanLinks {
		urlToPrint = defangURL(urlToPrint)
		domainToPrint = defangURL(domainToPrint)
		originalURL = defangURL(originalURL)
	}

	log.Printf("%s => %s", domainToPrint, urlToPrint)
	origDom, _ := extractDomain(originalURL)
	if domainToPrint != origDom {
		log.Printf("  original URL: %s\n", originalURL)
	}

	if resultAd.Advertiser != "" {
		log.Printf("  advertiser name: %s\n  advertiser location: %s\n", resultAd.Advertiser, resultAd.Location)
	}
	fmt.Println()
}

// PrintFlags prints the current values of the command-line arguments
func PrintFlags() {
	log.Println("Configuration Flags:")
	log.Printf("  ConfigFilePath: %s\n", ConfigFilePath)
	log.Printf("  ConcurrencyLevel: %d\n", ConcurrencyLevel)
	log.Printf("  ScreenshotPath: %s\n", ScreenshotPath)
	log.Printf("  PrintCleanLinks: %t\n", PrintCleanLinks)
	log.Printf("  EnableNotifications: %t\n", EnableNotifications)
	log.Printf("  PrintRedirectChain: %t\n", PrintRedirectChain)
	log.Printf("  UserAgentString: %s\n", UserAgentString)
	log.Printf("  EnableURLScan: %t\n", EnableURLScan)
	log.Printf("  OutputFilePath: %s\n", OutputFilePath)
	log.Printf("  NoRedirection: %t\n", NoRedirection)
	log.Printf("  HtmlPath: %s\n", HtmlPath)
	log.Printf("  Logger: %t\n\n", Logger)
}
