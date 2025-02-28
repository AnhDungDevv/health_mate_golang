package models

import (
	"time"

	"gorm.io/gorm"
)

type NotificationStatus string

const (
	Pending   NotificationStatus = "pending"
	Read      NotificationStatus = "read"
	Delivered NotificationStatus = "delivered"
)

type NotificationType string

const (
	Success NotificationType = "success"
	Error   NotificationType = "error"
	Info    NotificationType = "info"
)

type EventType string

const (
	PaymentSuccess EventType = "payment_success"
	PaymentFailure EventType = "payment_failure"
	Withdrawal     EventType = "withdrawal"
	ProfileUpdate  EventType = "profile_update"
	TaskCompleted  EventType = "task_completed"
	SystemAlert    EventType = "system_alert"
)

type Notification struct {
	gorm.Model
	Title            string             `gorm:"size:255;not null"`
	Content          string             `gorm:"type:text;not null"`
	Status           NotificationStatus `gorm:"type:varchar(20);not null"`
	Type             NotificationType   `gorm:"type:varchar(20);not null"`
	Event            EventType          `gorm:"type:varchar(20);not null"`
	UserID           uint               `gorm:"not null"`
	User             User               `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	NotificationDate time.Time          `gorm:"not null"`
}

// func SendNotification(userID uint, title, content string, notifType NotificationType, eventType EventType) error {
// 	notification := Notification{
// 		Title:            title,
// 		Content:          content,
// 		Status:           Pending,
// 		Type:             notifType,
// 		Event:            eventType,
// 		UserID:           userID,
// 		NotificationDate: time.Now(),
// 	}

// 	if err := db.Create(&notification).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
