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
		beego.NSRouter("/endpoints", &BashBoardController{}, "get:EndpRegxqury;post:EndpRegxqury"),
		beego.NSRouter("/endpointcounters", &BashBoardController{}, "get:CounterQuery;post:CounterQuery"),
	)
	hostgroup := beego.NewNamespace("/api/v1/hostgroup",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/query", &BashBoardController{}, "get:HostGroupQuery;post:HostGroupQuery"),
		beego.NSRouter("/hosts", &BashBoardController{}, "get:HostsQueryByHostGroups;post:HostsQueryByHostGroups"),
		beego.NSRouter("/hostgroupscounters", &BashBoardController{}, "get:CounterQueryByHostGroup;post:CounterQueryByHostGroup"),
	)
	beego.AddNamespace(dashboard)
	beego.AddNamespace(hostgroup)
}
