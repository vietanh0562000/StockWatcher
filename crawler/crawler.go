package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stockwatcher/config"
	"stockwatcher/database"
	"stockwatcher/models"
	"time"
)

type IResourceClientService interface {
	GetTrade(symbol string) *models.Trade
	GetQuote(symbol string) *models.Quote
}

type Crawler struct {
	config  *config.Config
	storage *database.Storage
	client  IResourceClientService
}

func NewCrawler(config *config.Config, storage *database.Storage, client IResourceClientService) *Crawler {
	return &Crawler{
		config:  config,
		storage: storage,
		client:  client,
	}
}

func CreateQuoteListener(quoteChan chan models.Quote, storage *database.Storage) {
	for {
		quote := <-quoteChan
		err := storage.QuoteRepo.CreateQuote(&quote)
		if err != nil {
			fmt.Println("Create quote failed ", err.Error())
		} else {
			fmt.Println("Create quite successful")
		}
	}
}

func (c *Crawler) StartCrawlCronJob() {
	// Start a new goroutine for the cron job
	symbols := getTier1HighActivitySymbols()
	for {
		for _, symbol := range symbols {
			c.CrawTrade(symbol)
			c.CrawlQuote(symbol)
		}
		// Wait for a minute before the next crawl
		time.Sleep(time.Minute)
	}
}

func (c *Crawler) CrawTrade(symbol string) {
	trade := c.client.GetTrade(symbol)
	if trade != nil {
		database.Instance.TradeRepo.CreateTrade(trade)
	}
}

func (c *Crawler) CrawlQuote(symbol string) {
	quote := c.client.GetQuote(symbol)
	if quote != nil {
		database.Instance.QuoteRepo.CreateQuote(quote)
	}
}

func fetchSymbols(exchange string, apiToken string) []models.Symbol {
	// create url
	url := fmt.Sprintf("https://finnhub.io/api/v1/stock/symbol?exchange=%s&token=%s", exchange, apiToken)

	// create http.Client with longer timeout
	client := &http.Client{Timeout: time.Second * 30}

	// call client.Get
	resp, err := client.Get(url)

	// check status
	if err != nil {
		fmt.Printf("Get symbol error: %s", err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get symbol %s", resp.StatusCode)
		return nil
	}

	// parse and save symbol into database
	var symbols []models.Symbol
	err = json.NewDecoder(resp.Body).Decode(&symbols)

	if err != nil {
		fmt.Printf("Can't parse body: %s", err.Error())
		return nil
	}

	return symbols
}

func getTier1HighActivitySymbols() []string {
	return []string{
		// These generate the most WebSocket messages
		"AAPL",  // Apple - Massive volume, constant trading
		"NVDA",  // NVIDIA - AI hype = constant price movement
		"AMZN",  // Amazon - Large cap with good volatility
		"META",  // Meta - Social media volatility
		"GOOGL", // Google - Search/AI developments
		"AMD",   // AMD - Semiconductor volatility
	}
}
