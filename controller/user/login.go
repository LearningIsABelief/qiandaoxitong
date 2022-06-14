package user

import (
	"github.com/gin-gonic/gin"
	"qiandao/pkg/app"
	"qiandao/service"
	"qiandao/viewmodel"
)

func Login(ctx *gin.Context) {
	var loginRequest viewmodel.LoginRequest
	if err := ctx.Bind(&loginRequest); err != nil {
		app.SendResponse(ctx, app.ErrBind, nil)
		return
	}
	login, err := service.Login(loginRequest, ctx)
	if err != nil {
		app.SendResponse(ctx, err, nil)
		return
	}
	app.SendResponse(ctx, nil, login)
}
