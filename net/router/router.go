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
	router.GET("/index", service.GetIndex)
	router.POST("/user/CreateUser", service.CreateUser)
	router.POST("/user/Login", service.Login)
	return router
}
