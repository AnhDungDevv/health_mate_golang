package models

import (
	"time"

	"gorm.io/gorm"
)

type LiveStreamStatus string

const (
	LiveStreamActive    LiveStreamStatus = "active"
	LiveStreamEnded     LiveStreamStatus = "ended"
	LiveStreamScheduled LiveStreamStatus = "scheduled"
)

type LiveStream struct {
	gorm.Model
	ConsultantID uint             `gorm:"not null;index"`
	Title        string           `gorm:"type:varchar(100);not null"`
	Description  string           `gorm:"type:text"`
	Status       LiveStreamStatus `gorm:"type:varchar(20);not null"`
	StartTime    time.Time        `gorm:"not null"`
	EndTime      time.Time        `gorm:"not null"`
	StreamURL    string           `gorm:"type:varchar(255)"`
	Consultant   User             `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE"`
}

type LiveStreamComment struct {
	gorm.Model
	LiveStreamID uint       `gorm:"not null;index"`
	CustomerID   uint       `gorm:"not null;index"`
	Comment      string     `gorm:"type:text;not null"`
	Timestamp    time.Time  `gorm:"not null"`
	LiveStream   LiveStream `gorm:"foreignKey:LiveStreamID;constraint:OnDelete:CASCADE"`
	Customer     User       `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}
type LiveStreamView struct {
	gorm.Model
	LiveStreamID uint `gorm:"not null;index"`
	CustomerID   uint `gorm:"not null;index"`
}
type LiveStreamStats struct {
	gorm.Model
	LiveStreamID uint      `gorm:"not null;index"`
	ViewCount    int       `gorm:"not null"`
	RecordedAt   time.Time `gorm:"not null"`
}
