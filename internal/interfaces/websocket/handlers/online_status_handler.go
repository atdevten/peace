package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	appjwt "github.com/atdevten/peace/internal/application/services/jwt"
	"github.com/atdevten/peace/internal/application/usecases"
	"github.com/atdevten/peace/internal/interfaces/http/handlers"
	httpmiddleware "github.com/atdevten/peace/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"bearer"},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// OnlineStatusHandler handles WebSocket connections for user online status
type OnlineStatusHandler struct {
	userOnlineStatusUC usecases.UserOnlineStatusUseCase
	jwtService         appjwt.Service
}

// NewOnlineStatusHandler creates a new OnlineStatusHandler
func NewOnlineStatusHandler(userOnlineStatusUC usecases.UserOnlineStatusUseCase, jwtService appjwt.Service) *OnlineStatusHandler {
	return &OnlineStatusHandler{
		userOnlineStatusUC: userOnlineStatusUC,
		jwtService:         jwtService,
	}
}

// HandleWebSocket handles WebSocket connection for online status
func (h *OnlineStatusHandler) HandleWebSocket(c *gin.Context) {
	// Prefer user from middleware (Authorization header)
	var userID, userEmail string
	if idVO, ok := httpmiddleware.GetUserIDFromGinContext(c); ok && idVO != nil {
		userID = idVO.String()
	}
	if emailVO, ok := httpmiddleware.GetUserEmailFromGinContext(c); ok && emailVO != nil {
		userEmail = emailVO.String()
	}

	// If not present, try Sec-WebSocket-Protocol subprotocol (e.g., ['bearer', token] or [token])
	if userID == "" || userEmail == "" {
		protocolHeader := c.Request.Header.Get("Sec-WebSocket-Protocol")
		if protocolHeader == "" {
			handlers.Error(c, "UNAUTHORIZED", "Missing credentials")
			return
		}
		parts := strings.Split(protocolHeader, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		var token string
		if len(parts) >= 2 && strings.EqualFold(parts[0], "bearer") {
			token = parts[1]
		} else if len(parts) >= 1 {
			token = parts[0]
		}
		if token == "" {
			handlers.Error(c, "UNAUTHORIZED", "Missing token in subprotocol")
			return
		}
		claims, err := h.jwtService.ValidateAccessToken(token)
		if err != nil {
			handlers.Error(c, "UNAUTHORIZED", "Invalid token in subprotocol")
			return
		}
		userID = claims.UserID
		userEmail = claims.Email
	}

	// Set user as online
	err := h.userOnlineStatusUC.SetUserOnline(c.Request.Context(), userID, userEmail)
	if err != nil {
		log.Printf("Failed to set user %s online: %v", userID, err)
		handlers.Error(c, "INTERNAL_ERROR", "Failed to set user online")
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed for user %s: %v", userID, err)
		if !c.Writer.Written() {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		_ = h.userOnlineStatusUC.SetUserOffline(c.Request.Context(), userID)
		return
	}
	defer func() {
		conn.Close()
		_ = h.userOnlineStatusUC.SetUserOffline(c.Request.Context(), userID)
		log.Printf("User %s disconnected from WebSocket", userID)
	}()

	// Send welcome message
	welcome := map[string]interface{}{
		"type": "connection_established",
		"data": map[string]interface{}{
			"user_id":  userID,
			"status":   "online",
			"endpoint": "/ws",
			"message":  "Single WebSocket endpoint for all operations",
			"supported_events": []string{
				"get_amount_online_users",
				"get_online_users_list",
				"ping",
				"pong",
			},
		},
	}
	if err := conn.WriteJSON(welcome); err != nil {
		log.Printf("Failed to write welcome message to user %s: %v", userID, err)
		return
	}

	const (
		readWait   = 60 * time.Second
		pingPeriod = 30 * time.Second
		writeWait  = 10 * time.Second
	)
	conn.SetReadLimit(1 << 20)
	_ = conn.SetReadDeadline(time.Now().Add(readWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(readWait))
	})

	// Ping ticker
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				deadline := time.Now().Add(writeWait)
				if err := conn.WriteControl(websocket.PingMessage, nil, deadline); err != nil {
					log.Printf("Ping write failed for user %s: %v", userID, err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	ctx := c.Request.Context()

	// Read loop (handle message types)
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if messageType != websocket.TextMessage {
			continue
		}

		var req struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data,omitempty"`
		}
		if err := json.Unmarshal(message, &req); err != nil {
			_ = conn.WriteJSON(map[string]interface{}{
				"type": "error",
				"data": map[string]interface{}{
					"code":    "INVALID_MESSAGE",
					"message": "Invalid JSON format",
				},
			})
			continue
		}

		switch req.Type {
		case "ping":
			if err := h.userOnlineStatusUC.UpdateUserLastSeen(ctx, userID); err != nil {
				log.Printf("Update last seen failed for %s: %v", userID, err)
			}
			_ = conn.WriteJSON(map[string]interface{}{
				"type": "pong",
				"data": map[string]interface{}{
					"ts": time.Now().Unix(),
				},
			})

		case "get_amount_online_users":
			n, err := h.userOnlineStatusUC.GetOnlineCount(ctx)
			if err != nil {
				_ = conn.WriteJSON(map[string]interface{}{
					"type": "error",
					"data": map[string]interface{}{
						"code":    "INTERNAL_ERROR",
						"message": "Failed to get online users",
					},
				})
				continue
			}
			_ = conn.WriteJSON(map[string]interface{}{
				"type": "amount_online_users",
				"data": map[string]interface{}{
					"count": n,
					"ts":    time.Now().Unix(),
				},
			})

		case "get_online_users_list":
			users, err := h.userOnlineStatusUC.GetOnlineUsers(ctx)
			if err != nil {
				_ = conn.WriteJSON(map[string]interface{}{
					"type": "error",
					"data": map[string]interface{}{
						"code":    "INTERNAL_ERROR",
						"message": "Failed to get online users",
					},
				})
				continue
			}
			var userDTOs []map[string]interface{}
			for _, u := range users {
				userDTOs = append(userDTOs, map[string]interface{}{
					"user_id":    u.UserID().String(),
					"user_email": u.UserEmail().String(),
					"is_online":  u.IsOnline(),
					"last_seen":  u.LastSeen().Format(time.RFC3339),
				})
			}
			_ = conn.WriteJSON(map[string]interface{}{
				"type": "online_users_list",
				"data": map[string]interface{}{
					"users": userDTOs,
					"count": len(userDTOs),
					"ts":    time.Now().Unix(),
				},
			})

		default:
			_ = conn.WriteJSON(map[string]interface{}{
				"type": "error",
				"data": map[string]interface{}{
					"code":    "UNKNOWN_TYPE",
					"message": "Unsupported message type",
				},
			})
		}
	}

	close(done)
}
