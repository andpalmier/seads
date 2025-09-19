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

	// OLD allAdResults, notifications, submitToURLScan, err := internal.RunAdSearch(config)
	allAdResults, notifications, err := internal.RunAdSearch(config)
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

	fmt.Println()
}
