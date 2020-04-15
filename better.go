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

func newscorerank() {
	db, err := sql.Open("mysql", "root:123456@/gaokao")
	if db.Ping() != nil {
		println("初始化-数据库连接出错", err)
	}

	getinfo, err := db.Prepare(`
		select 专业名称, 录取最低分, 最低位次, 平均分, 平均位次 from ws where 年份=? 
	`)
	writein, err := db.Prepare(`
		insert into ws19 (name, minscore, minrank, avescore, averank) values (?, ?, ?, ?, ?)
	`)
	if err != nil {
		println("创建表达式", err.Error())
	}

	if err != nil {
		println(err.Error())
	}
	rows, err := getinfo.Query(2019)
	if err != nil {
		println(err.Error())
	}
	for rows.Next() {
		var name string
		var minscore, minrank, avescore, averank int
		err := rows.Scan(&name, &minscore, &minrank, &avescore, &averank)
		if err != nil {
			println("读取出错", err.Error())
		}
		_, err = writein.Exec(name, minscore, minrank, avescore, averank)
		if err != nil {
			println("写入出错", err.Error())
		}
	}
}

//统计性别比例
func sex() {
	db, err := sql.Open("mysql", "root:123456@/gaokao")
	if db.Ping() != nil {
		println("初始化-数据库连接出错", err)
	}

	//getman, err := db.Prepare("select 专业,count(性别) from gaokao.19totaldata where 性别='男' and 科类代码='理工' group by 专业")
	getwoman, err := db.Prepare("select count(性别) from gaokao.18totaldata where 性别='女' and 科类代码='文史' and 专业=? group by 专业")
	writein, err := db.Prepare("update ws18 set sex=? where name=?")

	mans, err := db.Query("select 专业,count(性别) from gaokao.18totaldata where 性别='男' and 科类代码='文史' group by 专业")
	if err != nil {
		println("读取出错", err.Error())
	}
	for mans.Next() {
		var name string
		var score, womanscore, point float32
		mans.Scan(&name, &score)

		womans := getwoman.QueryRow(name)
		womans.Scan(&womanscore)

		if womanscore == 0 {
			point = 1
		} else {
			point = score / (womanscore + score)
		}

		_, err := writein.Exec(point, name)
		if err != nil {
			println("写入出错", err.Error())
		}
	}

}
