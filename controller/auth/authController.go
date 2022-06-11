package auth

import (
	"github.com/gin-gonic/gin"
	"qiandao/pkg/app"
	"qiandao/service"
)

// Code
// @Description: 获取验证码图片
// @Author YangXuZheng 2022-06-11 10:49
func Code(ctx *gin.Context) {
	result, err := service.CreateCodeService()
	if err != nil {
		app.SendResponse(ctx, err, nil)
		return
	}
	app.SendResponse(ctx, nil, result)
}
