package models

import "time"

type Conversation struct {
	ID        uint      `gorm:"primaryKey"`
	IsGroup   bool      `gorm:"not null"`
	Name      *string   `gorm:"default:null"`
	Messages  []Message `gorm:"foreignKey:ConversationID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Users     []User `gorm:"many2many:conversation_users;"`
}
