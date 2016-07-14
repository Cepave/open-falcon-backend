package portal

import events "github.com/Cepave/open-falcon-backend/modules/fe/model/falcon_portal"

func (this *PortalController) WhenStrategyUpdated() {
	baseResp := this.BasicRespGen()
	id, _ := this.GetInt("id", 0)
	if id == 0 {
		this.ResposeError(baseResp, "id is missing")
		return
	}
	err, resp := events.WhenStrategyUpdated(id)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["affectedRows"] = resp
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) WhenStrategyDeleted() {
	baseResp := this.BasicRespGen()
	id, _ := this.GetInt("id", 0)
	if id == 0 {
		this.ResposeError(baseResp, "id is missing")
		return
	}
	err, resp := events.WhenStrategyDeleted(id)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["affectedRows"] = resp
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) WhenTempleteDeleted() {
	baseResp := this.BasicRespGen()
	templateId, _ := this.GetInt("templateId", 0)
	if templateId == 0 {
		this.ResposeError(baseResp, "templateId is missing")
		return
	}
	err, resp := events.WhenTempleteDeleted(templateId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["affectedRows"] = resp
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) WhenTempleteUnbind() {
	baseResp := this.BasicRespGen()
	hostgroupId, _ := this.GetInt("hostgroupId", 0)
	templateId, _ := this.GetInt("templateId", 0)
	if templateId == 0 && hostgroupId == 0 {
		this.ResposeError(baseResp, "templateId or hostgroupId  is missing")
		return
	}
	err, resp := events.WhenTempleteUnbind(templateId, hostgroupId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["affectedRows"] = resp
	this.ServeApiJson(baseResp)
	return
}

func (this *PortalController) WhenEndpointUnbind() {
	baseResp := this.BasicRespGen()
	hostgroupId, _ := this.GetInt("hostgroupId", 0)
	hostId, _ := this.GetInt("hostId", 0)
	if hostId == 0 && hostgroupId == 0 {
		this.ResposeError(baseResp, "hostId or hostgroupId is missing")
		return
	}
	err, resp := events.WhenEndpointUnbind(hostId, hostgroupId)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	baseResp.Data["affectedRows"] = resp
	this.ServeApiJson(baseResp)
	return
}
