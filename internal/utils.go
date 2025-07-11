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

// processSearchResults handles post-search processing of ads and returns unique Ads
func processSearchResults(ads []AdResult, userAgent string, noRedirection bool) ([]AdResult, error) {
	// Remove duplicates
	uniqueAds, err := removeDuplicateAds(ads)
	if err != nil {
		return nil, fmt.Errorf("failed to remove duplicates: %v", err)
	}

	// Follow redirects if enabled
	if !noRedirection {
		for i := range uniqueAds {
			redirectChain, _ := findRedirectionChain(uniqueAds[i].OriginalAdURL, userAgent)
			uniqueAds[i].RedirectChain = redirectChain
		}
	}

	return uniqueAds, nil
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

// processAdResults processes the ad results and updates the respective lists
func processAdResults(adResults []AdResult, expectedDomainList []string, allAdResults *[]AdResult, notifications *[]AdResult, config Config) error {
	// Iterate over each ad result
	for i := range adResults {
		adResult := adResults[i]

		if !IsExpectedDomain(adResult.FinalDomainURL, expectedDomainList) {
			if Logger {
				safePrintf(nil, "\nURL's domain not on expectedDomain: %s not in '%s'\n", adResult.FinalDomainURL, expectedDomainList)
			}
			printDomainInfo(adResult, false)
			adResult.ExpectedDomains = false

			// Submit original advertisement URL to URLScan if enabled
			if EnableURLScan {
				urlScanResult, err := SubmitURLScan(config, adResult.OriginalAdURL)
				if err != nil {
					log.Printf("Error submitting to URLScan: %v\n", err)
				} else {
					adResult.URLScan = urlScanResult
				}
			}

			// Append the ad result to the notifications list if notifications are enabled
			if EnableNotifications {
				*notifications = append(*notifications, adResult)
			}

			// Print the redirection chain if enabled
			if PrintRedirectChain {
				if err := printRedirectionChain(adResult.RedirectChain); err != nil {
					return fmt.Errorf("failed to print redirection chain: %w", err)
				}
			}
		} else {
			// add is in the expected domain list
			printDomainInfo(adResult, true)
			adResult.ExpectedDomains = true
		}
		// Append the ad result to the allAdResults list
		*allAdResults = append(*allAdResults, adResult)
	}
	return nil
}

/* OLD
// processAdResults processes the ad results and updates the respective lists
func processAdResults(adResults []AdResult, expectedDomainList []string, allAdResults *[]AdResult, notifications *[]AdResult, submitToURLScan *[]AdResult) error {
	// Iterate over each ad result
	for _, adResult := range adResults {
		// Append the ad result to the allAdResults list
		*allAdResults = append(*allAdResults, adResult)

		if !IsExpectedDomain(adResult.FinalDomainURL, expectedDomainList) {
			if Logger {
				safePrintf(nil, "\nURL's domain not on expectedDomain: %s not in '%s'\n", adResult.FinalDomainURL, expectedDomainList)
			}
			printDomainInfo(adResult, false)

			// Append the ad result to submitToURLScan list if enabled
			if EnableURLScan {
				*submitToURLScan = append(*submitToURLScan, adResult)
			}

			// Append the ad result to the notifications list if notifications are enabled
			if EnableNotifications {
				*notifications = append(*notifications, adResult)
			}

			// Print the redirection chain if enabled
			if PrintRedirectChain {
				if err := printRedirectionChain(adResult.RedirectChain); err != nil {
					return fmt.Errorf("failed to print redirection chain: %w", err)
				}
			}
		} else {
			// add is in the expected domain list
			printDomainInfo(adResult, true)
		}

	}
	return nil
}
*/
