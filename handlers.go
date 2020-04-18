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
	maxScore := make([]int, 0, 3)
	var minscore, minrank, avescore, maxscore int

	//返回的密度数据
	axisX := make([]int, 0, 100)

	//执行sql查询
	switch postInfor.Type {
	case "理工":
		//预编译表达式，但是不够完美，写预编译纯属为了好看

		//获取三年平均最低折线图
		getExplg, err := db.Prepare(`
			select minscore,minrank,avescore,maxscore from gaokao.lg17 where name=?
			union
			select minscore,minrank,avescore,maxscore from gaokao.lg18 where name=?
			union
			select minscore,minrank,avescore,maxscore from gaokao.lg19 where name=?
		`)
		if err != nil {
			println("预编译表达式出错", err.Error())
		}
		rows, err := getExplg.Query(postInfor.Profession, postInfor.Profession, postInfor.Profession)
		if err != nil {
			println("执行sql出错", err.Error())
		}

		//rows结果，分别是17 18 19
		for rows.Next() {
			err := rows.Scan(&minscore, &minrank, &avescore, &maxscore)
			if err != nil {
				println("执行sql后写入数据出错", err.Error())
			}
			defer rows.Close()
			//此专业数据残缺（没有最高分），返回421
			if maxscore == 0 {
				ctx.StatusCode(421)
				println("421-数据残缺")
				return
			}
			aveScore = append(aveScore, avescore)
			minScore = append(minScore, minscore)
			minRank = append(minRank, minrank)
			maxScore = append(maxScore, maxscore)
		}

		//没有这个专业数据 返回 404
		if len(aveScore) == 0 {
			ctx.StatusCode(404)
			println("404-找不到专业数据")
			return
		}

		//获取性别比例
		getSex, err := db.Prepare("select sex from gaokao.lg18 where name=?")
		var sexnum float32
		err = getSex.QueryRow(postInfor.Profession).Scan(&sexnum)
		if err != nil {
			println("执行获取性别sql出错", err.Error())
		}

		//获取分数密度
		getThisScore17, err := db.Prepare("select 成绩,count(*) from gaokao.17totaldata where 专业=? and 成绩 between ? and ? group by 成绩 order by 成绩")
		getAllScore17, err := db.Prepare("select count(*) from gaokao.17totaldata where 专业=?")
		getThisScore18, err := db.Prepare("select 成绩,count(*) from gaokao.18totaldata where 专业=? and 成绩 between ? and ? group by 成绩 order by 成绩")
		getAllScore18, err := db.Prepare("select count(*) from gaokao.18totaldata where 专业=?")
		getThisScore19, err := db.Prepare("select 成绩,count(*) from gaokao.19totaldata where 专业=? and 成绩 between ? and ? group by 成绩 order by 成绩")
		getAllScore19, err := db.Prepare("select count(*) from gaokao.19totaldata where 专业=?")
		if err != nil {
			println("分数密度预编译表达式出错", err.Error())
		}

		//获取X轴数据，[三年最低分，三年最高]
		maxdata := getMax(maxScore)
		mindata := getMin(minScore)
		//写入X轴数据
		for i := mindata; i < maxdata; i++ {
			axisX = append(axisX, i)
		}

		//17级的Y轴百分比
		var allNum17 float64 //分母-专业录取人数
		err = getAllScore17.QueryRow(postInfor.Profession).Scan(&allNum17)
		if err != nil {
			println("获取17X轴出错", err.Error())
		}

		thisScoreRows17, err := getThisScore17.Query(postInfor.Profession, mindata, maxdata)
		if err != nil {
			println("获取17X轴对应Y比例出错", err.Error())
		}
		ratio17 := make([]float64, 0, 100)

		thisScoreRows17.Next()
		for i := mindata - 1; i < maxdata; i++ {
			var thisNum, result float64
			var xnum int
			thisScoreRows17.Scan(&xnum, &thisNum)
			if xnum == i {
				result = thisNum / allNum17
				ratio17 = append(ratio17, result)
			} else {
				ratio17 = append(ratio17, 0)
				continue
			}
			thisScoreRows17.Next()
		}

		//18级的Y轴百分比
		var allNum18 float64 //分母-专业录取人数
		err = getAllScore18.QueryRow(postInfor.Profession).Scan(&allNum18)
		if err != nil {
			println("获取18X轴出错", err.Error())
		}

		thisScoreRows18, err := getThisScore18.Query(postInfor.Profession, mindata, maxdata)
		if err != nil {
			println("获取17X轴对应Y比例出错", err.Error())
		}
		ratio18 := make([]float64, 0, 100)

		thisScoreRows18.Next()
		for i := mindata - 1; i < maxdata; i++ {
			var thisNum, result float64
			var xnum int
			thisScoreRows18.Scan(&xnum, &thisNum)
			if xnum == i {
				result = thisNum / allNum17
				ratio18 = append(ratio18, result)
			} else {
				ratio18 = append(ratio18, 0)
				continue
			}
			thisScoreRows18.Next()
		}

		//19级的Y轴百分比
		var allNum19 float64 //分母-专业录取人数
		err = getAllScore19.QueryRow(postInfor.Profession).Scan(&allNum19)
		if err != nil {
			println("获取17X轴出错", err.Error())
		}

		thisScoreRows19, err := getThisScore19.Query(postInfor.Profession, mindata, maxdata)
		if err != nil {
			println("获取17X轴对应Y比例出错", err.Error())
		}
		ratio19 := make([]float64, 0, 100)

		thisScoreRows19.Next()
		for i := mindata - 1; i < maxdata; i++ {
			var thisNum, result float64
			var xnum int
			thisScoreRows19.Scan(&xnum, &thisNum)
			if xnum == i {
				result = thisNum / allNum17
				ratio19 = append(ratio19, result)
			} else {
				ratio19 = append(ratio19, 0)
				continue
			}
			thisScoreRows19.Next()
		}

		resMap := make(map[string]interface{})
		resMap["avescore"] = aveScore
		resMap["minscore"] = minScore
		resMap["minrank"] = minRank
		resMap["sex"] = sexnum
		resMap["axisx"] = axisX
		resMap["axis17"] = ratio17
		resMap["axis18"] = ratio18
		resMap["axis19"] = ratio19

		_, err = ctx.JSON(resMap)
		if err != nil {
			println("打包返回json失败", err.Error())
		} else {
			fmt.Println("成功返回三年数据", aveScore, minScore, minRank, sexnum)
		}

		getExplg.Close()

	case "文史":

		//获取三年数据
		getExplg, err := db.Prepare(`
			select minscore,minrank,avescore from gaokao.ws17 where name=?
			union
			select minscore,minrank,avescore from gaokao.ws18 where name=?
			union
			select minscore,minrank,avescore from gaokao.ws19 where name=?
		`)
		if err != nil {
			println("预编译表达式出错", err.Error())
		}
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
		//没有这个专业数据 返回 404
		if len(aveScore) == 0 {
			ctx.StatusCode(404)
			println("404-找不到专业数据")
			return
		}
		//获取性别比例
		getSex, err := db.Prepare("select sex from gaokao.ws18 where name=?")
		var sexnum float32
		err = getSex.QueryRow(postInfor.Profession).Scan(&sexnum)
		if err != nil {
			println("执行获取性别sql出错", err.Error())
		}

		resMap := make(map[string]interface{})
		resMap["avescore"] = aveScore
		resMap["minscore"] = minScore
		resMap["minrank"] = minRank
		resMap["sex"] = sexnum

		_, err = ctx.JSON(resMap)
		if err != nil {
			println("打包返回json失败", err.Error())
		} else {
			fmt.Println("成功返回三年数据", aveScore, minScore, minRank, sexnum)
		}
		getExplg.Close()
	}

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
