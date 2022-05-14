package store

import (
	"github.com/lexkong/log"
	"qiandao/model"
	"qiandao/pkg/app"
	"qiandao/viewmodel"
)

// GetByClassNameMapper 查询要创建的班级是否在数据库中存在 mapper
func GetByClassNameMapper(className string) bool {
	isClassName := DB.Self.Where("class_name = ?", className).First(&model.Class{})
	if isClassName.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

// CreateClassMapper 创建班级 mapper
func CreateClassMapper(class *model.Class) error {
	result := DB.Self.Create(&class)
	if result.Error != nil {
		log.Errorf(result.Error, "创建班级：创建班级失败")
		return app.InternalServerError
	}
	log.Infof("创建班级：成功创建 %v 条记录", result.RowsAffected)
	return nil
}

// GetAllClassPageMapper 分页获取班级列表 mapper
func GetAllClassPageMapper(offset, limit int) ([]viewmodel.ClassInfo, uint64, error) {
	if limit == 0 {
		limit = 50
	}
	classes := make([]*model.Class, 0)

	responseClassInfo := make([]viewmodel.ClassInfo, limit)
	var count uint64
	classCount := DB.Self.Table("class").Count(&count)
	if classCount.Error != nil {
		log.Errorf(classCount.Error, "获取所有用户数量出错")
		return responseClassInfo, count, app.InternalServerError
	}
	log.Infof("获取所有班级：一共有：%v个用户", count)
	result := DB.Self.Select([]string{"class_id", "class_name"}).Offset(offset).Limit(limit).Order("created_at desc").Find(&classes)
	if result.Error != nil {
		log.Errorf(result.Error, "分页获取用户信息出错")
		return responseClassInfo, count, app.InternalServerError
	}
	for k, _ := range classes {
		responseClassInfo[k].ClassId = classes[k].ClassId
		responseClassInfo[k].ClassName = classes[k].ClassName
	}
	log.Info("分页获取班级列表成功")
	return responseClassInfo, count, nil
}

// GetAllClassMapper 获取班级列表 mapper
func GetAllClassMapper() ([]viewmodel.ClassInfo, uint64, error) {
	classes := make([]*model.Class, 0)
	var count uint64
	classCount := DB.Self.Table("class").Count(&count)
	if classCount.Error != nil {
		log.Errorf(classCount.Error, "获取所有用户数量出错")
		return []viewmodel.ClassInfo{}, count, app.InternalServerError
	}
	log.Infof("获取所有班级：一共有：%v个用户", count)
	responseClassInfo := make([]viewmodel.ClassInfo, count)
	result := DB.Self.Select([]string{"class_id", "class_name"}).Order("created_at desc").Find(&classes)
	if result.Error != nil {
		log.Errorf(result.Error, "获取用户信息出错")
		return responseClassInfo, count, app.InternalServerError
	}
	for k, _ := range classes {
		responseClassInfo[k].ClassId = classes[k].ClassId
		responseClassInfo[k].ClassName = classes[k].ClassName
	}
	log.Info("获取所有班级列表成功")
	return responseClassInfo, count, nil
}
