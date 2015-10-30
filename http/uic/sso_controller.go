package uic

import (
	"github.com/Cepave/fe/http/base"
	"github.com/Cepave/fe/model/uic"
	"github.com/Cepave/fe/utils"
	"time"
)

type SsoController struct {
	base.BaseController
}

func (this *SsoController) Sig() {
	this.Ctx.Output.Body([]byte(utils.GenerateUUID()))
}

func (this *SsoController) User() {
	sig := this.Ctx.Input.Param(":sig")
	if sig == "" {
		this.NotFound("sig is blank")
		return
	}

	sessionObj := uic.ReadSessionBySig(sig)
	if sessionObj == nil {
		this.NotFound("no such sig")
		return
	}

	if int64(sessionObj.Expired) < time.Now().Unix() {
		uic.RemoveSessionByUid(sessionObj.Uid)
		this.SessionExpired()
		return
	}

	u := uic.ReadUserById(sessionObj.Uid)
	if u == nil {
		this.NotFound("no such user")
		return
	}

	this.Data["json"] = map[string]interface{}{
		"user": u,
	}
	this.ServeJSON()
}

func (this *SsoController) Logout() {
	sig := this.Ctx.Input.Param(":sig")
	if sig == "" {
		this.ServeErrJson("sig is blank")
		return
	}

	sessionObj := uic.ReadSessionBySig(sig)
	if sessionObj != nil {
		uic.RemoveSessionByUid(sessionObj.Uid)
	}

	this.ServeOKJson()
}
