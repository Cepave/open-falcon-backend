package dashboard

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/base"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/dashboard"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/uic"
	"github.com/SlyMarbo/rss"
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

type xmlEntry struct {
	ID      string `xml:"id"`
	Updated string `xml:"updated"`
}

type xmlData struct {
	EntryList []xmlEntry `xml:"entry"`
}

func (this *DashBoardController) LatestPlugin() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}

	v := xmlData{}
	if resp, err := http.Get(g.Config().AtomAddr); err != nil {
		// handle error.
		log.Println("Error retrieving resource:", err)
	} else {
		defer resp.Body.Close()
		xml.NewDecoder(resp.Body).Decode(&v)
	}

	baseResp.Data["EntryList"] = v.EntryList
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

//endpoints query by counter
func (this *DashBoardController) EndpointQuery() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	}
	counter := this.GetString("counter", "")
	limitNum, _ := this.GetInt("limit", 0)
	metricQuery := this.GetString("metricQuery", "")
	negateMatch, _ := this.GetBool("negateMatch", false)
	if metricQuery == "" && counter == "" {
		this.ResposeError(baseResp, "query string && query pattern are both empty, please check it")
		return
	}
	endpoints, err := dashboard.QueryEndpointsByCounter(counter, limitNum, metricQuery, negateMatch)
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
	log.Println(baseResp)
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
		enpRow, _ := dashboard.QueryEndpintByNameRegxForOps(queryStr)
		enp := gitInfoAdapter(enpRow)
		if len(enp) > 0 {
			var ips []string
			this.Data["Endpoints"] = enp
			this.Data["Len"] = len(enp)
			for _, en := range enp {
				if en.Ip != "" {
					ips = append(ips, en.Ip)
				}
			}
			this.Data["IP"] = strings.Join(ips, ",")
		} else {
			this.Data["Endpoints"] = []string{}
			this.Data["Len"] = 0
			this.Data["IP"] = ""
		}
	}
	this.TplName = "dashboard/endpoints.html"
}

var commitsInfo []*rss.Item

func gitInfoAdapter(enpRow []dashboard.Hosts) (enp []dashboard.GitInfo) {
	feed, err := rss.Fetch(g.Config().AtomAddr)
	if err != nil {
		log.Println(err)
	}

	commitsInfo = append(commitsInfo, feed.Items...)
	log.Debugln("commit atom feed is:", feed.Items)
	log.Debugln("commitsInfo is:", commitsInfo)
	for _, host := range enpRow {
		gitInfo := dashboard.GitInfo{Hostname: host.Hostname,
			Ip:            host.Ip,
			AgentVersion:  host.AgentVersion,
			PluginVersion: host.PluginVersion,
			Valid:         false}
		for _, item := range commitsInfo {
			titleArray := strings.Split(item.ID, "/")
			hash := strings.TrimSpace(titleArray[len(titleArray)-1])
			if hash == host.PluginVersion {
				// copy Title and Date column
				gitInfo.Date = item.Date
				gitInfo.Title = item.Title
				gitInfo.Valid = true
				break
			}
		}
		enp = append(enp, gitInfo)
	}

	return
}
