package uic

import (
	"github.com/Cepave/fe/http/base"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func ConfigRoutes() {

	//open-falcon's routes
	beego.Router("/root", &UserController{}, "get:CreateRoot")
	beego.Router("/auth/login", &AuthController{}, "get:LoginGet;post:LoginPost")
	beego.Router("/auth/login/:token", &AuthController{}, "get:LoginWithToken")
	beego.Router("/auth/third-party", &AuthController{}, "post:LoginThirdParty")
	beego.Router("/auth/register", &AuthController{}, "get:RegisterGet;post:RegisterPost")

	beego.Router("/sso/sig", &SsoController{}, "get:Sig")
	beego.Router("/sso/user/:sig:string", &SsoController{}, "get:User")
	beego.Router("/sso/logout/:sig:string", &SsoController{}, "get:Logout")

	beego.Router("/user/query", &UserController{}, "get:Query")
	beego.Router("/user/in", &UserController{}, "get:In")
	beego.Router("/user/qrcode/:id:int", &UserController{}, "get:QrCode")
	beego.Router("/about/:name:string", &UserController{}, "get:About")

	beego.Router("/team/users", &TeamController{}, "get:Users")
	beego.Router("/team/query", &TeamController{}, "get:Query")
	beego.Router("/team/all", &TeamController{}, "get:All")

	//owl-protal-routes
	apins := beego.NewNamespace("/api/v1/auth",
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSRouter("/register", &AuthApiController{}, "post:Register"),
		beego.NSRouter("/login", &AuthApiController{}, "post:Login"),
		beego.NSRouter("/sessioncheck", &AuthApiController{}, "get:AuthSession;post:AuthSession"),
		beego.NSRouter("/logout", &AuthApiController{}, "get:Logout;post:Logout"),
		beego.NSRouter("/user", &AuthApiController{}, "post:GetUser"),
		beego.NSRouter("/user/update", &AuthApiController{}, "put:UpdateUser;post:UpdateUser"),
		beego.NSRouter("/team/count", &AuthApiController{}, "get:CountNumOfTeam;post:CountNumOfTeam"),
	)
	beego.AddNamespace(apins)

	loginRequired :=
		beego.NewNamespace("/me",
			beego.NSCond(func(ctx *context.Context) bool {
				return true
			}),
			beego.NSBefore(base.FilterLoginUser),
			beego.NSRouter("/logout", &AuthController{}, "*:Logout"),
			beego.NSRouter("/info", &UserController{}, "get:Info"),
			beego.NSRouter("/profile", &UserController{}, "get:ProfileGet;post:ProfilePost"),
			beego.NSRouter("/chpwd", &UserController{}, "*:ChangePassword"),
			beego.NSRouter("/users", &UserController{}, "get:Users"),
			beego.NSRouter("/user/c", &UserController{}, "get:CreateUserGet;post:CreateUserPost"),
			beego.NSRouter("/teams", &TeamController{}, "get:Teams"),
			beego.NSRouter("/team/c", &TeamController{}, "get:CreateTeamGet;post:CreateTeamPost"),
		)

	beego.AddNamespace(loginRequired)

	targetUserRequired :=
		beego.NewNamespace("/target-user",
			beego.NSCond(func(ctx *context.Context) bool {
				return true
			}),
			beego.NSBefore(base.FilterLoginUser, base.FilterTargetUser),
			beego.NSRouter("/delete", &UserController{}, "*:DeleteUser"),
			beego.NSRouter("/edit", &UserController{}, "get:EditGet;post:EditPost"),
			beego.NSRouter("/chpwd", &UserController{}, "post:ResetPassword"),
			beego.NSRouter("/role", &UserController{}, "*:Role"),
		)

	beego.AddNamespace(targetUserRequired)

	targetTeamRequired :=
		beego.NewNamespace("/target-team",
			beego.NSCond(func(ctx *context.Context) bool {
				return true
			}),
			beego.NSBefore(base.FilterLoginUser, base.FilterTargetTeam),
			beego.NSRouter("/delete", &TeamController{}, "*:DeleteTeam"),
			beego.NSRouter("/edit", &TeamController{}, "get:EditGet;post:EditPost"),
		)

	beego.AddNamespace(targetTeamRequired)

}
