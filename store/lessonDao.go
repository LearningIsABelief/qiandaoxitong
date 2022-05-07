package store

import (
	"qiandao/model"
)

// InsertLesson 插入课程信息
func InsertLesson(lesson *model.Lesson, classLesson []model.ClassLesson) error {
	err := DB.Self.Create(&lesson).Error;if err != nil{
		return err
	}

	for _, v := range classLesson{
		err = DB.Self.Create(&v).Error;if err != nil{
			return err
		}
	}
	return err
}
