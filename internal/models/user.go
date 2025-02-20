package models

import (
	"strings"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"golang.org/x/crypto/bcrypt"
)

// User model for GORM
type User struct {
	UserID      uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName   string     `gorm:"type:varchar(30);not null"`
	LastName    string     `gorm:"type:varchar(30);not null"`
	Email       string     `gorm:"type:varchar(60);uniqueIndex"`
	Password    string     `gorm:"type:varchar(255);not null"`
	Role        *string    `gorm:"type:varchar(10)"`
	About       *string    `gorm:"type:varchar(1024)"`
	Avatar      *string    `gorm:"type:varchar(512)"`
	PhoneNumber *string    `gorm:"type:varchar(20)"`
	Address     *string    `gorm:"type:varchar(250)"`
	City        *string    `gorm:"type:varchar(24)"`
	Country     *string    `gorm:"type:varchar(24)"`
	Gender      *string    `gorm:"type:varchar(10)"`
	Postcode    *int       `gorm:"type:int"`
	Birthday    *time.Time `gorm:"type:date"`
	CreatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
	LoginDate   *time.Time `gorm:"type:timestamp"`
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
	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)

	}
	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	return nil

}

//Prepare for update

func (u *User) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	}
	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
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
