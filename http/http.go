package http

import (
	"github.com/Cepave/fe/g"
	"github.com/Cepave/fe/http/dashboard"
	"github.com/Cepave/fe/http/home"
	"github.com/Cepave/fe/http/uic"
	uic_model "github.com/Cepave/fe/model/uic"
	"github.com/astaxie/beego"
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

	beego.AddFuncMap("member", uic_model.MembersByTeamId)
	beego.Run(addr)
}
