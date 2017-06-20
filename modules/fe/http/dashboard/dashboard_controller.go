package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Cepave/open-falcon-backend/modules/fe/http/base"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/dashboard"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/uic"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/sys"
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

func (this *DashBoardController) CounterRegxQuery() {
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
	counters, err := dashboard.QueryCounterByNameRegx(queryStr, limitNum)
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	if len(counters) > 0 {
		baseResp.Data["counters"] = counters
	} else {
		baseResp.Data["counters"] = []string{}
	}
	this.ServeApiJson(baseResp)
	return
}

func gitLsRemote(gitRepo string, refs string) (string, error) {
	// This function depends on git command
	if resultStr, err := sys.CmdOut("git", "ls-remote", gitRepo, refs); err != nil {
		return "", err
	} else {
		// resultStr should be:
		// cb7a2998571cb25693867afcb24a7331f597768e        refs/heads/master
		strList := strings.Fields(resultStr)
		return strList[0], nil
	}
}

func (this *DashBoardController) LatestPlugin() {
	baseResp := this.BasicRespGen()
	if c, q_err := dashboard.QueryConfig("git_repo"); q_err != nil {
		log.Errorln("Query error: ", q_err)
		this.ResposeError(baseResp, "Error when getting git_repo address from database.")
		return
	} else {
		log.Debugln("git_repo address value is: ", c.Value)
		if hash, glr_err := gitLsRemote(c.Value, "refs/heads/master"); glr_err != nil {
			this.ResposeError(baseResp, "Error when git ls-remote GIT_ADDR")
			return
		} else {
			baseResp.Data["latestCommitHash"] = hash
		}
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
	metricQuery := this.GetString("metricQuery", "")
	counters, err := dashboard.QueryCounterByEndpoints(endpointsArr, limitNum, metricQuery)
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

// endpoint query by counter
func (this *DashBoardController) EndpointsQuery() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	counters := this.GetString("counters", "")
	rexstr, _ := regexp.Compile("^\\s*\\[\\s*|\\s*\\]\\s*$")
	countersArr := strings.Split(rexstr.ReplaceAllString(counters, ""), ",")
	limitNum, _ := this.GetInt("limit", 0)
	// We need a string pattern to filter outputs.
	// default filter pattern is .+
	filter := this.GetString("filter", ".+")
	if counters == "" {
		this.ResposeError(baseResp, "query string counters is empty, please check it")
		return
	}
	endpoints, err := dashboard.QueryEndpointsByCounter(countersArr, limitNum, filter)
	switch {
	case err != nil:
		this.ResposeError(baseResp, err.Error())
		return
	case len(endpoints) == 0:
		baseResp.Data["endpoints"] = []string{}
	default:
		baseResp.Data["endpoints"] = endpoints
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
		metricQuery := this.GetString("metricQuery", "")
		counters, err := dashboard.QueryCounterByEndpoints(endpoints, limitNum, metricQuery)
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

func (this *DashBoardController) EndpRunningPlugin() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}

	addr := this.GetString("addr", "")
	resp, AgentErr := http.Get(addr)
	baseResp.Data["requestAddr"] = addr
	log.Debugln("response from Agent: ", resp)
	log.Debugln("error message from Agent: ", AgentErr)
	if AgentErr != nil {
		baseResp.Data["errorFromAgent"] = AgentErr.Error()
	} else {
		defer resp.Body.Close()
		data := map[string]interface{}{}
		json.NewDecoder(resp.Body).Decode(&data)
		baseResp.Data["msgFromAgent"] = data["msg"]
		baseResp.Data["dataFromAgent"] = data["data"]
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *DashBoardController) EndpRegxquryForPlugin() {
	baseResp := this.BasicRespGen()
	session, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	var username *uic.User
	if session.Uid <= 0 {
		baseResp.Data["SessionFlag"] = true
		baseResp.Data["ErrorMsg"] = "Session is not vaild"
	} else {
		baseResp.Data["SessionFlag"] = false
		username = uic.SelectUserById(session.Uid)
		if username.Name != "root" {
			baseResp.Data["SessionFlag"] = true
			baseResp.Data["ErrorMsg"] = "You don't have permission to access this page"
		}
	}
	queryStr := ".+"
	if baseResp.Data["SessionFlag"] == false {
		enp, _ := dashboard.QueryEndpintByNameRegxForOps(queryStr)
		if len(enp) > 0 {
			baseResp.Data["Endpoints"] = enp
		} else {
			baseResp.Data["Endpoints"] = []string{}
		}
	}
	log.Debugln(baseResp)
	this.ServeApiJson(baseResp)
	return
}
