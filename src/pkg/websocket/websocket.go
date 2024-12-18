package websocket

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"sync"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
	ID   uint // User ID
}

type Message struct {
	Type    string          `json:"type"`
	RoomID  uint            `json:"room_id,omitempty"`
	Content json.RawMessage `json:"content"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	rooms      map[uint]map[*Client]bool
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[uint]map[*Client]bool),
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
				close(client.Send)

				// Remove from all rooms
				h.mu.Lock()
				for _, room := range h.rooms {
					delete(room, client)
				}
				h.mu.Unlock()
			}

		case message := <-h.broadcast:
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			switch msg.Type {
			case "room":
				h.mu.RLock()
				room := h.rooms[msg.RoomID]
				h.mu.RUnlock()

				for client := range room {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(room, client)
					}
				}

			case "private":
				var privateMsg struct {
					ToUserID uint `json:"to_user_id"`
				}
				if err := json.Unmarshal(msg.Content, &privateMsg); err != nil {
					continue
				}

				for client := range h.clients {
					if client.ID == privateMsg.ToUserID {
						select {
						case client.Send <- message:
						default:
							close(client.Send)
							delete(h.clients, client)
						}
						break
					}
				}
			}
		}
	}
}

func (h *Hub) JoinRoom(client *Client, roomID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
}

func (h *Hub) LeaveRoom(client *Client, roomID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		delete(room, client)
	}
}
