package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
)

type sqlobj struct {
	Code        string `json:"code"`
	Collegename string `json:"collegename"`
	Majorcode   string `json:"majorcode"`
	Majorname   string `json:"majorname"`
	Minscore    int    `json:"minscore"`
	Minrank     int    `json:"minrank"`
	Avescore    int    `json:"avescore"`
	Year        int    `json:"year"`
}

//查询3届数据
func rankQuery(ctx iris.Context, db *sql.DB) {
	//测试
	ctx.ContentType("application/javascript")
	if db.Ping() != nil {
		println("handler-数据库连接出错")
	} else {
		println("handler-连接数据库成功")
	}

	//从数据库获取数据，计划采用和spring一样的数据结构
	/*responseData := make(map[string]interface{})
	scoreData := make(map[string]interface{})
	singleScore := make(map[string]interface{})

	//预编译sql: 传入lg/ws 专业名称 年份(2017)
	getMinGrade,err := db.Prepare("select 录取最低分 from ? where 专业名称=? and 年份=?")
	getMinRank,err := db.Prepare("select 最低位次 from ? where 专业名称=? and 年份=?")
	getAverage,err := db.Prepare("select 平均分 from ? where 专业名称=? and 年份=?")
	*/

	rows, err := db.Query("SELECT * FROM `gaokao`.`lg` LIMIT 1")
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
