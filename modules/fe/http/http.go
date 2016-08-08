package http

import (
	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/dashboard"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/fastweb"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/home"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/portal"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/uic"
	uic_model "github.com/Cepave/open-falcon-backend/modules/fe/model/uic"
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

	home.ConfigRoutes()
	uic.ConfigRoutes()
	dashboard.ConfigRoutes()
	portal.ConfigRoutes()
	fastweb.ConfigRoutes()

	beego.AddFuncMap("member", uic_model.MembersByTeamId)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
	}))
	beego.Run(addr)
}
