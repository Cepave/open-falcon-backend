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
	status := this.GetString("status", "ALL")
	processStatus := this.GetString("process_status", "ALL")
	metrics := this.GetString("metrics", "ALL")

	username := this.GetString("cName", "")
	limitNum, _ := this.GetInt("limit", 0)
	elimit, _ := this.GetInt("elimit", 0)
	events, err := event.GetEventCases(startTime, endTime, prioprity, status, processStatus, limitNum, elimit, username, metrics)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["eventCases"] = events
	this.ServeApiJson(baseResp)
	return
}

//will deprecated
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

func (this *PortalController) AddNote() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	username := this.GetString("cName", "")
	note := this.GetString("note", "")
	id := this.GetString("id", "xxx")
	status := this.GetString("status", "")
	caseId := this.GetString("caseId", "")
	switch {
	case id == "xxx":
		this.ResposeError(baseResp, "You dosen't pick any event id")
		return
	case note == "":
		this.ResposeError(baseResp, "You can not skip closed note")
		return
	}
	err = event.AddNote(username, note, id, status, caseId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) NotesGet() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	id := this.GetString("id", "xxx")
	limitNum, _ := this.GetInt("limit", 0)
	if id == "xxx" {
		this.ResposeError(baseResp, "You dosen't pick any event id")
		return
	}
	notes, err := event.GetNotes(id, limitNum)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["notes"] = notes
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) GetNote() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	id, _ := this.GetInt64("id", 0)
	if id == 0 {
		this.ResposeError(baseResp, "You dosen't pick any note id")
		return
	}
	note, err := event.GetNote(id)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["note"] = note
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
	status := this.GetString("status", "ALL")
	limitNum, _ := this.GetInt("limit", 0)
	events, err := event.GetEvents(startTime, endTime, status, limitNum)
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
