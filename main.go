package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
)

func main3() {
	slice1 := make([]int, 0, 3)
	slice1 = append(slice1, 8, 5, 10)
	fmt.Println(getMax(slice1))
}

func main() {
	//初始化数据库连接
	db, err := sql.Open("mysql", "root:123456@/gaokao")
	if db.Ping() != nil {
		println("初始化-数据库连接出错", err)
	}

	//初始化iris框架
	app := iris.New()

	//后端api接口
	app.Post("numschange", func(ctx iris.Context) {
		numsChange(ctx, db)
	})

	app.Listen(":8090")
}
