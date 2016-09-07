package portal

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func ConfigRoutes() {
	//owl-protal-routes
	portal := beego.NewNamespace("/api/v1/portal",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/eventcases/get", &PortalController{}, "get:GetEventCases;post:GetEventCases"),
		beego.NSRouter("/events/get", &PortalController{}, "get:GetEvent;post:GetEvent"),
		beego.NSRouter("/eventcases/close", &PortalController{}, "get:ColseCase;post:ColseCase;put:ColseCase"),
		beego.NSRouter("/eventcases/addnote", &PortalController{}, "get:AddNote;post:AddNote;put:AddNote"),
		beego.NSRouter("/eventcases/notes", &PortalController{}, "get:GetNotes;post:GetNotes;put:GetNotes"),
		beego.NSRouter("/eventcases/note", &PortalController{}, "get:GetNote;post:GetNote;put:GetNote"),
		beego.NSRouter("/tpl/count", &PortalController{}, "get:CountNumOfTlp;post:CountNumOfTlp"),
	)

	portalv2 := beego.NewNamespace("/api/v2/portal",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/eventcases/get", &PortalController{}, "get:GetEventCasesV3;post:GetEventCasesV3"),
		beego.NSRouter("/eventcases/feed", &PortalController{}, "get:OnTimeFeeding;post:OnTimeFeeding"),
	)

	alarmAdjust := beego.NewNamespace("/api/v1/alarmadjust",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/whenstrategyupdated", &PortalController{}, "get:WhenStrategyUpdated;post:WhenStrategyUpdated"),
		beego.NSRouter("/whenstrategydeleted", &PortalController{}, "get:WhenStrategyDeleted;post:WhenStrategyDeleted"),
		beego.NSRouter("/whentempletedeleted", &PortalController{}, "get:WhenTempleteDeleted;post:WhenTempleteDeleted"),
		beego.NSRouter("/whentempleteunbind", &PortalController{}, "get:WhenTempleteUnbind;post:WhenTempleteUnbind"),
		beego.NSRouter("/whenendpointunbind", &PortalController{}, "get:WhenEndpointUnbind;post:WhenEndpointUnbind"),
		beego.NSRouter("/whenendpointonmaintain", &PortalController{}, "get:WhenEndpointOnMaintain;post:WhenEndpointOnMaintain"),
	)

	beego.AddNamespace(portal)
	beego.AddNamespace(portalv2)
	beego.AddNamespace(alarmAdjust)
}
