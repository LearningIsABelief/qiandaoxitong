package model

import "time"

type Connection struct {
	ClassRoomId string `json:"class_room_id"   gorm:"primary_key;column:class_room_id;not null"`
	ClassId     string `json:"class_id"     gorm:"column:class_id"`
	UserId      string `json:"user_id"  gorm:"column:user_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
