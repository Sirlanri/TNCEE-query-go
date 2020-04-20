package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iris-contrib/middleware/cors"
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
	//允许跨域请求
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //允许通过的主机名称
		AllowCredentials: true,
	})
	//后端api接口

	go1 := app.Party("/go", crs).AllowMethods(iris.MethodOptions)
	{
		go1.Post("/numschange", func(ctx iris.Context) {
			numsChange(ctx, db)
		})
		go1.Post("/recommend", func(ctx iris.Context) {
			recommend(ctx, db)
		})
	}
	app.Listen(":8090")
}
