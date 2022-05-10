package model

import (
	"time"
)

type Checkin struct {
	CheckinID   string    `json:"checkin_id" gorm:"primaryKey"`
	CreatorID   string    `json:"creator"`
	LessonID    string    `json:"lesson_id"`
	BeginTime   time.Time `json:"begin_time"`
	EndTime     time.Time `json:"end_time"`
	CheckinCode string    `json:"checkin_code"`

	CreatedAt time.Time  `gorm:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at"`
	DeletedAt *time.Time `gorm:"deleted_at"`
}
