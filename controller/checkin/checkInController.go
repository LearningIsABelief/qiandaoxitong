package checkin

import (
	"github.com/gin-gonic/gin"
	"qiandao/pkg/app"
	"qiandao/service"
	"qiandao/viewmodel"
)

func CreateCheckin(ctx *gin.Context) {
	checkin := new(viewmodel.CreateCheckin)
	err := ctx.ShouldBindJSON(checkin)
	if err != nil {
		app.SendResponse(ctx, app.ErrBind, nil)
		return
	}
	if checkin.CreatorID == "" || checkin.CheckinCode == "" || checkin.Duration < 0 || checkin.LessonID == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	createCheckin, err := service.CreateCheckin(checkin)
	if err != nil {
		app.SendResponse(ctx, app.InternalServerError, nil)
		return
	}
	app.SendResponse(ctx, app.OK, createCheckin.CheckinID)
}

func CheckIn(ctx *gin.Context) {
	checkIn := new(viewmodel.Checkin)
	err := ctx.ShouldBindJSON(checkIn)
	if err != nil {
		app.SendResponse(ctx, app.ErrBind, nil)
		return
	}
	if checkIn.CheckinID == "" || checkIn.UserID == "" || checkIn.UserName == "" || checkIn.CheckinCode == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	checkinResponse, err := service.StuCheckin(checkIn)
	if err != nil {
		app.SendResponse(ctx, app.InternalServerError, checkinResponse)
		return
	}
	app.SendResponse(ctx, app.OK, checkinResponse)
}

func GetCheckinDetails(ctx *gin.Context) {
	checkinID := ctx.Query("checkin_id")
	if checkinID == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	checkinDetailsResponse, err := service.GetCheckInDetails(checkinID)
	if err != nil {
		app.SendResponse(ctx, app.InternalServerError, nil)
		return
	}
	app.SendResponse(ctx, app.OK, checkinDetailsResponse)
}

func GetCreatedCheckinList(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	response, err := service.GetCreatedCheckInList(userID)
	if err != nil {
		app.SendResponse(ctx, app.InternalServerError, nil)
		return
	}
	app.SendResponse(ctx, app.OK, response)
}

func GetCheckinRecList(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		app.SendResponse(ctx, app.ErrParamNull, nil)
		return
	}
	response, err := service.GetShouldCheckInList(userID)
	if err != nil {
		app.SendResponse(ctx, app.InternalServerError, nil)
		return
	}
	app.SendResponse(ctx, app.OK, response)
}
