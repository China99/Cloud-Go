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
		"INSERT into tb_user(`username`,`password`,`status`)values (?,?,1)")
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
