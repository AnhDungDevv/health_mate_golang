package chat

import (
	"github.com/gin-gonic/gin"
)

// Interface cho Conversation API
type ConversationHandler interface {
	GetConversations() gin.HandlerFunc
	GetConversation() gin.HandlerFunc
	CreateConversation() gin.HandlerFunc
	DeleteConversation() gin.HandlerFunc
}

// Interface cho Message API
type MessageHandler interface {
	GetMessages() gin.HandlerFunc
	SendMessage() gin.HandlerFunc
	UpdateMessage() gin.HandlerFunc
	DeleteMessage() gin.HandlerFunc
	UnsendMessage() gin.HandlerFunc
}

type ChatHandler interface {
	ConversationHandler
	MessageHandler
}
