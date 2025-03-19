package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	ConsultantID uuid.UUID `gorm:"not null;unique"`
	Balance      float64   `gorm:"type:numeric(10,2);not null;default:0"`
}

type TransactionType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
	Received TransactionType = "received"
)

type WalletTransaction struct {
	gorm.Model
	WalletID        uuid.UUID       `gorm:"not null;index"`
	Amount          float64         `gorm:"type:numeric(10,2);not null"`
	TransactionType TransactionType `gorm:"type:varchar(20);not null"`
	TransactionDate time.Time       `gorm:"not null"`
	ReferenceID     uuid.UUID       `gorm:"index"`
	Wallet          Wallet          `gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE"`
}
