package websocket

import (
	"health_backend/internal/chat"

	"github.com/gin-gonic/gin"
)

func MapChatRoutes(chatGroup *gin.RouterGroup, h chat.WebSocketHandler) {
	// WebSocket Route (Real-time chat)
	chatGroup.GET("/ws", h.HandleWebSocket())
}
