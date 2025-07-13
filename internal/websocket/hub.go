package websocket

import (
	"sync"
)

type Client struct {
	Username string
	Conn     Connection
}

type Connection interface {
	WriteJSON(v interface{}) error
	ReadJSON(v interface{}) error
	Close() error
}

type Hub struct {
	clients map[string]map[Connection]struct{} // 用户名 -> 连接集合
	lock    sync.RWMutex
}

var hub = &Hub{
	clients: make(map[string]map[Connection]struct{}),
}

// 用户上线，注册连接
func (h *Hub) Register(username string, conn Connection) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.clients[username] == nil {
		h.clients[username] = make(map[Connection]struct{})
	}
	h.clients[username][conn] = struct{}{}
}

// 用户下线，移除连接
func (h *Hub) Unregister(username string, conn Connection) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if conns, ok := h.clients[username]; ok {
		delete(conns, conn)
		if len(h.clients[username]) == 0 {
			delete(h.clients, username)
		}
	}
}

// 向某个用户的所有在线端推送消息
func (h *Hub) SendToUser(username string, msg interface{}) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for conn := range h.clients[username] {
		conn.WriteJSON(msg)
	}
}

// IsUserOnline checks if a user has at least one active connection.
func (h *Hub) IsUserOnline(username string) bool {
	h.lock.RLock()
	defer h.lock.RUnlock()
	connections, ok := h.clients[username]
	return ok && len(connections) > 0
}

// GetHub provides access to the global hub instance.
// This is a simple way to allow other packages to access the hub.
func GetHub() *Hub {
	return hub
}
