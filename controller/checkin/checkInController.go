package checkin

import (
	"github.com/gin-gonic/gin"
	"qiandao/pkg/app"
	"qiandao/service"
	"qiandao/viewmodel"
)

// CreateCheckin
// @Description: 创建签到
// @Author zhandongyang 2022-05-11 22:02:22
// @Param ctx
func CreateCheckin(ctx *gin.Context) {
	viewCheckin := new(viewmodel.CreateCheckin)
	err := ctx.ShouldBindJSON(viewCheckin)
	if err != nil {
		app.SendResponse(ctx, app.ErrBind, nil)
		return
	}
	if viewCheckin.CreatorID == "" || viewCheckin.CheckinCode == "" || viewCheckin.Duration <= 0 ||
		viewCheckin.LessonID == "" || viewCheckin.Longitude == "" || viewCheckin.Latitude == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	createCheckin, err := service.CreateCheckin(viewCheckin)
	if err != nil {
		app.SendResponse(ctx, err, nil)
		return
	}
	app.SendResponse(ctx, app.OK, createCheckin.CheckinID)
}

// StuCheckIn
// @Description: 学生签到
// @Author zhandongyang 2022-05-11 22:02:30
// @Param ctx
func StuCheckIn(ctx *gin.Context) {
	checkIn := new(viewmodel.Checkin)
	err := ctx.ShouldBindJSON(checkIn)
	if err != nil {
		app.SendResponse(ctx, app.ErrBind, nil)
		return
	}
	if checkIn.CheckinID == "" || checkIn.UserID == "" || checkIn.CheckinCode == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	checkinResponse, err := service.StuCheckin(checkIn)
	if err != nil {
		app.SendResponse(ctx, err, checkinResponse)
		return
	}
	app.SendResponse(ctx, app.OK, checkinResponse)
}

// GetCheckinDetails
// @Description: 获取签到详情
// @Author zhandongyang 2022-05-11 22:02:37
// @Param ctx
func GetCheckinDetails(ctx *gin.Context) {
	checkinID := ctx.Query("checkin_id")
	if checkinID == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	checkinDetailsResponse, err := service.GetCheckinDetails(checkinID)
	if err != nil {
		app.SendResponse(ctx, err, nil)
		return
	}
	app.SendResponse(ctx, app.OK, checkinDetailsResponse)
}

// GetCreatedCheckinLst
// @Description: 获取已创建签到列表
// @Author zhandongyang 2022-05-11 22:02:47
// @Param ctx
func GetCreatedCheckinLst(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	response, err := service.GetCreatedCheckInList(userID)
	if err != nil {
		app.SendResponse(ctx, err, nil)
		return
	}
	app.SendResponse(ctx, app.OK, response)
}

// GetCheckinRecLst
// @Description: 获取签到记录列表
// @Author zhandongyang 2022-05-11 22:02:54
// @Param ctx
func GetCheckinRecLst(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	response, err := service.GetCheckinRecList(userID)
	if err != nil {
		app.SendResponse(ctx, err, nil)
		return
	}
	app.SendResponse(ctx, app.OK, response)
}
