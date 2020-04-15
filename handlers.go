package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
)

type sqlobj struct {
	code        string `json:"code"`
	collegename string `json:"collegename"`
	majorcode   string `json:"majorcode"`
	majorname   string `json:"majorname"`
	minscore    int    `json:"minscore"`
	minrank     int    `json:"minrank"`
	avescore    int    `json:"avescore"`
	year        int    `json:"year"`
}

//查询3届数据
func rankQuery(ctx iris.Context, db *sql.DB) {
	ctx.ContentType("application/javascript")
	if db.Ping() != nil {
		println("handler-数据库连接出错")
	} else {
		println("handler-连接成功")
	}
	rows, err := db.Query("SELECT * FROM `gaokao`.`lg` LIMIT 10")
	if err != nil {
		println("handler-数据库测试出错", err)
	}

	sqls := make([]sqlobj, 0, 1000)

	for rows.Next() {
		var sqlnow sqlobj
		err := rows.Scan(&sqlnow)
		if err != nil {
			println("遍历出错", err.Error())
		}
		_, err = ctx.JSON(sqlnow)
		if err != nil {
			println("json出错", err.Error())
		}
		sqls = append(sqls, sqlnow)
	}
	fmt.Println(sqls)

}
