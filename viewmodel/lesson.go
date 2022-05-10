package viewmodel

import "time"

// Lesson 课程
type Lesson struct {
	LessonName 	    string    `json:"lesson_name"`   	  // 课程名称
	LessonCreator 	string    `json:"user_id"`           // 课程发起者
	ClassList       []string  `json:"class_list"`       // 班级id列表
}

// ListObj 用户所创建的课程响应实体
type ListObj struct {
	LessonId		string        `json:"lesson_id"`			// 课程id
	LessonName   	string	      `json:"lesson_name"`	        // 课程名称
	CreatedAt    	time.Time     `json:"created_at"`           //  创建时间
	ClassName       []string      `json:"class_name_list"`      //  班级名称列表
}

// ClassObj 班级
type ClassObj struct {
	ClassName  string `json:"class_name"`    // 班级名
	ClassId    string `json:"class_id"`     //  班级id
}

// LessonEditor 编辑班级信息
type LessonEditor struct {
	LessonID      string    `json:"lesson_id"`
	LessonName    string   `json:"lesson_name"`
	ClassIdList []string   `json:"class_id_list"`
}




