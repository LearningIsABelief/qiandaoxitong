package service

import (
	"github.com/lexkong/log"
	"qiandao/model"
	"qiandao/pkg/app"
	"qiandao/pkg/util"
	"qiandao/store"
	"qiandao/viewmodel"
)

// CreateUser 创建用户 service
func CreateUser(registerRequest *viewmodel.RegisterRequest) error {
	// 加密密码
	password, err := util.Encrypt(registerRequest.Password)
	if err != nil {
		log.Errorf(app.ErrEncrypt, "密码加密出错")
		return app.ErrEncrypt
	}
	// 判断注册的手机号是否在数据库中存在
	if phone := store.GetByPhoneMapper(registerRequest.Phone); phone {
		log.Errorf(app.ErrPhoneExist, "手机号已被注册")
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
	if err := store.CreateUserMapper(&user, registerRequest.ClassId); err != nil {
		return err
	}
	return nil
}

// UpdateUser 修改用户信息 service
func UpdateUser(updateUserInfo viewmodel.UpdateUserInfoRequest) error {
	if err := store.UpdateUserMapper(updateUserInfo); err != nil {
		return err
	}
	return nil
}

// UpdateEmail 修改邮箱 service
func UpdateEmail(updateUserEmail viewmodel.UpdateEmailRequest) error {
	if err := store.UpdateEmailMapper(updateUserEmail); err != nil {
		return err
	}
	return nil
}

// UpdateNickName 修改昵称 service
func UpdateNickName(updateUserInfo viewmodel.UpdateNickNameRequest) error {
	if err := store.UpdateNickNameMapper(updateUserInfo); err != nil {
		return err
	}
	return nil
}
