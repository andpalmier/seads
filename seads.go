package main

import (
	"flag"
	"fmt"
	"github.com/andpalmier/seads/internal"
	"log"
)

func main() {

	flag.Parse()
	fmt.Println(internal.AsciiArt)
	if internal.EnableURLScan && internal.Logger {
		log.Printf("URLScan enabled\n")
	}

	config, err := internal.ParseConfig(internal.ConfigFilePath)
	if err != nil {
		if internal.ConfigFilePath == "config.yaml" {
			fmt.Printf("no config file found...")
			internal.ShowHelp()
		}
		log.Fatalf("error parsing config file: %v\n", err)
	}

	internal.PrintFlags()

	allAdResults, notifications, submitToURLScan, err := internal.RunAdSearch(config)
	if err != nil {
		log.Fatalf("error running ad search: %v\n", err)
	}

	if internal.OutputFilePath != "" && len(allAdResults) > 0 {
		fmt.Println()
		err := internal.ExportAdResults(internal.OutputFilePath, allAdResults)
		if err == nil {
			fmt.Println("file exported successfully at: " + internal.OutputFilePath)
		} else {
			fmt.Println("error exporting file: " + err.Error())
		}
	}

	if internal.EnableNotifications && len(notifications) > 0 {
		fmt.Println()
		config.SendNotifications(notifications)
	}

	// Submit domain to URLScan
	if internal.EnableURLScan {
		if len(submitToURLScan) > 0 {
			fmt.Println("Total URLs to be submitted to URLScan: ", len(submitToURLScan))
			config.SubmitURLScan(submitToURLScan)
		} else {
			fmt.Println("URLScan enabled, but no submissions")
		}
	}

	fmt.Println()
}

/*

package main

import (
	"flag"
	"fmt"
	"github.com/andpalmier/seads/internal"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	fmt.Println(internal.AsciiArt)

	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	internal.PrintFlags()
	logConfiguration()

	// Pass config pointer directly since loadConfig already returns *internal.Config
	searchResult, err := internal.RunAdSearch(*config)
	if err != nil {
		return fmt.Errorf("search error: %w", err)
	}

	if err := processResults(config, searchResult); err != nil {
		return fmt.Errorf("results processing error: %w", err)
	}

	return nil
}

func loadConfig() (*internal.Config, error) {
	config, err := internal.ParseConfig(internal.ConfigFilePath)
	if err != nil {
		if internal.ConfigFilePath == "config.yaml" {
			fmt.Printf("No config file found...\n")
			internal.ShowHelp()
		}
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}
	return &config, nil
}

func logConfiguration() {
	if internal.EnableURLScan && internal.Logger {
		log.Printf("URLScan enabled")
	}
}

func processResults(config *internal.Config, results internal.SearchResult) error {
	// Export all ads
	if err := exportResults(results.AllAds); err != nil {
		return err
	}

	// Handle notifications for unexpected domains
	if err := handleNotifications(config, results.Notifications); err != nil {
		return err
	}

	// Handle URLScan for suspicious ads
	if err := handleURLScan(config, results.URLScanAds); err != nil {
		return err
	}

	fmt.Println()
	return nil
}

func exportResults(results []internal.AdResult) error {
	if internal.OutputFilePath == "" || len(results) == 0 {
		return nil
	}

	fmt.Println()
	if err := internal.ExportAdResults(internal.OutputFilePath, results); err != nil {
		return fmt.Errorf("error exporting results: %w", err)
	}

	fmt.Printf("File exported successfully at: %s\n", internal.OutputFilePath)
	return nil
}

func handleNotifications(config *internal.Config, notifAds []internal.AdResult) error {
	if !internal.EnableNotifications || len(notifAds) == 0 {
		return nil
	}

	notifications := make([]string, 0, len(notifAds))
	for _, ad := range notifAds {
		notifications = append(notifications, fmt.Sprintf("Unexpected domain found: %s", ad.FinalDomainURL))
	}

	fmt.Println()
	config.SendNotifications(notifications)
	return nil
}

func handleURLScan(config *internal.Config, scanAds []internal.AdResult) error {
	if !internal.EnableURLScan {
		return nil
	}

	urlsToScan := make([]string, 0, len(scanAds))
	for _, ad := range scanAds {
		urlsToScan = append(urlsToScan, ad.FinalRedirectURL)
	}

	if len(urlsToScan) > 0 {
		fmt.Printf("Total URLs to be submitted to URLScan: %d\n", len(urlsToScan))
		config.SubmitURLScan(urlsToScan)
	} else {
		fmt.Println("URLScan enabled, but no submissions")
	}

	return nil
}

*/
