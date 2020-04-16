package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
)

//post请求-三年分数图
type numsChangePost struct {
	Profession string `json:"profession"`
	Province   string `json:"province"`
	Type       string `json:"type"`
}

//handler-3年分数位次曲线
func numsChange(ctx iris.Context, db *sql.DB) {
	ctx.ContentType("application/javascript")
	if db.Ping() != nil {
		println("handler-数据库连接出错")
	}

	//解析post数据
	var postInfor numsChangePost
	if err := ctx.ReadJSON(&postInfor); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("非法的post请求格式" + err.Error())
		println("非法请求格式", err.Error())
		return
	}

	aveScore := make([]int, 0, 3)
	minScore := make([]int, 0, 3)
	minRank := make([]int, 0, 3)
	var minscore, minrank, avescore int

	getExplg, err := db.Prepare(`
		select minscore,minrank,avescore from gaokao.lg17 where name=?
		union
		select minscore,minrank,avescore from gaokao.lg18 where name=?
		union
		select minscore,minrank,avescore from gaokao.lg19 where name=?
	`)
	if err != nil {
		println("预编译表达式出错", err.Error())
	}

	//执行sql查询
	switch postInfor.Type {
	case "理工":
		rows, err := getExplg.Query(postInfor.Profession, postInfor.Profession, postInfor.Profession)
		if err != nil {
			println("执行sql出错", err.Error())
		}
		//rows结果，分别是17 18 19
		for rows.Next() {
			err := rows.Scan(&minscore, &minrank, &avescore)
			if err != nil {
				println("执行sql后写入数据出错", err.Error())
			}
			aveScore = append(aveScore, avescore)
			minScore = append(minScore, minscore)
			minRank = append(minRank, minrank)
		}

		resMap := make(map[string]interface{})
		resMap["avescore"] = aveScore
		resMap["minscore"] = minScore
		resMap["minrank"] = minRank

		_, err = ctx.JSON(resMap)
		if err != nil {
			println("打包返回json失败", err.Error())
		} else {
			fmt.Println("成功返回三年数据", aveScore, minScore, minRank)
		}
	}
	getExplg.Close()
}
func sexPro(ctx iris.Context, db *sql.DB) {
}

//查询3届数据
/*
func rankQuery(ctx iris.Context, db *sql.DB) {
	//测试
	ctx.ContentType("application/javascript")
	if db.Ping() != nil {
		println("handler-数据库连接出错")
	} else {
		println("handler-连接数据库成功")
	}

	//从数据库获取数据
	responseData := make(map[string]interface{})
	scoreData := make(map[string]interface{})
	singleScore := make(map[string]interface{})

	aveScore := make([]int, 0, 3)
	minScore := make([]int, 0, 3)
	minRank := make([]int, 0, 3)

	//预编译sql: 传入lg/ws 专业名称 年份(2017)
	getMinGrade, err := db.Prepare("select 录取最低分 from ? where 专业名称=? and 年份=?")
	getMinRank, err := db.Prepare("select 最低位次 from ? where 专业名称=? and 年份=?")
	getAverage, err := db.Prepare("select 平均分 from ? where 专业名称=? and 年份=?")

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
*/
