package dashboard

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func ConfigRoutes() {
	//owl-protal-routes
	dashboard := beego.NewNamespace("/api/v1/dashboard",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/endpoints", &BashBoardController{}, "post:EndpRegxqury"),
	)
	beego.AddNamespace(dashboard)
}
