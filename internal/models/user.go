package models

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	Customer   = "customer"
	Consultant = "consultant"
	Admin      = "admin"
)

type User struct {
	ID        uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Role     string  `gorm:"size:20;not null" json:"role" binding:"required,oneof=customer consultant" validate:"required,oneof=customer consultant"`
	Name     string  `gorm:"size:100;not null" json:"name" binding:"required,min=2,max=100" validate:"required,min=2,max=100"`
	Email    string  `gorm:"size:100;uniqueIndex;not null" json:"email" binding:"required,email" validate:"required,email"`
	Password string  `gorm:"not null" json:"password" binding:"required,min=6" validate:"required,min=6"`
	Phone    *string `gorm:"size:15" json:"phone" binding:"omitempty,e164" validate:"omitempty,e164"`
	Active   bool    `gorm:"default:true" json:"active" binding:"omitempty" validate:"omitempty"`

	City       *string    `gorm:"size:100" json:"city,omitempty"`
	Country    *string    `gorm:"size:100" json:"country,omitempty"`
	Expertiese Expertiese `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE" json:"expertiese,omitempty"`

	Interests []Interest `gorm:"type:json" json:"interests,omitempty"`

	Followers []Follow  `gorm:"foreignKey:FollowingID" json:"-" binding:"-" validate:"-"`
	Following []Follow  `gorm:"foreignKey:FollowerID" json:"-" binding:"-" validate:"-"`
	Post      []Post    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-" binding:"-" validate:"-"`
	Comment   []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-" binding:"-" validate:"-"`
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
