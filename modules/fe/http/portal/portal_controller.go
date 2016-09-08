package portal

import (
	"github.com/Cepave/open-falcon-backend/modules/fe/http/base"
	event "github.com/Cepave/open-falcon-backend/modules/fe/model/falcon_portal"
)

type PortalController struct {
	base.BaseController
}

func (this *PortalController) GetEventCases() {
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
	caseId := this.GetString("caseId", "")
	includeEvents, _ := this.GetBool("includeEvents", false)
	events, err := event.GetEventCases(includeEvents, startTime, endTime, prioprity, status, processStatus, limitNum, elimit, username, metrics, caseId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["eventCases"] = events
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) GetEventCasesV2() {
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
	processStatus := this.GetString("process", "ALL")
	metrics := this.GetString("metric", "ALL")

	username := this.GetString("cName", "")
	limitNum, _ := this.GetInt("limit", 0)
	elimit, _ := this.GetInt("elimit", 0)
	caseId := this.GetString("caseId", "")
	showAll, _ := this.GetBool("show_all", false)
	includeEvents, _ := this.GetBool("includeEvents", false)
	events, err := event.GetEventCases(includeEvents, startTime, endTime, prioprity, status, processStatus, limitNum, elimit, username, metrics, caseId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	alertTmpStore, endpoints, err := event.AlertsConvert(events)
	alertswithInfo := event.GetAlertInfo(alertTmpStore, endpoints, showAll)
	alertswithNote := event.GetAlertsNotes(alertswithInfo)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["eventCases"] = alertswithNote
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) GetEventCasesV3() {
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
	processStatus := this.GetString("process", "ALL")
	metrics := this.GetString("metric", "ALL")

	username := this.GetString("cName", "")
	limitNum, _ := this.GetInt("limit", 0)
	elimit, _ := this.GetInt("elimit", 0)
	caseId := this.GetString("caseId", "")
	showAll, _ := this.GetBool("show_all", false)
	includeEvents, _ := this.GetBool("includeEvents", false)
	events, err := event.GetEventCases(includeEvents, startTime, endTime, prioprity, status, processStatus, limitNum, elimit, username, metrics, caseId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	alerts, endpoints, err := event.AlertsConvert(events)
	alerts2 := event.GetAlertInfoFromDB(alerts, endpoints, showAll)
	alerts3 := event.GetAlertsNotes(alerts2)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["eventCases"] = alerts3
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) GetEvent() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	startTime, _ := this.GetInt64("startTime", 0)
	endTime, _ := this.GetInt64("endTime", 0)
	status := this.GetString("status", "ALL")
	limit, _ := this.GetInt("limit", 0)
	caseId := this.GetString("caseId", "")
	events, err := event.GetEvents(startTime, endTime, status, limit, caseId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["events"] = events
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
	bossId := this.GetString("caseId", "")
	limit, _ := this.GetInt("limit", 0)
	switch {
	case id == "xxx":
		this.ResposeError(baseResp, "You dosen't pick any event id")
		return
	case note == "":
		this.ResposeError(baseResp, "You can not skip closed note")
		return
	case status == "":
		this.ResposeError(baseResp, "You can not skip status of note")
	}
	err = event.AddNote(username, note, id, status, bossId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	notes, err := event.GetNotes(id, limit, 0, 0, false)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["notes"] = notes
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) BatchUpdateNote() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	username := this.GetString("cName", "")
	note := this.GetString("note", "")
	ids := this.GetString("ids", "[]")
	status := this.GetString("status", "ignored")
	if status == "ignored" && note == "" {
		note = "ignored by ignored api."
	}
	bossId := this.GetString("caseIds", "")
	switch {
	case ids == "[]":
		this.ResposeError(baseResp, "You dosen't pick any event id")
		return
	case note == "":
		this.ResposeError(baseResp, "You can not skip closed note")
		return
	}
	err = event.AddNote(username, note, ids, status, bossId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) GetNotes() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	id := this.GetString("id", "xxx")
	limitNum, _ := this.GetInt("limit", 0)
	startTime, _ := this.GetInt64("startTime", 0)
	endTime, _ := this.GetInt64("endTime", 0)
	filterIgnored, _ := this.GetBool("filterIgnored", false)
	if id == "xxx" {
		this.ResposeError(baseResp, "You dosen't pick any event id")
		return
	}
	notes, err := event.GetNotes(id, limitNum, startTime, endTime, filterIgnored)
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
