package portal

import (
	"github.com/Cepave/fe/http/base"
	event "github.com/Cepave/fe/model/falcon_portal"
)

type PortalController struct {
	base.BaseController
}

func (this *PortalController) EventCasesGet() {
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

	username := this.GetString("cName", "")
	limitNum, _ := this.GetInt("limit", 0)
	elimit, _ := this.GetInt("elimit", 0)
	events, err := event.GetEventCases(startTime, endTime, prioprity, status, limitNum, elimit, username)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["eventCases"] = events
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
	closedNote := this.GetString("closedNote", "")
	id := this.GetString("id", "xxx")
	switch {
	case id == "xxx":
		this.ResposeError(baseResp, "You dosen't pick any event id")
		return
	case closedNote == "":
		this.ResposeError(baseResp, "You can not skip closed note")
		return
	}
	err = event.CloseEvent(username, closedNote, id)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	this.ServeApiJson(baseResp)
	return
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

	limitNum, _ := this.GetInt("limit", 0)
	events, err := event.GetEvents(startTime, endTime, limitNum)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["events"] = events
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
