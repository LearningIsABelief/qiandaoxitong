package viewmodel

import "time"

// CreateCheckin 创建签到请求结构体
type CreateCheckin struct {
	LessonID    string `json:"lesson_id" form:"lesson_id"`
	Duration    int    `json:"duration" form:"duration"`
	CreatorID   string `json:"creator_id" form:"creator_id"`
	CheckinCode string `json:"checkin_code" form:"checkin_code"`
}

// Checkin 学生签到请求结构体
type Checkin struct {
	CheckinID   string `json:"checkin_id" form:"checkin_id"`
	UserID      string `json:"user_id" form:"user_id"`
	UserName    string `json:"user_name" form:"user_name"`
	CheckinCode string `json:"checkin_code" form:"checkin_code"`
}

type List struct {
	ClassName string `json:"class_name"`
	UserName  string `json:"user_name"`
}

// CheckinDetailsResponse 签到详情响应结构体
type CheckinDetailsResponse struct {
	TotalList     []List `json:"total_list"`
	CheckedInList []List `json:"checkin_list"`
	NotCheckList  []List `json:"not_check_list"`
}

// ListResponse 已创建签到列表/签到记录列表响应结构体
type ListResponse struct {
	CheckinID  string    `json:"checkin_id"`
	LessonName string    `json:"lesson_name"`
	BeginTime  time.Time `json:"begin_time"`
	State      bool      `json:"state"`
}
