package dashboard

import (
	"regexp"
	"strings"

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

//counter query by endpoints
func (this *BashBoardController) CounterQuery() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	endpoints := this.GetString("endpoints", "")
	endpointcheck, _ := regexp.Compile("^\\s*\\[\\s*\\]\\s*$")
	if endpoints == "" || endpointcheck.MatchString(endpoints) {
		this.ResposeError(baseResp, "query string is empty, please it")
		return
	}
	rexstr, _ := regexp.Compile("^\\s*\\[\\s*|\\s*\\]\\s*$")
	endpointsArr := strings.Split(rexstr.ReplaceAllString(endpoints, ""), ",")
	counters, err := dashboard.QueryCounterByEndpoints(endpointsArr)
	switch {
	case err != nil:
		this.ResposeError(baseResp, err.Error())
		return
	case len(counters) == 0:
		baseResp.Data["counters"] = []string{}
	default:
		baseResp.Data["counters"] = counters
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *BashBoardController) HostGroupQuery() {
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

	hostgroupList, err := dashboard.QueryHostGroupByNameRegx(queryStr)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	if len(hostgroupList) > 0 {
		baseResp.Data["hostgroups"] = hostgroupList
	} else {
		baseResp.Data["hostgroups"] = []string{}
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *BashBoardController) HostsQueryByHostGroup() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	hostgroup := this.GetString("hostgroup", "")
	if hostgroup == "" {
		this.ResposeError(baseResp, "query string is empty, please it")
		return
	}

	hosts, err := dashboard.GetHostsByHostGroupName(hostgroup)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	if len(hosts) > 0 {
		baseResp.Data["hosts"] = hosts
	} else {
		baseResp.Data["hosts"] = []string{}
	}
	this.ServeApiJson(baseResp)
	return
}
