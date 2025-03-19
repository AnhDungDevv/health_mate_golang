package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ConversationID uuid.UUID  `gorm:"type:uuid;index;not null" json:"conversation_id"`
	SenderID       uuid.UUID  `gorm:"type:uuid;index;not null" json:"sender_id"`
	ReceiverID     uuid.UUID  `gorm:"type:uuid;index;not null" json:"receiver_id"`
	Content        string     `gorm:"type:text;not null" json:"content"`
	Metadata       string     `gorm:"type:jsonb" json:"metadata,omitempty"`
	IsRead         bool       `gorm:"default:false" json:"is_read"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
