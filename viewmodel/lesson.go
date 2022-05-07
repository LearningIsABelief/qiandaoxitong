package viewmodel

import "time"

type Lesson struct {
	LessonName 	    string    `json:"lesson_name"`   	  // 课程名称
	LessonCreator 	string    `json:"user_id"`           // 课程发起者
	ClassList       []string  `json:"class_list"`       // 班级id列表
}

type  ListObj struct {
	LessonName   	string	      `json:"lesson_name"`	        // 课程名称
	CreatedAt    	time.Time     `json:"created_at"`           //  创建时间
	ClassName       []string      `json:"class_name_list"`      //  课程名称列表
}

type ClassObj struct {
	ClassName  string `json:"class_name"`    // 班级名
	ClassId    string `json:"class_id"`     //  班级id
}



