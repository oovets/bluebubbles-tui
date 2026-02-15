package ws

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/bluebubbles-tui/models"
)

type Client struct {
	baseURL  string
	password string
	conn     *websocket.Conn
	Events   chan models.WSEvent
	done     chan struct{}
}

func NewClient(baseURL, password string) *Client {
	return &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		password: password,
		Events:   make(chan models.WSEvent, 10),
		done:     make(chan struct{}),
	}
}

// Connect dials the WebSocket endpoint
func (c *Client) Connect() error {
	// Convert https to wss, http to ws
	wsURL := c.baseURL
	wsURL = strings.ReplaceAll(wsURL, "https://", "wss://")
	wsURL = strings.ReplaceAll(wsURL, "http://", "ws://")

	// Append Socket.IO endpoint with EIO=4 for raw WebSocket transport
	u, err := url.Parse(fmt.Sprintf("%s/socket.io/?EIO=4&transport=websocket&guid=%s", wsURL, url.QueryEscape(c.password)))
	if err != nil {
		return err
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		NetDialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip TLS verification for self-signed certs
		},
	}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %v", err)
	}

	c.conn = conn

	// Start read loop in goroutine
	go c.readLoop()

	return nil
}

// readLoop handles incoming WebSocket messages
func (c *Client) readLoop() {
	defer close(c.Events)

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			// Connection closed or error
			return
		}

		msg := string(raw)

		switch {
		case strings.HasPrefix(msg, "0"):
			// Socket.IO open frame - server handshake, ignore
			continue

		case strings.HasPrefix(msg, "2"):
			// Socket.IO ping - respond with pong (3)
			c.conn.WriteMessage(websocket.TextMessage, []byte("3"))
			continue

		case strings.HasPrefix(msg, "42"):
			// Socket.IO event frame: 42[eventName, eventData]
			payload := msg[2:]

			// Parse as JSON array: [eventName, eventData]
			var arr []json.RawMessage
			if err := json.Unmarshal([]byte(payload), &arr); err != nil {
				continue
			}

			if len(arr) < 1 {
				continue
			}

			// Extract event type (first element)
			var eventType string
			if err := json.Unmarshal(arr[0], &eventType); err != nil {
				continue
			}

			// Extract data (second element if exists)
			var eventData json.RawMessage
			if len(arr) > 1 {
				eventData = arr[1]
			}

			// Send to events channel
			select {
			case c.Events <- models.WSEvent{Type: eventType, Data: eventData}:
			case <-c.done:
				return
			}

		default:
			// Unknown frame type, ignore
			continue
		}
	}
}

// Close closes the WebSocket connection
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	close(c.done)
	return c.conn.Close()
}
