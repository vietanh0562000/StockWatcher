package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"stockwatcher/config"
	"stockwatcher/models"

	"github.com/gorilla/websocket"
)

type FinnhubSocket struct {
	config *config.Config
	conn   *websocket.Conn
	done   chan struct{}
	quote  chan models.Quote
}

type SubscribeMessage struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
}

func NewFinnhubSocket(config *config.Config, quote chan models.Quote) *FinnhubSocket {
	return &FinnhubSocket{
		config: config,
		quote:  quote,
	}
}

func (fs *FinnhubSocket) OpenSocket() error {
	url := fs.config.FinnhubWURL + "?token=" + fs.config.FinnhubAPI
	fmt.Println("Connecting to websocket URL:", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		fmt.Println("Error when websocket connect :", err.Error())
		return err
	}

	fs.conn = conn
	fmt.Println("Connect socket successfully")
	return nil
}

func (fs *FinnhubSocket) Subscribe(symbol string) error {
	msg := SubscribeMessage{
		Type:   "subscribe",
		Symbol: symbol,
	}

	if err := fs.conn.WriteJSON(msg); err != nil {
		fmt.Println("Failed to send message by socket: ", err.Error())
		return err
	}

	fmt.Println("Succeed to send message by socket")
	return nil
}

func (fs *FinnhubSocket) Unsubscribe(symbol string) error {
	msg := SubscribeMessage{
		Type:   "unsubscribe",
		Symbol: symbol,
	}

	if err := fs.conn.WriteJSON(msg); err != nil {
		fmt.Println("Failed to send message by socket: ", err.Error())
		return err
	}

	fmt.Println("Succeed to send message by socket")
	return nil
}

func (fs *FinnhubSocket) Listen() {
	defer close(fs.done)

	for {
		select {
		case <-fs.done:
			return
		default:
			var message json.RawMessage

			// Program block here until a message come OR connection is closed
			err := fs.conn.ReadJSON(&message)
			if err != nil {
				log.Printf("Read error: %v", err)
				return
			}

			// Parse the message to determine its type
			var baseMsg struct {
				Type string `json:"type"`
			}

			if err := json.Unmarshal(message, &baseMsg); err != nil {
				log.Printf("Error parsing message type: %v", err)
				continue
			}

			switch baseMsg.Type {
			case "trade":
				var quote models.Quote
				if err := json.Unmarshal(message, &quote); err != nil {
					log.Printf("Error parsing trade data: %v", err)
					continue
				}
				fs.handleQuote(&quote)
			case "ping":
				// Respond to ping with pong
				pong := map[string]string{"type": "pong"}
				if err := fs.conn.WriteJSON(pong); err != nil {
					log.Printf("Error sending pong: %v", err)
				}
			default:
				log.Printf("Received message: %s", string(message))
			}
		}
	}
}

func (fs *FinnhubSocket) handleQuote(quote *models.Quote) {
	fs.quote <- *quote
}

func (fs *FinnhubSocket) CloseSocket() {
	if fs.conn != nil {
		log.Println("Closing connection...")
		fs.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		fs.conn.Close()
	}
	close(fs.done)
}
