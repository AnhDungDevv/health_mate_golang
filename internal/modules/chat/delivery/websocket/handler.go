package websocket

import (
	"context"
	"encoding/json"
	"health_backend/config"
	"health_backend/internal/chat/interfaces"
	"health_backend/pkg/logger"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

// ChatMessage represents the structure of a chat message
type ChatMessage struct {
	Type      string          `json:"type"`      // "message", "typing", "read", etc.
	From      string          `json:"from"`      // Sender ID
	To        string          `json:"to"`        // Receiver ID
	Content   string          `json:"content"`   // Message content
	Timestamp time.Time       `json:"timestamp"` // Message timestamp
	Metadata  map[string]any  `json:"metadata,omitempty"` // Additional metadata
}

var (
	clients sync.Map
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type HandleWebSocket struct {
	cfg         *config.Config
	chatUC      interfaces.UseCase
	logger      logger.Logger
	kafkaWriter *kafka.Writer
	redisClient *redis.Client
}

// Constructor cho WebSocketHandler
// Constructor WebSocketHandler
func NewWebsocketHandler(cfg *config.Config, uc interfaces.UseCase, log logger.Logger, kafkaWriter *kafka.Writer, redisClient *redis.Client) interfaces.WebSocketHandler {
	return &HandleWebSocket{
		cfg:         cfg,
		chatUC:      uc,
		logger:      log,
		kafkaWriter: kafkaWriter,
		redisClient: redisClient,
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

		//Match user as online in redis
		_, err = h.redisClient.Get(ctx, "online:"+clientID).Result()
		if err == redis.Nil {
			h.redisClient.Set(ctx, "online:"+clientID, "1", 10*time.Minute).Err()
		}
		if err != nil && err != redis.Nil {
			h.logger.Error("Failed to set Redis key:", err)
		}

		// Add to list
		clients.Store(clientID, conn)
		h.logger.Info("New WebSocket connection: ", clientID)

		// Cleanup when disconnected
		defer func() {
			clients.Delete(clientID)
			h.redisClient.Del(ctx, "online:"+clientID).Err()
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

			// Parse the message
			var chatMsg ChatMessage
			if err := json.Unmarshal(rawMsg, &chatMsg); err != nil {
				h.logger.Error("Failed to parse message:", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
				continue
			}

			// Validate message
			if chatMsg.To == "" || chatMsg.Content == "" {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid message: missing required fields"))
				continue
			}

			// Add metadata
			chatMsg.From = clientID
			chatMsg.Timestamp = time.Now()

			// Convert message back to JSON
			msgBytes, err := json.Marshal(chatMsg)
			if err != nil {
				h.logger.Error("Failed to marshal message:", err)
				continue
			}

			h.logger.Info("Received message: ", string(msgBytes))

			// Publish message to Kafka
			err = h.kafkaWriter.WriteMessages(ctx, kafka.Message{
				Key:   []byte(clientID),
				Value: msgBytes,
			})
			if err != nil {
				h.logger.Error("Failed to publish message to Kafka:", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to deliver message"))
			}
		}
	}
}

func SendToClient(clientID string, message []byte) {
	if conn, exists := clients.Load(clientID); exists {
		wsConn := conn.(*websocket.Conn)
		if err := wsConn.WriteMessage(websocket.TextMessage, message); err != nil {
			// Remove client if connection is broken
			clients.Delete(clientID)
			wsConn.Close()
		}
	}
}
