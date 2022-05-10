package model

import "time"

type CheckedIn struct {
	ID        string `json:"id" gorm:"primaryKey"`
	CheckinID string `json:"checkin_id" sql:"index"`
	UserID    string `json:"user_id" sql:"index"`
	UserName  string `json:"user_name"`
	State     bool   `json:"state"`

	CreatedAt time.Time  `gorm:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at"`
	DeletedAt *time.Time `gorm:"deleted_at"`
}
