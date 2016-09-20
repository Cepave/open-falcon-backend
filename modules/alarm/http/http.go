package http

import (
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	_ "net/http/pprof"

	"github.com/Cepave/open-falcon-backend/modules/alarm/g"
	log "github.com/Sirupsen/logrus"

	"github.com/astaxie/beego"
)

func configRoutes() {
	beego.Router("/", &MainController{}, "get:Index")
	beego.Router("/version", &MainController{}, "get:Version")
	beego.Router("/health", &MainController{}, "get:Health")
	beego.Router("/workdir", &MainController{}, "get:Workdir")
	beego.Router("/config/reload", &MainController{}, "get:ConfigReload")
	beego.Router("/event", &MainController{}, "get:Event")
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
	beego.SetLogger("console", `{"color":false}`)
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
		beego.BConfig.RunMode = "dev"
	} else {
		beego.BConfig.RunMode = "prod"
	}

	if g.Config().Http.ViewPath != "" {
		log.Infof("set http view_path in %v", g.Config().Http.ViewPath)
		beego.SetViewsPath(g.Config().Http.ViewPath)
	}
	if g.Config().Http.StaticPath != "" {
		log.Infof("set http static_path in %v", g.Config().Http.StaticPath)
		beego.SetStaticPath("/static", g.Config().Http.StaticPath)
	}

	beego.Run(addr)

	log.Println("http listening", addr)
}
