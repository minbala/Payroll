package main

import (
	_ "payroll/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	sqlConn, err := beego.AppConfig.String("sqlconn")
	if err != nil {

	}
	orm.RegisterDataBase("default", "mysql", sqlConn)
	err = orm.RunSyncdb("default", false, false)
	if err != nil {

	}
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
