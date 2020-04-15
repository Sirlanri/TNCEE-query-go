package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
)

type sqlobj struct {
	code        string
	collegename string
	majorcode   string
	majorname   string
	minscore    int
	minrank     int
	avescore    int
	year        int
}

//查询3届数据
func rankQuery(ctx iris.Context, db *sql.DB) {
	ctx.Text("链接数据库成功~")

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
		sqlnow := sqlobj{}
		err := rows.Scan(&sqlnow.code, &sqlnow.collegename, &sqlnow.majorcode, &sqlnow.majorname, &sqlnow.minscore,
			&sqlnow.minrank, &sqlnow.avescore, &sqlnow.year)
		if err != nil {
			println("遍历出错", err.Error())
		}
		js, _ := json.Marshal(sqlnow)
		_, err = ctx.JSON(js)
		if err != nil {
			println("json出错", err.Error())
		}
		sqls = append(sqls, sqlnow)
	}
	fmt.Println(sqls)

}
