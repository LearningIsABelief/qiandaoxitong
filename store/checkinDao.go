package store

import (
	"github.com/jinzhu/gorm"
	"qiandao/model"
)

type Tx struct {
	tx *gorm.DB
}

func GetTx() *Tx {
	return &Tx{tx: DB.Self}
}

func (tx *Tx) Begin() {
	tx.tx = tx.tx.Begin()
}

func (tx *Tx) RollBack() {
	tx.tx.Rollback()
}

func (tx *Tx) Commit() (err error) {
	err = tx.tx.Commit().Error
	return
}

// CreateCheckin
// @Description: 创建一个签到
// @Author zhandongyang 2022-05-08 22:44:50 ${time}
// @Param checkin
// @Return err
func (tx *Tx) CreateCheckin(checkin *model.Checkin) (err error) {
	err = tx.tx.Create(checkin).Error
	return
}

// GetCheckinById
// @Description: 根据签到id获取一个签到
// @Author zhandongyang 2022-05-08 22:46:45 ${time}
// @Param checkinID
// @Return checkin
// @Return err
func (tx *Tx) GetCheckinById(checkinID string) (checkin model.Checkin, err error) {
	err = tx.tx.Where("checkin_id = ?", checkinID).First(&checkin).Error
	return
}

// GetCheckinByCreator
// @Description: 根据用户id获取一个签到
// @Author zhandongyang 2022-05-09 16:58:08
// @Param creator
// @Return checkinList
// @Return err
func (tx *Tx) GetCheckinByCreator(creatorID string) (checkinList []model.Checkin, err error) {
	err = tx.tx.Where("creator_id = ?", creatorID).Find(&checkinList).Error
	return
}

//func (tx *Tx) GetACheckin(field, fieldValue string) (checkinLst []model.Checkin, err error) {
//	err = tx.tx.Where(fmt.Sprintf("%v = ?", field), fieldValue).Find(&checkinLst).Error
//	return
//}

// GetLessonByID
// @Description: 根据课程id获取课程
// @Author zhandongyang 2022-05-09 17:13:17
// @Param lessonID
// @Return lesson
// @Return err
func (tx *Tx) GetLessonByID(lessonID string) (lesson model.Lesson, err error) {
	err = tx.tx.Where("lesson_id = ?", lessonID).First(&lesson).Error
	return
}

// GetClassLstByLessonID
// @Description: 根据课程id获取需要签到的班级列表
// @Author zhandongyang 2022-05-08 23:12:01
// @Param lessonID
// @Return classList
// @Return err
func (tx *Tx) GetClassLstByLessonID(lessonID string) (classList []model.Class, err error) {
	err = tx.tx.Raw("select * from class where class_id IN (select class_id from class_lesson where lesson_id = ? and class_lesson.deleted_at is null) and class.deleted_at is null", lessonID).Scan(&classList).Error
	return
}

// GetStuLstByClassID
// @Description: 根据班级id获取需要签到的学生列表
// @Author zhandongyang 2022-05-08 23:02:59
// @Param classID
// @Param userID
// @Return stu
// @Return err
func (tx *Tx) GetStuLstByClassID(classID string) (stuList []model.User, err error) {
	err = tx.tx.Raw("select * from user where class_id = ? and user.deleted_at is null", classID).Scan(&stuList).Error
	return
}

// AddCheckinRec
// @Description: 添加学生签到记录
// @Author zhandongyang 2022-05-08 22:58:34 ${time}
// @Param stuCheckin
// @Return err
func (tx *Tx) AddCheckinRec(stuCheckin *model.CheckinRec) (err error) {
	err = tx.tx.Create(stuCheckin).Error
	return
}

// UpdateCheckinRecStateByID
// @Description: 根据签到记录id更新学生签到状态
// @Author zhandongyang 2022-05-09 17:32:27
// @Param stuCheckin
// @Return err
func (tx *Tx) UpdateCheckinRecStateByID(checkinRecID string, checkinRecState int) (err error) {
	err = tx.tx.Model(&model.CheckinRec{}).Where("checkin_rec_id = ?", checkinRecID).Update("state", checkinRecState).Error
	return
}

// GetCheckinRecByID
// @Description: 根据签到记录id获取一个签到记录
// @Author zhandongyang 2022-05-08 22:59:20
// @Param checkedID
// @Return checkedIn
// @Return err
func (tx *Tx) GetCheckinRecByID(checkinRecID string) (checkinRec model.CheckinRec, err error) {
	err = tx.tx.Where("checkin_rec_id = ?", checkinRecID).First(&checkinRec).Error
	return
}

// GetCheckinRecLstByCheckinID
// @Description: 根据签到id获取签到记录列表
// @Author zhandongyang 2022-05-10 16:11:54
// @Param checkinID
// @Return checkedInList
// @Return err
func (tx *Tx) GetCheckinRecLstByCheckinID(checkinID string) (checkedInList []model.CheckinRec, err error) {
	err = tx.tx.Where("checkin_id = ?", checkinID).Find(&checkedInList).Error
	return
}

// GetCheckinRecByUserID
// @Description: 根据用户id获取某个用户需要签到的列表
// @Author zhandongyang 2022-05-09 15:46:01
// @Param checkedID
// @Return checkedIn
// @Return err
func (tx *Tx) GetCheckinRecByUserID(userID string) (checkedInList []model.CheckinRec, err error) {
	err = tx.tx.Where("user_id = ?", userID).Order("end_time desc").Find(&checkedInList).Error
	return
}

func (tx *Tx) GetClassByUserID(userID string) (class model.Class, err error) {
	err = tx.tx.Raw("select * from class where class_id = (select class_id from connection where user_id = ? and connection.deleted_at is null) and class.deleted_at is null", userID).Scan(&class).Error
	return
}
