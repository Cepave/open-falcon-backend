package http

import (
	"strings"

	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/boss"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/dashboard"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/home"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/portal"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/uic"
	uic_model "github.com/Cepave/open-falcon-backend/modules/fe/model/uic"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func Start() {
	if !g.Config().Http.Enabled {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	switch strings.ToLower(g.Config().Log) {
	case "info":
		beego.SetLevel(beego.LevelInformational)
	case "debug":
		beego.SetLevel(beego.LevelDebug)
	case "warn":
		beego.SetLevel(beego.LevelWarning)
	case "error":
		beego.SetLevel(beego.LevelError)
	}

	home.ConfigRoutes()
	uic.ConfigRoutes()
	dashboard.ConfigRoutes()
	portal.ConfigRoutes()
	boss.ConfigRoutes()

	beego.SetLogger("console", `{"color":false}`)
	beego.AddFuncMap("member", uic_model.MembersByTeamId)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
	}))
	if g.Config().Http.ViewPath != "" {
		log.Infof("set http view_path in %v", g.Config().Http.ViewPath)
		beego.SetViewsPath(g.Config().Http.ViewPath)
	}
	if g.Config().Http.StaticPath != "" {
		log.Infof("set http static_path in %v", g.Config().Http.StaticPath)
		beego.SetStaticPath("/static", g.Config().Http.StaticPath)
	}
	log.Infof("current beego verion: %v", beego.VERSION)
	beego.Run(addr)
}
