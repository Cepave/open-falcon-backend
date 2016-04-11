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
		beego.NSRouter("/endpoints", &DashBoardController{}, "get:EndpRegxqury;post:EndpRegxqury"),
		beego.NSRouter("/endpointcounters", &DashBoardController{}, "get:CounterQuery;post:CounterQuery"),
	)
	hostgroup := beego.NewNamespace("/api/v1/hostgroup",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/query", &DashBoardController{}, "get:HostGroupQuery;post:HostGroupQuery"),
		beego.NSRouter("/hosts", &DashBoardController{}, "get:HostsQueryByHostGroups;post:HostsQueryByHostGroups"),
		beego.NSRouter("/hostgroupscounters", &DashBoardController{}, "get:CounterQueryByHostGroup;post:CounterQueryByHostGroup"),
		beego.NSRouter("/count", &DashBoardController{}, "get:CountNumOfHostGroup;post:CountNumOfHostGroup"),
	)
	beego.Router("/ops/endpoints", &DashBoardController{}, "get:EndpRegxquryForOps")
	beego.AddNamespace(dashboard)
	beego.AddNamespace(hostgroup)
}
