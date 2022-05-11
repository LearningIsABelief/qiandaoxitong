package service

import (
	"fmt"
	"math"
	"qiandao/model"
	"qiandao/pkg/util"
	"qiandao/store"
	"qiandao/viewmodel"
	"reflect"
	"strconv"
	"time"
)

// CreateCheckin 创建签到
func CreateCheckin(viewCreatCheckin *viewmodel.CreateCheckin) (*model.Checkin, error) {
	tx := store.GetTx()
	tx.Begin()
	checkin := &model.Checkin{
		CheckinID:   util.GetUUID(),
		CreatorID:   viewCreatCheckin.CreatorID,
		LessonID:    viewCreatCheckin.LessonID,
		BeginTime:   time.Now(),
		EndTime:     time.Now().UTC().Add(time.Duration(viewCreatCheckin.Duration) * time.Minute),
		CheckinCode: viewCreatCheckin.CheckinCode,
		Longitude:   viewCreatCheckin.Longitude,
		Latitude:    viewCreatCheckin.Latitude,
	}
	err := tx.CreateCheckin(checkin)
	if err != nil {
		tx.RollBack()
		return nil, err
	}
	shouldCheckInClass, err := tx.GetShouldCheckInClass(checkin.LessonID)
	if len(shouldCheckInClass) == 0 || err != nil {
		tx.RollBack()
		return nil, err
	}
	for i := range shouldCheckInClass {
		class := shouldCheckInClass[i]
		shouldCheckInStu, err := tx.GetShouldCheckInStu(class.ClassId)
		if err != nil {
			tx.RollBack()
			return nil, err
		}
		for j := range shouldCheckInStu {
			stu := shouldCheckInStu[j]
			stuCheckedIn := model.CheckedIn{
				ID:        checkin.CheckinID + stu.UserId,
				CheckinID: checkin.CheckinID,
				UserID:    stu.UserId,
				UserName:  stu.RealName,
				State:     2,
				EndTime:   checkin.EndTime,
			}
			err := tx.AddCheckedIn(&stuCheckedIn)
			if err != nil {
				tx.RollBack()
				return nil, err
			}
		}
	}
	tx.Commit()
	return checkin, nil
}

// StuCheckin
// @Description: 学生签到
// @Author zhandongyang 2022-05-08 23:14:00
// @Param viewCheckin
// @Return res 1:签到成功 2:签到失败 3:重复的签到 4:非法的签到 5:签到已过期 6:签到码错误 7:超出签到范围
// @Return err
func StuCheckin(viewCheckin *viewmodel.Checkin) (res int, err error) {
	tx := store.GetTx()
	checkin, err := tx.GetCheckinById(viewCheckin.CheckinID)
	if err != nil {
		fmt.Println(err)
		return 4, err
	}
	// 获取正确的已签到信息
	rightCheckedIn, err := tx.GetCheckedIn(viewCheckin.CheckinID + viewCheckin.UserID)
	if err != nil {
		fmt.Println(err)
		// 检查签到合法性
		if err.Error() == "record not found" {
			return 4, err
		}
		return 2, err
	}
	// 检查签到合法性
	if reflect.DeepEqual(rightCheckedIn, model.CheckedIn{}) {
		return 4, err
	}
	// 检查签到码
	if checkin.CheckinCode != viewCheckin.CheckinCode {
		return 6, nil
	}
	// 检查签到时间是否过期
	if checkin.EndTime.Before(time.Now()) {
		return 5, nil
	}
	// 检查是否重复签到
	if rightCheckedIn.State == 2 {
		return 3, err
	}
	// 检查签到位置是否超出签到范围
	distance := func() float64 {
		lng1, _ := strconv.ParseFloat(checkin.Longitude, 64)
		lat1, _ := strconv.ParseFloat(checkin.Latitude, 64)
		lng2, _ := strconv.ParseFloat(viewCheckin.Longitude, 64)
		lat2, _ := strconv.ParseFloat(viewCheckin.Latitude, 64)
		const PI float64 = 3.141592653589793
		radlat1 := PI * lat1 / 180
		radlat2 := PI * lat2 / 180

		theta := lng1 - lng2
		radtheta := PI * theta / 180

		dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

		if dist > 1 {
			dist = 1
		}

		dist = math.Acos(dist)
		dist = dist * 180 / PI
		dist = dist * 60 * 1.1515
		dist = dist * 1.609344
		return dist
	}()
	if distance > 100 {
		return 7, nil
	}
	// 创建已签到记录
	checkedIn := &model.CheckedIn{
		ID:        viewCheckin.CheckinID + viewCheckin.UserID,
		CheckinID: viewCheckin.CheckinID,
		UserID:    viewCheckin.UserID,
		UserName:  viewCheckin.UserName,
		State:     1,
	}
	err = tx.UpdateCheckedIn(checkedIn)
	if err != nil {
		return 2, err
	}
	return 1, nil
}

// GetCheckInDetails
// @Description: 获取签到详情
// @Author zhandongyang 2022-05-09 15:38:01
// @Param checkinID
// @Return checkinDetails
// @Return err
func GetCheckInDetails(checkinID string) (checkinDetails *viewmodel.CheckinDetailsResponse, err error) {
	tx := store.GetTx()
	if err != nil {
		return nil, err
	}
	// 获取需要签到的学生列表
	shouldCheckInStuList, err := tx.GetAllCheckedInByCheckinID(checkinID)
	if err != nil {
		return nil, err
	}
	// 获取所有需要签到、已经签到、没有签到的 数据响应列表
	totalStuList := make([]viewmodel.List, len(shouldCheckInStuList))
	var checkedInStuList, noCheckedInStuList []viewmodel.List
	for i := range shouldCheckInStuList {
		checkedIn := shouldCheckInStuList[i]
		state := func() string {
			if checkedIn.State == 1 {
				return "已签到"
			}
			return "未签到"
		}()
		class, err := tx.GetClassByUserID(checkedIn.UserID)
		if err != nil {
			return nil, err
		}
		totalStuList[i] = viewmodel.List{
			ClassName: class.ClassName,
			UserName:  checkedIn.UserName,
			State:     state,
		}
		if checkedIn.State == 1 {
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
// @Description: 获取已创建的签到
// @Author zhandongyang 2022-05-09 16:59:29
// @Param userID
// @Return createdCheckInList
// @Return err
func GetCreatedCheckInList(creatorID string) (listResponse []viewmodel.ListResponse, err error) {
	tx := store.GetTx()
	createdCheckInList, err := tx.GetCheckinByCreator(creatorID)
	if err != nil {
		return nil, err
	}
	listResponse = make([]viewmodel.ListResponse, len(createdCheckInList))
	for i := range createdCheckInList {
		lesson, err := tx.GetLessonById(createdCheckInList[i].LessonID)
		if err != nil {
			return nil, err
		}
		endTime := createdCheckInList[i].EndTime
		checkinState := 2
		if time.Now().Before(endTime) {
			checkinState = 1
		}
		listResponse[i] = viewmodel.ListResponse{
			CheckinID:    createdCheckInList[i].CheckinID,
			LessonName:   lesson.LessonName,
			BeginTime:    createdCheckInList[i].BeginTime,
			CheckinState: checkinState,
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
	tx := store.GetTx()
	checkedInList, err := tx.GetAllCheckedInByUserID(userID)
	if err != nil {
		return nil, err
	}
	shouldCheckInList = make([]viewmodel.ListResponse, len(checkedInList))
	for i := range checkedInList {
		checkedIn := checkedInList[i]
		checkin, err := tx.GetCheckinById(checkedIn.CheckinID)
		endTime := checkin.EndTime
		checkinState := 2
		if time.Now().Before(endTime) {
			checkinState = 1
		}
		if err != nil {
			return nil, err
		}
		lesson, err := tx.GetLessonById(checkin.LessonID)
		if err != nil {
			return nil, err
		}
		shouldCheckInList[i] = viewmodel.ListResponse{
			CheckinID:    checkedIn.CheckinID,
			LessonName:   lesson.LessonName,
			BeginTime:    checkin.BeginTime,
			State:        checkedIn.State,
			CheckinState: checkinState,
		}
	}
	return
}
