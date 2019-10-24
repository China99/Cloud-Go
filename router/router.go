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
	//不需要验证可以直接访问
	/*	router.GET("/user/signup", handler.SignUpHandler)
		router.POST("/user/signup", handler.DoSignUpHandler)
		router.GET("/user/signin", handler.SignInHandler)
		router.POST("/user/signin", handler.DoSignInHandler)*/
	router.GET("/user/getbalance", handler.GetBalance)
	router.GET("/user/getbalances", handler.GetBalances)
	//加入中间件，用于校验token的拦截器
	router.Use(handler.HTTPInterceptor())

	return router
}
