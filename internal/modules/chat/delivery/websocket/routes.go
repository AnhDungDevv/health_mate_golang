package websocket

import (
	"health_backend/internal/chat/interfaces"

	"github.com/gin-gonic/gin"
)

func MapChatRoutes(chatGroup *gin.RouterGroup, h interfaces.WebSocketHandler) {
	// WebSocket Route (Real-time chat)
	chatGroup.GET("/ws", h.HandleWebSocket())
}
