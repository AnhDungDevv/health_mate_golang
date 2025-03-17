package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string    `gorm:"varchar(255);not null"`
	Content  string    `gorm:"type:text"`
	ImageURL string    `gorm:"type:text"`
	UserID   uint      `gorm:"not null"`
	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comment  []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
