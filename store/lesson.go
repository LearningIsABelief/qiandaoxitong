package store

import (
	"qiandao/model"
	"qiandao/viewmodel"
)

// InsertLesson 插入课程信息
func InsertLesson(lesson *model.Lesson, classLesson []model.ClassLesson) error {
	tx := DB.Self.Begin()
	err := tx.Create(&lesson).Error;if err != nil{
		tx.Rollback()
		return err
	}
	for _, v := range classLesson{
		tx.Rollback()
		err = tx.Create(&v).Error;if err != nil{
			return err
		}
	}
	tx.Commit()
	return err
}


// GetLessonList 获取课程列表
func GetLessonList(userId string) ([]*viewmodel.ListObj,error){
	  // 返回给前端的最终结果集
	  var res []*viewmodel.ListObj
	  // 获取课程信息:课程名、创建时间
	  lesson := make([]*model.Lesson,0)
	  // 存储每个课程对应的班级
	  var lessonClass  = make(map[string][]string)
	  //  班级实体
	  classEntity := make([]viewmodel.ClassObj,0)

	  // 查询数据库，获取课程信息，存入lesson
	  err := DB.Self.Table("lesson").Select([]string{`lesson_id`,`lesson_name`,`created_at`}).Where("lesson_creator = ? ",userId).Find(&lesson).Error
	  if err != nil{
	  	return nil,err
	  }

	  // 获取每个课程对应的班级，存入lessonClass
	  for _,v := range lesson {
	 //  查询数据库，根据上述查询的课程id获取班级名称。
	  	err = DB.Self.Table("class").Select([]string{`class_name`,`class.class_id`}).Joins("inner join class_lesson on class.class_id = class_lesson.class_id").
	  		Where("lesson_id = ?",v.LessonID).Find(&classEntity).Error
	  	if err != nil{
	  		return nil,err
		}
	  // 存储每个课程对应的班级
	  	var tmp []string
	  	for _,v1 := range classEntity{
	  		tmp = append(tmp,v1.ClassName)
		}
		lessonClass[v.LessonID] = tmp
	  }
	  // 存入最终结果集
	  for _,v := range lesson{
	  	val := &viewmodel.ListObj{
	  		LessonName: v.LessonName,
	  		CreatedAt: v.CreatedAt,
	  		ClassName: lessonClass[v.LessonID],
		}
		res = append(res,val)
	  }
	return res,err
}

// LessonIsExist 查询当前用户是否创建重复的课程
func LessonIsExist(lessonParam *viewmodel.Lesson)(err error,ok bool) {
	var lesson model.Lesson
	err = DB.Self.Select(`lesson_name`).Where("lesson_creator = ?",lessonParam.LessonCreator).Find(&lesson).Error
	if err != nil {
		return err,false
	}
	return err,lesson.LessonName == lessonParam.LessonName
}
