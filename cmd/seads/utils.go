package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen)
	italic = color.New(color.Italic)
	red    = color.New(color.FgRed)
)

// removeDuplicateAds removes ads with same domain from the given list
func removeDuplicateAds(adLinks []AdLinkPair, noRedirectionFlag bool) ([]AdLinkPair, error) {
	var uniqueAdLinks []AdLinkPair
	seenDomains := make(map[string]struct{})
	seenURLs := make(map[string]struct{})

	// Normalize the URL to avoid duplicates
	var normalizedAdURL string
	for _, adLink := range adLinks {
		// With noredirection flag, pick the original url
		if noRedirectionFlag {
			normalizedAdURL = normalizeURL(adLink.OriginalAdURL)
		} else {
			normalizedAdURL = normalizeURL(adLink.FinalAdURL)
		}

		parsedURL, err := url.Parse(normalizedAdURL)
		if err != nil {
			return nil, err
		}

		adDomain := parsedURL.Host
		if _, seen := seenDomains[adDomain]; !seen {
			uniqueAdLinks = append(uniqueAdLinks, adLink)
			seenDomains[adDomain] = struct{}{}
		}

		// If no redirection flag is set, we only want to see the original ad URL
		if noRedirectionFlag {
			if _, seen := seenURLs[adLink.OriginalAdURL]; !seen {
				uniqueAdLinks = append(uniqueAdLinks, adLink)
				seenURLs[adLink.OriginalAdURL] = struct{}{}
			}
		}
	}
	return uniqueAdLinks, nil
}

// normalizeURL normalizes an ad URL by adding "https://" if missing
func normalizeURL(adURL string) string {
	if strings.HasPrefix(adURL, "https://") {
		return adURL
	}
	if strings.HasPrefix(adURL, "http://") {
		return strings.ReplaceAll(adURL, "http://", "https://")
	}
	return "https://" + adURL
}

// encodeSearchTerm encodes an input string
func encodeSearchTerm(inputString string) string {
	return url.QueryEscape(inputString)
}

// defangAdURL prevents a URL from being clickable
func defangAdURL(url string) string {
	return strings.ReplaceAll(url, ".", "[.]")
}

// extractDomain extracts domain name from a URL
func extractDomain(inputURL string) (string, error) {
	if !strings.Contains(inputURL, "https") {
		inputURL = "https://" + inputURL
	}
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	adHost := parsedURL.Host
	return strings.TrimPrefix(adHost, "www."), nil
}

// generateAdResults gets AdResult list from a list of ad
func generateAdResults(adLinks []AdLinkPair, searchKeyword string, searchEngineName string, time time.Time, noRedirectionFlag bool) ([]AdResult, error) {
	var adResults []AdResult
	uniqueAdLinks, err := removeDuplicateAds(adLinks, noRedirectionFlag)
	if err != nil {
		return adResults, err
	}
	for _, adLink := range uniqueAdLinks {
		domain, err := extractDomain(adLink.FinalAdURL)
		if err != nil {
			return nil, errors.New("cannot get domain from following URL: " + adLink.FinalAdURL)
		}
		if noRedirectionFlag {
			adResults = append(adResults, AdResult{searchEngineName, searchKeyword, adLink.OriginalAdURL, domain,
				adLink.FinalAdURL, nil, time})

		} else {
			redirectChain, _ := findRedirectionChain(adLink.OriginalAdURL, *userAgentString)
			adResults = append(adResults, AdResult{searchEngineName, searchKeyword, adLink.OriginalAdURL, domain,
				adLink.FinalAdURL, redirectChain, time})
		}

	}
	return adResults, nil
}

// isDomainExpected checks if the domainAd is in the expectedDomains list
func isDomainExpected(domainAd string, expectedDomains []string) bool {
	for _, domain := range expectedDomains {
		if domainAd == domain {
			return true
		}
	}
	return false
}

// printExpectedDomainInfo prints the expected domain
func printExpectedDomainInfo(resultAd AdResult) {
	green.Printf("  [+] expected domain: ")
	if *printCleanLinks {
		fmt.Printf("%s => %s\n", resultAd.FinalDomainURL, resultAd.FinalRedirectURL)
	} else {
		fmt.Printf("%s => %s\n", defangAdURL(resultAd.FinalDomainURL), defangAdURL(resultAd.FinalRedirectURL))
	}
}

// printUnexpectedDomainInfo prints the unexpected domain
func printUnexpectedDomainInfo(resultAd AdResult) {
	red.Printf("  [!] unexpected domain: ")
	if *printCleanLinks {
		fmt.Printf("%s => %s\n", resultAd.FinalDomainURL, resultAd.FinalRedirectURL)
	} else {
		fmt.Printf("%s => %s\n", defangAdURL(resultAd.FinalDomainURL), defangAdURL(resultAd.FinalRedirectURL))
	}
}

// exportAdResults exports results in a beautified JSON file
func exportAdResults(filepath string, allAds []AdResult) error {

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(allAds)
	if err != nil {
		return err
	}
	return nil
}

// Function to check if a string begins with "http" or "https"
func beginsWithHTTP(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}

func isAdsExpected(ads string, expectedDomains []string) bool {
	for _, expectedDomain := range expectedDomains {
		adsURL, _ := url.QueryUnescape(ads) // Unescape any encoded characters

		// Parse the URL
		parsedURL, err := url.Parse(adsURL)
		if err != nil {
			fmt.Printf("Skipping invalid URL: %s, Error: %v\n", adsURL, err)
			continue
		}

		// Extract query parameters
		queryParams := parsedURL.Query()

		// Check for HTTP values in query parameters
		for key, values := range queryParams {
			for _, value := range values {
				if beginsWithHTTP(value) {
					// Parse the URL
					currentParsedURL, err := url.Parse(value)
					if err != nil {
						fmt.Printf("Skipping invalid URL: %s, Error: %v\n", currentParsedURL, err)
						continue
					}

					// Ads host
					currentHost := currentParsedURL.Host
					currentHost = strings.TrimPrefix(currentHost, "www.")

					// If Ads host matches exceptedDomain, return function with true
					if currentHost == expectedDomain {
						return true
					}
				}

				// Exception on DDG ad_domain
				if key == "ad_domain" && strings.HasPrefix(ads, "https://duckduckgo.com") {
					if value == expectedDomain {
						return true
					}
				}
			}
		}
	}
	return false
}
