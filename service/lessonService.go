package service

import (
	"qiandao/model"
	"qiandao/store"
	"qiandao/viewmodel"
)

// CreateLesson 创建课程
func CreateLesson(lessonParam *viewmodel.Lesson) error{
	// 事务处理

	// 处理课程表实体并加入数据库
	lesson := &model.Lesson{
		LessonID:"",
		LessonName:    lessonParam.LessonName,
		LessonCreator: lessonParam.LessonCreator,
	}

	// 遍历班级id列表，创建中间表实体，加入切片
	 classLessonSlice := make([]model.ClassLesson,0)

	for _, v := range lessonParam.ClassList {
		classLesson := model.ClassLesson{
			ClassLessonID: "",
			ClassID:v,
			LessonID:lesson.LessonID,
		}
		classLessonSlice = append(classLessonSlice,classLesson)
	}

	// 存入数据库
	err := store.InsertLesson(lesson,classLessonSlice)
	if err != nil{
		return err
	}
	return nil
}
