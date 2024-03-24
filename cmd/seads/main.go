package main

import (
	"flag"
	"fmt"
	"log"
)

/*
	TODO
	- handle errors :(
	- export results as text
	- change User-Agent stri
*/

var (
	configPath     = flag.String("config", "config.yaml", "path to config file")
	consumers      = flag.Int("concurrency", 4, "number of concurrent headless browsers")
	screenshotPath = flag.String("screenshot", "", "path to store screenshots (if empty, the screenshot feature will be disabled)")
	cleanLinks     = flag.Bool("cleanlinks", false, "print clear links in output (links will remain defanged in notifications)")
	notify         = flag.Bool("notify", false, "notify if unexpected domains are found (requires notifications fields in config.yaml)")
	seURLs         = map[string]string{
		"Google":     "https://www.google.com/search?q=",
		"Bing":       "https://www.bing.com/search?q=",
		"Yahoo":      "https://search.yahoo.com/search?q=",
		"DuckDuckGo": "https://duckduckgo.com/?ia=web&?q=",
	}
	sf = []SearchFunctions{
		{Name: "Google", Function: getGoogleAds},
		{Name: "Bing", Function: getBingAds},
		{Name: "Yahoo", Function: getYahooAds},
		{Name: "DuckDuckGo", Function: getDuckDuckGoAds},
	}
)

// SearchFunc holds the search engine name and its corresponding function
type SearchFunctions struct {
	Name     string
	Function func(string) ([]string, error)
}

// ResultAd represents an ad result
type ResultAd struct {
	Domain string
	Link   string
}

// search return the ads found in the search engines for the specified config
func search(config Config) []string {
	var toNotify []string

	for _, query := range config.Queries {
		fmt.Printf("Searching for: '%s'\n", query.SearchTerm)

		for _, engine := range sf {
			ads := searchAds(engine.Function, query, engine.Name)
			resultAds, _ := GetResultAdsFromURLs(ads)
			fmt.Println()
			fmt.Printf("* Searching ads for '%s' on %s: ",
				query.SearchTerm, engine.Name)
			if len(resultAds) == 0 {
				italic.Println("no ads found")
			} else {
				fmt.Println()
				for _, resultAd := range resultAds {
					if isExpectedDomain(resultAd.Domain, query.ExpectedDomains) {
						printExpectedDomain(resultAd)
					} else {
						printUnexpectedDomain(resultAd)
						if *notify {
							toNotify = append(toNotify, formatNotification(engine.Name,
								query.SearchTerm, resultAd))
						}
					}
				}
			}
		}
		fmt.Println()
	}
	return toNotify
}

func main() {

	flag.Parse()

	config, err := parseConfig(*configPath)
	if err != nil {
		if *configPath == "config.yaml" {
			log.Fatalf("no config file found at config.yaml, please be sure to use " +
				"-config to specify the config file path")
		}
		log.Fatalf("error parsing config file: %v\n", err)
	}

	toNotify := search(config)

	if *notify && len(toNotify) > 0 {
		fmt.Println()
		config.notify(toNotify)

	}
	fmt.Println()
}
