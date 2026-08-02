package ws

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, verify origin
	},
}

type WSHandler struct {
	// Hub can be added later for full WebSocket support
}

func NewWSHandler() *WSHandler {
	return &WSHandler{}
}

func (h *WSHandler) Serve(c *gin.Context) {
	chatID := c.Query("chatId")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chatId required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		chatID: chatID,
		send:   make(chan []byte, 256),
		conn:   conn,
	}

	go client.writePump()
	go client.readPump()
}

// Client represents a WebSocket client with connection
type Client struct {
	chatID string
	send   chan []byte
	conn   *websocket.Conn
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// TODO: Handle incoming messages through hub
		_ = msg
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}
