package mysql

import (
	"Cloud-Go/config"
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", config.MysqlLink)
	db.SetMaxOpenConns(1000)

	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1) //退出
	}
	fmt.Println("连接成功")
}

func DBConn() *sql.DB {
	return db
}
