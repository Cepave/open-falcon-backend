package fastweb

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func ConfigRoutes() {
	//owl-protal-routes
	fastweb := beego.NewNamespace("/api/v1/fastweb",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/getplaftom", &FastWebController{}, "get:GetPlaftom;post:GetPlaftom"),
		beego.NSRouter("/getcontact", &FastWebController{}, "get:GetContact;post:GetContact"),
	)
	beego.AddNamespace(fastweb)
}
