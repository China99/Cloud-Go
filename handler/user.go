package handler

import (
	"Cloud-Go/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	pwdSalt   = "!(@*#&$^%"
	tokenSalt = "_tokenSalt"
)

//注册界面
func SignUpHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

//注册用户
func DoSignUpHandler(c *gin.Context) {
	//获取用户名和密码
	//c.Query("username")
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	if len(username) < 3 || len(password) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "请求参数无效",
			"code": util.StatusRegisterFailed,
		})
		return
	}
	//对密码进行加盐及取MD5值加密
	encPwd := util.MD5([]byte(password + pwdSalt))

	//4.向file表存储记录
	mydb.UserSignUp(username, encPwd)
}
