package home

import (
	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["Shortcut"] = g.Config().Shortcut
	this.TplName = "home/index.html"
}
