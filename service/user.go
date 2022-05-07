package service

import (
	"github.com/lexkong/log"
	"qiandao/model"
	"qiandao/pkg/app"
	"qiandao/pkg/util"
	"qiandao/store"
	"qiandao/viewmodel"
)

// CreateUser 创建用户service
func CreateUser(registerRequest *viewmodel.RegisterRequest) error {
	// 加密密码
	password, err := util.Encrypt(registerRequest.Password)
	if err != nil {
		log.Errorf(app.ErrEncrypt, "err：%v", app.ErrEncrypt)
		return app.ErrEncrypt
	}
	// 判断注册的手机号是否在数据库中存在
	if phone := store.GetByPhone(registerRequest.Phone); phone {
		log.Errorf(app.ErrPhoneExist, "err：%v", app.ErrPhoneExist)
		return app.ErrPhoneExist
	}
	user := model.User{
		UserId:   util.GetUUID(),
		Phone:    registerRequest.Phone,
		Password: password,
		Email:    registerRequest.Email,
		Role:     registerRequest.Role,
		ClassId:  registerRequest.ClassId,
	}
	// 创建用户并维护班级用户中间表
	if err := store.CreateUser(&user, registerRequest.ClassId); err != nil {
		return err
	}
	return nil
}

// UpdateUser 修改用户信息 service
func UpdateUser(updateUserInfo viewmodel.UpdateUserInfoRequest) error {
	if err := store.UpdateUser(updateUserInfo); err != nil {
		return err
	}
	return nil
}

// UpdateEmail 修改邮箱 service
func UpdateEmail(updateUserEmail viewmodel.UpdateEmailRequest) error {
	if err := store.UpdateEmail(updateUserEmail); err != nil {
		return err
	}
	return nil
}

// UpdateNickName 修改昵称 service
func UpdateNickName(updateUserInfo viewmodel.UpdateNickNameRequest) error {
	if err := store.UpdateNickName(updateUserInfo); err != nil {
		return err
	}
	return nil
}
