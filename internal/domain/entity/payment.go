package entity

import (
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string
type PaymentMethod string

const (
	PendingPayment   PaymentStatus = "pending"
	CompletedPayment PaymentStatus = "completed"
	FailedPayment    PaymentStatus = "failed"
)

const (
	CreditCard   PaymentMethod = "credit_card"
	PayPal       PaymentMethod = "paypal"
	BankTransfer PaymentMethod = "bank_transfer"
)

type Payment struct {
	gorm.Model
	ConsultationID uint          `gorm:"not null;index"`
	Amount         float64       `gorm:"type:numeric(10,2);not null"`
	Status         PaymentStatus `gorm:"type:varchar(20);not null"`
	PaymentDate    time.Time     `gorm:"not null"`
	PaymentMethod  PaymentMethod `gorm:"size:50"`
	TransactionID  string        `gorm:"unique"`
}
