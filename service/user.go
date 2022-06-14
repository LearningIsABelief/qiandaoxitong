package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"qiandao/model"
	"qiandao/pkg/app"
	"qiandao/pkg/token"
	"qiandao/pkg/util"
	"qiandao/store"
	"qiandao/viewmodel"
	"strings"
	"time"
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
		log.Errorf(app.ErrEmailExist, "该邮箱已被绑定，请换一个试试")
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
	if strings.Compare(updatePasswordRequest.NewPassword, updatePasswordRequest.NewConfirmPassword) != 0 {
		log.Errorf(app.ErrOldNewInconsistent, "请确保两次输入的密码一样")
		return app.ErrOldNewInconsistent
	}
	// 密码进行加密
	psw, err := util.Encrypt(updatePasswordRequest.NewConfirmPassword)
	if err != nil {
		log.Errorf(app.ErrEncrypt, "密码加密出错")
		return app.ErrEncrypt
	}
	if err := store.UpdatePasswordByFieldMapper("user_id", updatePasswordRequest.UserId, psw); err != nil {
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
	if strings.Compare(email, forgetPasswordRequest.Email) != 0 {
		log.Errorf(app.ErrPhoneBinEmail, "请输入手机号绑定的正确邮箱")
		return app.ErrPhoneBinEmail
	}
	// 密码进行加密
	psw, err := util.Encrypt("123456")
	if err != nil {
		log.Errorf(app.ErrEncrypt, "密码加密出错")
		return app.ErrEncrypt
	}
	if err := store.UpdatePasswordByFieldMapper("phone", forgetPasswordRequest.Phone, psw); err != nil {
		return err
	}
	return nil
}

// Login 用户登录 service
func Login(loginRequest viewmodel.LoginRequest, ctx *gin.Context) (viewmodel.LoginResponse, error) {
	// 判断验证码是否过期
	_, err := store.RedisDB.Self.Get("login-code-" + loginRequest.Uuid).Result()
	// 清除redis中的验证码的key
	util.RedisDel("login-code-" + loginRequest.Uuid)
	if err != nil {
		log.Errorf(err, "验证码已过期")
		return viewmodel.LoginResponse{}, app.ErrCodeExpired
	}
	if !util.VerifyCaptcha(loginRequest.Uuid, loginRequest.VerifyValue) {
		log.Errorf(err, "验证码错误")
		return viewmodel.LoginResponse{}, app.ErrCode
	}
	// 根据手机号拿到用户信息
	userInfo, err := store.GetUserInfoByPhone(loginRequest.Phone)
	if err != nil {
		return viewmodel.LoginResponse{}, err
	}
	// 验证密码正确与否
	if err := util.Decrypt(userInfo.Password, loginRequest.Password); err != nil {
		log.Errorf(app.ErrLoginPassword, "用户输入的原来的密码，和数据库中不一致")
		return viewmodel.LoginResponse{}, app.ErrLoginPassword
	}
	// 账号存在 密码正确 生成token
	userToken, err := token.Sign(nil, token.Context{ID: userInfo.UserId}, viper.GetString("jwt_secret"))
	if err != nil {
		return viewmodel.LoginResponse{}, app.ErrTokenCreate
	}

	// 自定义在线用户信息 有用户ID、用户账号、加密的用户颁发token、ip地址，地区， 登录时间
	onlineUserInfo := viewmodel.OnlineUserInfo{
		Id:          userInfo.UserId,
		UserAccount: loginRequest.Phone,
		// 将用户id转换为byte数组进行加密，再将加密后的byte数组转换为base64
		Token:     base64.StdEncoding.EncodeToString(util.DesEncrypt(util.StringToByteSlice(loginRequest.Uuid))),
		Ip:        util.GetRequestIP(ctx),
		Address:   util.GetAddressByLngAndLat(loginRequest.Longitude, loginRequest.Latitude),
		LoginTime: time.Now().Unix(),
	}
	// 将结构体序列化
	marshal, err := json.Marshal(onlineUserInfo)
	if err != nil {
		log.Errorf(err, "结构体序列化失败")
		return viewmodel.LoginResponse{}, errors.New("结构体序列化失败")
	}
	// 将通过验证的账号的信息存入redis
	err = util.RedisSet(viper.GetString("jwt.online_key")+userToken, util.ByteSliceToString(marshal))
	if err != nil {
		log.Errorf(err, "redis键设置失败")
		return viewmodel.LoginResponse{}, err
	}
	// 单用户模式：true，多用户模式：false
	// 为true需要踢掉之前登陆过的用户
	if viper.GetBool("login.single_login") {
		err := util.CheckLoginOnUser(loginRequest.Uuid, userToken)
		if err != nil {
			return viewmodel.LoginResponse{}, err
		}
	}
	return viewmodel.LoginResponse{
		Token: userToken,
		User:  userInfo,
	}, nil
}
