package internal

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// printRedirectionChain prints the chain of redirection URLs or returns an error if none are found
func printRedirectionChain(redirectionURLs []string) error {
	log.Printf("  redirect chain: ")
	if len(redirectionURLs) > 1 {
		fmt.Println()
		for i, url := range redirectionURLs {
			if !PrintCleanLinks {
				url = defangURL(url)
			}
			log.Printf("    %d) %s\n", i+1, url)
		}
	} else {
		log.Printf("no redirects found!\n")
		return fmt.Errorf("no redirects found in the chain")
	}
	fmt.Println()
	return nil
}

// createNoRedirectHTTPClient creates an HTTP client that prevents automatic redirects
func createNoRedirectHTTPClient(userAgent string) *http.Client {
	client := *http.DefaultClient
	client.CheckRedirect = requests.NoFollow
	client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	// setting custom user agent string, if provided
	if userAgent != "" {
		client.Transport = &userAgentTransport{
			userAgent: userAgent,
			transport: client.Transport,
		}
	}
	return &client
}

// userAgentTransport is custom transport that sets the User-Agent header
type userAgentTransport struct {
	userAgent string
	transport http.RoundTripper
}

// RoundTrip executes a single HTTP transaction
func (uat *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", uat.userAgent)
	return uat.transport.RoundTrip(req)
}

// findRedirectionChain retrieves the redirection chain starting from the given URL
func findRedirectionChain(initialURL string, userAgent string) ([]string, error) {
	client := createNoRedirectHTTPClient(userAgent)
	var redirectionChain []string
	currentURL := initialURL
	redirectionChain = append(redirectionChain, currentURL)

	for {
		redirectURL, err := fetchRedirectURL(client, currentURL)
		if err != nil {
			return redirectionChain, err
		}

		if isValidRedirect(redirectURL, initialURL) {
			redirectionChain = append(redirectionChain, redirectURL)
			currentURL = redirectURL
		} else {
			break
		}
	}
	return redirectionChain, nil
}

// fetchRedirectURL fetches the redirect location for a given URL using an HTTP client
func fetchRedirectURL(client *http.Client, url string) (string, error) {
	var redirectLocation string
	err := requests.URL(url).
		Client(client).
		CheckStatus(http.StatusFound).
		Handle(func(res *http.Response) error {
			redirectLocation = res.Header.Get("Location")
			return nil
		}).
		Fetch(context.Background())
	return redirectLocation, err
}

// isValidRedirect checks if a redirect URL is valid and different from the initial URL
func isValidRedirect(redirectURL, initialURL string) bool {
	return strings.HasPrefix(redirectURL, "http") && !strings.HasPrefix(redirectURL, initialURL)
}

// extractDestURL extracts a destination URL from a given query parameter key
func extractDestURL(rawURL, paramKey string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("Skipping invalid URL: %s, Error: %v", rawURL, err)
	}
	queryParams := parsedURL.Query()
	destURL := queryParams.Get(paramKey)
	test, err := url.Parse(destURL)
	if err != nil || test.Host == "" {
		return "", fmt.Errorf("Skipping invalid destination URL: %s, Error: %v", destURL, err)
	}
	return destURL, nil
}

// resolveAdURLByDomainRecursive recursively resolves ad URLs by following known resolver chains.
func resolveAdURLByDomainRecursive(adURL string, visited map[string]struct{}) (string, string) {
	if visited == nil {
		visited = make(map[string]struct{})
	}
	if _, seen := visited[adURL]; seen {
		// Prevent infinite loops
		return adURL, ""
	}
	visited[adURL] = struct{}{}

	adDomain, err := extractDomain(adURL)
	if err != nil {
		if Logger {
			safePrintf(nil, "Error extracting domain from URL: %s", adURL)
		}
		return adURL, ""
	}

	resolvers := map[string]func(string) (string, error){
		googledomain:      ResolveGoogleAdURL,
		adsenseadsdomain:  ResolveGoogleAdURL,
		syndicateddomain:  ResolveGoogleAdURL,
		bingdomain:        ResolveBingAdURL,
		ddgdomain:         ResolveDuckDuckGoAdURL,
		doubleclickdomain: ResolveDoubleClickAdURL,
		googleadsservices: ResolveGoogleAdURL,
		dadxio:            ResolveDadxioAdURL,
		dartsearch:        ResolveDartSearchAdURL,
		clickcease:        ResolveClickCeaseAdURL,
		agkn:              ResolveAgknAdURL,
	}

	if resolver, exists := resolvers[adDomain]; exists {
		if resolvedURL, err := resolver(adURL); err == nil && resolvedURL != adURL {
			return resolveAdURLByDomainRecursive(resolvedURL, visited)
		}
	}

	return adURL, adDomain
}

// ResolveAdUrl resolves the ad URL to its final destination and updates the AdResult.
func ResolveAdUrl(adURL string, currentAd *AdResult) {
	if currentAd.FinalRedirectURL != "" {
		return
	}
	redirectURL, finalDomain := resolveAdURLByDomainRecursive(adURL, nil)
	currentAd.FinalRedirectURL = redirectURL
	currentAd.FinalDomainURL = finalDomain
}

// ResolveDartSearchAdURL uses the generic extractor
func ResolveDartSearchAdURL(dartSearchAdURL string) (string, error) {
	return extractDestURL(dartSearchAdURL, "ds_dest_url")
}

// ResolveDadxioAdURL uses the generic extractor
func ResolveDadxioAdURL(dadxioAdURL string) (string, error) {
	return extractDestURL(dadxioAdURL, "xu")
}

// ResolveClickCeaseAdURL uses the generic extractor
func ResolveClickCeaseAdURL(ClickCeaseAdURL string) (string, error) {
	return extractDestURL(ClickCeaseAdURL, "url")
}

// ResolveDoubleClickAdURL uses the generic extractor
func ResolveDoubleClickAdURL(doubleClickURL string) (string, error) {
	return extractDestURL(doubleClickURL, "ds_dest_url")
}

// ResolveDoubleClickAdURL uses the generic extractor
func ResolveAgknAdURL(AgknAdURL string) (string, error) {
	return extractDestURL(AgknAdURL, "l0")
}
