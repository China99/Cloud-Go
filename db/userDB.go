package db

import (
	mydb "Cloud-Go/db/mysql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"log"
)

var db *gorm.DB

type User struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	CreateTime   string `json:"create_time"`
	LastEditTime string `json:"last_edit_time"`
	Status       int    `json:"status"`
}

func UserSignUp(username string, password string) bool {
	/*var user User
	err := db.Select("username").Where("username=? AND password=?", username, password).First(&user).Error
	if err != nil {
		return false
	}
	if username != "" || password != "" {
		return true
	}
	return false*/

	//代码冗余，待修改
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tb_user(`username`, `password`, `status`) values(?,?,1)")
	if err != nil {
		log.Fatalf("Faild to prepare statement,err", err)
		return false
	}
	defer stmt.Close()
	result, err := stmt.Exec(username, password)
	if err != nil {
		panic(err)
		return false
	}
	rows, err := result.RowsAffected()
	fmt.Printf("rows:%d\n", rows)
	if nil == err && rows > 0 {
		fmt.Printf("user %s creat OK\n", username)
		return true
	}
	return false
}

//sign in a user
func UserSignin(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select * from tb_user where username=? and password=? limit 1")
	if err != nil {
		fmt.Println("stmt err:", err.Error())
		return false
	}
	defer stmt.Close()

	//查询账号
	rows, err := stmt.Query(username, password)
	if err != nil {
		fmt.Println("query user err", err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found" + username)
		return false
	}
	fmt.Println("username:", username, "log in ok")
	return true

}

//创建令牌 或者更新令牌
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tb_user_token(username,user_token) values (?,?)")
	if err != nil {
		fmt.Println("stmt err", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println("update token err ", err.Error())
		return false
	}
	return true

}
