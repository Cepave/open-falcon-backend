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
	alarmAdjust := beego.NewNamespace("/api/v1/alarmadjust",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/whenstrategyupdated", &PortalController{}, "get:WhenStrategyUpdated;post:WhenStrategyUpdated"),
		beego.NSRouter("/whentempleteunbind", &PortalController{}, "get:WhenTempleteUnbind;post:WhenTempleteUnbind"),
	)
	beego.AddNamespace(portal)
	beego.AddNamespace(alarmAdjust)
}
