package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/url"
	"strings"
	"sync"
)

var (
	green  = color.New(color.FgGreen)
	italic = color.New(color.Italic)
	red    = color.New(color.FgRed)
)

// removeDuplicates removes ads with same domain from the given list
func removeDuplicates(ads []string) ([]string, error) {
	var results []string
	seen := make(map[string]struct{})

	for _, adURL := range ads {
		adURL = normalizeURL(adURL)
		parsedURL, err := url.Parse(adURL)
		if err != nil {
			return nil, err
		}
		domain := parsedURL.Host
		if _, ok := seen[domain]; !ok {
			results = append(results, adURL)
			seen[domain] = struct{}{}
		}
	}
	return results, nil
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

// EncodeString encodes an input string
func EncodeString(input string) string {
	return url.QueryEscape(input)
}

// DefangURL prevents a URL from being clickable
func DefangURL(url string) string {
	return strings.ReplaceAll(url, ".", "[.]")
}

// ExtractDomainFromURL extracts domain from a URL
func ExtractDomainFromURL(inputURL string) (string, error) {
	if !strings.Contains(inputURL, "https") {
		inputURL = "https://" + inputURL
	}
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	host := parsedURL.Host
	return strings.TrimPrefix(host, "www."), nil
}

// GetResultAdsFromURLs gets ResultAd list from a list of ads
func GetResultAdsFromURLs(ads []string) ([]ResultAd, error) {
	var results []ResultAd
	uniqueAds, err := removeDuplicates(ads)
	if err != nil {
		return results, err
	}
	for _, ad := range uniqueAds {
		domain, err := ExtractDomainFromURL(ad)
		if err != nil {
			return nil, errors.New("cannot get domain from following URL: " + ad)
		}
		results = append(results, ResultAd{domain, ad})
	}
	return results, nil
}

// isExpectedDomain checks if the domainAd is in the expectedDomains list
func isExpectedDomain(domainAd string, expectedDomains []string) bool {
	for _, domain := range expectedDomains {
		if domainAd == domain {
			return true
		}
	}
	return false
}

// printExpectedDomain prints the expected domain
func printExpectedDomain(resultAd ResultAd) {
	green.Printf("[+] expected domain: ")
	if *cleanLinks {
		fmt.Printf("%s => %s\n", resultAd.Domain, resultAd.Link)
	} else {
		fmt.Printf("%s => %s\n", DefangURL(resultAd.Domain), DefangURL(resultAd.Link))
	}
}

// printUnexpectedDomain prints the unexpected domain
func printUnexpectedDomain(resultAd ResultAd) {
	red.Printf("[!] unexpected domain: ")
	if *cleanLinks {
		fmt.Printf("%s => %s\n", resultAd.Domain, resultAd.Link)
	} else {
		fmt.Printf("%s => %s\n", DefangURL(resultAd.Domain), DefangURL(resultAd.Link))
	}
}

// searchAds searches ads using a specific engine function
func searchAds(engineFunc func(string) ([]string, error), query Query, engineName string) []string {
	encoded := EncodeString(query.SearchTerm)

	adsFoundChan := make(chan []string, *consumers)
	var wg sync.WaitGroup

	searchFunc := func(i int) {
		defer wg.Done()
		ads, err := engineFunc(encoded)
		if err != nil {
			log.Printf("Error searching %s ad: %v", engineName, err)
			return
		}
		adsFoundChan <- ads
	}

	for i := 0; i < *consumers; i++ {
		wg.Add(1)
		go searchFunc(i)
	}

	wg.Wait()
	close(adsFoundChan)

	adsFound := make([]string, 0)
	for ads := range adsFoundChan {
		adsFound = append(adsFound, ads...)
	}

	return adsFound
}
