package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qiandao/controller/checkin"

	"qiandao/controller/class"
	"qiandao/controller/lesson"
	"qiandao/controller/sd"
	"qiandao/controller/user"
)

func Load(engine *gin.Engine, handlerFunc ...gin.HandlerFunc) *gin.Engine {
	engine.Use(gin.Recovery())
	engine.Use(handlerFunc...)

	// NoRoute()是默认情况下都返回404代码
	engine.NoRoute(func(context *gin.Context) {
		// 将给定的字符串写入响应正文
		context.String(http.StatusNotFound, "API路由错误")
	})

	userAPI := engine.Group("/api/user")
	{
		userAPI.POST("/register", user.Register)
		userAPI.POST("/login", user.Login)
		userAPI.PUT("/update-user", user.UpdateUserInfo)
		userAPI.PUT("/update-email", user.UpdateEmail)
		userAPI.PUT("/update-nick-name", user.UpdateNickName)
	}

	classAPI := engine.Group("/api/class")
	{
		classAPI.POST("", class.Create)
		classAPI.GET("", class.GetAllClass)
	}

	// 检查http健康的路由组
	svcd := engine.Group("/api/sd")
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
		lessonApi.POST("", lesson.CreateLesson)
		// 获取创建的课程列表
		lessonApi.GET("/user", lesson.GetCreateLessonList)
		//获取加入的课程列表
		lessonApi.GET("/join", lesson.GetJoinLessonList)
	}

	// 签到
	checkInApi := engine.Group("/api/checkin")
	{
		// 创建签到
		checkInApi.POST("createCheckin", checkin.CreateCheckin)
		checkInApi.POST("checkin", checkin.CheckIn)
		checkInApi.GET("getCheckinDetails", checkin.GetCheckinDetails)
		checkInApi.GET("getCheckinRecList", checkin.GetCheckinRecList)
		checkInApi.GET("getCreatedCheckinList", checkin.GetCreatedCheckinList)
	}

	return engine
}
