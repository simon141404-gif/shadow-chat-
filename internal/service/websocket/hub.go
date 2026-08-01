package websocket

import (
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	logger     *zap.Logger
	mu         sync.RWMutex
}

type Message struct {
	Type      string          `json:"type"`
	ChatID    string          `json:"chatId,omitempty"`
	SenderID  string          `json:"senderId,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type Client struct {
	hub      *Hub
	userID   string
	chatID   string
	send     chan []byte
}

func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.chatID] == nil {
				h.clients[client.chatID] = make(map[*Client]bool)
			}
			h.clients[client.chatID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.chatID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.chatID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.clients[message.ChatID]; ok {
				for client := range clients {
					select {
					case client.send <- mustMarshal(message):
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(chatID, senderID string, payload interface{}) {
	msg := Message{
		Type:     "message",
		ChatID:   chatID,
		SenderID: senderID,
	}
	if data, ok := payload.(json.RawMessage); ok {
		msg.Payload = data
	}
	h.broadcast <- &msg
}

func (h *Hub) GetClientCount(chatID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[chatID])
}

func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
