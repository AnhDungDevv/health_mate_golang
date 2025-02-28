package models

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50;unique; not null"`
	User []User `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE;"`
}
