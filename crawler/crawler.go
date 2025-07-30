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

type Crawler struct {
	config  *config.Config
	storage *database.Storage
}

func NewCrawler(config *config.Config, storage *database.Storage) *Crawler {
	return &Crawler{
		config:  config,
		storage: storage,
	}
}

func (c *Crawler) CrawlOne() {
	symbols, err := c.storage.SymbolRepo.GetAllSymbols()
	if err != nil {
		fmt.Printf("Get all symbol error: %s", err.Error())
		return
	}

	// Process symbols with rate limiting
	for i, symbol := range symbols {
		// Add delay between requests to avoid rate limiting
		if i > 0 {
			time.Sleep(1 * time.Second) // 1 second delay between requests
		}

		quote := fetchQuote(symbol.Symbol, c.config.FinnhubAPI)
		if quote != nil {
			// Set the symbol ID for the quote
			quote.SymbolID = symbol.ID
			quote.Timestamp = time.Now()

			err := c.storage.QuoteRepo.CreateQuote(quote)
			if err != nil {
				fmt.Printf("Create quote error for %s: %s", symbol.Symbol, err.Error())
			} else {
				fmt.Printf("Get quote successful for %s\n", symbol.Symbol)
			}
		} else {
			fmt.Printf("Failed to fetch quote for %s\n", symbol.Symbol)
		}

		// Add extra delay every 10 requests to be extra safe
		if (i+1)%10 == 0 {
			time.Sleep(5 * time.Second)
			fmt.Printf("Processed %d symbols, taking a short break...\n", i+1)
		}
	}
}

func updateSymbols(c *Crawler) {
	symbols := fetchSymbols("US", c.config.FinnhubAPI)
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
