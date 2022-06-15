package util

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"qiandao/model"
	"time"
)

// BulkInsert 中间表内容批量插入
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