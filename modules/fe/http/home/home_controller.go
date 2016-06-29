package home

import (
	"github.com/astaxie/beego"
	"github.com/Cepave/fe/g"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["Shortcut"] = g.Config().Shortcut
	this.TplName = "home/index.html"
}
