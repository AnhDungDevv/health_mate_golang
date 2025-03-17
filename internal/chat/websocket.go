package chat

import "github.com/gin-gonic/gin"

type WebSocketHandler interface {
	HandleWebSocket() gin.HandlerFunc
}
