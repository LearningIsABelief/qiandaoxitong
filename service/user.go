package service

import (
	"github.com/lexkong/log"
	"qiandao/model"
	"qiandao/pkg/app"
	"qiandao/pkg/util"
	"qiandao/store"
	"qiandao/viewmodel"
	"strings"
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
	if phone := store.IsExistUser("phone", registerRequest.Phone, &model.User{}); phone {
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
	if isExits := store.IsExistUser("email", updateUserEmail.Email, &model.User{}); isExits {
		log.Errorf(app.ErrEmailExist, "邮箱已存在")
		return app.ErrEmailExist
	}
	if err := store.UpdateEmailMapper(updateUserEmail); err != nil {
		return err
	}
	return nil
}

// UpdateNickName 修改昵称 service
func UpdateNickName(updateUserInfo viewmodel.UpdateNickNameRequest) error {
	if isExits := store.IsExistUser("nick_name", updateUserInfo.NickName, &model.User{}); isExits {
		log.Errorf(app.ErrNickNameExist, "昵称已存在")
		return app.ErrNickNameExist
	}
	if err := store.UpdateNickNameMapper(updateUserInfo); err != nil {
		return err
	}
	return nil
}

// UpdatePassword 修改密码 service
func UpdatePassword(updatePasswordRequest viewmodel.UpdatePasswordRequest) error {
	// 根据用户id找到该用户数据库中的密码
	password, err := store.GetPasswordById(updatePasswordRequest.UserId)
	if err != nil {
		return err
	}
	// 将数据库中的密码和用户的密码进行比对
	if err := util.Decrypt(password, updatePasswordRequest.OldPassword); err != nil {
		log.Errorf(app.ErrPassword, "用户输入的原来的密码，和数据库中不一致")
		return app.ErrPassword
	}
	// 判断新密码和确认输入的新密码是否相等
	if updatePasswordRequest.NewPassword != updatePasswordRequest.NewConfirmPassword {
		log.Errorf(app.ErrOldNewInconsistent, "请确保两次输入的密码一样")
		return app.ErrOldNewInconsistent
	}
	// 密码进行加密
	psw, err := util.Encrypt(updatePasswordRequest.NewConfirmPassword)
	if err != nil {
		log.Errorf(app.ErrEncrypt, "密码加密出错")
		return app.ErrEncrypt
	}
	if err := store.UpdatePasswordMapper("user_id", updatePasswordRequest.UserId, psw); err != nil {
		return err
	}
	return nil
}

// ForgetPassword 忘记密码 service
func ForgetPassword(forgetPasswordRequest viewmodel.ForgetPasswordRequest) error {
	// 判断手机号是否在数据库中存在
	if phone := store.IsExistUser("phone", forgetPasswordRequest.Phone, &model.User{}); !phone {
		log.Errorf(app.ErrPhoneDoesNotExist, "手机号不存在")
		return app.ErrPhoneDoesNotExist
	}
	// 判断邮箱是否是当前账号(手机号)下的
	email, err2 := store.GetEmailByPhone(forgetPasswordRequest.Phone)
	if err2 != nil {
		return err2
	}
	if strings.Compare(email, forgetPasswordRequest.Email) == 1 {
		log.Errorf(app.ErrPhoneBinEmail, "请输入手机号绑定的正确邮箱")
		return app.ErrPhoneBinEmail
	}
	// 密码进行加密
	psw, err := util.Encrypt("123456")
	if err != nil {
		log.Errorf(app.ErrEncrypt, "密码加密出错")
		return app.ErrEncrypt
	}
	if err := store.UpdatePasswordMapper("phone", forgetPasswordRequest.Phone, psw); err != nil {
		return err
	}
	return nil
}

func Login(loginRequest viewmodel.LoginRequest) (viewmodel.UserInfo, error) {
	// 根据手机号拿到用户信息
	store.GetUserInfoByPhone(loginRequest.Phone)
}
