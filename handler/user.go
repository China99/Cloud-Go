package handler

import (
	"Cloud-Go/util"
	"fmt"
	"time"
)

const (
	pwdSalt   = "!(@*#&$^%"
	tokenSalt = "_tokenSalt"
)

/*
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
	encPassword := util.MD5([]byte(password + pwdSalt))

	//4.向file表存储记录
	suc := db.UserSignUp(username, encPassword)

	if suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "注册成功",
			"code": util.StatusOK,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "注册失败",
			"code": util.StatusRegisterFailed,
		})

	}
	return

}

//登录界面
func SignInHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

//登录处理
func DoSignInHandler(c *gin.Context) {
	//1.获取账号密码
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	//2.判断参数是否合法
	if len(username) < 3 || len(password) < 5 {
		//还没有验证特殊字符，前端可以判断去掉
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "请求参数无效",
			"code": util.StatusLoginFailed, //登录失败
		})
		return
	}
	//3 对密码进行加盐及取MD5值加密
	encPassword := util.MD5([]byte(password + pwdSalt))

	//4.检查用户名以及密码是否在db中
	suc := db.UserSignin(username, encPassword)
	if !suc {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "用户名或者密码错误",
			"code": util.StatusLoginFailed, //登录错误

		})
		return
	}

	//5.生成token并且保存到数据库中
	//token := GenToken(username)
	//db.UpdateToken
}
*/
//生成用户令牌
func GenToken(username string) string {
	timestamp := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + timestamp + tokenSalt))
	token := tokenPrefix + timestamp[:8]
	fmt.Printf("username:%s Token: %s\n", username, token)
	return token
}

//检查token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}

	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}
