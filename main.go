package main

import (
	"fmt"
	"log"
	"os"
	"stockwatcher/config"
	"stockwatcher/crawler"
	"stockwatcher/database"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [api|crawl]")
		os.Exit(1)
	}

	config := config.LoadConfig()

	db, err := database.Connect(config)
	if err != nil {
		log.Fatal("Failed to connect database: ", err.Error())
		return
	}

	database.GetStorageInstance(db)

	switch os.Args[1] {
	case "api":
		startAPI(config)
	case "crawl":
		startCrawler(config, database.Instance)
	default:
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println("Usage: go run main.go [api|crawl]")
		os.Exit(1)
	}

}

func startAPI(config *config.Config) {
	// get default routes

	// set up cors

	// sett routes

	// start listen
}

func startCrawler(config *config.Config, storage *database.Storage) {
	resourceClient := crawler.NewAlpacaClient(config)
	crawler := crawler.NewCrawler(config, storage, resourceClient)
	crawler.CrawlOne()
}
