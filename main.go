package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
)

func main() {
	//初始化数据库连接
	db, err := sql.Open("mysql", "root:123456@/gaokao")
	if db.Ping() != nil {
		println("初始化-数据库连接出错", err)
	}

	//初始化iris框架
	app := iris.New()
	app.Get("/rankQuery", func(ctx iris.Context) {
		ctx.Text("连接成功")
		rankQuery(&ctx, db)
	})

	app.Listen(":8080")
}
