package internal

import (
	"fmt"
	"github.com/fatih/color"
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
