package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	configPath     = flag.String("config", "config.yaml", "path to config file")
	consumers      = flag.Int("concurrency", 4, "number of concurrent headless browsers")
	screenshotPath = flag.String("screenshot", "", "path to store screenshots (if empty, the screenshot feature will be disabled)")
	cleanLinks     = flag.Bool("cleanlinks", false, "print clear links in output (links will remain defanged in notifications)")
	notify         = flag.Bool("notify", false, "notify if unexpected domains are found (requires notifications fields in config.yaml)")
	userAgent      = flag.String("ua", "", "User-Agent string to be used to click on ads")
	output         = flag.String("out", "", "path of file containing links of gathered ads")
	seURLs         = map[string]string{
		"Google":     "https://www.google.com/search?q=",
		"Bing":       "https://www.bing.com/search?q=",
		"Yahoo":      "https://search.yahoo.com/search?q=",
		"DuckDuckGo": "https://duckduckgo.com/?ia=web&q=",
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
	Function func(string, string) ([]string, error)
}

// ResultAd represents an ad result
type ResultAd struct {
	Engine string    `json:"engine"`
	Query  string    `json:"query"`
	Domain string    `json:"domain"`
	Link   string    `json:"link"`
	Time   time.Time `json:"time"`
}

// search return the ads found in the search engines for the specified config
func search(config Config) ([]ResultAd, []ResultAd, error) {
	var adsToNotify []ResultAd
	var allAds []ResultAd
	for _, query := range config.Queries {
		fmt.Printf("Searching for: '%s'\n", query.SearchTerm)

		for _, engine := range sf {
			resultAds, err := searchAds(engine.Function, query, engine.Name, *userAgent)
			if err != nil {
				return nil, nil, err
			}
			fmt.Println()
			fmt.Printf("* Searching ads for '%s' on %s: ",
				query.SearchTerm, engine.Name)
			if len(resultAds) == 0 {
				italic.Println("no ads found")
			} else {
				fmt.Println()
				for _, resultAd := range resultAds {
					allAds = append(allAds, resultAd)
					if isExpectedDomain(resultAd.Domain, query.ExpectedDomains) {
						printExpectedDomain(resultAd)
					} else {
						printUnexpectedDomain(resultAd)
						if *notify {
							adsToNotify = append(adsToNotify, resultAd)
						}
					}
				}
			}
		}
		fmt.Println()
	}
	return allAds, adsToNotify, nil
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

	allAds, adsToNotify, err := search(config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if *output != "" && len(allAds) > 0 {
		fmt.Println()
		err := export(*output, allAds)
		if err == nil {
			fmt.Println("file exported successfully at: " + *output)
		} else {
			fmt.Println("error exporting file: " + err.Error())
		}
	}

	if *notify && len(adsToNotify) > 0 {
		fmt.Println()
		config.notify(adsToNotify)
	}

	fmt.Println()
}
