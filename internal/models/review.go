package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	CustomerID   uuid.UUID    `gorm:"not null" `
	ConsultantID uuid.UUID    `gorm:"not null"`
	Rating       float64 `gorm:"not null"`
	Comment      string  `gorm:"type:text"`

	Customer   User `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Consultant User `gorm:"foreignKey;ConsultantID;constraint:OnDelete:CASCADE"`
}
