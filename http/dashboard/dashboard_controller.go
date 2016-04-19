package dashboard

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Cepave/fe/g"
	"github.com/Cepave/fe/http/base"
	"github.com/Cepave/fe/model/dashboard"
	"github.com/Cepave/fe/model/uic"
)

type DashBoardController struct {
	base.BaseController
}

func (this *DashBoardController) EndpRegxqury() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	queryStr := this.GetString("queryStr", "")
	if queryStr == "" {
		this.ResposeError(baseResp, "query string is empty, please check it")
		return
	}
	limitNum, _ := this.GetInt("limit", 0)
	enp, err := dashboard.QueryEndpintByNameRegx(queryStr, limitNum)
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
func (this *DashBoardController) CounterQuery() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	endpoints := this.GetString("endpoints", "")
	endpointcheck, _ := regexp.Compile("^\\s*\\[\\s*\\]\\s*$")
	if endpoints == "" || endpointcheck.MatchString(endpoints) {
		this.ResposeError(baseResp, "query string is empty, please check it")
		return
	}
	rexstr, _ := regexp.Compile("^\\s*\\[\\s*|\\s*\\]\\s*$")
	endpointsArr := strings.Split(rexstr.ReplaceAllString(endpoints, ""), ",")
	limitNum, _ := this.GetInt("limit", 0)
	counters, err := dashboard.QueryCounterByEndpoints(endpointsArr, limitNum)
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

func (this *DashBoardController) HostGroupQuery() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	queryStr := this.GetString("queryStr", "")
	if queryStr == "" {
		this.ResposeError(baseResp, "query string is empty, please check it")
		return
	}
	limitNum, _ := this.GetInt("limit", 0)
	hostgroupList, err := dashboard.QueryHostGroupByNameRegx(queryStr, limitNum)
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

func (this *DashBoardController) HostsQueryByHostGroups() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}

	hostgroups := this.GetString("hostgroups", "")
	hostgroupscheck, _ := regexp.Compile("^\\s*\\[\\s*\\]\\s*$")
	if hostgroups == "" || hostgroupscheck.MatchString(hostgroups) {
		this.ResposeError(baseResp, "query string is empty, please check it")
		return
	}
	rexstr, _ := regexp.Compile("^\\s*\\[\\s*|\\s*\\]\\s*$")
	hostgroupsArr := strings.Split(rexstr.ReplaceAllString(hostgroups, ""), ",")
	hosts_resp, err := dashboard.GetHostsByHostGroupName(hostgroupsArr)

	if len(hosts_resp) > 0 {
		baseResp.Data["hosts"] = hosts_resp
	} else {
		baseResp.Data["hosts"] = []string{}
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *DashBoardController) CounterQueryByHostGroup() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}

	hostgroups := this.GetString("hostgroups", "")
	hostgroupscheck, _ := regexp.Compile("^\\s*\\[\\s*\\]\\s*$")
	if hostgroups == "" || hostgroupscheck.MatchString(hostgroups) {
		this.ResposeError(baseResp, "query string is empty, please check it")
		return
	}
	rexstr, _ := regexp.Compile("^\\s*\\[\\s*|\\s*\\]\\s*$")
	hostgroupsArr := strings.Split(rexstr.ReplaceAllString(hostgroups, ""), ",")

	hosts, err := dashboard.GetHostsByHostGroupName(hostgroupsArr)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	if len(hosts) > 0 {
		var endpoints []string
		for _, v := range hosts {
			endpoints = append(endpoints, fmt.Sprintf("\"%v\"", v.Hostname))
		}
		limitNum, _ := this.GetInt("limit", 0)
		counters, err := dashboard.QueryCounterByEndpoints(endpoints, limitNum)
		if err != nil {
			this.ResposeError(baseResp, err.Error())
			return
		} else if len(counters) > 0 {
			baseResp.Data["counters"] = counters
		} else {
			baseResp.Data["counters"] = []string{}
		}
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *DashBoardController) CountNumOfHostGroup() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()

	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	} else {
		numberOfteam, err := dashboard.CountNumOfHostGroup()
		if err != nil {
			this.ResposeError(baseResp, err.Error())
			return
		}
		baseResp.Data["count"] = numberOfteam
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *DashBoardController) EndpRegxquryForOps() {
	this.Data["Shortcut"] = g.Config().Shortcut
	sig := this.Ctx.GetCookie("sig")
	session := uic.ReadSessionBySig(sig)
	var username *uic.User
	if sig == "" || session.Uid <= 0 {
		this.Data["SessionFlag"] = true
		this.Data["ErrorMsg"] = "Session is not vaild"
	} else {
		this.Data["SessionFlag"] = false
		username = uic.SelectUserById(session.Uid)
		if username.Name != "root" {
			this.Data["SessionFlag"] = true
			this.Data["ErrorMsg"] = "You don't have permission to access this page"
		}
	}
	queryStr := this.GetString("queryStr", "")
	this.Data["QueryCondstion"] = queryStr
	if queryStr == "" || this.Data["SessionFlag"] == true {
		this.Data["Init"] = true
	} else {
		enp, _ := dashboard.QueryEndpintByNameRegxForOps(queryStr)
		if len(enp) > 0 {
			var ips []string
			this.Data["Endopints"] = enp
			this.Data["Len"] = len(enp)
			for _, en := range enp {
				if en.Ip != "" {
					ips = append(ips, en.Ip)
				}
			}
			this.Data["IP"] = strings.Join(ips, ",")
		} else {
			this.Data["Endopints"] = []string{}
			this.Data["Len"] = 0
			this.Data["IP"] = ""
		}
	}
	this.TplName = "dashboard/endpoints.html"
}
