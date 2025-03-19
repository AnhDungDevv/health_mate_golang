package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Follow struct {
	FollowerID  uuid.UUID      `gorm:"primaryKey;"`
	FollowingID uuid.UUID      `gorm:"primaryKey;"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	// Foreignkey
	Follower  User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE"`
	Following User `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE"`
}
