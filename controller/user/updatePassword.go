package user

import (
	"github.com/gin-gonic/gin"
	"qiandao/pkg/app"
	"qiandao/service"
	"qiandao/viewmodel"
)

func UpdatePassword(ctx *gin.Context) {
	var updatePasswordRequest viewmodel.UpdatePasswordRequest
	if err := ctx.Bind(&updatePasswordRequest); err != nil {
		app.SendResponse(ctx, app.ErrBind, nil)
		return
	}
	if err := service.UpdatePassword(updatePasswordRequest); err != nil {
		app.SendResponse(ctx, err, nil)
		return
	}
	app.SendResponse(ctx, nil, nil)
}
