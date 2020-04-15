package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
)

type sqlobj struct {
	Code        string
	Collegename string
	Majorcode   string
	Majorname   string
	Minscore    int
	Minrank     int
	Avescore    int
	Year        int
}

//查询3届数据
func rankQuery(ctx iris.Context, db *sql.DB) {
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
		err := rows.Scan(&sqlnow.Code, &sqlnow.Collegename, &sqlnow.Majorcode, &sqlnow.Majorname, &sqlnow.Minscore,
			&sqlnow.Minrank, &sqlnow.Avescore, &sqlnow.Year)
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
