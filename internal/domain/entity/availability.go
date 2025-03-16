package entity

import (
	"time"

	"gorm.io/gorm"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "Monday"
	Tuesday   DayOfWeek = "Tuesday"
	Wednesday DayOfWeek = "Wednesday"
	Thursday  DayOfWeek = "Thursday"
	Friday    DayOfWeek = "Friday"
	Saturday  DayOfWeek = "Saturday"
	Sunday    DayOfWeek = "Sunday"
)

type Availability struct {
	gorm.Model
	DayOfWeek DayOfWeek `gorm:"type:day_of_week_enum;not null"` // Dùng enum
	StartTime time.Time `gorm:"type:time;not null"`
	EndTime   time.Time `gorm:"type:time;not null"`
}
