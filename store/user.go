package store

import (
	"github.com/lexkong/log"
	"qiandao/model"
	"qiandao/pkg/app"
	"qiandao/viewmodel"
)

// GetByPhoneMapper 查询手机号是否在数据库中存在 mapper
func GetByPhoneMapper(phone string) bool {
	isPhone := DB.Self.Where("phone = ?", phone).First(&model.User{})
	if isPhone.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

// CreateUserMapper 创建用户 mapper
func CreateUserMapper(user *model.User, classId string) error {
	tx := DB.Self.Begin()
	// 创建用户
	db := tx.Create(&user)
	if db.Error != nil {
		tx.Rollback()
		log.Errorf(db.Error, "创建班级失败")
		return app.InternalServerError
	}
	log.Infof("创建班级：成功创建:%v条记录", db.RowsAffected)
	// 维护用户班级中间表
	intermediateTable := tx.Create(&model.Connection{
		ClassId: classId,
		UserId:  user.UserId,
	})
	if intermediateTable.Error != nil {
		tx.Rollback()
		log.Errorf(intermediateTable.Error, "维护中间表失败")
		return app.InternalServerError
	}
	log.Infof("维护班级用户中间表：成功创建:%v条记录", intermediateTable.RowsAffected)
	tx.Commit()
	return nil
}

// UpdateUserMapper 修改用户信息 mapper
func UpdateUserMapper(updateUser viewmodel.UpdateUserInfoRequest) error {
	result := DB.Self.Model(model.User{}).Where("user_id = ?", updateUser.UserId).Updates(model.User{
		Email:    updateUser.Email,
		RealName: updateUser.RealName,
		Hobby:    updateUser.Hobby,
		Address:  updateUser.Address,
		Sex:      updateUser.Sex,
		Age:      updateUser.Age,
	})
	if result.Error != nil {
		log.Errorf(result.Error, "修改用户信息失败")
		return app.InternalServerError
	}
	log.Infof("修改用户信息：成功修改 %v 条记录", result.RowsAffected)
	return nil
}

// UpdateEmailMapper 修改邮箱 mapper
func UpdateEmailMapper(updateEmail viewmodel.UpdateEmailRequest) error {
	result := DB.Self.Model(&model.User{}).Where("user_id = ?", updateEmail.UserId).Update("email", updateEmail.Email)
	if result.Error != nil {
		log.Errorf(result.Error, "修改邮箱失败失败")
		return app.InternalServerError
	}
	log.Infof("修改邮箱：成功修改 %v 条记录", result.RowsAffected)
	return nil
}

// UpdateNickNameMapper 修改昵称 mapper
func UpdateNickNameMapper(updateNickName viewmodel.UpdateNickNameRequest) error {
	result := DB.Self.Model(&model.User{}).Where("user_id = ?", updateNickName.UserId).Update("nick_name", updateNickName.NickName)
	if result.Error != nil {
		log.Errorf(result.Error, "修改昵称失败")
		return app.InternalServerError
	}
	log.Infof("修改昵称：成功修改 %v 条记录", result.RowsAffected)
	return nil
}
