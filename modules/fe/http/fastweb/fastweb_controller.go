package fastweb

import (
	"github.com/Cepave/open-falcon-backend/modules/fe/http/base"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/fastweb"
)

type FastWebController struct {
	base.BaseController
}

func (this *FastWebController) GetPlaftom() {
	baseResp := this.BasicRespGen()
	res, err := fastweb.GetPlatformASJSON()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["res"] = res
	this.ServeApiJson(baseResp)
	return
}

func (this *FastWebController) GetContact() {
	baseResp := this.BasicRespGen()
	platform := this.GetString("platform", "")
	if platform == "" {
		this.ResposeError(baseResp, "platfrom is empty")
		return
	}
	res, err := fastweb.QueryContact(platform)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["res"] = res
	this.ServeApiJson(baseResp)
	return
}
