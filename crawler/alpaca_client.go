package crawler

import (
	"encoding/json"
	"fmt"
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

type AlpacaAuthMessage struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type AlpacaSubscribeMessage struct {
	Action string   `json:"action"`
	Trades []string `json:"trades,omitempty"`
	Quotes []string `json:"quotes,omitempty"`
	Bars   []string `json:"bars,omitempty"`
}

type AlpacaMessage struct {
	T   string      `json:"T"`           // Message type
	S   string      `json:"S,omitempty"` // Symbol
	Msg string      `json:"msg,omitempty"`
	Raw interface{} `json:"-"` // Store raw message for debugging
}

type AlpacaTradeResponse struct {
	Symbol string       `json:"symbol"`
	Trade  models.Trade `json:"trade"`
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

func (client *AlpacaClient) Connect() error {
	url := client.config.ResourceWURL
	fmt.Println("Connecting to websocket URL:", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		fmt.Println("Error when websocket connect :", err.Error())
		return err
	}

	client.conn = conn
	client.isAuthenticated = false
	fmt.Println("Connect socket successfully")
	return nil
}

func (client *AlpacaClient) Auth() error {
	authMsg := AlpacaAuthMessage{
		Action: "auth",
		Key:    client.config.ResourceAPIKey,
		Secret: client.config.ResourceSecretKey,
	}

	if err := client.conn.WriteJSON(authMsg); err != nil {
		return err
	}

	return nil
}

func (client *AlpacaClient) Subscribe(dataType string, symbols []string) error {
	msg := AlpacaSubscribeMessage{
		Action: "subscribe",
	}

	switch dataType {
	case "Trade":
		msg.Trades = append(msg.Trades, symbols...)
	case "Quote":
		msg.Quotes = append(msg.Quotes, symbols...)
	case "Bar":
		msg.Bars = append(msg.Bars, symbols...)
	}

	fmt.Println(msg)

	if err := client.conn.WriteJSON(msg); err != nil {
		fmt.Println("Failed to send message by socket: ", err.Error())
		return err
	}

	fmt.Println("Succeed to send message by socket")
	return nil
}

func (client *AlpacaClient) Unsubscribe(dataType string, symbols []string) error {
	msg := AlpacaSubscribeMessage{
		Action: "unsubscribe",
	}

	switch dataType {
	case "Trade":
		msg.Trades = append(msg.Trades, symbols...)
	case "Quote":
		msg.Quotes = append(msg.Quotes, symbols...)
	case "Bar":
		msg.Bars = append(msg.Bars, symbols...)
	}

	if err := client.conn.WriteJSON(msg); err != nil {
		fmt.Println("Failed to send message by socket: ", err.Error())
		return err
	}

	fmt.Println("Succeed to send message by socket")
	return nil
}

func (client *AlpacaClient) Listen() error {
	defer close(client.done)

	for {
		select {
		case <-client.done:
			return nil
		default:
			_, message, err := client.conn.ReadMessage()

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			processMessage(message)

		}
	}
}

func processMessage(message []byte) {
	var messages []json.RawMessage
	if err := json.Unmarshal(message, &messages); err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, message := range messages {
		processSingleMessage(message)
	}
}

func processSingleMessage(message json.RawMessage) {
	var alpacaMsg AlpacaMessage
	if err := json.Unmarshal(message, &alpacaMsg); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(alpacaMsg)

	switch alpacaMsg.T {
	case "success":
		fmt.Printf("✅ Success: %s\n", alpacaMsg.Msg)
	case "error":
		fmt.Printf("❌ Error: %s\n", alpacaMsg.Msg)
	case "subscription":
		fmt.Printf("📡 Subscription confirmed: %s\n", alpacaMsg.Msg)
	case "t": // Trade data
		fmt.Printf("✅ Trade: %s\n", alpacaMsg.Msg)
	case "q": // Quote data
		fmt.Printf("✅ Quote: %s\n", alpacaMsg.Msg)
	default:
		fmt.Printf("Unknown message type '%s': %s\n", alpacaMsg.T, string(message))
	}
}

func (client *AlpacaClient) Close() {
	if client.conn != nil {
		log.Println("Closing connection...")
		client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		client.conn.Close()
	}
	close(client.done)
}
