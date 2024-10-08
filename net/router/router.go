package router

import (
	"tiny_talk/docs"
	"tiny_talk/net/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 设置页面
	router.GET("/", service.GetIndex)
	router.GET("/Register", service.Register)
	router.GET("/home", service.Home)

	// 静态资源
	router.Static("/static", "/root/tiny_talk/webui")

	//设置服务
	router.POST("/user/CreateUser", service.CreateUser)
	router.POST("/user/Login", service.UserLogin)
	router.POST("/user/TestToken", service.TestToken)
	router.GET("/friend/GetFriendList", service.GetFriendList)
	router.GET("/ws", service.HandleWebSocket)
	return router
}
