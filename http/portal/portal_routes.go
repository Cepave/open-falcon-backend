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
		beego.NSRouter("/eventcases/get", &PortalController{}, "get:EventCasesGet;post:EventCasesGet"),
		beego.NSRouter("/events/get", &PortalController{}, "get:EventGet;post:EventGet"),
		beego.NSRouter("/eventcases/close", &PortalController{}, "get:ColseCase;post:ColseCase;put:ColseCase"),
		beego.NSRouter("/eventcases/addnote", &PortalController{}, "get:AddNote;post:AddNote;put:AddNote"),
		beego.NSRouter("/eventcases/notes", &PortalController{}, "get:NotesGet;post:NotesGet;put:NotesGet"),
		beego.NSRouter("/eventcases/note", &PortalController{}, "get:GetNote;post:GetNote;put:GetNote"),
		beego.NSRouter("/tpl/count", &PortalController{}, "get:CountNumOfTlp;post:CountNumOfTlp"),
	)
	beego.AddNamespace(portal)
}
