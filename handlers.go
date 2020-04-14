package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
)

//查询3届数据
func rankQuery(ctx *iris.Context, db *sql.DB) {

	if db.Ping() != nil {
		println("handler-数据库连接出错")
	} else {
		println("handler-连接成功")
	}
	rows, err := db.Query("SELECT * FROM `gaokao`.`lg` LIMIT 10")
	if err != nil {
		println("handler-数据库测试出错", err)
	}
	for rows.Next() {
		text1 := ""
		rows.Scan(text1)
		fmt.Println(text1)
	}
}
