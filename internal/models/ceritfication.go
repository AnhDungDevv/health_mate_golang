package models

import "time"

type Certification struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ConsultantID uint      `gorm:"not null;index" json:"consultant_id" binding:"required" validate:"required"`
	Name         string    `gorm:"size:255;not null" json:"name" binding:"required" validate:"required"`
	ImageURL     string    `gorm:"type:text" json:"image_url" binding:"omitempty,url" validate:"omitempty,url"`
	IssuedBy     string    `gorm:"size:255" json:"issued_by" binding:"omitempty" validate:"omitempty"`
	IssuedDate   time.Time `gorm:"type:timestamp" json:"issued_date" binding:"required" validate:"required"`

	Consultant *User `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE" json:"-" binding:"-" validate:"-"`
}
