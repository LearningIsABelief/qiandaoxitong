package service

import (
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"qiandao/pkg/app"
	"qiandao/pkg/util"
	"qiandao/store"
	"qiandao/viewmodel"
	"time"
)

// CreateCodeService
// @Description: 获取验证码图片service层
// @Author YangXuZheng 2022-06-11 14:52
func CreateCodeService() (*viewmodel.CodeInfoResponse, error) {
	id, base64, err := util.CreateCode()
	if err != nil {
		log.Errorf(err, "验证码创建失败")
		return &viewmodel.CodeInfoResponse{}, app.InternalServerError
	}
	// 获取验证码结果
	verificationCode := util.GetCodeAnswer(id)
	// 获取完验证码将验证码和答案 以key-value方式存入redis,设置过期时间2分钟
	store.RedisDB.Self.Set("login-code-"+id, verificationCode, viper.GetDuration("code.expiration")*time.Minute)

	log.Info("验证码：" + verificationCode)
	return &viewmodel.CodeInfoResponse{
		Uuid:   id,
		Base64: base64,
	}, nil
}
