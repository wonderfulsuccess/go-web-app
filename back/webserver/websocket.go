package webserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/wonderfulsuccess/go-web-app/back/logger"
)

// WSMessage represents the envelope shared between server and clients.
type WSMessage struct {
	Sender    string          `json:"sender"`
	Receiver  string          `json:"receiver"`
	Timestamp time.Time       `json:"timestamp"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}

// Hub orchestrates WebSocket clients and message routing.
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan WSMessage
	incoming   chan WSMessage
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan WSMessage, 32),
		incoming:   make(chan WSMessage, 32),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				if msg.Receiver != "" && msg.Receiver != client.id && msg.Receiver != "*" {
					continue
				}
				select {
				case client.send <- msg:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}

// SendMessage allows other packages to emit WebSocket messages.
func (h *Hub) SendMessage(msg WSMessage) {
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now().UTC()
	}
	h.broadcast <- msg
}

// Incoming exposes server-side visibility into messages pushed by clients.
func (h *Hub) Incoming() <-chan WSMessage {
	return h.incoming
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocket upgrades an HTTP request to a WebSocket connection.
func (h *Hub) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("failed to upgrade websocket: %v", err)
		return
	}

	clientID := c.Query("clientId")
	if clientID == "" {
		clientID = c.ClientIP()
	}

	client := &Client{
		id:   clientID,
		hub:  h,
		conn: conn,
		send: make(chan WSMessage, 16),
	}

	h.register <- client

	go client.writePump()
	go client.readPump()
}

// Client represents an active websocket connection.
type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan WSMessage
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(5120)
	_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	for {
		var msg WSMessage
		if err := c.conn.ReadJSON(&msg); err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				return
			}
			logger.Errorf("websocket read error: %v", err)
			return
		}

		if msg.Timestamp.IsZero() {
			msg.Timestamp = time.Now().UTC()
		}
		if msg.Sender == "" {
			msg.Sender = c.id
		}

		select {
		case c.hub.incoming <- msg:
		default:
		}

		c.hub.SendMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(msg); err != nil {
				logger.Errorf("websocket write error: %v", err)
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ID() string {
	return c.id
}
