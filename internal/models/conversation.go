package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Conversation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	IsGroup   bool      `gorm:"not null;default:false" json:"is_group"`     // false = 1-1, true = Nhóm
	Name      *string   `gorm:"default:null" json:"name,omitempty"`         // Chỉ dùng khi là nhóm
	Users     []User    `gorm:"many2many:conversation_users;" json:"users"` // Thành viên cuộc trò chuyện
	Messages  []Message `gorm:"foreignKey:ConversationID" json:"messages"`  // Danh sách tin nhắn
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
