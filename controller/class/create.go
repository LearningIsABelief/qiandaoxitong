package class

import (
	"github.com/gin-gonic/gin"
	"qiandao/pkg/app"
	"qiandao/service"
	"qiandao/viewmodel"
)

// Create 创建班级 controller
func Create(ctx *gin.Context) {
	var createClassRequest viewmodel.CreateClassRequest
	if err := ctx.Bind(&createClassRequest); err != nil {
		app.SendResponse(ctx, app.ErrBind, nil)
		return
	}
	if err := service.CreateClass(createClassRequest); err != nil {
		app.SendResponse(ctx, err, nil)
		return
	}
	app.SendResponse(ctx, nil, nil)
}
