package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Expertiese struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	ConsultantID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"consultant_id" binding:"required" validate:"required"`
	Category      string    `gorm:"size:255;not null" json:"category" binding:"required" validate:"required"`
	VideoURL      string    `gorm:"type:text" json:"video_url,omitempty" binding:"omitempty,url" validate:"omitempty,url"`
	IdentityProof string    `gorm:"size:255" json:"identity_proof,omitempty" binding:"omitempty,url" validate:"omitempty,url"`

	Consultant *User `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE" json:"-"`
}
