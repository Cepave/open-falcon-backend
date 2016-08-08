package boss

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func ConfigRoutes() {
	//owl-protal-routes
	boss := beego.NewNamespace("/api/v1/boss",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/getplaftom", &BossController{}, "get:GetPlaftom;post:GetPlaftom"),
		beego.NSRouter("/getcontact", &BossController{}, "get:GetContact;post:GetContact"),
	)
	beego.AddNamespace(boss)
}
