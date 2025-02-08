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
	searchEngineURLs    = map[string]string{
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
	SearchFunction func(string, string) ([]AdLinkPair, error)
}

// AdResult contains information regarding an ad found
type AdResult struct {
	Engine           string    `json:"engine"`
	Query            string    `json:"query"`
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
func performAdSearch(config Config) ([]AdResult, []AdResult, error) {
	var notifications []AdResult
	var allAdResults []AdResult
	for _, searchQuery := range config.Queries {
		fmt.Printf("\nSearching for: '%s'\n", searchQuery.SearchTerm)

		for _, engine := range searchEnginesFunctions {
			adResults, err := searchAdsWithEngine(engine.SearchFunction, searchQuery, engine.EngineName, *userAgentString)
			if err != nil {
				return nil, nil, err
			}
			fmt.Printf("\n* %s ads for '%s': ",
				engine.EngineName, searchQuery.SearchTerm)
			if len(adResults) == 0 {
				italic.Println("no ads found")
			} else {
				fmt.Printf("\n")
				for _, adResult := range adResults {
					allAdResults = append(allAdResults, adResult)
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
	return allAdResults, notifications, nil
}

func main() {

	flag.Parse()

	config, err := parseConfig(*configFilePath)
	if err != nil {
		if *configFilePath == "config.yaml" {
			log.Fatalf("no config file found at config.yaml, please be sure to use " +
				"-config to specify the config file path")
		}
		log.Fatalf("error parsing config file: %v\n", err)
	}

	allAdResults, notifications, err := performAdSearch(config)
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

	fmt.Println()
}
