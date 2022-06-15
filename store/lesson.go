package store

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"qiandao/model"
	"qiandao/pkg/app"
	//"qiandao/pkg/util"
	"qiandao/viewmodel"
	"strings"
	"time"
)

// InsertLesson 插入课程信息
func InsertLesson(lesson *model.Lesson, classLesson []model.ClassLesson) error {
//  开启事务
	tx := DB.Self.Begin()
//  插入课程表
	err := tx.Create(&lesson).Error; if err != nil{
		log.Errorf(err,"插入课程记录失败")
		tx.Rollback()
		return app.InternalServerError
	}
//  批量插入中间表
	err = BulkInsert(DB.Self,classLesson)
	if err != nil {
		log.Errorf(err,"插入课程记录失败")
		tx.Rollback()
		log.Errorf(err,"批量插入失败")
	}
	tx.Commit()
	return err
}

// GetLessonList 获取当前用户创建课程列表
// 主要逻辑:根据课程创建者获取其所创建的所有课程id列表，然后根据课程id列表中间表查询，获取最后所需要的返回结果。
func GetLessonList(userId string) ([]*viewmodel.ListObj,error){
	  var lessonList []*viewmodel.ListObj 	 			        // 返回给前端的最终结果集
	  var queryData []*viewmodel.LessonClass           // 连表查询需要的结果集
	  var lessonClassMap  = make(map[string][]string) // 去重lessonID，记录每个课程所拥有的班级
	  var lessonObjList = make([]model.Lesson,0)		 // 根据课程创建者,获取所有课程对象

	  // 1.根据userID获取该用户所创建的课程
	  db := DB.Self.Table("lesson").Select([]string{`lesson_id`}).Where("lesson_creator = ?",userId).Find(&lessonObjList)
	 // 创建课程id列表
	 var lessonIDList = make([]string,len(lessonObjList))
	  // 存入id列表
	 for _,v := range lessonObjList{
		 lessonIDList = append(lessonIDList,v.LessonID)
	 }
	  // 错误处理
     if db.RowsAffected == 0 {
		log.Errorf(db.Error,"GetLessonList 查询lessonID列表失败")
		return lessonList,app.ErrRecordNotExist
	 }

	  // 2.根据课程id连表查询出最后返回的结果集
	// 2.根据课程id列表,获取所有最后结果集
	db = DB.Self.Table("class_lesson").Select([]string{`lesson_id`,`lesson_name`,`class_name`,`created_at`}).
		Where("lesson_id IN (?) AND deleted_at is null ",lessonIDList).Find(&queryData)
	//  错误处理
		if db.RowsAffected == 0 {
			log.Errorf(db.Error,"GetLessonList 连表查询失败")
			return lessonList,app.ErrRecordNotExist
		}

	  // 3.获取每个课程对应的班级，存入lessonClass
	  for _,v := range queryData {
	  // key是课程id+课程名+创建时间
		  mapKey := v.LessonId+","+v.LessonName+","+v.CreatedAt.String()
		  lessonClassMap[mapKey] = append(lessonClassMap[mapKey],v.ClassName)
	  }
	  // 4.存入最终结果集
	  for key,val := range lessonClassMap{
	    vals := strings.Split(key,",")   // 分割
		createdAt,_ := time.Parse("2006-01-02 15:04:05",vals[2]) // 字符串转成日期
	  	lessonObj := &viewmodel.ListObj{
	  		LessonId : vals[0],
	  		LessonName:vals[1],
	  		CreatedAt: createdAt,
	  		ClassName:val,
		}
		lessonList = append(lessonList,lessonObj)
	  }
	 log.Infof("查询用户创建列表成功%v",lessonList)
	 return lessonList,nil
}

// GetJoinLessonList 查询当前用户加入的课程
func GetJoinLessonList(classId string) ([]*viewmodel.ListObj,error) {
	var joinList []*viewmodel.ListObj			  // 返回结果
	var lesson []model.Lesson		  		 	 // 创建课程实体
    var queryData []*viewmodel.LessonClass      // 连表查询需要的结果集
	classLessonMap := make(map[string][]string) // 存入每个课堂对应的班级

	// 1.根据班级id查询出当前班级所加入的所有课程,获取了课程id列表
	db := DB.Self.Table("class_lesson").Select([]string{"lesson.lesson_id"}).
		Joins("INNER JOIN lesson ON lesson.lesson_id = class_lesson.lesson_id").
		Where("class_lesson.class_id = ?",classId).Find(&lesson)
	// 错误处理
	if db.RowsAffected == 0 {
		log.Errorf(db.Error,"查询用户所在班级加入的课堂失败")
		return nil,app.ErrRecordNotExist
	}
	// 课程id列表
	 lessonIDList := make([]string,len(lesson))
	 for _,v := range lesson {
		 lessonIDList = append(lessonIDList,v.LessonID)
	 }

	// 2.根据课程id列表,获取所有最后结果集
	db = DB.Self.Table("class_lesson").Select([]string{`lesson_id`,`lesson_name`,`class_name`,`created_at`}).
		Where("lesson_id IN (?) AND deleted_at is null ",lessonIDList).Find(&queryData)

	//  错误处理
	if db.RowsAffected == 0 {
		log.Errorf(db.Error,"GetLessonList 连表查询失败")
		return joinList,app.ErrRecordNotExist
	}

	// 3.处理相应课程对应的班级
	for _,v := range queryData {
		// key是课程id+课程名+创建时间
		mapKey := v.LessonId+","+v.LessonName+","+v.CreatedAt.String()
		classLessonMap[mapKey] = append(classLessonMap[mapKey],v.ClassName)
	}

	// 4.存入joinList
	for key,val := range classLessonMap{
		vals := strings.Split(key,",")   // 分割
		createdAt,_ := time.Parse("2006-01-02 15:04:05",vals[2]) // 字符串转成日期
		val := &viewmodel.ListObj{
			LessonId : vals[0],
			LessonName:vals[1],
			CreatedAt: createdAt,
			ClassName:val,
		}
		joinList = append(joinList,val)
	}

	log.Infof("获取用户加入课程成功")

	return joinList,nil
}

// UpdateLessonName  更新课程名称
func UpdateLessonName(lesson *viewmodel.LessonEditor)(err error) {
	tx := DB.Self.Begin()
//	更新课程名
    tx.Table("lesson").Model(&model.Lesson{}).Where("lesson_id = ?",lesson.LessonID).
		Update("lesson_name",lesson.LessonName)
	if tx.Error !=  nil{
		log.Errorf(tx.Error,"更新失课程名失败")
		tx.Rollback()
		return app.ErrUpdated
	}
	tx.Commit()
	log.Infof("更新课程名称成功%s",lesson.LessonName)
	return nil
}

// InsertClassLesson 插入中间表信息
func InsertClassLesson(classLessonSlice []model.ClassLesson)(err error) {
	tx := DB.Self.Begin()
	err = BulkInsert(tx,classLessonSlice)
	if err != nil{
		log.Errorf(tx.Error,"插入中间表信息失败")
		tx.Rollback()
		return err
	}
	tx.Commit()
	log.Infof("插入中间表信息成功%v",classLessonSlice)
	return nil
}

// RemoveLesson 移除课程
func RemoveLesson(lesson *viewmodel.LessonRemove)(err error){
// 移除课程表中的数据
	tx := DB.Self.Begin()
	db := DB.Self.Table("lesson").
		Where("lesson_id = ? and lesson_creator = ?",lesson.LessonID,lesson.LessonCreator).
		Delete(&model.Lesson{})
	if db.RowsAffected == 0 {
		log.Errorf(db.Error,"删除课程记录失败")
		tx.Rollback()
		return app.ErrDeleted
	}
//	移除中间表中的数据
    db = DB.Self.Table("class_lesson").
		Where("lesson_id = ?",lesson.LessonID).
		Delete(&model.ClassLesson{})
	if db.RowsAffected == 0 {
		log.Errorf(db.Error,"删除中间表记录失败")
		tx.Rollback()
		return app.ErrDeleted
	}
	tx.Commit()
    log.Infof("移除课程记录成功%d",db.RowsAffected)
	return nil
}

// LessonCreatorIsExist 查询lesson_creator是否和创建的课程匹配，防止误删除
func LessonCreatorIsExist(lesson *viewmodel.LessonRemove) error {
	db := DB.Self.Table("lesson").
		Find(&model.Lesson{},"lesson_creator = ? and lesson_id = ?",lesson.LessonCreator,lesson.LessonID)
	if db.Error != nil {
		log.Errorf(db.Error,"查询课程记录失败")
		return app.ErrRecordNotExist
	}
	log.Infof("查询指定用户创建的课程记录成功%d",db.RowsAffected)
	return nil
}
// LessonIsExist 查询当前用户是否创建重复的课程
func LessonIsExist(lessonParam *viewmodel.Lesson)(err error,ok bool) {
	var lesson model.Lesson
	db := DB.Self.Select(`lesson_name`).Where("lesson_creator = ?",lessonParam.LessonCreator).
		Find(&lesson)
	if db.RowsAffected == 0 {
		log.Errorf(db.Error,"查询课程记录失败")
		return app.ErrRecordNotExist,false
	}
	log.Infof("查询当前用户是否创建重复课程操作成功%d",db.RowsAffected)
	return nil,lesson.LessonName == lessonParam.LessonName
}

// GetLessonInfoByLessonId 根据课程id查找课程信息
func GetLessonInfoByLessonId(lessonId string) (err error,lesson *model.Lesson){
	// 初始化参数，不然会报空指针异常
	lesson = &model.Lesson{}
	db := DB.Self.Table("lesson").Find(&lesson,"lesson_id = ?",lessonId)
	if db.RowsAffected == 0 {
		log.Errorf(db.Error,"查询课程记录失败")
		return app.ErrRecordNotExist,nil
	}
	log.Infof("查询课程信息成功%d",db.RowsAffected)
	return nil,lesson
}

// DeleteClassIdByLessonId 根据课程id删除班级id
func DeleteClassIdByLessonId(lessonId string)(err error){
	tx := DB.Self.Begin()
	tx.Table("class_lesson").Where("lesson_id = ?",lessonId).
		Delete(&model.ClassLesson{})
	if tx.Error != nil {
		log.Errorf(tx.Error,"删除班级id失败")
		tx.Rollback()
		return app.ErrDeleted
	}
	tx.Commit()
	log.Infof("删除班级id是啊比%d",tx.RowsAffected)
	return nil
}

func BulkInsert(db *gorm.DB,data []model.ClassLesson) error{
	// 声明buffer缓冲器
	var buffer bytes.Buffer
	sql := "insert into `class_lesson` (`class_lesson_id`,`class_id`,`lesson_id`,`class_name`,`lesson_name`,`created_at`,`updated_at`) values"
	// 将字符串放到缓冲器的尾部
	if _,err := buffer.WriteString(sql);err != nil{
		return err
	}
	// 拼接字符串,构成多行插入语句
	for i,val := range data{
		if i == len(data)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s');",val.ClassLessonID,val.ClassID,val.LessonID,val.ClassName,val.LessonName,time.Now().Format("2006-01-02 15:04:05.000"),time.Now().Format("2006-01-02 15:04:05.000")))
		}else{
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s'),",val.ClassLessonID,val.ClassID,val.LessonID,val.ClassName,val.LessonName,time.Now().Format("2006-01-02 15:04:05.000"),time.Now().Format("2006-01-02 15:04:05.000")))
		}
	}
	// 执行插入
	return db.Exec(buffer.String()).Error

}








