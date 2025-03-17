package models

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	ConsultantID uint    `gorm:"not null;unique"`
	Balance      float64 `gorm:"type:numeric(10,2);not null;default:0"`
}

type TransactionType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
	Received TransactionType = "received"
)

type WalletTransaction struct {
	gorm.Model
	WalletID        uint            `gorm:"not null;index"`
	Amount          float64         `gorm:"type:numeric(10,2);not null"`
	TransactionType TransactionType `gorm:"type:varchar(20);not null"`
	TransactionDate time.Time       `gorm:"not null"`
	ReferenceID     uint            `gorm:"index"`
	Wallet          Wallet          `gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE"`
}
