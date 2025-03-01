package models

type Profile struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	ConsultantID uint    `gorm:"uniqueIndex;not null" json:"consultant_id" binding:"required" validate:"required"`
	Profession   string  `gorm:"size:255;not null" json:"profession" binding:"required" validate:"required"`
	Experience   int     `gorm:"not null" json:"experience" binding:"required,gte=0" validate:"required,gte=0"`
	Rating       float32 `gorm:"type:float" json:"rating" binding:"omitempty,gte=0,lte=5" validate:"omitempty,gte=0,lte=5"`
	TotalReviews int     `gorm:"type:int" json:"total_reviews" binding:"omitempty,gte=0" validate:"omitempty,gte=0"`

	Consultant *User `gorm:"foreignKey:ConsultantID;references:ID;constraint:OnDelete:CASCADE" json:"-" binding:"-" validate:"-"`
}
