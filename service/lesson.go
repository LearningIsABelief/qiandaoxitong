package service

import (
	"github.com/lexkong/log"
	"qiandao/model"
	"qiandao/pkg/app"
	"qiandao/pkg/util"
	"qiandao/store"
	"qiandao/viewmodel"
)

// CreateLesson 创建课程
func CreateLesson(lessonParam *viewmodel.Lesson) error{
	// 1.当前创建者的课程名不能重复
	_,ok := store.LessonIsExist(lessonParam);if ok{
		log.Errorf(app.ErrLessonExist,"当前用户创建重复课程")
		return app.ErrLessonExist
	}
	// 处理课程表实体并加入数据库
	lesson := &model.Lesson{
		LessonID:		util.GetUUID(),
		LessonName:    lessonParam.LessonName,
		LessonCreator: lessonParam.LessonCreator,
	}

	// 遍历班级id列表，创建中间表实体，加入切片
    classLessonSlice := make([]model.ClassLesson,0)
	for _, v := range lessonParam.ClassList {
		classLesson := model.ClassLesson{
			ClassLessonID:util.GetUUID(),
			ClassID:v.ClassId,
			LessonID:lesson.LessonID,
			ClassName: v.ClassName,
			LessonName:lesson.LessonName,
		}
		// 追加到最终结果集中
		classLessonSlice = append(classLessonSlice,classLesson)
	}

	// 存入数据库
	err := store.InsertLesson(lesson,classLessonSlice)
	if err != nil{
		return err
	}
	return err
}

// GetCreateLessonList 获取当前用户创建的所有课程
func GetCreateLessonList(userId string)(lessonList []*viewmodel.ListObj,err error) {
	// 根据userId查询数据库,获取相应的数据
	lessonList,_ = store.GetLessonList(userId)
	if err != nil {
		return nil,err
	}
	return lessonList,err
}

// GetJoinLessonList 获取当前用户加入的所有课程
func GetJoinLessonList(classId string)(lessonList []*viewmodel.ListObj,err error) {
	lessonList,err = store.GetJoinLessonList(classId)
	if err != nil {
		return nil,err
	}
	return lessonList,err
}

// EditorLesson 编辑课程信息
func EditorLesson(lesson *viewmodel.LessonEditor)(err error){
	// 查询课程名是否更改
	err,OldLesson := store.GetLessonInfoByLessonId(lesson.LessonID)
	if err != nil {
		return err
	}
	if OldLesson.LessonName != lesson.LessonName{
		// 更新课程名称
		err = store.UpdateLessonName(lesson)
		if err != nil {
			return err
		}
	}
	// 删除该课程对应的班级
	err = store.DeleteClassIdByLessonId(lesson.LessonID)
	if err != nil {
		return err
	}
	// 存储需要插入中间表的数据
	classLessonSlice := make([]model.ClassLesson,0)
	for _,v := range lesson.ClassList{
		classLesson := model.ClassLesson{
			ClassLessonID:util.GetUUID(),
			ClassID: v.ClassId,
			LessonID: lesson.LessonID,
			ClassName: v.ClassName,
			LessonName: lesson.LessonName,
		}
		classLessonSlice = append(classLessonSlice,classLesson)
	}
	// 调用插入语句重新插入
	err = store.InsertClassLesson(classLessonSlice)
	if err != nil{
		return err
	}
	return nil
}
// RemoveLesson 移除所选课程
func RemoveLesson(lesson *viewmodel.LessonRemove)(err error){
//  判定传入参数是否不匹配
	err = store.LessonCreatorIsExist(lesson)
	if err != nil {
		return err
	}
//  调用移除功能
	err = store.RemoveLesson(lesson)
	if err != nil {
		return err
	}
	return nil
}

