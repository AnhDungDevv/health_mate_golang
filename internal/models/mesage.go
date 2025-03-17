package models

import "time"

type Message struct {
	ID             uint   `gorm:"primaryKey"`
	ConversationID uint   `gorm:"not null"`
	SenderID       string `gorm:"not null"`
	Content        string `gorm:"not null"`
	CreatedAt      time.Time
	IsRead         bool `gorm:"default:false"`
}
