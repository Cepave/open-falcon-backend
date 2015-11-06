package home

import (
	"github.com/astaxie/beego"
	"github.com/freedomkk-qfeng/feg"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["Shortcut"] = g.Config().Shortcut
	this.TplNames = "home/index.html"
}
