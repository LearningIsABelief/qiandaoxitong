package lesson

import (
	"github.com/gin-gonic/gin"
	"qiandao/pkg/app"
	"qiandao/service"
	"qiandao/viewmodel"
)

// CreateLesson 创建课程
func CreateLesson(ctx *gin.Context)  {
//	 1.绑定参数
	lesson := new(viewmodel.Lesson)
	err := ctx.ShouldBindJSON(lesson)
	if err != nil {
		app.SendResponse(ctx,app.ErrBind,nil)
		return
	}
//	 2.调用业务逻辑
	if lesson.LessonName == "" {
		app.SendResponse(ctx,app.ErrParamNull,nil)
		return
	}
	if len(lesson.ClassList) == 0 {
		app.SendResponse(ctx,app.ErrParamNull,nil)
		return
	}
	err = service.CreateLesson(lesson)
	if err != nil {
		return
	}

//	 3.返回响应
	app.SendResponse(ctx,app.OK,nil)
}

// GetCreateLessonList 获取创建的课程列表
func GetCreateLessonList(ctx *gin.Context)  {
	
}

// GetJoinLessonList 获取加入的课程列表
func GetJoinLessonList(ctx *gin.Context)  {

}
