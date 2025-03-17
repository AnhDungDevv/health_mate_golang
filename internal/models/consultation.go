package models

import (
	"time"

	"gorm.io/gorm"
)

type CallType string
type CallStatus string

const (
	Calling   CallStatus = "calling"
	Completed CallStatus = "completed"
	Cancelled CallStatus = "cancelled"
)

type Consultation struct {
	gorm.Model
	CustomerID   uint        `gorm:"not null;index"`
	ConsultantID uint        `gorm:"not null;index"`
	PaymentID    uint        `gorm:"not null;index"`
	StartTime    time.Time   `gorm:"not null"`
	EndTime      time.Time   `gorm:"not null"`
	ServiceType  ServiceType `gorm:"type:varchar(20);not null"`
	Status       CallStatus  `gorm:"type:varchar(20);not null"`
	Payment      Payment     `gorm:"foreignKey:PaymentID;constraint:OnDelete:CASCADE"`
	Customer     User        `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Consultant   User        `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE"`
	Price        Pricing     `gorm:"foreignKey:ServiceType;constraint:OnDelete:CASCADE"`
}
