package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Message string `gorm:"type:text not null"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	PostID  uint   `gorm:"not null"`
	Post    Post   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
