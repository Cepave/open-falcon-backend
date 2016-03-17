package portal

import (
	"github.com/Cepave/fe/http/base"
	"github.com/Cepave/fe/model/event"
)

type PortalController struct {
	base.BaseController
}

func (this *PortalController) EventGet() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	startTime, _ := this.GetInt64("startTime", 0)
	endTime, _ := this.GetInt64("endTime", 0)
	prioprity, _ := this.GetInt("prioprity", -1)
	status := this.GetString("status", "PROBLEM")

	events, err := event.GetEvent(startTime, endTime, prioprity, status)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["events"] = events
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) ColseCase() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	username := this.GetString("cName", "")
	id := this.GetString("id", "xxx")
	if id == "xxx" {
		this.ResposeError(baseResp, "You dosen't pick any event id")
		return
	}
	err = event.CloseEvent(username, id)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) CountNumOfTlp() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()

	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	} else {
		numberOfteam, err := event.CountNumOfTlp()
		if err != nil {
			this.ResposeError(baseResp, err.Error())
			return
		}
		baseResp.Data["count"] = numberOfteam
	}
	this.ServeApiJson(baseResp)
	return
}
