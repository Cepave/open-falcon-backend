package uic

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"github.com/gin-gonic/gin"
)

var db config.DBPool

const badstatus = http.StatusBadRequest

func Routes(r *gin.Engine) {
	db = config.Con()
	//session
	u := r.Group("/api/v1/user")
	u.GET("/auth_session", AuthSession)
	u.POST("/login", Login)
	u.GET("/logout", Logout)

	//user modify
	u.POST("/create", CreateUser)
	authapi := r.Group("/api/v1/user")
	authapi.Use(utils.AuthSessionMidd)
	authapi.GET("/current", UserInfo)
	authapi.GET("/u/:uid", GetUser)
	authapi.PUT("/update", UpdateUser)
	authapi.PUT("/cgpasswd", ChangePassword)
	authapi.GET("/users", UserList)
	adminapi := r.Group("/api/v1/admin")
	adminapi.Use(utils.AuthSessionMidd)
	adminapi.PUT("/change_user_role", ChangeRuleOfUser)
	adminapi.PUT("/change_user_passwd", AdminChangePassword)
	adminapi.DELETE("/delete_user", AdminUserDelete)

	//team
	authapi_team := r.Group("/api/v1")
	authapi_team.Use(utils.AuthSessionMidd)
	authapi_team.GET("/team", Teams)
	authapi_team.GET("/team/:team_id", GetTeam)
	authapi_team.POST("/team", CreateTeam)
	authapi_team.PUT("/team", UpdateTeam)
	authapi_team.DELETE("/team/:team_id", DeleteTeam)

	//third party
	third_party := r.Group("/api/v1")
	//third_party.GET("/third-party/forwarding", ForwardToBossLoginPage)
	//third_party.GET("/third-party/auth/login/:utoken", BossRedirectLogin)
	third_party.GET("/third-party/auth_session", GetBossUserInfoByCookie)
}
