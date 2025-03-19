package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ServiceType string

const (
	AudioCall ServiceType = "audio_call"
	VideoCall ServiceType = "video_call"
	Chat      ServiceType = "chat"
)

type Pricing struct {
	gorm.Model
	ConsultantID uuid.UUID        `gorm:"not null;index"`
	ServiceType  ServiceType `gorm:"type:varchar(20);not null"`
	Price        float64     `gorm:"type:numeric(10,2);not null"`
	IsFree       bool        `gorm:"default:false"`
}
