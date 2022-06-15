package viewmodel

import "time"

// Lesson 创建课程
type Lesson struct {
	LessonName 	    string           `form:"lesson_name" json:"lesson_name"`   	    // 课程名称
	LessonCreator 	string           `form:"user_id"     json:"user_id"`           // 课程发起者
	ClassList       []ClassContext   `form:"class_list"  json:"class_list"`            // 班级对象列表
}
// ClassContext 班级信息
type ClassContext struct {
	ClassId    string 		`form:"class_id"    json:"class_id"`     			// 班级id
	ClassName  string 		`form:"class_name"  json:"class_name"`    		   // 班级名
}

// ListObj 用户所创建的课程响应实体
type ListObj struct {
	LessonId		string        `form:"lesson_id"        json:"lesson_id"`			//课程id
	LessonName   	string	      `form:"lesson_name"      json:"lesson_name"`	       // 课程名称
	CreatedAt    	time.Time     `form:"created_at"       json:"created_at"`         //  创建时间
	ClassName       []string      `form:"class_name_list"  json:"class_name_list"`   // 班级名称列表
}

// LessonClass 用户创建的所有课程的相应
type LessonClass struct {
	LessonId		string        `form:"lesson_id"        json:"lesson_id"`		    // 课程id
	LessonName   	string	      `form:"lesson_name"      json:"lesson_name"`	       // 课程名称
	CreatedAt    	time.Time     `form:"created_at"       json:"created_at"`         //  创建时间
	ClassName       string        `form:"class_name"       json:"class_name"`        //  班级名
}

// ClassObj 用户创建的班级、用户加入的班级
type ClassObj struct {
	ClassName  string 		`form:"class_name"  json:"class_name"`    		// 班级名
	ClassId    string 		`form:"class_id"    json:"class_id"`     	   // 班级id
	DeletedAt   *time.Time  `gorm:"deleted_at" sql:"index"`	      		  // 删除时间

}

// LessonEditor 编辑班级信息
type LessonEditor struct {
	LessonID      string        `form:"lesson_id"     json:"lesson_id"`			 // 课程id
	LessonName    string        `form:"lesson_name"   json:"lesson_name"`		//  课程名
    ClassList    []ClassContext `form:"class_list"         json:"class_list"`            // 班级对象列表
}

// LessonRemove 移除课程
type LessonRemove struct {
	LessonID      string `form:"lesson_id"      json:"lesson_id"`        // 课程id
	LessonCreator string `form:"lesson_creator" json:"lesson_creator"`  // 课程创建者
}




