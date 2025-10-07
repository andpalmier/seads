package internal

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"sync"
)

var printMutex sync.Mutex

// safePrintln println in a thread-safe manner
func safePrintln(a ...interface{}) {
	printMutex.Lock()
	fmt.Println(a...)
	printMutex.Unlock()
}

// safePrintf printf in a thread-safe manner
func safePrintf(c *color.Color, format string, a ...interface{}) {
	printMutex.Lock()
	defer printMutex.Unlock()
	if c != nil {
		c.Printf(format, a...)
	} else {
		fmt.Printf(format, a...)
	}
}

// printDomainInfo logs domain information based on whether it is expected or unexpected
func printDomainInfo(resultAd AdResult, expected bool) {

	domainToPrint := resultAd.FinalDomainURL
	urlToPrint := resultAd.FinalRedirectURL
	originalURL := resultAd.OriginalAdURL

	if PrintCleanLinks {
		urlToPrint = defangURL(urlToPrint)
		domainToPrint = defangURL(domainToPrint)
		originalURL = defangURL(originalURL)
	}

	if expected {
		safePrintf(green, "  [+] expected domain: ")
	} else {
		safePrintf(red, "  [!] unexpected domain: ")
	}

	safePrintf(nil, "  %s => %s\n", domainToPrint, urlToPrint)
	origDom, _ := extractDomain(originalURL)
	if domainToPrint != origDom {
		safePrintf(nil, "  original URL: %s\n", originalURL)
	}

	if resultAd.Advertiser != "" {
		safePrintf(nil, "  advertiser name: %s\n  advertiser location: %s\n", resultAd.Advertiser, resultAd.Location)
	}

	safePrintf(nil, "\n")
}

// printStringFlag prints a label and its value, using italics for "none" values
func printStringFlag(label, value string) {
	if value == "" {
		safePrintf(nil, "  %s", label)
		safePrintf(italic, "%s\n", "none")
	} else {
		safePrintf(nil, "  %s%s\n", label, value)
	}
}

// PrintEngines prints list of searchEngines
func PrintEngines() {
	var names []string
	for _, se := range searchEnginesFunctions {
		names = append(names, se.EngineName)
	}
	safePrintf(nil, "  Engines (%d): %s\n", len(names), strings.Join(names, ", "))
}

// PrintQueryKeywords prints list of query keywords from config file
func PrintQueryKeywords(config Config) {
	if DirectQuery == "" {
		if config.Queries == nil {
			fmt.Println("  No queries defined in config")
		} else {
			var queries []string
			for _, searchQuery := range config.Queries {
				queries = append(queries, searchQuery.SearchTerm)
			}
			if len(queries) == 0 {

			} else {
				safePrintf(nil, "  Queries: %s\n", strings.Join(queries, ", "))
			}
		}
	} else {
		safePrintf(nil, "  Direct search query: %s\n", DirectQuery)
	}
}

// PrintTotalGlobalExclusions prints total global exclusion list
func PrintTotalGlobalExclusions(config Config) {
	globalDomainExclusionList := config.GlobalDomainExclusion.GlobalDomainExclusionList
	fmt.Printf("  Global Domain Exclusion List: %d\n", len(globalDomainExclusionList))
}

// PrintConfigOverview config overview after banner
func PrintConfigOverview(config Config) {
	safePrintln("Search Engines and Keywords:")
	PrintEngines()
	PrintQueryKeywords(config)
	PrintTotalGlobalExclusions(config)
}

// PrintFlags prints the current values of the command-line arguments
func PrintFlags() {
	safePrintln("Configuration Flags:")
	printStringFlag("ConfigFilePath: ", ConfigFilePath)
	safePrintf(nil, "  ConcurrencyLevel: %d\n", ConcurrencyLevel)
	printStringFlag("ScreenshotPath: ", ScreenshotPath)
	safePrintf(nil, "  PrintCleanLinks: %t\n", PrintCleanLinks)
	safePrintf(nil, "  EnableNotifications: %t\n", EnableNotifications)
	safePrintf(nil, "  PrintRedirectChain: %t\n", PrintRedirectChain)
	printStringFlag("UserAgentString: ", UserAgentString)
	safePrintf(nil, "  EnableURLScan: %t\n", EnableURLScan)
	printStringFlag("OutputFilePath: ", OutputFilePath)
	safePrintf(nil, "  NoRedirection: %t\n", NoRedirection)
	printStringFlag("HTMLFilePath: ", HtmlPath)
	printStringFlag("DirectQuery: ", DirectQuery)
	safePrintf(nil, "  Logger: %t\n\n", Logger)
}
