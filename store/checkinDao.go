package store

import (
	"qiandao/model"
)

// CreateCheckin
// @Description: 创建一个签到
// @Author zhandongyang 2022-05-08 22:44:50 ${time}
// @Param checkin
// @Return err
func CreateCheckin(checkin *model.Checkin) (err error) {
	err = DB.Self.Create(checkin).Error
	return
}

// GetCheckinById
// @Description: 根据签到id获取一个签到
// @Author zhandongyang 2022-05-08 22:46:45 ${time}
// @Param checkinID
// @Return checkin
// @Return err
func GetCheckinById(checkinID string) (checkin model.Checkin, err error) {
	err = DB.Self.Where("checkin_id = ?", checkinID).First(&checkin).Error
	return
}

// GetCheckinByCreator
// @Description: 根据用户id获取一个签到
// @Author zhandongyang 2022-05-09 16:58:08
// @Param creator
// @Return checkinList
// @Return err
func GetCheckinByCreator(creatorID string) (checkinList []model.Checkin, err error) {
	err = DB.Self.Find(&checkinList).Where("creator_id = ? and deleted_at is null", creatorID).Error
	return
}

// GetShouldCheckInClass
// @Description: 根据课程id获取需要签到的班级列表
// @Author zhandongyang 2022-05-08 23:12:01
// @Param lessonID
// @Return classList
// @Return err
func GetShouldCheckInClass(lessonID string) (classList []model.Class, err error) {
	err = DB.Self.Raw("select * from class where class_id IN (select class_id from class_lesson where lesson_id = ?)", lessonID).Scan(&classList).Error
	return
}

// GetShouldCheckInStu
// @Description: 获取需要签到的学生列表
// @Author zhandongyang 2022-05-08 23:02:59
// @Param classID
// @Param userID
// @Return stu
// @Return err
func GetShouldCheckInStu(classID string) (stuList []model.User, err error) {
	err = DB.Self.Raw("select * from user where class_id = ? and deleted_at is null", classID).Scan(&stuList).Error
	return
}

// GetLessonById
// @Description: 根据课程id获取课程
// @Author zhandongyang 2022-05-09 17:13:17
// @Param lessonID
// @Return lesson
// @Return err
func GetLessonById(lessonID string) (lesson model.Lesson, err error) {
	err = DB.Self.First(&lesson).Where("lesson_id = ?", lessonID).Error
	return
}

// AddCheckedIn
// @Description: 添加学生签到信息
// @Author zhandongyang 2022-05-08 22:58:34 ${time}
// @Param stuCheckin
// @Return err
func AddCheckedIn(stuCheckin *model.CheckedIn) (err error) {
	err = DB.Self.Create(stuCheckin).Error
	return
}

// UpdateCheckedIn
// @Description: 更新学生签到信息
// @Author zhandongyang 2022-05-09 17:32:27
// @Param stuCheckin
// @Return err
func UpdateCheckedIn(stuCheckedIn *model.CheckedIn) (err error) {
	err = DB.Self.Model(stuCheckedIn).Update("state", stuCheckedIn.State).Error
	return
}

// GetCheckedIn
// @Description: 获取单个已签到信息
// @Author zhandongyang 2022-05-08 22:59:20
// @Param checkedID
// @Return checkedIn
// @Return err
func GetCheckedIn(checkedID string) (checkedIn model.CheckedIn, err error) {
	err = DB.Self.First(&checkedIn).Where("id = ?", checkedID).Error
	return
}

// GetAllCheckedInByCheckinID
// @Description: 根据签到id获取需要签到的学生列表
// @Author zhandongyang 2022-05-10 16:11:54
// @Param checkinID
// @Return checkedInList
// @Return err
func GetAllCheckedInByCheckinID(checkinID string) (checkedInList []model.CheckedIn, err error) {
	err = DB.Self.Where("checkin_id = ?", checkinID).Find(&checkedInList).Error
	return
}

// GetAllCheckedInByUserID
// @Description: 获取某个学生需要签到的列表
// @Author zhandongyang 2022-05-09 15:46:01
// @Param checkedID
// @Return checkedIn
// @Return err
func GetAllCheckedInByUserID(userID string) (checkedInList []model.CheckedIn, err error) {
	err = DB.Self.Where("user_id = ?", userID).Find(&checkedInList).Error
	return
}

func GetClassByUserID(userID string) (class model.Class, err error) {
	err = DB.Self.Raw("select * from class where class_id = (select class_id from connection where user_id = ?)", userID).Scan(&class).Error
	return
}
