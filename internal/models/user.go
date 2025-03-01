package models

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	Customer   = "customer"
	Consultant = "consultant"
)

// User model for GORM
type User struct {
	gorm.Model
	Role           string          `gorm:"size:20;not null" json:"role" binding:"required,oneof=customer consultant" validate:"required,oneof=customer consultant"`
	Name           string          `gorm:"size:100;not null" json:"name" binding:"required,min=2,max=100" validate:"required,min=2,max=100"`
	Email          string          `gorm:"size:100;uniqueIndex;not null" json:"email" binding:"required,email" validate:"required,email"`
	Password       string          `gorm:"not null" json:"password" binding:"required,min=6" validate:"required,min=6"`
	Phone          *string         `gorm:"size:15" json:"phone" binding:"omitempty,e164" validate:"omitempty,e164"`
	Avatar         string          `gorm:"type:text;null" json:"avatar" binding:"omitempty" validate:"omitempty"`
	Status         string          `gorm:"size:20" json:"status" binding:"omitempty,oneof=active inactive banned" validate:"omitempty,oneof=active inactive banned"`
	Profile        *Profile        `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE" json:"profile,omitempty" binding:"-" validate:"-"`
	Certifications []Certification `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE" json:"certifications,omitempty" binding:"-" validate:"-"`
	Followers      []Follow        `gorm:"foreignKey:FollowingID" json:"-" binding:"-" validate:"-"`
	Following      []Follow        `gorm:"foreignKey:FollowerID" json:"-" binding:"-" validate:"-"`
	Post           []Post          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-" binding:"-" validate:"-"`
	Comment        []Comment       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-" binding:"-" validate:"-"`
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
