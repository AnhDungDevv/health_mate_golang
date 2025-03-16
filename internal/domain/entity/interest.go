package entity

import (
	"time"
)

type Interest struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Name      string `gorm:"size:255;not null;unique" json:"name"`
	ImageURL  string `gorm:"type:text" json:"image_url"`
	Customers []User `gorm:"many2many:user_interests;" json:"-"`
}
