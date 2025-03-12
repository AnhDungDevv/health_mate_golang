package models

import "time"

type Expertiese struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	ConsultantID  uint   `gorm:"not null;index" json:"consultant_id" binding:"required" validate:"required"`
	Category      string `gorm:"size:255;not null" json:"category" binding:"required" validate:"required"`
	VideoURL      string `gorm:"type:text" json:"video_url" binding:"omitempty,url" validate:"omitempty,url"`
	IdentityProof string `gorm:"size:255" json:"identity_proof" binding:"omitempty,url" validate:"omitempty,url"`

	Consultant *User `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE" json:"-" binding:"-" validate:"-"`
}
