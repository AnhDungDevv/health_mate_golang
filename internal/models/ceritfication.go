package models

import "time"

type Certification struct {
	ID           uint   `gorm:"primaryKey"`
	ConsultantID uint   `gorm:"not null;index"`
	Name         string `gorm:"size:255;not null"`
	ImageURL     string `gorm:"type:text"`
	IssuedBy     string `gorm:"size:255"`
	IssuedDate   time.Time

	Consultant User `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE;"`
}
