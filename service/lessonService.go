package service

import (
	"errors"
	"qiandao/model"
	"qiandao/pkg/util"
	"qiandao/store"
	"qiandao/viewmodel"
)

// CreateLesson 创建课程
func CreateLesson(lessonParam *viewmodel.Lesson) error {
	// 1.当前创建者的课程名不能重复
	_, ok := store.LessonIsExist(lessonParam)
	if ok {
		return errors.New("课程名称重复，不可创建")
	}
	// 处理课程表实体并加入数据库
	lesson := &model.Lesson{
		LessonID:      util.GetUUID(),
		LessonName:    lessonParam.LessonName,
		LessonCreator: lessonParam.LessonCreator,
	}

	// 遍历班级id列表，创建中间表实体，加入切片
	classLessonSlice := make([]model.ClassLesson, 0)

	for _, v := range lessonParam.ClassList {
		classLesson := model.ClassLesson{
			ClassLessonID: util.GetUUID(),
			ClassID:       v,
			LessonID:      lesson.LessonID,
		}
		classLessonSlice = append(classLessonSlice, classLesson)
	}

	// 存入数据库
	err := store.InsertLesson(lesson, classLessonSlice)
	if err != nil {
		return err
	}
	return nil
}

// GetCreateLessonList 获取当前用户创建的所有课程
func GetCreateLessonList(userId string) (lessonList []*viewmodel.ListObj, err error) {
	// 根据userId查询数据库,获取相应的数据
	lessonList, _ = store.GetLessonList(userId)
	if err != nil {
		return nil, err
	}
	return lessonList, err
}

func GetJoinLessonList(classId string) (lessonList []*viewmodel.ListObj, err error) {
	lessonList, err = store.GetJoinLessonList(classId)
	if err != nil {
		return nil, err
	}
	return lessonList, err
}
