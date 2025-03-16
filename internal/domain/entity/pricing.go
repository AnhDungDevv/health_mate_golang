package entity

import "gorm.io/gorm"

type ServiceType string

const (
	AudioCall ServiceType = "audio_call"
	VideoCall ServiceType = "video_call"
	Chat      ServiceType = "chat"
)

type Pricing struct {
	gorm.Model
	ConsultantID uint        `gorm:"not null;index"`
	ServiceType  ServiceType `gorm:"type:varchar(20);not null"`
	Price        float64     `gorm:"type:numeric(10,2);not null"`
	IsFree       bool        `gorm:"default:false"`
}
