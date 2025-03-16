package http

import (
	"health_backend/config"
	"health_backend/internal/chat/interfaces"
	"health_backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type chatHandlers struct {
	cfg    *config.Config
	chatUC interfaces.UseCase
	logger logger.Logger
}

// ✅ Thay vì trả về `ConversationHandler`, trả về `ChatHandler`
func NewChathandler(cfg *config.Config, chatUC interfaces.UseCase, log logger.Logger) interfaces.ChatHandler {
	return &chatHandlers{
		cfg: cfg, chatUC: chatUC, logger: log,
	}
}

// Implement ConversationHandler
func (c *chatHandlers) GetConversations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Logic xử lý
	}
}
func (c *chatHandlers) GetConversation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Logic xử lý
	}
}
func (c *chatHandlers) CreateConversation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Logic xử lý
	}
}
func (c *chatHandlers) DeleteConversation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Logic xử lý
	}
}

// Implement MessageHandler
func (c *chatHandlers) GetMessages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Logic xử lý
	}
}
func (c *chatHandlers) SendMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Logic xử lý
	}
}
func (c *chatHandlers) UpdateMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Logic xử lý
	}
}
func (c *chatHandlers) DeleteMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Logic xử lý
	}
}
func (c *chatHandlers) UnsendMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Logic xử lý
	}
}
