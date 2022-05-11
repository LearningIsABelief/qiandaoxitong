package model

import (
	"time"
)

// Checkin 签到结构体
type Checkin struct {
	CheckinID   string    `json:"checkin_id" gorm:"primaryKey"`
	CreatorID   string    `json:"creator"`
	LessonID    string    `json:"lesson_id"`
	BeginTime   time.Time `json:"begin_time"`
	EndTime     time.Time `json:"end_time"`
	CheckinCode string    `json:"checkin_code"`
	Longitude   string    `json:"longitude"`
	Latitude    string    `json:"latitude"`

	CreatedAt time.Time  `gorm:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at"`
	DeletedAt *time.Time `gorm:"deleted_at"`
}
