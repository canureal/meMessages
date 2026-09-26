package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu		sync.RWMutex
	clients	map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
	}
}

/*
	We need mutex, because if we don't use the mutex,
	we could get race conditions from it.
	And it is bad, really.
	So when we using that concurrency,
	1: Lock it
	2: Do what you do
	3: Unlock the mutex(with defer(keep in mind that there is a LIFO(Last In First Out) mechanism))
*/

func (h *Hub) Register(userId string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userId] = conn
}

func (h *Hub) Unregister(userId string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, userId)
}

// we store the message as a byte array for performance.
func (h *Hub) SendTo(userId string, message []byte) error {
	h.mu.RLock()
	conn, ok := h.clients[userId]
	defer h.mu.RUnlock()

	if !ok {
		return nil
	}
	return conn.WriteMessage(websocket.TextMessage, message)
}


