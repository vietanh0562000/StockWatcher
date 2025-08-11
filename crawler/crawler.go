package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stockwatcher/config"
	"stockwatcher/database"
	"stockwatcher/models"
	"time"
)

type IResourceClientService interface {
	GetTrade(symbol string) *models.Trade
	Connect() error
	Auth() error
	Subscribe(dataType string, symbols []string) error
	Unsubscribe(dataType string, symbols []string) error
	Listen() error
	Close()
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

func (c *Crawler) CrawlRealtime() {
	var client IResourceClientService
	client = NewAlpacaClient(c.config)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	if err := client.Auth(); err != nil {
		log.Fatal(err)
		return
	}

	symbols := getTier1HighActivitySymbols()

	client.Subscribe("Quote", symbols)
	client.Subscribe("Trade", symbols)

	// Handle graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go client.Listen()

	// Wait for interrupt signal
	<-interrupt

	// Unsubscribe from all symbols
	client.Unsubscribe("Quote", symbols)
	client.Unsubscribe("Trade", symbols)

	// Close the connection
	client.Close()
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
			c.CrawlOne(symbol)
		}
		// Wait for a minute before the next crawl
		time.Sleep(time.Minute)
	}
}

func (c *Crawler) CrawlOne(symbol string) {
	trade := c.client.GetTrade(symbol)
	if trade != nil {
		database.Instance.TradeRepo.CreateTrade(trade)
	}
}

func updateSymbols(c *Crawler) {
	symbols := fetchSymbols("US", c.config.ResourceAPIKey)
	for _, symbol := range symbols {
		err := c.storage.SymbolRepo.CreateSymbol(&symbol)

		if err != nil {
			fmt.Printf("Create symbol error: %s", err.Error())
		}
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

func fetchQuote(symbol string, apiToken string) *models.Quote {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, apiToken)

	client := http.Client{Timeout: time.Second * 30}

	resp, err := client.Get(url)

	if err != nil {
		fmt.Printf("Fail to fetch quote for %s: %s", symbol, err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Fail to fetch quote for %s: %s", symbol, resp.Status)
		return nil
	}

	var quote models.Quote

	err = json.NewDecoder(resp.Body).Decode(&quote)

	if err != nil {
		fmt.Printf("Decode failed for %s: %s", symbol, err.Error())
		return nil
	}

	return &quote
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
