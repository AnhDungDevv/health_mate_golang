package models

type Profile struct {
	ID           uint    `gorm:"primaryKey"`
	ConsultantID uint    `gorm:"uniqueIndex;not null"`
	Profession   string  `gorm:"size:255;not null"`
	Experience   int     `gorm:"not null"`
	Rating       float32 `gorm:"type:float"`
	TotalReviews int

	Consultant User `gorm:"foreignKey:ConsultantID;constraint:OnDelete:CASCADE;"`
}
