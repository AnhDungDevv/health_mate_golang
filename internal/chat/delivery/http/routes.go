package http

import (

	"github.com/gin-gonic/gin"
)

func MapChatRoutes(chatGroup *gin.RouterGroup, h chat.ChatHandler) {
	// Conversation Routes
	chatGroup.GET("/conversations", h.GetConversations())
	chatGroup.POST("/conversations", h.CreateConversation())
	chatGroup.GET("/conversations/:conversation_id", h.GetConversation())
	chatGroup.DELETE("/conversations/:conversation_id", h.DeleteConversation())

	// Message Routes
	chatGroup.GET("/conversations/:conversation_id/messages", h.GetMessages())
	chatGroup.POST("/conversations/:conversation_id/messages", h.SendMessage())
	chatGroup.PUT("/messages/:message_id", h.UpdateMessage())
	chatGroup.POST("/messages/:message_id/unsend", h.UnsendMessage())
}
