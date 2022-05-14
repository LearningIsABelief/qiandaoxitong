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

// CreateClass 创建班级 service
func CreateClass(createClassRequest viewmodel.CreateClassRequest) error {
	// 查看要创建的班级是否存在
	if className := store.GetByClassNameMapper(createClassRequest.ClassName); className {
		log.Errorf(app.ErrClassExist, "创建班级失败")
		return app.ErrClassExist
	}

	if err := store.CreateClassMapper(&model.Class{
		ClassId:       util.GetUUID(),
		ClassName:     createClassRequest.ClassName,
		ClassCapacity: createClassRequest.Capacity,
		CreateId:      "1",
	}); err != nil {
		return err
	}
	return nil
}

// GetAllClass 获取班级列表 service
func GetAllClass(page util.PageRequest) (viewmodel.GetClassListResponse, error) {
	var infos []viewmodel.ClassInfo
	var count uint64
	var err error
	if strings.Compare(page.Logo, "register") == 0 {
		infos, count, err = store.GetAllClassMapper()
	} else if strings.Compare(page.Logo, "universal") == 0 {
		infos, count, err = store.GetAllClassPageMapper(page.Offset, page.Limit)
	}
	return viewmodel.GetClassListResponse{
		TotalCount: count,
		Class:      infos,
	}, err
}
