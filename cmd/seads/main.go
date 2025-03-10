package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	configFilePath      = flag.String("config", "config.yaml", "path to config file")
	concurrencyLevel    = flag.Int("concurrency", 4, "number of concurrent headless browsers")
	screenshotPath      = flag.String("screenshot", "", "path to store screenshots (if empty, the screenshot feature will be disabled)")
	printCleanLinks     = flag.Bool("cleanlinks", false, "print clear links in output (links will remain defanged in notifications)")
	enableNotifications = flag.Bool("notify", false, "notify if unexpected domains are found (requires notifications fields in config.yaml)")
	printRedirectChain  = flag.Bool("redirect", false, "print redirection chain")
	userAgentString     = flag.String("ua", "", "User-Agent string to be used to click on ads")
	outputFilePath      = flag.String("out", "", "path of file containing links of gathered ads")
	enableURLScan       = flag.Bool("urlscan", false, "submit url to urlscan.io for analysis")
	noRedirection       = flag.Bool("noredirection", false, "do not follow redirection, if URLScan submit link to resolve by URLScan instead")
	htmlPath            = flag.String("html", "", "path to store search engine result html page (if empty, the htmlPath feature will be disabled)")

	searchEngineURLs = map[string]string{
		"Google":     "https://www.google.com/search?q=",
		"Bing":       "https://www.bing.com/search?form=QBLH&q=",
		"Yahoo":      "https://search.yahoo.com/search?q=",
		"DuckDuckGo": "https://duckduckgo.com/?ia=web&q=",
		"Syndicated": "https://syndicatedsearch.goog/afs/ads?adsafe=medium&adtest=off&adpage=1&channel=ch1&client=amg-informationvine&r=m&hl=en&ie=utf-8&adrep=5&oe=utf-8&type=0&format=p5%7Cn5&ad=n5p5&output=uds_ads_only&v=3&bsl=8&pac=0&u_his=5&uio=--&cont=text-ad-block-0%7Ctext-ad-block-1&rurl=https%3A%2F%2Fwww.ask.com%2Fweb%3F%26o%3D0%26an%3Dorganic%26ad%3DOther%2BSEO%26capLimitBypass%3Dfalse%26qo%3DserpSearchTopBox%26q&q=",
	}
	searchEnginesFunctions = []SearchEngineFunction{
		{EngineName: "Google", SearchFunction: searchGoogleAds},
		{EngineName: "Bing", SearchFunction: searchBingAds},
		{EngineName: "Yahoo", SearchFunction: searchYahooAds},
		{EngineName: "DuckDuckGo", SearchFunction: searchDuckDuckGoAds},
		{EngineName: "Syndicated", SearchFunction: searchSyndicatedAds},
	}
)

// SearchEngineFunction holds the search engine name and its corresponding function
type SearchEngineFunction struct {
	EngineName     string
	SearchFunction func(string, string, bool) ([]AdLinkPair, error)
}

// AdResult contains information regarding an ad found
type AdResult struct {
	Engine           string    `json:"engine"`
	Query            string    `json:"query"`
	OriginalAdURL    string    `json:"OriginalAdURL"`
	FinalDomainURL   string    `json:"final-domain-url"`
	FinalRedirectURL string    `json:"final-redirect-url"`
	RedirectChain    []string  `json:"redirect-chain"`
	Time             time.Time `json:"time"`
}

// AdLinkPair represents a pair of the original ad link and the final redirected URL
type AdLinkPair struct {
	OriginalAdURL string
	FinalAdURL    string
}

// performAdSearch return the ads found in the search engines for the specified config
func performAdSearch(config Config) ([]AdResult, []AdResult, []AdResult, error) {
	var notifications []AdResult
	var allAdResults []AdResult
	var submitToURLScan []AdResult

	for _, searchQuery := range config.Queries {
		log.Printf("\nSearching for: '%s'\n", searchQuery.SearchTerm)

		for _, engine := range searchEnginesFunctions {
			log.Printf("Search Engine Lookup using '%s' for keyword '%s'", engine.EngineName, searchQuery.SearchTerm)
			adResults, err := searchAdsWithEngine(engine.SearchFunction, searchQuery, engine.EngineName, *userAgentString, *noRedirection)
			if err != nil {
				return nil, nil, nil, err
			}
			if len(adResults) == 0 {
				italic.Println("no ads found")
			} else {
				for _, adResult := range adResults {
					allAdResults = append(allAdResults, adResult)
					if *enableURLScan {
						if *noRedirection {
							// parse ads
							if !isAdsExpected(adResult.OriginalAdURL, searchQuery.ExpectedDomains) {
								log.Printf("\nURL's domain not on expectedDomain: '%s'\n", searchQuery.ExpectedDomains)
								submitToURLScan = append(submitToURLScan, adResult)
							}
						} else {
							submitToURLScan = append(submitToURLScan, adResult)
						}
					} else {
						if isDomainExpected(adResult.FinalDomainURL, searchQuery.ExpectedDomains) {
							printExpectedDomainInfo(adResult)
						} else {
							printUnexpectedDomainInfo(adResult)
							if *enableNotifications {
								notifications = append(notifications, adResult)
							}

						}
						if *printRedirectChain {
							printRedirectionChain(adResult.RedirectChain)
						}
					}
				}
			}
		}
	}
	return allAdResults, notifications, submitToURLScan, nil
}

func main() {

	flag.Parse()
	if *enableURLScan {
		fmt.Println("URLScan Enable")
	}

	config, err := parseConfig(*configFilePath)
	if err != nil {
		if *configFilePath == "config.yaml" {
			log.Fatalf("no config file found at config.yaml, please be sure to use " +
				"-config to specify the config file path")
		}
		log.Fatalf("error parsing config file: %v\n", err)
	}

	allAdResults, notifications, submitToURLScan, err := performAdSearch(config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if *outputFilePath != "" && len(allAdResults) > 0 {
		fmt.Println()
		err := exportAdResults(*outputFilePath, allAdResults)
		if err == nil {
			fmt.Println("file exported successfully at: " + *outputFilePath)
		} else {
			fmt.Println("error exporting file: " + err.Error())
		}
	}

	if *enableNotifications && len(notifications) > 0 {
		fmt.Println()
		config.sendNotifications(notifications)
	}

	// Submit domain to URLScan
	if *enableURLScan {
		if len(submitToURLScan) > 0 {
			fmt.Println("Total URLs for submission: ", len(submitToURLScan))
			config.submitURLScan(submitToURLScan)
		} else {
			fmt.Println("URLScan enabled with no possible submission")
		}
	}

	fmt.Println()
}
