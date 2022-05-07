package model

import "time"

type Class struct {
	ClassId       string `json:"class_id"       gorm:"primary_key;column:class_id;not null"`
	ClassName     string `json:"class_name"     gorm:"column:class_name"`
	ClassCapacity string `json:"class_capacity" gorm:"column:class_capacity"`
	CreateId      string `json:"create_id"      gorm:"column:create_id"`

	CreatedAt time.Time  `gorm:"column:createdAt" json:"-"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"-"`
	DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index" json:"-"`
}
