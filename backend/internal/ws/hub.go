package ws

import (
	"chatapp/internal/auth"
	"chatapp/internal/config"
	"chatapp/internal/constants"
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Hub struct {
	mu       sync.RWMutex
	byRoom   map[string]map[*websocket.Conn]struct{}
	upgrader websocket.Upgrader
	pub      *Publisher
}

func BuildHub() *Hub {
	return &Hub{
		byRoom: make(map[string]map[*websocket.Conn]struct{}),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *Hub) WithPublisher(p *Publisher) *Hub {
	h.pub = p
	return h
}

func (h *Hub) RegisterRoutes(r *gin.Engine, cfg config.AppConfig) {
	group := r.Group(constants.APIv1 + "/ws")
	group.Use(auth.WebSocketAuthMiddleware(cfg.JWTSecret))
	group.GET("", h.handleWS)
}

func (h *Hub) handleWS(c *gin.Context) {
	roomID := c.Query("roomId")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roomId is required"})
		return
	}
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	h.addConn(roomID, conn)
	uid := c.GetString("uid")
	log.Printf("ws connected uid=%s room=%s", uid, roomID)

	go h.readLoop(roomID, conn, uid)
}

func (h *Hub) addConn(roomID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	set, ok := h.byRoom[roomID]
	if !ok {
		set = make(map[*websocket.Conn]struct{})
		h.byRoom[roomID] = set
	}
	set[conn] = struct{}{}
}

func (h *Hub) removeConn(roomID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.byRoom[roomID]; ok {
		delete(set, conn)
		if len(set) == 0 {
			delete(h.byRoom, roomID)
		}
	}
}

// Broadcast sends JSON payload to all connections in a room.
func (h *Hub) Broadcast(roomID string, payload any) {
	h.mu.RLock()
	set := h.byRoom[roomID]
	h.mu.RUnlock()
	for conn := range set {
		_ = conn.WriteJSON(payload)
	}
}

// readLoop currently drains frames and closes on error; producer/consumer wiring next.
type clientSubmit struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (h *Hub) readLoop(roomID string, conn *websocket.Conn, userID string) {
	defer func() {
		h.removeConn(roomID, conn)
		_ = conn.Close()
	}()
	for {
		var in clientSubmit
		if err := conn.ReadJSON(&in); err != nil {
			break
		}
		if in.Type == "submit" && h.pub != nil {
			_ = h.pub.SubmitMessage(context.Background(), roomID, userID, in.Text)
		}
	}
}
