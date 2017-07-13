package boss

import (
	"github.com/Cepave/open-falcon-backend/modules/fe/http/base"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/boss"
)

type BossController struct {
	base.BaseController
}

func (this *BossController) GetPlaftom() {
	baseResp := this.BasicRespGen()
	res, err := boss.GetPlatformASJSON()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["res"] = res
	this.ServeApiJson(baseResp)
	return
}

func (this *BossController) GetContact() {
	baseResp := this.BasicRespGen()
	platform := this.GetString("platform", "")
	if platform == "" {
		this.ResposeError(baseResp, "platform is empty")
		return
	}
	res, err := boss.QueryContact(platform)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["res"] = res
	this.ServeApiJson(baseResp)
	return
}
