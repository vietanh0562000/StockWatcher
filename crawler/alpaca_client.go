package crawler

import (
	"encoding/json"
	"log"
	"net/http"
	"stockwatcher/config"
	"stockwatcher/models"

	"github.com/gorilla/websocket"
)

type AlpacaClient struct {
	config          *config.Config
	conn            *websocket.Conn
	isAuthenticated bool
	done            chan struct{}
}

type AlpacaTradeResponse struct {
	Symbol string       `json:"symbol"`
	Trade  models.Trade `json:"trade"`
}

type AlpacaQuoteResponse struct {
	Symbol string       `json:"symbol"`
	Quote  models.Quote `json:"quote"`
}

func NewAlpacaClient(config *config.Config) *AlpacaClient {
	return &AlpacaClient{
		config: config,
		done:   make(chan struct{}),
	}
}

func (client *AlpacaClient) GetTrade(symbol string) *models.Trade {
	url := client.config.StockURL + "/" + symbol + "/trades/latest"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Apca-Api-Key-Id", client.config.ResourceAPIKey)
	req.Header.Set("Apca-Api-Secret-Key", client.config.ResourceSecretKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error making request:", err)
		return nil
	}

	defer response.Body.Close()

	var tradeResponse AlpacaTradeResponse
	if response.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %d from Alpaca API\n", response.Status)
		return nil
	}

	if err := json.NewDecoder(response.Body).Decode(&tradeResponse); err != nil {
		log.Println("Error decoding response:", err)
		return nil
	}

	tradeResponse.Trade.SymbolName = tradeResponse.Symbol
	return &tradeResponse.Trade
}

func (client *AlpacaClient) GetQuote(symbol string) *models.Quote {
	url := client.config.StockURL + "/" + symbol + "/quotes/latest"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Apca-Api-Key-Id", client.config.ResourceAPIKey)
	req.Header.Set("Apca-Api-Secret-Key", client.config.ResourceSecretKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error making request:", err)
		return nil
	}

	defer response.Body.Close()

	var quoteResponse AlpacaQuoteResponse
	if response.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %d from Alpaca API\n", response.Status)
		return nil
	}

	if err := json.NewDecoder(response.Body).Decode(&quoteResponse); err != nil {
		log.Println("Error decoding response:", err)
		return nil
	}

	quoteResponse.Quote.SymbolName = quoteResponse.Symbol
	return &quoteResponse.Quote
}
