package models

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:50;unique; not null"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text;not null"`
	ImageURL    string `gorm:"type:text"`
	User        []User `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE;"`
}
