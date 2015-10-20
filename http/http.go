package http

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"github.com/Cepave/alarm/g"
	"github.com/astaxie/beego"
	"log"
	_ "net/http/pprof"

	"github.com/astaxie/beego"
	"github.com/open-falcon/alarm/g"
)

func configRoutes() {
	beego.Router("/", &MainController{}, "get:Index")
	beego.Router("/version", &MainController{}, "get:Version")
	beego.Router("/health", &MainController{}, "get:Health")
	beego.Router("/workdir", &MainController{}, "get:Workdir")
	beego.Router("/config/reload", &MainController{}, "get:ConfigReload")
	beego.Router("/event/solve", &MainController{}, "post:Solve")
}

func Duration(now, before int64) string {
	d := now - before
	if d <= 60 {
		return "just now"
	}

	if d <= 120 {
		return "1 minute ago"
	}

	if d <= 3600 {
		return fmt.Sprintf("%d minutes ago", d/60)
	}

	if d <= 7200 {
		return "1 hour ago"
	}

	if d <= 3600*24 {
		return fmt.Sprintf("%d hours ago", d/3600)
	}

	if d <= 3600*24*2 {
		return "1 day ago"
	}

	return fmt.Sprintf("%d days ago", d/3600/24)
}

func init() {
	configRoutes()
	beego.AddFuncMap("duration", Duration)
}

func Start() {
	if !g.Config().Http.Enabled {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	if g.Config().Debug {
		beego.RunMode = "dev"
	} else {
		beego.RunMode = "prod"
	}

	account := g.Config().Database.Account
	password := g.Config().Database.Password
	host := g.Config().Database.Host
	port := g.Config().Database.Port
	database := g.Config().Database.Db
	str := account + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8"

	orm.RegisterDriver("mysql", orm.DR_MySQL)
	maxIdle := 30
	maxConn := 30
	orm.RegisterDataBase("default", "mysql", str, maxIdle, maxConn)
	orm.RegisterModel(new(Session))

	beego.Run(addr)

	log.Println("http listening", addr)
}
