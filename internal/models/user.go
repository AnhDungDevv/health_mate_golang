package models

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User model for GORM
type User struct {
	gorm.Model
	RoleID   uint    `gorm:"index;not null"`
	Role     Role    `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	Name     string  `gorm:"size:100;not null"`
	Email    string  `gorm:"size:100;uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Phone    *string `gorm:"size:15"`
	Avatar   string  `gorm:"type:text;null"`
	Bio      string  `gorm:"type:text"`
	Status   string  `gorm:"size:20"`

	Followers []Follow `gorm:"foreignKey:FollowingID"`
	Following []Follow `gorm:"foreignKey:FollowerID"`

	Post    []Post    `gorm:"foreignKey;UserID;constraint:OnDelete:CASCADE"`
	Comment []Comment `gorm:"foreignKey;UserID;constraint:OnDelete:CASCADE"`
}

// Hash user password with bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Coompare password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
}

// Prepare for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}
	if u.Phone != nil {
		*u.Phone = strings.TrimSpace(*u.Phone)
	}
	return nil
}

// Prepare for update
func (u *User) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	if u.Phone != nil {
		*u.Phone = strings.TrimSpace(*u.Phone)
	}
	if u.Password != "" && !strings.HasPrefix(u.Password, "$2a$") {
		if err := u.HashPassword(); err != nil {
			return err
		}
	}
	return nil
}

// All Users response
type UsersList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*User `json:"users"`
}
type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
