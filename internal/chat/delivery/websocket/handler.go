package websocket

import (
	"context"
	"encoding/json"
	"health_backend/config"
	"health_backend/internal/chat"
	"health_backend/internal/models"
	"health_backend/pkg/logger"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

// ChatMessage represents the structure of a chat message
type ChatMessage struct {
	Type      string         `json:"type"`               // "message", "typing", "read", etc.
	From      string         `json:"from"`               // Sender ID
	To        string         `json:"to"`                 // Receiver ID
	Content   string         `json:"content"`            // Message content
	Timestamp time.Time      `json:"timestamp"`          // Message timestamp
	Metadata  map[string]any `json:"metadata,omitempty"` // Additional metadata
}

var (
	clients sync.Map
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type HandleWebSocket struct {
	cfg    *config.Config
	chatUC chat.UseCase
	logger logger.Logger
}

// Constructor WebSocketHandler
func NewWebsocketHandler(cfg *config.Config, uc chat.UseCase, log logger.Logger) chat.WebSocketHandler {
	return &HandleWebSocket{
		cfg:    cfg,
		chatUC: uc,
		logger: log,
	}
}

// HandleWebSocket implements interfaces.WebSocketHandler.
func (h *HandleWebSocket) HandleWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			h.logger.Error("WebSocket upgrade failed: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to establish WebSocket connection"})
			return
		}
		defer conn.Close()

		// Setup ping/pong handlers
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			return nil
		})

		// Get client id
		clientID := c.Query("user_id")
		if clientID == "" {
			h.logger.Error("Missing user_id in query")
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Missing user_id"))
			return
		}

		ctx := context.Background()
		// Start ping ticker
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
						h.logger.Error("ping error:", err)
						return
					}
				}
			}
		}()

		// Add to list
		clients.Store(clientID, conn)
		h.logger.Info("New WebSocket connection: ", clientID)
		userUUID, err := uuid.FromString(clientID)
		if err != nil {
			h.logger.Error("Invalid UUID format", err, "ClientID:", clientID)
			return
		}

		err = h.chatUC.SetUserOnlineStatus(ctx, userUUID, true)
		if err != nil {
			h.logger.Error("Failed to set user online status", err, "UserID:", userUUID)
		}
		h.chatUC.NotifyUserOnline(ctx, clientID)

		// Cleanup when disconnected
		defer func() {
			clients.Delete(clientID)
			h.logger.Info("WebSocket disconnected: ", clientID)
		}()

		for {
			// Read message from client
			_, rawMsg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					h.logger.Error("Websocket read error:", err)
				}
				break
			}

			var chatMsg ChatMessage
			if err := json.Unmarshal(rawMsg, &chatMsg); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
				continue
			}

			if chatMsg.To == "" || chatMsg.Content == "" {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid message: missing required fields"))
				continue
			}

			fromUUID, err := uuid.FromString(chatMsg.From)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid sender UUID"))
				return
			}

			toUUID, err := uuid.FromString(chatMsg.To)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid receiver UUID"))
				return
			}

			// ✅ Chuyển `Metadata` từ `map[string]any` sang JSON string
			metadataJSON, err := json.Marshal(chatMsg.Metadata)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid metadata format"))
				return
			}

			// ✅ Lấy ID cuộc trò chuyện
			conversationID, err := h.chatUC.GetConversationID(ctx, fromUUID, toUUID)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to get conversation ID"))
				return
			}

			// ✅ Tạo tin nhắn mới với kiểu dữ liệu đúng
			message := &models.Message{
				ID:             uuid.Must(uuid.NewV4()),
				SenderID:       fromUUID,
				ConversationID: conversationID,
				ReceiverID:     toUUID,
				Content:        chatMsg.Content,
				Metadata:       string(metadataJSON),
				IsRead:         false,
				ReadAt:         nil,
				CreatedAt:      time.Now(),
			}
			messageJSON, err := json.Marshal(message)
			if err != nil {
				h.logger.Error("Failed to marshal message to JSON:", err)
			} else {
				h.logger.Info("Message JSON:", string(messageJSON))
			}

			// Gửi tin nhắn qua UseCase
			err = h.chatUC.SendMessage(ctx, conversationID, message)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to send message"))
				continue
			}

		}
	}
}

func SendToClient(clientID string, message []byte) {
	if conn, exists := clients.Load(clientID); exists {
		wsConn, ok := conn.(*websocket.Conn)
		if !ok {
			clients.Delete(clientID)
			return
		}

		err := wsConn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
		if err != nil {
			clients.Delete(clientID)
			wsConn.Close()
			return
		}

		if err := wsConn.WriteMessage(websocket.TextMessage, message); err != nil {
			clients.Delete(clientID)
			wsConn.Close()
		}
	}
}
