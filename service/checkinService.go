package service

import (
	"github.com/lexkong/log"
	"math"
	"qiandao/model"
	"qiandao/pkg/app"
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
		log.Errorf(app.ErrCheckinCreate, "用户'%v'创建签到失败", checkin.CreatorID)
		tx.RollBack()
		return nil, err
	}
	shouldCheckInClassLst, err := tx.GetClassLstByLessonID(checkin.LessonID)
	if len(shouldCheckInClassLst) == 0 || err != nil {
		log.Errorf(app.ErrCheckinClassGet, "签到'%v'的课程'%v'的签到班级无法获取", checkin.CheckinID, checkin.LessonID)
		tx.RollBack()
		return nil, err
	}
	// 根据需要签到的班级列表获取所有需要签到的学生
	for i := range shouldCheckInClassLst {
		class := shouldCheckInClassLst[i]
		shouldCheckInStuLst, err := tx.GetStuLstByClassID(class.ClassId)
		if err != nil {
			log.Errorf(app.ErrCheckinStuGet, "签到'%v'根据班级'%v'获取学生列表失败", checkin.CheckinID, class.ClassId)
			tx.RollBack()
			return nil, err
		}
		for j := range shouldCheckInStuLst {
			stu := shouldCheckInStuLst[j]
			checkinRec := model.CheckinRec{
				CheckinRecID: checkin.CheckinID + stu.UserId,
				CheckinID:    checkin.CheckinID,
				UserID:       stu.UserId,
				UserName:     stu.RealName,
				State:        2,
				EndTime:      checkin.EndTime,
			}
			err := tx.AddCheckinRec(&checkinRec)
			if err != nil {
				log.Errorf(app.ErrCheckinRecCreate, "签到'%v'添加签到记录'%v'失败", checkin.CheckinID, checkinRec)
				tx.RollBack()
				return nil, err
			}
		}
	}
	log.Infof("用户'%v'创建签到'%v'成功", checkin.CreatorID, checkin.CheckinID)
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return checkin, nil
}

// StuCheckin
// @Description: 学生签到
// @Author zhandongyang 2022-05-08 23:14:00
// @Param viewCheckin
// @Return res 1:签到成功 2:签到失败 3:重复的签到 4:非法的签到 5:签到已过期 6:签到码错误 7:超出签到范围
// @Return err
func StuCheckin(viewStuCheckin *viewmodel.Checkin) (res int, err error) {
	tx := store.GetTx()
	checkin, err := tx.GetCheckinById(viewStuCheckin.CheckinID)
	if err != nil {
		log.Errorf(app.ErrCheckinGet, "用户'%v'获取签到'%v'失败", viewStuCheckin.UserID, viewStuCheckin.CheckinID)
		return 4, app.ErrCheckinGet
	}
	// 获取签到记录
	checkinRec, err := tx.GetCheckinRecByID(viewStuCheckin.CheckinID + viewStuCheckin.UserID)
	if err != nil {
		// 检查签到合法性
		if err.Error() == "record not found" {
			log.Errorf(app.ErrCheckinRecNotExist, "签到'%v'的应签到列表中无此学生'%v'", viewStuCheckin.CheckinID, viewStuCheckin.UserID)
			return 4, app.ErrCheckinRecNotExist
		}
		log.Errorf(app.ErrCheckinRecNotExist, "获取签到记录'%v'失败", viewStuCheckin.CheckinID+viewStuCheckin.UserID)
		return 2, app.ErrCheckinRecNotExist
	}
	// 检查签到合法性
	if reflect.DeepEqual(checkinRec, model.CheckinRec{}) {
		log.Errorf(app.ErrCheckinRecNotExist, "签到'%v'的应签到列表中无此学生'%v'", viewStuCheckin.CheckinID, viewStuCheckin.UserID)
		return 4, app.ErrCheckinRecNotExist
	}
	// 检查签到码
	if checkin.CheckinCode != viewStuCheckin.CheckinCode {
		log.Errorf(app.ErrCheckinCode, "用户'%v'签到'%v'时的签到码'%v'错误", viewStuCheckin.UserID, viewStuCheckin.CheckinID, viewStuCheckin.CheckinCode)
		return 6, app.ErrCheckinCode
	}
	// 检查签到时间是否过期
	if checkin.EndTime.Before(time.Now()) {
		log.Errorf(app.ErrCheckinExpired, "用户'%v'签到'%v'时签到过期，过期时间为'%v'", viewStuCheckin.UserID, viewStuCheckin.CheckinID, checkin.EndTime)
		return 5, app.ErrCheckinExpired
	}
	// 检查是否重复签到
	if checkinRec.State == 1 {
		log.Errorf(app.ErrCheckinRepeat, "用户'%v'签到'%v'时重复签到", viewStuCheckin.UserID, viewStuCheckin.CheckinID)
		return 3, app.ErrCheckinRepeat
	}
	// 检查签到位置是否超出签到范围
	distance := func() float64 {
		lng1, _ := strconv.ParseFloat(checkin.Longitude, 64)
		lat1, _ := strconv.ParseFloat(checkin.Latitude, 64)
		lng2, _ := strconv.ParseFloat(viewStuCheckin.Longitude, 64)
		lat2, _ := strconv.ParseFloat(viewStuCheckin.Latitude, 64)
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
		log.Errorf(app.ErrCheckinOutOfRng, "用户'%v'签到'%v'时超出签到范围，用户签到范围为'%v','%v'，签到范围为'%v','%v'",
			viewStuCheckin.UserID, viewStuCheckin.CheckinID, viewStuCheckin.Longitude, viewStuCheckin.Latitude,
			checkin.Longitude, checkin.Latitude)
		return 7, app.ErrCheckinOutOfRng
	}
	err = tx.UpdateCheckinRecStateByID(viewStuCheckin.CheckinID+viewStuCheckin.UserID, 1)
	if err != nil {
		log.Errorf(app.ErrCheckinUpdateState, "更新用户'%v'的签到状态失败,签到为'%v'", viewStuCheckin.UserID, viewStuCheckin.CheckinID)
		return 2, app.ErrCheckinUpdateState
	}
	log.Infof("用户'%v'签到成功，签到为'%v'", viewStuCheckin.UserID, viewStuCheckin.CheckinID)
	return 1, nil
}

// GetCheckinDetails
// @Description: 获取签到详情
// @Author zhandongyang 2022-05-09 15:38:01
// @Param checkinID
// @Return checkinDetails
// @Return err
func GetCheckinDetails(checkinID string) (checkinDetails *viewmodel.CheckinDetailsResponse, err error) {
	tx := store.GetTx()
	// 获取签到记录列表
	checkinRecLst, err := tx.GetCheckinRecLstByCheckinID(checkinID)
	if err != nil {
		log.Errorf(app.ErrCheckinRecGet, "签到'%v'的签到记录获取失败", checkinID)
		return nil, app.ErrCheckinRecGet
	}
	// 获取所有需要签到、已经签到、没有签到的 数据响应列表
	totalStuList := make([]viewmodel.List, len(checkinRecLst))
	var checkedInStuList, noCheckedInStuList []viewmodel.List
	for i := range checkinRecLst {
		checkinRec := checkinRecLst[i]
		state := func() string {
			if checkinRec.State == 1 {
				return "已签到"
			}
			return "未签到"
		}()
		class, err := tx.GetClassByUserID(checkinRec.UserID)
		if err != nil {
			log.Errorf(app.ErrCheckinClassGet, "获取用户'%v'的班级失败", checkinRec.UserID)
			return nil, app.ErrCheckinClassGet
		}
		totalStuList[i] = viewmodel.List{
			ClassName: class.ClassName,
			UserName:  checkinRec.UserName,
			State:     state,
		}
		if checkinRec.State == 1 {
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
	log.Infof("获取签到'%v'成功", checkinID)
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
	checkinList, err := tx.GetCheckinByCreator(creatorID)
	if err != nil {
		log.Errorf(app.ErrCheckinGet, "获取用户'%v'创建的签到失败", creatorID)
		return nil, app.ErrCheckinGet
	}
	listResponse = make([]viewmodel.ListResponse, len(checkinList))
	for i := range checkinList {
		lesson, err := tx.GetLessonByID(checkinList[i].LessonID)
		if err != nil {
			log.Errorf(app.ErrCheckinLessonGet, "获取用户'%v'创建的课程'%v'失败", creatorID, lesson.LessonID)
			return nil, app.ErrCheckinLessonGet
		}
		endTime := checkinList[i].EndTime
		checkinState := 2
		if time.Now().Before(endTime) {
			checkinState = 1
		}
		listResponse[i] = viewmodel.ListResponse{
			CheckinID:    checkinList[i].CheckinID,
			LessonName:   lesson.LessonName,
			BeginTime:    checkinList[i].BeginTime.Format("2006/01/02 15:04"),
			CheckinState: checkinState,
		}
	}
	log.Infof("获取用户'%v'已创建的签到成功", creatorID)
	return
}

// GetCheckinRecList
// @Description: 获取签到记录
// @Author zhandongyang 2022-05-09 21:16:44
// @Param userID
// @Return shouldCheckInList
// @Return err
func GetCheckinRecList(userID string) (shouldCheckInList []viewmodel.ListResponse, err error) {
	tx := store.GetTx()
	checkinRecLst, err := tx.GetCheckinRecByUserID(userID)
	if err != nil {
		log.Errorf(app.ErrCheckinRecGet, "获取用户'%v'的签到记录失败", userID)
		return nil, app.ErrCheckinRecGet
	}
	shouldCheckInList = make([]viewmodel.ListResponse, len(checkinRecLst))
	for i := range checkinRecLst {
		checkinRec := checkinRecLst[i]
		checkin, err := tx.GetCheckinById(checkinRec.CheckinID)
		if err != nil {
			log.Errorf(app.ErrCheckinGet, "获取签到'%v'失败", checkinRec.CheckinID)
			return nil, app.ErrCheckinGet
		}
		endTime := checkin.EndTime
		checkinState := 2
		if time.Now().Before(endTime) {
			checkinState = 1
		}
		lesson, err := tx.GetLessonByID(checkin.LessonID)
		if err != nil {
			log.Errorf(app.ErrCheckinLessonGet, "获取课程'%v'失败", checkin.LessonID)
			return nil, app.ErrCheckinLessonGet
		}
		shouldCheckInList[i] = viewmodel.ListResponse{
			CheckinID:    checkinRec.CheckinID,
			LessonName:   lesson.LessonName,
			BeginTime:    checkin.BeginTime.Format("2006/01/02 15:04"),
			State:        checkinRec.State,
			CheckinState: checkinState,
		}
	}
	log.Infof("获取用户'%v'的签到记录成功", userID)
	return
}
