package models

import "time"

type Follow struct {
	FollowerID  uint      `gorm:"primaryKey;"`
	FollowingID uint      `gorm:"primaryKey;"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	// Foreignkey
	Follower  User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE"`
	Following User `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE"`
}
