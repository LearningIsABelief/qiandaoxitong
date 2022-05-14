package class

import (
	"github.com/gin-gonic/gin"
	"qiandao/pkg/app"
	"qiandao/pkg/util"
	"qiandao/service"
)

// GetAllClass 获取班级列表 controller
func GetAllClass(ctx *gin.Context) {
	var pageCondition util.PageRequest
	if err := ctx.Bind(&pageCondition); err != nil {
		app.SendResponse(ctx, app.ErrBind, nil)
		return
	}
	class, err := service.GetAllClass(pageCondition)
	if err != nil {
		app.SendResponse(ctx, err, nil)
		return
	}
	app.SendResponse(ctx, nil, class)
}
