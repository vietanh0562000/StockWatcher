package crawler

import (
	"fmt"
	"log"
	"stockwatcher/config"
	"stockwatcher/models"
	"strings"

	"github.com/gorilla/websocket"
)

type PolygonClient struct {
	config          *config.Config
	conn            *websocket.Conn
	isAuthenticated bool
	done            chan struct{}
	quote           chan models.Quote
}

type PolygonActionMessage struct {
	Action string `json:"action"`
	Params string `json:"params"`
}

type PolygonResponseMessage struct {
	// Common fields
	Ev  string `json:"ev"`  // Event type
	Sym string `json:"sym"` // Symbol

	// Trade specific
	X *int     `json:"x,omitempty"` // Exchange
	P *float64 `json:"p,omitempty"` // Price
	S *int     `json:"s,omitempty"` // Size
	T *int64   `json:"t,omitempty"` // Timestamp

	// Quote specific
	BP *float64 `json:"bp,omitempty"` // Bid Price
	AP *float64 `json:"ap,omitempty"` // Ask Price
	BS *int     `json:"bs,omitempty"` // Bid Size
	AS *int     `json:"as,omitempty"` // Ask Size

	// Aggregate specific
	O *float64 `json:"o,omitempty"` // Open
	C *float64 `json:"c,omitempty"` // Close
	H *float64 `json:"h,omitempty"` // High
	L *float64 `json:"l,omitempty"` // Low
	V *int     `json:"v,omitempty"` // Volume

	// Status specific
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

const AUTHENTICATED = "authenticated"

func NewPolygonClient(config *config.Config, quote chan models.Quote) *PolygonClient {
	return &PolygonClient{
		config: config,
		quote:  quote,
	}
}

func (client *PolygonClient) Connect() error {
	url := client.config.ResourceWURL + "?apiKey=" + client.config.ResourceAPIKey
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

func (client *PolygonClient) Auth() error {
	authMsg := map[string]interface{}{
		"action": "auth",
		"params": client.config.ResourceAPIKey,
	}

	if err := client.conn.WriteJSON(authMsg); err != nil {
		return err
	}

	return nil
}

func (client *PolygonClient) Subscribe(dataType string, symbols []string) error {
	msg := PolygonActionMessage{
		Action: "subscribe",
		Params: createParams(dataType, symbols),
	}

	if err := client.conn.WriteJSON(msg); err != nil {
		fmt.Println("Failed to send message by socket: ", err.Error())
		return err
	}

	fmt.Println("Succeed to send message by socket")
	return nil
}

func (client *PolygonClient) Unsubscribe(dataType string, symbols []string) error {
	msg := PolygonActionMessage{
		Action: "unsubscribe",
		Params: createParams(dataType, symbols),
	}

	if err := client.conn.WriteJSON(msg); err != nil {
		fmt.Println("Failed to send message by socket: ", err.Error())
		return err
	}

	fmt.Println("Succeed to send message by socket")
	return nil
}

func createParams(dataType string, symbols []string) string {
	var params []string
	for _, symbol := range symbols {
		params = append(params, fmt.Sprintf("%s.%s", dataType, symbol))
	}

	return strings.Join(params, ",")
}

func (client *PolygonClient) Listen() error {
	defer close(client.done)

	for {
		select {
		case <-client.done:
			return nil
		default:
			var responseMsgs []PolygonResponseMessage

			// Program block here until a message come OR connection is closed
			err := client.conn.ReadJSON(&responseMsgs)
			if err != nil {
				log.Printf("Read error: %v", err)
				return err
			}

			for _, responseMsg := range responseMsgs {
				fmt.Println(responseMsg)
				switch responseMsg.Ev {
				case "Q":
					if !client.isAuthenticated {
						continue
					}
					var quote models.Quote

					client.handleQuote(&quote)
				case "status":
					if responseMsg.Message == AUTHENTICATED {
						client.isAuthenticated = true
					}
				default:
					log.Printf("Received message: %v", responseMsg)
				}
			}
		}
	}
}

func (client *PolygonClient) handleQuote(quote *models.Quote) {
	client.quote <- *quote
}

func (client *PolygonClient) Close() {
	if client.conn != nil {
		log.Println("Closing connection...")
		client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		client.conn.Close()
	}
	close(client.done)
}
