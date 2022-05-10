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
		err = tx.Create(&v).Error;if err != nil{
			tx.Rollback()
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
	  	err = DB.Self.Table("class").Select([]string{`class_name`,`class.class_id`}).
	  		Joins("inner join class_lesson on class.class_id = class_lesson.class_id").
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
	  		LessonId : v.LessonID,
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

func GetJoinLessonList(classId string) ([]*viewmodel.ListObj,error) {
	// 返回结果
	var resListObj []*viewmodel.ListObj
	// 创建课程实体
	var lesson []model.Lesson
	// 存入每个课堂对应的班级
	classLessonMap := make(map[string][]string)
	// 创建班级实体
	var classLesson []viewmodel.ClassObj
	// 根据中间表关联查询到当前班级加入的课堂
	err := DB.Self.Table("class_lesson").Select([]string{`lesson.lesson_name`,`lesson.created_at`,`lesson.lesson_id`}).
		Joins("inner join lesson on lesson.lesson_id = class_lesson.lesson_id").
		Joins("inner join class on class.class_id = class_lesson.class_id").Where("class_lesson.class_id = ?",classId).Find(&lesson).Error
	if err != nil {
		return nil,err
	}
	// 根据查询出的课堂id,去反查询，得到加入该课堂的相应班级
	for _,v := range lesson{
		err = DB.Self.Table("class").Select([]string{`class_name`}).Joins("inner join class_lesson on class.class_id = class_lesson.class_id").
			Where("class_lesson.lesson_id = ?",v.LessonID).Find(&classLesson).Error
		if err != nil {
			return nil,err
		}
		var tmp []string
		for _,v1 := range classLesson{
			tmp = append(tmp,v1.ClassName)
		}
		classLessonMap[v.LessonID] = tmp
	}
	for _,v := range lesson{
		vobj := &viewmodel.ListObj{
			LessonName: v.LessonName,
			CreatedAt: v.CreatedAt,
			ClassName: classLessonMap[v.LessonID],

		}
		resListObj = append(resListObj,vobj)
	}
	return resListObj,err
}

// UpdateLessonName  更新课程名称
func UpdateLessonName(lesson *viewmodel.LessonEditor)(err error) {
	tx := DB.Self.Begin()
//	首先查库，判断课程名是否修改
	err,oldLesson := GetLessonInfoByLessonId(lesson.LessonID)
	if err != nil {
		tx.Rollback()
		return err
	}
//	课程名修改，更新课程名
	if oldLesson.LessonName != lesson.LessonName{
		err = DB.Self.Table("lesson").Model(oldLesson).Where("lesson_id = ?",lesson.LessonID).
			Update("lesson_name",lesson.LessonName).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	DB.Self.Table("lesson").Find(&model.Lesson{},"lesson_id = ?",lesson.LessonID)
	tx.Commit()
	return nil
}

// GetLessonInfoByLessonId 根据课程id查找课程信息
func GetLessonInfoByLessonId(lessonId string) (err error,lesson *model.Lesson){
	// 初始化参数，不然会报空指针异常
	lesson = &model.Lesson{}
	err = DB.Self.Table("lesson").Find(&lesson,"lesson_id = ?",lessonId).Error
	if err != nil {
		return err,nil
	}
	return nil,lesson
}

// DeleteClassIdByLessonId 根据课程id删除班级id
func DeleteClassIdByLessonId(lessonId string)(err error){
	tx := DB.Self.Begin()
	err = DB.Self.Table("class_lesson").Where("lesson_id = ?",lessonId).
		Delete(&model.ClassLesson{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
// InsertClassLesson 插入中间表信息
func InsertClassLesson(classLessonSlice []model.ClassLesson)(err error) {
	tx := DB.Self.Begin()
	for _, v := range classLessonSlice{
		err = tx.Create(&v).Error;if err != nil{
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}







