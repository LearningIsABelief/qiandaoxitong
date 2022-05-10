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
	tx.tx.Begin()
}

func (tx *Tx) RollBack() {
	tx.tx.Rollback()
}

func (tx *Tx) Commit() {
	tx.tx.Commit()
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

// GetShouldCheckInClass
// @Description: 根据课程id获取需要签到的班级列表
// @Author zhandongyang 2022-05-08 23:12:01
// @Param lessonID
// @Return classList
// @Return err
func (tx *Tx) GetShouldCheckInClass(lessonID string) (classList []model.Class, err error) {
	err = tx.tx.Raw("select * from class where class_id IN (select class_id from class_lesson where lesson_id = ? and class_lesson.deleted_at is null) and class.deleted_at is null", lessonID).Scan(&classList).Error
	return
}

// GetShouldCheckInStu
// @Description: 获取需要签到的学生列表
// @Author zhandongyang 2022-05-08 23:02:59
// @Param classID
// @Param userID
// @Return stu
// @Return err
func (tx *Tx) GetShouldCheckInStu(classID string) (stuList []model.User, err error) {
	err = tx.tx.Raw("select * from user where class_id = ? and user.deleted_at is null", classID).Scan(&stuList).Error
	return
}

// GetLessonById
// @Description: 根据课程id获取课程
// @Author zhandongyang 2022-05-09 17:13:17
// @Param lessonID
// @Return lesson
// @Return err
func (tx *Tx) GetLessonById(lessonID string) (lesson model.Lesson, err error) {
	err = tx.tx.Where("lesson_id = ?", lessonID).First(&lesson).Error
	return
}

// AddCheckedIn
// @Description: 添加学生签到信息
// @Author zhandongyang 2022-05-08 22:58:34 ${time}
// @Param stuCheckin
// @Return err
func (tx *Tx) AddCheckedIn(stuCheckin *model.CheckedIn) (err error) {
	err = tx.tx.Create(stuCheckin).Error
	return
}

// UpdateCheckedIn
// @Description: 更新学生签到信息
// @Author zhandongyang 2022-05-09 17:32:27
// @Param stuCheckin
// @Return err
func (tx *Tx) UpdateCheckedIn(stuCheckedIn *model.CheckedIn) (err error) {
	err = tx.tx.Model(stuCheckedIn).Update("state", stuCheckedIn.State).Error
	return
}

// GetCheckedIn
// @Description: 获取单个已签到信息
// @Author zhandongyang 2022-05-08 22:59:20
// @Param checkedID
// @Return checkedIn
// @Return err
func (tx *Tx) GetCheckedIn(checkedID string) (checkedIn model.CheckedIn, err error) {
	err = tx.tx.Where("id = ?", checkedID).First(&checkedIn).Error
	return
}

// GetAllCheckedInByCheckinID
// @Description: 根据签到id获取需要签到的学生列表
// @Author zhandongyang 2022-05-10 16:11:54
// @Param checkinID
// @Return checkedInList
// @Return err
func (tx *Tx) GetAllCheckedInByCheckinID(checkinID string) (checkedInList []model.CheckedIn, err error) {
	err = tx.tx.Where("checkin_id = ?", checkinID).Find(&checkedInList).Error
	return
}

// GetAllCheckedInByUserID
// @Description: 获取某个学生需要签到的列表
// @Author zhandongyang 2022-05-09 15:46:01
// @Param checkedID
// @Return checkedIn
// @Return err
func (tx *Tx) GetAllCheckedInByUserID(userID string) (checkedInList []model.CheckedIn, err error) {
	err = tx.tx.Where("user_id = ?", userID).Find(&checkedInList).Error
	return
}

func (tx *Tx) GetClassByUserID(userID string) (class model.Class, err error) {
	err = tx.tx.Raw("select * from class where class_id = (select class_id from connection where user_id = ? and connection.deleted_at is null) and class.deleted_at is null", userID).Scan(&class).Error
	return
}
