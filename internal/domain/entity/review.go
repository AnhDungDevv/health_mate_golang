package entity

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	CustomerID   uint    `gorm:"not null" `
	ConsultantID uint    `gorm:"not null"`
	Rating       float64 `gorm:"not null"`
	Comment      string  `gorm:"type:text"`

	Customer   User `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Consultant User `gorm:"foreignKey;ConsultantID;constraint:OnDelete:CASCADE"`
}
