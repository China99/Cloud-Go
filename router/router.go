package router

import (
	"Cloud-Go/handler"
	"github.com/gin-gonic/gin"
)

//Router:路由表配置
func Router() *gin.Engine {
	//gin framework
	router := gin.Default()

	//处理静态资源
	router.Static("/static", "./static")

	// 目录[handler/user.go
	router.GET("/user/signup", handler.SignUpHandler)
	router.POST("/user/signup", handler.DoSignUpHandler)
	router.GET("/user/signin", handler.SignInHandler)
	router.POST("/user/signin", handler.DoSignInHandler)
}
