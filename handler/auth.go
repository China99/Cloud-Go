package handler

import (
	//"Cloud-Go/config"
	"github.com/gin-gonic/gin"
	//"net/http"
)

//权限认证
func HTTPInterceptor() gin.HandlerFunc {
	return func(context *gin.Context) {
		username := context.Request.FormValue("username")
		token := context.Request.FormValue("token")
		if len(username) < 3 || !IsTokenValid(token) {

		}

	}

}
