package http

import (
	"github.com/astaxie/beego"
	"github.com/freedomkk-qfeng/fe/g"
	"github.com/freedomkk-qfeng/fehttp/home"
	"github.com/freedomkk-qfeng/fe/http/uic"
	uic_model "github.com/freedomkk-qfeng/femodel/uic"
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

	beego.AddFuncMap("member", uic_model.MembersByTeamId)
	beego.Run(addr)
}
