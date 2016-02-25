package dashboard

import (
	"github.com/Cepave/fe/http/base"
	"github.com/Cepave/fe/model/dashboard"
)

type BashBoardController struct {
	base.BaseController
}

func (this *BashBoardController) EndpRegxqury() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	queryStr := this.GetString("queryStr", "")
	if queryStr == "" {
		this.ResposeError(baseResp, "query string is empty, please it")
		return
	}
	enp, err := dashboard.QueryEndpintByNameRegx(queryStr)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	if len(enp) > 0 {
		baseResp.Data["endpoints"] = enp
	} else {
		baseResp.Data["endpoints"] = []string{}
	}
	this.ServeApiJson(baseResp)
	return
}
