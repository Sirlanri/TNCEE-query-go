package main

import (
	"database/sql"
)

//优化数据库,运行一次就完事儿
func getscores() {
	db, err := sql.Open("mysql", "root:123456@/gaokao")
	if db.Ping() != nil {
		println("初始化-数据库连接出错", err)
	}

	getrank, err := db.Prepare("select 文史科类累计人数 from rank17 where 成绩分数段=?")
	insertRank, err := db.Prepare("update ws set 平均位次=? where 平均分=? and 年份=2017")

	if err != nil {
		println(err.Error())
	}
	rows, err := db.Query("select 平均分 from ws where 年份=2017")
	if err != nil {
		println(err.Error())
	}
	for rows.Next() {
		score := 0
		rows.Scan(&score)
		ranksour := getrank.QueryRow(score)
		rank := ""
		ranksour.Scan(&rank)
		if err != nil {
			println("处理rows出错", err.Error())
		}

		//将数据写入平均位次
		insertRank.Exec(rank, score)
	}
}
