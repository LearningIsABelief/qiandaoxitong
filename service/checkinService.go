package service

import (
	"fmt"
	"qiandao/model"
	"qiandao/pkg/util"
	"qiandao/store"
	"qiandao/viewmodel"
	"reflect"
	"time"
)

// CreateCheckin 创建签到
func CreateCheckin(viewCreatCheckin *viewmodel.CreateCheckin) (*model.Checkin, error) {
	checkin := &model.Checkin{
		CheckinID:   util.GetUUID(),
		CreatorID:   viewCreatCheckin.CreatorID,
		LessonID:    viewCreatCheckin.LessonID,
		BeginTime:   time.Now(),
		EndTime:     time.Now().UTC().Add(time.Duration(viewCreatCheckin.Duration) * time.Minute),
		CheckinCode: viewCreatCheckin.CheckinCode,
	}
	err := store.CreateCheckin(checkin)
	shouldCheckInClass, err := store.GetShouldCheckInClass(checkin.LessonID)
	if err != nil {
		return nil, err
	}
	for i := range shouldCheckInClass {
		class := shouldCheckInClass[i]
		shouldCheckInStu, err := store.GetShouldCheckInStu(class.ClassId)
		if err != nil {
			return nil, err
		}
		for j := range shouldCheckInStu {
			stu := shouldCheckInStu[j]
			stuCheckedIn := model.CheckedIn{
				ID:        checkin.CheckinID + stu.UserId,
				CheckinID: checkin.CheckinID,
				UserID:    stu.UserId,
				UserName:  stu.RealName,
				State:     false,
			}
			err := store.AddCheckedIn(&stuCheckedIn)
			if err != nil {
				return nil, err
			}
		}
	}
	return checkin, nil
}

// StuCheckin
// @Description: 学生签到
// @Author zhandongyang 2022-05-08 23:14:00
// @Param viewCheckin
// @Return res 1:签到成功 2:签到失败 3:重复的签到 4:非法的签到 5:签到已过期
// @Return err
func StuCheckin(viewCheckin *viewmodel.Checkin) (res int, err error) {
	fmt.Println(0)
	checkin, err := store.GetCheckinById(viewCheckin.CheckinID)
	if err != nil {
		return 2, err
	}
	fmt.Println(1)
	// 获取正确的已签到信息
	rightCheckedIn, err := store.GetCheckedIn(viewCheckin.CheckinID + viewCheckin.UserID)
	if err != nil {
		return 2, err
	}
	fmt.Println(2)
	// 检查签到合法性
	if reflect.DeepEqual(rightCheckedIn, model.CheckedIn{}) {
		return 4, err
	}
	fmt.Println(3)
	// 检查签到码
	if checkin.CheckinCode != viewCheckin.CheckinCode {
		return 5, nil
	}
	fmt.Println(4)
	// 检查签到时间是否过期
	if checkin.EndTime.Before(time.Now()) {
		return 2, nil
	}
	fmt.Println(5)
	// 检查是否重复签到
	if rightCheckedIn.State == true {
		return 3, err
	}
	fmt.Println(6)
	// 创建已签到记录
	checkedIn := &model.CheckedIn{
		ID:        viewCheckin.CheckinID + viewCheckin.UserID,
		CheckinID: viewCheckin.CheckinID,
		UserID:    viewCheckin.UserID,
		UserName:  viewCheckin.UserName,
		State:     true,
	}
	err = store.UpdateCheckedIn(checkedIn)
	if err != nil {
		return 2, err
	}
	fmt.Println(7)
	return 1, nil
}

// GetCheckInDetails
// @Description: 获取签到详情
// @Author zhandongyang 2022-05-09 15:38:01
// @Param checkinID
// @Return checkinDetails
// @Return err
func GetCheckInDetails(checkinID string) (checkinDetails *viewmodel.CheckinDetailsResponse, err error) {
	if err != nil {
		return nil, err
	}
	// 获取需要签到的学生列表
	shouldCheckInStuList, err := store.GetAllCheckedInByCheckinID(checkinID)
	if err != nil {
		return nil, err
	}
	// 获取所有需要签到、已经签到、没有签到的 数据响应列表
	totalStuList := make([]viewmodel.List, len(shouldCheckInStuList))
	var checkedInStuList, noCheckedInStuList []viewmodel.List
	for i := range shouldCheckInStuList {
		checkedIn := shouldCheckInStuList[i]
		class, err := store.GetClassByUserID(checkedIn.UserID)
		if err != nil {
			return nil, err
		}
		totalStuList[i] = viewmodel.List{
			ClassName: class.ClassName,
			UserName:  checkedIn.UserName,
		}
		if checkedIn.State == true {
			checkedInStuList = append(checkedInStuList, totalStuList[i])
		} else {
			noCheckedInStuList = append(noCheckedInStuList, totalStuList[i])
		}
	}
	checkinDetails = &viewmodel.CheckinDetailsResponse{
		TotalList:     totalStuList,
		CheckedInList: checkedInStuList,
		NotCheckList:  noCheckedInStuList,
	}
	return
}

// GetCreatedCheckInList
// @Description: 获取创建的签到
// @Author zhandongyang 2022-05-09 16:59:29
// @Param userID
// @Return createdCheckInList
// @Return err
func GetCreatedCheckInList(creatorID string) (listResponse []viewmodel.ListResponse, err error) {
	createdCheckInList, err := store.GetCheckinByCreator(creatorID)
	if err != nil {
		return nil, err
	}
	listResponse = make([]viewmodel.ListResponse, len(createdCheckInList))
	for i := range createdCheckInList {
		lesson, err := store.GetLessonById(createdCheckInList[i].LessonID)
		if err != nil {
			return nil, err
		}
		endTime := createdCheckInList[i].EndTime
		state := false
		if time.Now().Before(endTime) {
			state = true
		}
		listResponse[i] = viewmodel.ListResponse{
			CheckinID:  createdCheckInList[i].CheckinID,
			LessonName: lesson.LessonName,
			BeginTime:  createdCheckInList[i].BeginTime,
			State:      state,
		}
	}
	return
}

// GetShouldCheckInList
// @Description: 获取需要签到的列表
// @Author zhandongyang 2022-05-09 21:16:44
// @Param userID
// @Return shouldCheckInList
// @Return err
func GetShouldCheckInList(userID string) (shouldCheckInList []viewmodel.ListResponse, err error) {
	checkedInList, err := store.GetAllCheckedInByUserID(userID)
	if err != nil {
		return nil, err
	}
	shouldCheckInList = make([]viewmodel.ListResponse, len(checkedInList))
	for i := range checkedInList {
		checkedIn := checkedInList[i]
		checkin, err := store.GetCheckinById(checkedIn.CheckinID)
		if err != nil {
			return nil, err
		}
		lesson, err := store.GetLessonById(checkin.LessonID)
		if err != nil {
			return nil, err
		}
		shouldCheckInList[i] = viewmodel.ListResponse{
			CheckinID:  checkedIn.CheckinID,
			LessonName: lesson.LessonName,
			BeginTime:  checkin.BeginTime,
			State:      checkedIn.State,
		}
	}
	return
}
