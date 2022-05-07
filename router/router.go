package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qiandao/controller/lesson"
	"qiandao/controller/sd"
)

func Load(engine *gin.Engine, handlerFunc ...gin.HandlerFunc) *gin.Engine {
	engine.Use(gin.Recovery())
	engine.Use(handlerFunc...)

	// NoRoute()是默认情况下都返回404代码
	engine.NoRoute(func(context *gin.Context) {
		// 将给定的字符串写入响应正文
		context.String(http.StatusNotFound, "API路由错误")
	})

	// 检查http健康的路由组
	svcd := engine.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	// 课程
	lessonApi := engine.Group("/api/lesson")
	{
		// 创建课程
		lessonApi.POST("",lesson.CreateLesson)
		// 获取创建的课程列表
		lessonApi.GET("/user",lesson.GetCreateLessonList)
       //获取加入的课程列表
		lessonApi.GET("/join",lesson.GetJoinLessonList)
	}
	return engine
}
