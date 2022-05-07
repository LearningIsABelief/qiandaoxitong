package viewmodel

type Lesson struct {
	LessonName 	    string    `json:"lesson_name"`   	 // 课程名称
	LessonCreator 	string    `json:"user_id"`          // 课程发起者
	ClassList       []string  `json:"class_list"`      // 班级id列表
}
