package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	"github.com/Cepave/open-falcon-backend/modules/query/graph"
	"github.com/Cepave/open-falcon-backend/modules/query/proc"
	log "github.com/sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	"github.com/bitly/go-simplejson"
)

type Tag struct {
	StrategyId int
	Name       string
	Value      string
	CreateAt   string
	UpdateAt   string
}

type Remark struct {
	User string
	Sig  string
	Host string
	Text string
}

/**
 * @function name:   func postByJson(rw http.ResponseWriter, req *http.Request, url string)
 * @description:     This function sends a POST request in JSON format.
 * @related issues:  OWL-171
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @param:           url string
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/12/2015
 * @last modified:   11/13/2015
 * @called by:       func queryInfo(rw http.ResponseWriter, req *http.Request)
 *                   func queryHistory(rw http.ResponseWriter, req *http.Request)
 */
func postByJSON(rw http.ResponseWriter, req *http.Request, url string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	s := buf.String()
	reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
	if err != nil {
		log.Errorf("Error = %v", err.Error())
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		log.Errorf("Error = %v", err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

/**
 * @function name:   func queryInfo(rw http.ResponseWriter, req *http.Request)
 * @description:     This function handles /graph/info API request.
 * @related issues:  OWL-171
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/12/2015
 * @last modified:   11/13/2015
 * @called by:       func configApiRoutes()
 */
func queryInfo(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Query + "/graph/info"
	postByJSON(rw, req, url)
}

/**
 * @function name:   func queryHistory(rw http.ResponseWriter, req *http.Request)
 * @description:     This function handles /graph/history API request.
 * @related issues:  OWL-171
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/12/2015
 * @last modified:   11/13/2015
 * @called by:       func configApiRoutes()
 */
func queryHistory(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Query + "/graph/history"
	postByJSON(rw, req, url)
}

/**
 * @function name:   func getRequest(rw http.ResponseWriter, url string)
 * @description:     This function sends GET request to given URL.
 * @related issues:  OWL-159
 * @param:           rw http.ResponseWriter
 * @param:           url string
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/24/2015
 * @last modified:   11/24/2015
 * @called by:       func dashboardEndpoints(rw http.ResponseWriter, req *http.Request)
 *                    in query/http/api.go
 * @called by:       func dashboardEndpoints(rw http.ResponseWriter, req *http.Request)
 *                    in query/http/api.go
 */
func getRequest(rw http.ResponseWriter, url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Error = %v", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error = %v", err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

/**
 * @function name:   func dashboardEndpoints(rw http.ResponseWriter, req *http.Request)
 * @description:     This function handles /api/endpoints API request.
 * @related issues:  OWL-159, OWL-171
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/12/2015
 * @last modified:   11/24/2015
 * @called by:       func configApiRoutes()
 */
func dashboardEndpoints(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Dashboard + req.URL.RequestURI()
	getRequest(rw, url)
}

/**
 * @function name:   func postByForm(rw http.ResponseWriter, req *http.Request, url string)
 * @description:     This function sends a POST request in Form format.
 * @related issues:  OWL-171
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @param:           url string
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/12/2015
 * @last modified:   11/13/2015
 * @called by:       func dashboardCounters(rw http.ResponseWriter, req *http.Request)
 *                   func dashboardChart(rw http.ResponseWriter, req *http.Request)
 */
func postByForm(rw http.ResponseWriter, req *http.Request, url string) {
	req.ParseForm()
	client := &http.Client{}
	resp, err := client.PostForm(url, req.PostForm)
	if err != nil {
		log.Errorf("Error = %v", err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

/**
 * @function name:   func dashboardCounters(rw http.ResponseWriter, req *http.Request)
 * @description:     This function handles /api/counters API request.
 * @related issues:  OWL-171
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/13/2015
 * @last modified:   11/13/2015
 * @called by:       func configApiRoutes()
 */
func dashboardCounters(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Dashboard + "/api/counters"
	postByForm(rw, req, url)
}

/**
 * @function name:   func dashboardChart(rw http.ResponseWriter, req *http.Request)
 * @description:     This function handles /api/chart API request.
 * @related issues:  OWL-171
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/13/2015
 * @last modified:   11/13/2015
 * @called by:       func configApiRoutes()
 */
func dashboardChart(rw http.ResponseWriter, req *http.Request) {
	url := g.Config().Api.Dashboard + "/chart"
	postByForm(rw, req, url)
}

func getAgentAliveData(hostnames []string, versions map[string]string, result map[string]interface{}) []cmodel.GraphLastResp {
	var queries []cmodel.GraphLastParam
	o := orm.NewOrm()
	var hosts []*Host
	_, err := o.Raw("SELECT hostname, agent_version FROM falcon_portal.host ORDER BY hostname ASC").QueryRows(&hosts)
	if err != nil {
		setError(err.Error(), result)
	} else {
		for _, host := range hosts {
			var query cmodel.GraphLastParam
			if !strings.Contains(host.Hostname, ".") && strings.Contains(host.Hostname, "-") {
				hostnames = append(hostnames, host.Hostname)
				versions[host.Hostname] = host.Agent_version
				query.Endpoint = host.Hostname
				query.Counter = "agent.alive"
				queries = append(queries, query)
			}
		}
	}
	s, err := json.Marshal(queries)
	if err != nil {
		setError(err.Error(), result)
	}
	url := g.Config().Api.Query + "/graph/last"
	reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
	if err != nil {
		setError(err.Error(), result)
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		setError(err.Error(), result)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := []cmodel.GraphLastResp{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		setError(err.Error(), result)
	}
	return data
}

func processAgentAliveData(data []cmodel.GraphLastResp, hostnames []string, versions map[string]string, result map[string]interface{}) {
	name := ""
	version := ""
	status := ""
	alive := 0
	countOfNormal := 0
	countOfWarn := 0
	countOfDead := 0
	anomalies := []interface{}{}
	items := []interface{}{}
	for key, row := range data {
		name = row.Endpoint
		var diff int64
		diff = 0
		var timestamp int64
		timestamp = 0
		status = "dead"
		alive = 0
		if name == "" {
			name = hostnames[key]
		} else {
			alive = int(row.Value.Value)
			timestamp = row.Value.Timestamp
			now := time.Now().Unix()
			diff = now - timestamp
		}
		version = versions[name]
		if alive > 0 {
			if diff > 3600 {
				status = "warn"
				countOfWarn++
			} else {
				status = "normal"
				countOfNormal++
			}
		} else {
			countOfDead++
		}
		item := map[string]interface{}{}
		item["id"] = strconv.Itoa(key + 1)
		item["hostname"] = name
		item["agent_version"] = version
		item["alive"] = alive
		item["timestamp"] = timestamp
		item["diff"] = diff
		item["status"] = status
		items = append(items, item)
		if diff > 60*60*24 && timestamp > 0 {
			anomalies = append(anomalies, item)
		}
	}
	var count = make(map[string]interface{})
	count["all"] = len(data)
	count["normal"] = countOfNormal
	count["warn"] = countOfWarn
	count["dead"] = countOfDead
	result["count"] = count
	result["items"] = items
}

func getAlive(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors

	data := []cmodel.GraphLastResp{}
	hostnames := []string{}
	var versions = make(map[string]string)
	data = getAgentAliveData(hostnames, versions, result)
	processAgentAliveData(data, hostnames, versions, result)
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func setStrategyTags(rw http.ResponseWriter, req *http.Request) {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}

	var nodes = make(map[string]interface{})
	nodes, _ = json.Map()
	strategyID := ""
	tagName := ""
	tagValue := ""
	if value, ok := nodes["strategyID"]; ok {
		strategyID = value.(string)
		delete(nodes, "strategyID")
	}
	if value, ok := nodes["tagName"]; ok {
		tagName = value.(string)
		delete(nodes, "tagName")
	}
	if value, ok := nodes["tagValue"]; ok {
		tagValue = value.(string)
		delete(nodes, "tagValue")
	}

	if len(strategyID) > 0 && len(tagName) > 0 && len(tagValue) > 0 {
		strategyIDint, err := strconv.Atoi(strategyID)
		if err != nil {
			setError(err.Error(), result)
		} else {
			o := orm.NewOrm()
			var tag Tag
			sqlcmd := "SELECT * FROM falcon_portal.tags WHERE strategy_id=?"
			err = o.Raw(sqlcmd, strategyIDint).QueryRow(&tag)
			if err == orm.ErrNoRows {
				sql := "INSERT INTO tags(strategy_id, name, value, create_at) VALUES(?, ?, ?, ?)"
				res, err := o.Raw(sql, strategyIDint, tagName, tagValue, getNow()).Exec()
				if err != nil {
					setError(err.Error(), result)
				} else {
					num, _ := res.RowsAffected()
					log.Debugf("mysql row affected nums = %v", num)
					result["strategyID"] = strategyID
					result["action"] = "create"
				}
			} else if err != nil {
				setError(err.Error(), result)
			} else {
				sql := "UPDATE tags SET name = ?, value = ? WHERE strategy_id = ?"
				res, err := o.Raw(sql, tagName, tagValue, strategyIDint).Exec()
				if err != nil {
					setError(err.Error(), result)
				} else {
					num, _ := res.RowsAffected()
					log.Debugf("mysql row affected nums = %v", num)
					result["strategyID"] = strategyID
					result["action"] = "update"
				}
			}
		}
	} else {
		setError("Input value errors.", result)
	}
	nodes["result"] = result
	setResponse(rw, nodes)
}

func getTemplateStrategies(rw http.ResponseWriter, req *http.Request) {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}
	countOfStrategies := 0
	arguments := strings.Split(req.URL.Path, "/")
	if arguments[len(arguments)-1] == "strategies" {
		templateID, err := strconv.Atoi(arguments[len(arguments)-2])
		if err != nil {
			setError(err.Error(), result)
		}
		o := orm.NewOrm()
		var strategyIDs []int64
		num, err := o.Raw("SELECT id FROM falcon_portal.strategy WHERE tpl_id = ? ORDER BY id ASC", templateID).QueryRows(&strategyIDs)
		if err != nil {
			setError(err.Error(), result)
		} else if num > 0 {
			countOfStrategies = int(num)
			var strategies = make(map[string]interface{})
			sids := ""
			for key, strategyID := range strategyIDs {
				sid := strconv.Itoa(int(strategyID))
				item := map[string]string{}
				item["templateID"] = strconv.Itoa(templateID)
				item["strategyID"] = sid
				strategies[sid] = item
				if key == 0 {
					sids = sid
				} else {
					sids += ", " + sid
				}
			}
			sqlcmd := "SELECT strategy_id, name, value FROM falcon_portal.tags WHERE strategy_id IN ("
			sqlcmd += sids
			sqlcmd += ") ORDER BY strategy_id ASC"
			var tags []*Tag
			_, err = o.Raw(sqlcmd).QueryRows(&tags)
			if err != nil {
				setError(err.Error(), result)
			} else {
				for _, tag := range tags {
					strategyID := strconv.Itoa(int(tag.StrategyId))
					strategy := strategies[strategyID].(map[string]string)
					strategy["tagName"] = tag.Name
					strategy["tagValue"] = tag.Value
					strategies[strategyID] = strategy
				}
			}
			for _, strategy := range strategies {
				items = append(items, strategy)
			}
		}
	}
	result["items"] = items
	result["count"] = countOfStrategies
	var nodes = make(map[string]interface{})
	nodes["result"] = result
	setResponse(rw, nodes)
}

func getPlatformJSON(nodes map[string]interface{}, result map[string]interface{}) {
	fcname := g.Config().Api.Name
	fctoken := getFctoken()
	url := g.Config().Api.Map + "/fcname/" + fcname + "/fctoken/" + fctoken
	url += "/show_active/yes/hostname/yes/pop_id/yes/ip/yes/show_ip_type/yes.json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		setError(err.Error(), result)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		setError(err.Error(), result)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &nodes); err != nil {
		setError(err.Error(), result)
	}
}

func queryHostsData(result map[string]interface{}) []map[string]string {
	hosts := []map[string]string{}
	var rows []orm.Params
	o := orm.NewOrm()
	o.Using("boss")
	sql := "SELECT hostname, activate, platform, ip FROM `boss`.`hosts`"
	sql += " WHERE exist = 1 AND platform != '' ORDER BY hostname ASC"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
		return hosts
	} else if num > 0 {
		for _, row := range rows {
			hostname := row["hostname"].(string)
			IP := row["ip"].(string)
			if IP != "" && IP == getIPFromHostname(hostname, result) {
				host := map[string]string{
					"hostname": row["hostname"].(string),
					"platform": row["platform"].(string),
					"ip":       row["ip"].(string),
					"activate": row["activate"].(string),
				}
				hosts = append(hosts, host)
			}
		}
	}
	return hosts
}

func queryIPsData(result map[string]interface{}) []map[string]string {
	IPs := []map[string]string{}
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT ip, hostname, platform, status FROM `boss`.`ips`"
	sql += " WHERE exist = 1 AND hostname != '' AND platform != ''"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
		return IPs
	} else if num > 0 {
		for _, row := range rows {
			item := map[string]string{
				"hostname": row["hostname"].(string),
				"ip":       row["ip"].(string),
				"platform": row["platform"].(string),
				"activate": row["status"].(string),
			}
			IPs = append(IPs, item)
		}
	}
	return IPs
}

func mergeIPsOfHost(data []map[string]string, result map[string]interface{}) (map[string][]map[string]string, []string, []string) {
	platforms := map[string][]map[string]string{}
	platformNames := []string{}
	hostnames := []string{}
	platformName := ""
	for _, host := range data {
		hostname := host["hostname"]
		platformName = host["platform"]
		hostnames = appendUniqueString(hostnames, hostname)
		platformNames = appendUniqueString(platformNames, platformName)
		if platform, ok := platforms[platformName]; ok {
			platform = append(platform, host)
			platforms[platformName] = platform
		} else {
			platforms[platformName] = []map[string]string{
				host,
			}
		}
	}
	sort.Strings(hostnames)
	sort.Strings(platformNames)

	hostsMap := map[string]string{}
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT hostname, activate FROM `boss`.`hosts`"
	sql += " WHERE hostname IN ('" + strings.Join(hostnames, "','")
	sql += "') AND exist = 1"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			hostname := row["hostname"].(string)
			activate := row["activate"].(string)
			hostsMap[hostname] = activate
		}
	}
	for _, groupName := range platformNames {
		platform := []map[string]string{}
		itemsMap := map[string][]map[string]string{}
		itemNames := []string{}
		group := platforms[groupName]
		for _, agent := range group {
			hostname := agent["hostname"]
			itemNames = appendUniqueString(itemNames, hostname)
			if item, ok := itemsMap[hostname]; ok {
				item = append(item, agent)
			} else {
				itemsMap[hostname] = []map[string]string{
					agent,
				}
			}
		}
		for _, itemName := range itemNames {
			slice := itemsMap[itemName]
			index := 0
			for key, item := range slice {
				hostname := item["hostname"]
				ip := item["ip"]
				if ip == getIPFromHostname(hostname, result) {
					index = key
				}
			}
			host := slice[index]
			if val, ok := hostsMap[itemName]; ok {
				host["activate"] = val
			}
			platform = append(platform, host)
		}
		platforms[groupName] = platform
	}
	return platforms, platformNames, hostnames
}

func setGraphQueries(hostnames []string, hostnamesExisted []string, result map[string]interface{}) (queries []*cmodel.GraphLastParam, versions map[string]map[string]string) {
	o := orm.NewOrm()
	o.Using("default")
	var hosts []*Host
	versions = make(map[string]map[string]string)
	hostnamesStr := strings.Join(hostnames, "','")
	sqlcommand := "SELECT hostname, agent_version, plugin_version FROM `falcon_portal`.`host` WHERE hostname IN ('"
	sqlcommand += hostnamesStr + "') ORDER BY hostname ASC"
	_, err := o.Raw(sqlcommand).QueryRows(&hosts)
	if err != nil {
		setError(err.Error(), result)
	} else {
		for _, host := range hosts {
			var query cmodel.GraphLastParam
			if !strings.Contains(host.Hostname, ".") && strings.Contains(host.Hostname, "-") {
				hostnamesExisted = append(hostnamesExisted, host.Hostname)
				version := map[string]string{
					"agent":  host.Agent_version,
					"plugin": host.Plugin_version,
				}
				versions[host.Hostname] = version
				query.Endpoint = host.Hostname
				query.Counter = "agent.alive"
				queries = append(queries, &query)
			}
		}
	}
	return
}

func queryAgentAlive(queries []*cmodel.GraphLastParam, reqHost string, result map[string]interface{}) []cmodel.GraphLastResp {
	data := []cmodel.GraphLastResp{}
	if len(queries) > 0 {
		if strings.Index(g.Config().Api.Query, reqHost) >= 0 {
			proc.LastRequestCnt.Incr()
			for _, param := range queries {
				if param == nil {
					continue
				}
				last, err := graph.Last(*param)
				if err != nil {
					setError("graph.last fail, err: "+err.Error(), result)
					return data
				}
				if last == nil {
					continue
				}
				data = append(data, *last)
			}
			proc.LastRequestItemCnt.IncrBy(int64(len(data)))
		} else {
			s, err := json.Marshal(queries)
			if err != nil {
				setError(err.Error(), result)
			}
			url := g.Config().Api.Query + "/graph/last"
			reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
			if err != nil {
				setError(err.Error(), result)
			}
			reqPost.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(reqPost)
			if err != nil {
				setError(err.Error(), result)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, &data)
			if err != nil {
				setError(err.Error(), result)
			}
		}
	}
	return data
}

func classifyAgentAliveResponse(data []cmodel.GraphLastResp, hostnamesExisted []string, versions map[string]map[string]string) (out_versions map[string]map[string]string) {
	name := ""
	status := ""
	alive := 0
	var diff int64
	var timestamp int64
	out_versions = make(map[string]map[string]string)
	for key, row := range data {
		name = row.Endpoint
		alive = 0
		diff = 0
		timestamp = 0
		status = "error"
		if name == "" {
			name = hostnamesExisted[key]
		} else {
			alive = int(row.Value.Value)
			timestamp = row.Value.Timestamp
			now := time.Now().Unix()
			diff = now - timestamp
		}
		version := versions[name]
		if alive > 0 {
			if diff > 3600 {
				status = "warn"
			} else {
				status = "normal"
			}
		}
		item := map[string]string{
			"status":  status,
			"version": version["agent"],
			"plugin":  version["plugin"],
		}
		out_versions[name] = item
	}
	return
}

func getAnomalies(errorHosts []map[string]string, result map[string]interface{}) map[string]interface{} {
	anomalies := map[string]interface{}{}
	provinces := map[string][]map[string]string{}
	provinceNames := []string{}
	for _, errorHost := range errorHosts {
		provinceName := errorHost["province"]
		if provinceName == "特区" {
			provinceName = errorHost["city"]
		}
		if province, ok := provinces[provinceName]; ok {
			province = append(province, errorHost)
			provinces[provinceName] = province
		} else {
			provinces[provinceName] = []map[string]string{
				errorHost,
			}
			provinceNames = append(provinceNames, provinceName)
		}
	}
	sort.Strings(provinceNames)
	for _, provinceName := range provinceNames {
		province := provinces[provinceName]
		anomalies[provinceName] = map[string]interface{}{
			"count": len(province),
			"hosts": province,
		}
	}
	return anomalies
}

func completeAgentAliveData(groups map[string][]map[string]string, groupNames []string, versions map[string]map[string]string, result map[string]interface{}) {
	errorHosts := []map[string]string{}
	platforms := []interface{}{}
	count := map[string]int{}
	countOfNormalSum := 0
	countOfWarnSum := 0
	countOfErrorSum := 0
	countOfMissSum := 0
	countOfDeactivatedSum := 0
	hostId := 1
	hostname := ""
	activate := ""
	version := ""
	plugin := ""
	status := ""
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT idc, city, province FROM `boss`.`hosts` WHERE hostname = ?"
	for _, groupName := range groupNames {
		platform := map[string]interface{}{}
		hosts := []interface{}{}
		count := map[string]int{}
		countOfNormal := 0
		countOfWarn := 0
		countOfError := 0
		countOfMiss := 0
		countOfDeactivated := 0
		group := groups[groupName]
		for _, agent := range group {
			hostname = agent["hostname"]
			activate = agent["activate"]
			version = ""
			plugin = ""
			if item, ok := versions[hostname]; ok {
				version = item["version"]
				plugin = item["plugin"]
			}
			status = ""
			if activate == "1" {
				if item, ok := versions[hostname]; ok {
					status = item["status"]
				} else {
					status = "miss"
					countOfMiss++
				}
			} else {
				status = "deactivated"
				countOfDeactivated++
			}
			if status == "normal" {
				countOfNormal++
			} else if status == "warn" {
				countOfWarn++
			} else if status == "error" {
				countOfError++
			}
			host := map[string]string{
				"id":       strconv.Itoa(hostId),
				"name":     hostname,
				"platform": groupName,
				"status":   status,
				"ip":       agent["ip"],
				"version":  version,
				"plugin":   plugin,
			}
			if host["status"] == "error" {
				num, err := o.Raw(sql, hostname).Values(&rows)
				if err != nil {
					setError(err.Error(), result)
				} else if num > 0 {
					row := rows[0]
					host["idc"] = row["idc"].(string)
					host["city"] = row["city"].(string)
					host["province"] = row["province"].(string)
				}
				errorHosts = append(errorHosts, host)
			}
			hosts = append(hosts, host)
			hostId++
		}
		count["normal"] = countOfNormal
		count["warn"] = countOfWarn
		count["error"] = countOfError
		count["miss"] = countOfMiss
		count["deactivated"] = countOfDeactivated
		count["all"] = countOfNormal + countOfWarn + countOfError + countOfMiss + countOfDeactivated
		platform["platformName"] = groupName
		platform["platformCount"] = count
		platform["hosts"] = hosts
		platforms = append(platforms, platform)
		countOfNormalSum += countOfNormal
		countOfWarnSum += countOfWarn
		countOfErrorSum += countOfError
		countOfMissSum += countOfMiss
		countOfDeactivatedSum += countOfDeactivated
	}
	count["normal"] = countOfNormalSum
	count["warn"] = countOfWarnSum
	count["error"] = countOfErrorSum
	count["miss"] = countOfMissSum
	count["deactivated"] = countOfDeactivatedSum
	count["all"] = countOfNormalSum + countOfWarnSum + countOfErrorSum + countOfMissSum + countOfDeactivatedSum
	result["count"] = count
	result["anomalies"] = getAnomalies(errorHosts, result)
	result["items"] = platforms
}

func getAgentAlivePlatforms(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	data := queryHostsData(result)
	platforms := map[string][]map[string]string{}
	platformNames := []string{}
	hostnames := []string{}
	for _, host := range data {
		hostname := host["hostname"]
		platformName := host["platform"]
		hostnames = appendUniqueString(hostnames, hostname)
		platformNames = appendUniqueString(platformNames, platformName)
		if platform, ok := platforms[platformName]; ok {
			platform = append(platform, host)
			platforms[platformName] = platform
		} else {
			platforms[platformName] = []map[string]string{
				host,
			}
		}
	}
	hostnamesExisted := []string{}
	queries, versions := setGraphQueries(hostnames, hostnamesExisted, result)
	agentAliveData := queryAgentAlive(queries, req.Host, result)
	status_versions := classifyAgentAliveResponse(agentAliveData, hostnamesExisted, versions)
	completeAgentAliveData(platforms, platformNames, status_versions, result)
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getMetricsByMetricType(metricType string) []string {
	metrics := []string{}
	if metricType == "bandwidths" {
		metrics = []string{
			"net.if.in.bits/iface=eth_all",
			"net.if.out.bits/iface=eth_all",
		}
	} else if metricType == "cpu" {
		metrics = []string{
			"cpu.idle",
			"cpu.system",
			"cpu.softirq",
			"cpu.user",
		}
	} else if metricType == "resources" {
		metrics = []string{
			"cpu.idle",
			"disk.io.util.max",
			"load.1min",
			"mem.memfree.percent",
			"mem.swapused.percent",
			"vmstat.procs.b",
			"vmstat.procs.r",
		}
	} else if metricType == "services" {
		metrics = []string{
			"http.response.time",
			"https.response.time",
			"ss.close.wait",
			"ss.established",
			"ss.syn.recv",
			"vfcc.squid.response.time",
		}
	} else if metricType == "all" || metricType == "aggregate" {
		metrics = []string{
			"net.if.in.bits/iface=eth_all",
			"net.if.out.bits/iface=eth_all",
			"cpu.idle",
			"disk.io.util.max",
			"load.1min",
			"mem.memfree.percent",
			"mem.swapused.percent",
			"http.response.time",
			"https.response.time",
			"ss.close.wait",
			"ss.established",
			"ss.syn.recv",
			"vfcc.squid.response.time",
			"vmstat.procs.b",
			"vmstat.procs.r",
		}
	}
	return metrics
}

func getDisksIOMetrics(hostname string, metricType string) []string {
	metrics := []string{}
	keyword := ""
	if metricType == "disks" {
		keyword = "disk.device.used.percent/device=%"
	} else if metricType == "io" {
		keyword = "disk.io.util/device=%"
	} else {
		return metrics
	}
	o := orm.NewOrm()
	var rows []orm.Params
	endpointID := ""
	sqlcmd := "SELECT id FROM `graph`.`endpoint` WHERE endpoint = ? LIMIT 1"
	num, err := o.Raw(sqlcmd, hostname).Values(&rows)
	if err != nil {
		log.Errorf("Error = %v", err.Error())
	} else if num > 0 {
		endpointID = rows[0]["id"].(string)
		sqlcmd = "SELECT counter FROM `graph`.`endpoint_counter` WHERE endpoint_id = ? AND counter LIKE ?"
		num, err = o.Raw(sqlcmd, endpointID, keyword).Values(&rows)
		if err != nil {
			log.Errorf("Error = %v", err.Error())
		} else if num > 0 {
			for _, row := range rows {
				metrics = append(metrics, row["counter"].(string))
			}
		}
	}
	return metrics
}

func getTimestampFromString(timeString string, result map[string]interface{}) int64 {
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		location = time.Local
	}
	if timeString == "" {
		return time.Now().In(location).Unix()
	}
	timestamp := int64(0)
	timestamp, err = strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		timeFormat := "2006-01-02 15:04:05"
		date, err := time.ParseInLocation(timeFormat, timeString, location)
		if err != nil {
			setError(err.Error(), result)
		} else {
			timestamp = date.Unix()
		}
	}
	return timestamp
}

func convertDurationToPoint(duration string, result map[string]interface{}) (int64, int64) {
	timestampFrom := int64(0)
	timestampTo := int64(0)
	if strings.Index(duration, ",") > -1 {
		from := strings.Split(duration, ",")[0]
		to := strings.Split(duration, ",")[1]
		timestampFrom = getTimestampFromString(from, result)
		timestampTo = getTimestampFromString(to, result)
		if timestampFrom >= timestampTo {
			log.Debugf("duration = %v", duration)
			log.Debugf("timestampFrom = %v", timestampFrom)
			log.Debugf("timestampTo = %v", timestampTo)
			setError("Value of timestampFrom should be less than value of timestampTo.", result)
		}
		return timestampFrom, timestampTo
	} else if strings.Index(duration, "d") > -1 || strings.Index(duration, "min") > -1 {
		unit := ""
		seconds := int64(0)
		if strings.Index(duration, "d") > -1 {
			unit = "d"
			seconds = int64(86400)
		} else {
			unit = "min"
			seconds = int64(60)
		}
		multiplier, err := strconv.Atoi(strings.Split(duration, unit)[0])
		if err != nil {
			setError(err.Error(), result)
		}
		offset := int64(multiplier) * seconds
		now := time.Now().Unix()
		timestampFrom = now - offset
		timestampTo = now + int64(5*60)
	}
	return timestampFrom, timestampTo
}

func getGraphQueryResponse(hostnames []string, metrics []string, duration string, consolidationFunction string, step int, result map[string]interface{}) []*cmodel.GraphQueryResponse {
	data := []*cmodel.GraphQueryResponse{}
	start, end := convertDurationToPoint(duration, result)
	for _, hostname := range hostnames {
		for _, metric := range metrics {
			request := cmodel.GraphQueryParam{
				Start:     start,
				End:       end,
				ConsolFun: consolidationFunction,
				Endpoint:  hostname,
				Counter:   metric,
			}
			if step >= 60 {
				request.Step = step
			}
			response, err := graph.QueryOne(request)
			if err != nil {
				log.Debugf("request = %v", request)
				log.Debugf("response = %v", response)
				log.Debugf("graph.queryOne fail, %v", request)
				setError("graph.queryOne fail, "+err.Error(), result)
				return data
			}
			data = append(data, response)
		}
	}
	return data
}

func getHostMetricValues(rw http.ResponseWriter, req *http.Request) {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}
	arguments := strings.Split(req.URL.Path, "/")
	hostname := ""
	metricType := ""
	duration := ""
	if len(arguments) == 6 {
		hostname = arguments[len(arguments)-3]
		metricType = arguments[len(arguments)-2]
		duration = arguments[len(arguments)-1]
	} else if len(arguments) == 5 {
		hostname = arguments[len(arguments)-2]
		metricType = arguments[len(arguments)-1]
		duration = "3d"
	}
	metrics := getMetricsByMetricType(metricType)
	if len(metrics) > 0 && strings.Index(duration, "d") > -1 {
		data := getGraphQueryResponse([]string{hostname}, metrics, duration, "AVERAGE", 1, result)

		for _, series := range data {
			values := []interface{}{}
			for _, rrdObj := range series.Values {
				value := []interface{}{
					rrdObj.Timestamp * 1000,
					rrdObj.Value,
				}
				values = append(values, value)
			}
			item := map[string]interface{}{
				"host":   series.Endpoint,
				"metric": series.Counter,
				"data":   values,
			}
			items = append(items, item)
		}
	}
	result["items"] = items
	var nodes = make(map[string]interface{})
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getExistedHostnames(hostnames []string, result map[string]interface{}) []string {
	hostnamesExisted := []string{}
	o := orm.NewOrm()
	var hosts []*Host
	hostnamesStr := strings.Join(hostnames, "','")
	sqlcommand := "SELECT hostname, agent_version FROM falcon_portal.host WHERE hostname IN ('"
	sqlcommand += hostnamesStr + "') ORDER BY hostname ASC"
	_, err := o.Raw(sqlcommand).QueryRows(&hosts)
	if err != nil {
		setError(err.Error(), result)
	} else {
		for _, host := range hosts {
			if !strings.Contains(host.Hostname, ".") && strings.Contains(host.Hostname, "-") {
				hostnamesExisted = append(hostnamesExisted, host.Hostname)
			}
		}
	}
	return hostnamesExisted
}

func appendUniqueString(slice []string, s string) []string {
	if len(s) == 0 {
		return slice
	}
	existed := false
	for _, val := range slice {
		if s == val {
			existed = true
		}
	}
	if !existed {
		slice = append(slice, s)
		sort.Strings(slice)
	}
	return slice
}

func appendUnique(slice []int, num int) []int {
	existed := false
	for _, element := range slice {
		if element == num {
			existed = true
		}
	}
	if !existed {
		slice = append(slice, num)
		sort.Ints(slice)
	}
	return slice
}

func getExistedHosts(hosts []map[string]string, hostnamesExisted []string, result map[string]interface{}) map[string]map[string]string {
	hostsExisted := map[string]map[string]string{}
	for key, hostname := range hostnamesExisted {
		host := map[string]string{
			"id":   strconv.Itoa(key + 1),
			"name": hostname,
		}
		hostsExisted[hostname] = host
	}
	for _, host := range hosts {
		hostname := host["name"]
		if _, ok := hostsExisted[hostname]; ok {
			hostExisted := hostsExisted[hostname]
			isp := strings.Split(hostname, "-")[0]
			province := strings.Split(hostname, "-")[1]
			hostExisted["isp"] = isp
			hostExisted["province"] = province
			hostExisted["idc"] = host["idc"]
			hostExisted["platform"] = host["platform"]
			hostsExisted[hostname] = hostExisted
		}
	}
	return hostsExisted
}

func getHostsLocations(hosts []map[string]string, hostnamesInput []string, result map[string]interface{}) ([]map[string]string, []string) {
	hostnames := []string{}
	hostsMap := map[string]map[string]string{}
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	hostnameStr := strings.Join(hostnamesInput, "','")
	sqlcmd := "SELECT hostname, idc, isp, province, city FROM `boss`.`hosts`"
	sqlcmd += " WHERE hostname IN ('" + hostnameStr + "')"
	sqlcmd += " AND exist = 1"
	sqlcmd += " ORDER BY hostname ASC"
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			hostname := row["hostname"].(string)
			if _, ok := hostsMap[hostname]; !ok {
				provinceCode := ""
				slice := strings.Split(hostname, "-")
				if len(slice) >= 4 {
					provinceCode = slice[1]
				}
				if len(provinceCode) > 3 {
					provinceCode = ""
				}
				host := map[string]string{
					"idc":          row["idc"].(string),
					"isp":          row["isp"].(string),
					"province":     row["province"].(string),
					"provinceCode": provinceCode,
					"city":         row["city"].(string),
				}
				hostsMap[hostname] = host
				hostnames = append(hostnames, hostname)
			}
		}
	}
	for key, host := range hosts {
		hostname := host["hostname"]
		host["name"] = hostname
		if val, ok := hostsMap[hostname]; ok {
			host["idc"] = val["idc"]
			host["isp"] = val["isp"]
			host["province"] = val["province"]
			host["provinceCode"] = val["provinceCode"]
			host["city"] = val["city"]
		}
		delete(host, "hostname")
		hosts[key] = host
	}
	return hosts, hostnames
}

func completeApolloFiltersData(hostsInput []map[string]string, result map[string]interface{}) {
	hosts := map[string]map[string]string{}
	keywords := map[string][]string{}
	for key, host := range hostsInput {
		id := strconv.Itoa(key + 1)
		platform := host["platform"]
		tags := []string{}
		tags = appendUniqueString(tags, platform)
		if _, ok := keywords[platform]; ok {
			keywords[platform] = appendUniqueString(keywords[platform], id)
		} else {
			keywords[platform] = []string{id}
		}
		host["platform"] = strings.Join(tags, ",")
		host["tag"] = strings.Join(tags, ",")
		delete(host, "activate")
		delete(host, "city")
		delete(host, "platform")
		delete(host, "isp")
		delete(host, "province")
		delete(host, "provinceCode")
		hosts[id] = host
	}
	result["hosts"] = hosts
	result["keywords"] = keywords
}

func getApolloFilters(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	count := 0
	data := queryIPsData(result)
	platforms, _, hostnames := mergeIPsOfHost(data, result)
	hosts := []map[string]string{}
	for _, platform := range platforms {
		for _, host := range platform {
			hosts = append(hosts, host)
		}
	}
	hosts, hostnames = getHostsLocations(hosts, hostnames, result)
	count = len(hosts)
	completeApolloFiltersData(hosts, result)
	nodes["count"] = count
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func addNetworkSpeedAndBondMode(metrics []string, hostnames []string, result map[string]interface{}) []string {
	metrics = append(metrics, "nic.bond.mode")
	o := orm.NewOrm()
	var rows []orm.Params
	hostname := strings.Join(hostnames, "','")
	sqlcmd := "SELECT id FROM graph.endpoint WHERE endpoint IN ('" + hostname + "')"
	num, err := o.Raw(sqlcmd).Values(&rows)
	endpointIDs := []string{}
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			endpointIDs = append(endpointIDs, row["id"].(string))
		}
	}
	tags := []string{}
	if len(endpointIDs) > 0 {
		sqlcmd = "SELECT DISTINCT tag FROM graph.tag_endpoint WHERE endpoint_id IN ('"
		sqlcmd += strings.Join(endpointIDs, "','") + "') AND tag LIKE 'device=%'"
		num, err := o.Raw(sqlcmd).Values(&rows)
		if err != nil {
			setError(err.Error(), result)
		} else if num > 0 {
			for _, row := range rows {
				tag := row["tag"].(string)
				if strings.Index(tag, "bond") > -1 || strings.Index(tag, "eth") > -1 {
					tags = append(tags, tag)
				}
			}
		}
	}
	for _, tag := range tags {
		speedMetric := "nic.default.out.speed/" + tag
		metrics = append(metrics, speedMetric)
	}
	return metrics
}

func addRecentData(data []*cmodel.GraphQueryResponse, dataRecent []*cmodel.GraphQueryResponse) []*cmodel.GraphQueryResponse {
	for key, item := range data {
		hostname := item.Endpoint
		metric := item.Counter
		latest := int64(0)
		if len(item.Values) > 0 {
			values := []*cmodel.RRDData{}
			for _, pair := range item.Values {
				if !math.IsNaN(float64(pair.Value)) && pair.Value > 0 {
					latest = pair.Timestamp
					values = append(values, pair)
				}
			}
			for _, itemRecent := range dataRecent {
				if itemRecent.Endpoint == hostname && itemRecent.Counter == metric {
					for _, pair := range itemRecent.Values {
						if pair.Timestamp > latest && !math.IsNaN(float64(pair.Value)) && pair.Value > 0 {
							values = append(values, pair)
						}
					}
				}
			}
			item.Values = values
		}
		data[key] = item
	}
	return data
}

func getApolloCharts(rw http.ResponseWriter, req *http.Request) {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}
	arguments := strings.Split(req.URL.Path, "/")
	metricType := arguments[4]
	hostnames := strings.Split(arguments[5], ",")
	duration := "1d"
	if len(arguments) > 6 {
		duration = arguments[6]
	}
	metrics := []string{}
	if metricType == "customized" {
		metrics = strings.Split(req.URL.Query()["metrics"][0], ",")
	} else {
		metrics = getMetricsByMetricType(metricType)
		if metricType == "bandwidths" {
			metrics = append(metrics, "nic.bond.mode")
			metrics = append(metrics, "nic.default.out.speed")
		}
	}
	data := []*cmodel.GraphQueryResponse{}
	step := 1
	start, end := convertDurationToPoint(duration, result)
	diff := math.Abs(float64(end - start))
	if diff > 3600 {
		step = 1200
	}
	if metricType == "disks" || metricType == "io" {
		dataOfHost := []*cmodel.GraphQueryResponse{}
		for _, hostname := range hostnames {
			metrics = getDisksIOMetrics(hostname, metricType)
			if diff < 300 {
				dataOfHost = getGraphLastData(metrics, []string{hostname}, result)
			} else {
				dataOfHost = getGraphQueryData(metrics, duration, []string{hostname}, step, result)
			}
			for _, series := range dataOfHost {
				data = append(data, series)
			}
		}
	} else {
		if diff < 300 {
			data = getGraphLastData(metrics, hostnames, result)
		} else {
			data = getGraphQueryData(metrics, duration, hostnames, step, result)
		}
	}
	for _, series := range data {
		metric := series.Counter
		if strings.Index(metric, "nic.default.out.speed") > -1 {
			if len(series.Values) > 0 && series.Values[0].Value > 0 {
				metric = "net.transmission.limit.80%"
				limit := series.Values[0].Value
				if series.Values[len(series.Values)-1].Value > 0 {
					limit = series.Values[len(series.Values)-1].Value
				}
				for _, item := range series.Values {
					value := item.Value
					if value > limit {
						limit = value
						break
					}
				}
				limit *= 1000 * 1000 * 0.8
				for key, _ := range series.Values {
					series.Values[key].Value = limit
				}
			} else {
				metric = ""
			}
		}
		if metric != "" {
			values := []interface{}{}
			for _, rrdObj := range series.Values {
				value := []interface{}{
					rrdObj.Timestamp * 1000,
					rrdObj.Value,
				}
				values = append(values, value)
			}
			item := map[string]interface{}{
				"host":   series.Endpoint,
				"metric": metric,
				"data":   values,
			}
			items = append(items, item)
		}
	}
	result["items"] = items
	var nodes = make(map[string]interface{})
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getApolloRemark(rw http.ResponseWriter, req *http.Request) {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := map[string]map[string]string{}
	arguments := strings.Split(req.URL.Path, "/")
	hostnames := strings.Split(arguments[4], ",")
	var rows []orm.Params
	o := orm.NewOrm()
	o.Using("apollo")
	sql := "SELECT remark, account, name, email, updated FROM `apollo`.`remarks`"
	sql += " WHERE hostname = ? ORDER BY updated DESC LIMIT 1"
	for _, hostname := range hostnames {
		num, err := o.Raw(sql, hostname).Values(&rows)
		if err != nil {
			setError(err.Error(), result)
		} else if num > 0 {
			row := rows[0]
			remark := map[string]string{
				"remark":  row["remark"].(string),
				"account": row["account"].(string),
				"name":    row["name"].(string),
				"email":   row["email"].(string),
				"updated": row["updated"].(string),
			}
			items[hostname] = remark
		}
	}
	result["items"] = items
	var nodes = make(map[string]interface{})
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func addApolloRemark(rw http.ResponseWriter, req *http.Request) {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}

	decoder := json.NewDecoder(req.Body)
	var remark Remark
	err := decoder.Decode(&remark)
	if err != nil {
		setError(err.Error(), result)
	}
	defer req.Body.Close()
	user := map[string]string{
		"account": remark.User,
		"userID":  "",
		"name":    "",
		"email":   "",
	}
	login := false

	o := orm.NewOrm()
	o.Using("default")
	var rows []orm.Params
	sql := "SELECT id, cnname, email FROM `uic`.`user` WHERE name = ? LIMIT 1"
	num, err := o.Raw(sql, remark.User).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		user["userID"] = rows[0]["id"].(string)
		user["name"] = rows[0]["cnname"].(string)
		user["email"] = rows[0]["email"].(string)
	}
	sql = "SELECT expired FROM `uic`.`session` WHERE sig = ? ORDER BY expired DESC LIMIT 1"
	num, err = o.Raw(sql, remark.Sig).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		expired := rows[0]["expired"].(string)
		expiredTimestamp, err := strconv.Atoi(expired)
		if err != nil {
			setError(err.Error(), result)
		} else {
			if int64(expiredTimestamp) > time.Now().Unix() {
				login = true
			}
		}
	}
	if login {
		text := remark.Text
		text = strings.Replace(text, "\r", " ", -1)
		text = strings.Replace(text, "\n", " ", -1)
		text = strings.Replace(text, "\t", " ", -1)
		text = strings.TrimSpace(text)
		re_inside_whiteSpaces := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
		text = re_inside_whiteSpaces.ReplaceAllString(text, " ")
		if len(text) > 500 {
			text = string([]rune(text)[0:250])
		}
		o.Using("apollo")
		sql = "INSERT INTO `apollo`.`remarks`(`hostname`, `remark`, `userid`,"
		sql += "`account`, `name`, `email`, `updated`) VALUES(?, ?, ?, ?, ?, ?, ?)"
		_, err := o.Raw(sql, remark.Host, text, user["userID"], remark.User,
			user["name"], user["email"], getNow()).Exec()
		if err != nil {
			setError(err.Error(), result)
		} else {
			sql = "SELECT remark, account, name, email, updated "
			sql += "FROM `apollo`.`remarks` WHERE hostname = ? "
			sql += "ORDER BY updated DESC LIMIT 1"
			num, err := o.Raw(sql, remark.Host).Values(&rows)
			if err != nil {
				setError(err.Error(), result)
			} else if num > 0 {
				item := map[string]string{}
				item["hostname"] = remark.Host
				item["remark"] = rows[0]["remark"].(string)
				item["account"] = rows[0]["account"].(string)
				item["name"] = rows[0]["name"].(string)
				item["email"] = rows[0]["email"].(string)
				item["updated"] = rows[0]["updated"].(string)
				items = append(items, item)
			}
		}
	} else {
		result["error"] = append(result["error"].([]string), "Please log in first.")
	}
	result["items"] = items
	var nodes = make(map[string]interface{})
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getIPFromHostname(hostname string, result map[string]interface{}) string {
	ip := ""
	fragments := strings.Split(hostname, "-")
	if len(fragments) == 6 {
		slice := []string{}
		fragments = fragments[2:]
		for _, fragment := range fragments {
			num, err := strconv.Atoi(fragment)
			if err != nil {
				setError(err.Error(), result)
			} else {
				slice = append(slice, strconv.Itoa(num))
			}
		}
		if len(slice) == 4 {
			ip = strings.Join(slice, ".")
		}
	}
	return ip
}

func getBandwidthsAverage(metricType string, duration string, hostnames []string, result map[string]interface{}) []interface{} {
	items := []interface{}{}
	sort.Strings(hostnames)
	metrics := getMetricsByMetricType(metricType)
	hostsMap := map[string]interface{}{}
	if len(metrics) > 0 && len(hostnames) > 0 {
		data := getGraphQueryResponse(hostnames, metrics, duration, "AVERAGE", 1, result)
		for _, series := range data {
			values := []interface{}{}
			for _, rrdObj := range series.Values {
				if !math.IsNaN(float64(rrdObj.Value)) {
					value := []interface{}{
						float64(rrdObj.Timestamp),
						float64(rrdObj.Value),
					}
					values = append(values, value)
				}
			}
			average := float64(0)
			sum := float64(0)
			divider := float64(0)
			timestamp := float64(0)
			for _, value := range values {
				timestamp = value.([]interface{})[0].(float64)
				sum += value.([]interface{})[1].(float64)
				divider++
			}
			if divider > 0 {
				average = sum / divider
			}
			item := map[string]interface{}{
				"host":             series.Endpoint,
				"ip":               getIPFromHostname(series.Endpoint, result),
				"net.in.bits.avg":  0,
				"net.out.bits.avg": 0,
				"time":             "",
			}
			if value, ok := hostsMap[series.Endpoint]; ok {
				item = value.(map[string]interface{})
			}
			if series.Counter == "net.if.in.bits/iface=eth_all" {
				item["net.in.bits.avg"] = int(average)
			} else if series.Counter == "net.if.out.bits/iface=eth_all" {
				item["net.out.bits.avg"] = int(average)
			}
			if item["time"] == "" && average > 0 {
				item["time"] = time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
			}
			hostsMap[series.Endpoint] = item
		}
	}
	for _, hostname := range hostnames {
		host := hostsMap[hostname]
		if host.(map[string]interface{})["host"].(string) != "" {
			items = append(items, host)
		}
	}
	return items
}

func getPlatformBandwidthsFiveMinutesAverage(platformName string, metricType string, rw http.ResponseWriter) map[string]interface{} {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	duration := "6min"
	var nodes = make(map[string]interface{})
	getPlatformJSON(nodes, result)
	hostnames := []string{}
	if nodes["status"] != nil && int(nodes["status"].(float64)) == 1 {
		hostname := ""
		for _, platform := range nodes["result"].([]interface{}) {
			groupName := platform.(map[string]interface{})["platform"].(string)
			if groupName == platformName {
				for _, device := range platform.(map[string]interface{})["ip_list"].([]interface{}) {
					hostname = device.(map[string]interface{})["hostname"].(string)
					ip := device.(map[string]interface{})["ip"].(string)
					if len(ip) > 0 && ip == getIPFromHostname(hostname, result) {
						if device.(map[string]interface{})["ip_status"].(string) == "1" {
							hostnames = append(hostnames, hostname)
						}
					}
				}
			}
		}
	}
	items := getBandwidthsAverage(metricType, duration, hostnames, result)
	if _, ok := nodes["info"]; ok {
		delete(nodes, "info")
	}
	if _, ok := nodes["status"]; ok {
		delete(nodes, "status")
	}
	result["items"] = items
	nodes["result"] = result
	nodes["count"] = len(items)
	nodes["platform"] = platformName
	return nodes
}

func getPlatformContact(platformName string, nodes map[string]interface{}) {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	var platformMap = make(map[string]interface{})
	fcname := g.Config().Api.Name
	fctoken := getFctoken()
	url := g.Config().Api.Contact
	params := map[string]string{
		"fcname":       fcname,
		"fctoken":      fctoken,
		"platform_key": platformName,
	}
	s, err := json.Marshal(params)
	if err != nil {
		setError(err.Error(), result)
	}
	reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
	if err != nil {
		setError(err.Error(), result)
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		setError(err.Error(), result)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &nodes)
		if err != nil {
			setError(err.Error(), result)
		} else if nodes["status"] != nil && int(nodes["status"].(float64)) == 1 {
			roles := []string{
				"principal",
				"backuper",
				"upgrader",
			}
			for _, name := range strings.Split(platformName, ",") {
				if platform, ok := nodes["result"].(map[string]interface{})[name].(map[string]interface{}); ok {
					items := map[string]map[string]string{}
					for _, role := range roles {
						if value, ok := platform[role].([]interface{}); ok {
							person := value[0]
							item := map[string]string{
								"name":  person.(map[string]interface{})["realname"].(string),
								"phone": person.(map[string]interface{})["cell"].(string),
								"email": person.(map[string]interface{})["email"].(string),
							}
							if role == "backuper" {
								items["deputy"] = item
							} else {
								items[role] = item
							}
						}
					}
					platformMap[name] = items
				}
			}
		}
	}
	if _, ok := nodes["info"]; ok {
		delete(nodes, "info")
	}
	if _, ok := nodes["status"]; ok {
		delete(nodes, "status")
	}
	result["items"] = platformMap
	nodes["result"] = result
	nodes["count"] = len(platformMap)
	nodes["platform"] = platformName
}

func getPlatforms(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []map[string]string{}
	var rows []orm.Params
	o := orm.NewOrm()
	o.Using("boss")
	sql := "SELECT platform, contacts FROM `boss`.`platforms`"
	sql += " ORDER BY platform ASC"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			item := map[string]string{
				"platform": row["platform"].(string),
				"contacts": row["contacts"].(string),
			}
			items = append(items, item)
		}
	}
	result["items"] = items
	nodes["result"] = result
	nodes["count"] = len(items)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getPlatformBandwidthsByISP(platformName string, duration string, nodes map[string]interface{}) {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}
	hostsMap := map[string]string{}
	ISPs := map[string][]string{}
	ISPNames := []string{
		"all",
		"cmb",
		"cnc",
		"ctl",
		"oth",
	}
	var rows []orm.Params
	o := orm.NewOrm()
	o.Using("boss")
	hostnames := []string{}
	sql := "SELECT DISTINCT hostname FROM `boss`.`ips`"
	sql += " WHERE platform = ? AND exist = 1 ORDER BY hostname ASC"
	num, err := o.Raw(sql, platformName).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			hostnames = append(hostnames, row["hostname"].(string))
		}
	}
	sql = "SELECT hostname, platform, isp FROM `boss`.`hosts`"
	sql += " WHERE hostname IN ('" + strings.Join(hostnames, "','")
	sql += "') AND exist = 1 ORDER BY isp ASC"
	num, err = o.Raw(sql).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			hostname := row["hostname"].(string)
			ISP := row["isp"].(string)
			if ISP == "cmb" || ISP == "cnc" || ISP == "ctl" {
				hostsMap[hostname] = ISP
				ISPs[ISP] = append(ISPs[ISP], hostname)
			} else {
				hostsMap[hostname] = "oth"
				ISPs["oth"] = append(ISPs["oth"], hostname)
			}
			ISPs["all"] = append(ISPs["all"], hostname)
		}
	}
	hostnames = ISPs["all"]
	metrics := getMetricsByMetricType("bandwidths")
	start, end := convertDurationToPoint(duration, result)
	diff := end - start
	data := getGraphQueryData(metrics, duration, hostnames, 1200, result)
	if len(data) > 0 {
		index := -1
		max := 0
		timestampLatest := int64(0)
		for key, item := range data {
			if max < len(item.Values) {
				max = len(item.Values)
				index = key
			}
			if len(item.Values) > 0 {
				timestamp := item.Values[len(item.Values)-1].Timestamp
				if timestampLatest < timestamp {
					timestampLatest = timestamp
				}
			}
		}
		unit := 5
		if diff < 1200 {
			unit = 1
		}
		tickers := []string{}
		tickersMap := map[string]float64{}
		if index >= 0 {
			for _, rrdObj := range data[index].Values {
				ticker := getTicker(rrdObj.Timestamp, unit)
				if _, ok := tickersMap[ticker]; !ok {
					if len(ticker) > 0 {
						tickersMap[ticker] = float64(0)
						tickers = append(tickers, ticker)
					}
				}
			}
		}
		ISPsMap := map[string]map[string]map[string]float64{}
		for _, ISP := range ISPNames {
			ISPsMap[ISP] = map[string]map[string]float64{}
			for _, metric := range metrics {
				ISPsMap[ISP][metric] = map[string]float64{}
			}
		}
		for _, series := range data {
			hostname := series.Endpoint
			ISP := hostsMap[hostname]
			metric := series.Counter
			for _, rrdObj := range series.Values {
				value := float64(rrdObj.Value)
				if !math.IsNaN(value) {
					timestamp := rrdObj.Timestamp
					ticker := getNearestTicker(float64(timestamp), tickers)
					if len(ticker) > 0 {
						if _, ok := ISPsMap[ISP][metric][ticker]; ok {
							ISPsMap[ISP][metric][ticker] += value
						} else {
							ISPsMap[ISP][metric][ticker] = value
						}
						if _, ok := ISPsMap["all"][metric][ticker]; ok {
							ISPsMap["all"][metric][ticker] += value
						} else {
							ISPsMap["all"][metric][ticker] = value
						}
					}
				}
			}
		}
		for _, ISP := range ISPNames {
			for _, metric := range metrics {
				max := float64(0)
				for _, ticker := range tickers {
					if value, ok := ISPsMap[ISP][metric][ticker]; ok {
						if max < value {
							max = value
						}
					}
				}
				threshold := max * 0.02
				data := [][]float64{}
				for _, ticker := range tickers {
					if value, ok := ISPsMap[ISP][metric][ticker]; ok {
						timestamp := getTimestampFromTicker(ticker)
						if value > threshold {
							datum := []float64{
								float64(timestamp * 1000),
								value,
							}
							data = append(data, datum)
						}
					}
				}
				if len(data) > 4 {
					value := data[len(data)-3][1]
					if (value < (data[len(data)-4][1])*0.9) && (value <= (data[len(data)-2][1])) {
						value = (data[len(data)-4][1] + data[len(data)-2][1]) / 2
						data[len(data)-3][1] = value
					}
				}
				if len(data) > 3 {
					value := data[len(data)-2][1]
					if (value < (data[len(data)-3][1])*0.9) && (value <= (data[len(data)-1][1])) {
						value = (data[len(data)-3][1] + data[len(data)-1][1]) / 2
						data[len(data)-2][1] = value
					}
				}
				item := map[string]interface{}{
					"data":       data,
					"hostsCount": len(ISPs[ISP]),
					"ISP":        ISP,
					"metric":     metric,
					"platform":   platformName,
				}
				items = append(items, item)
			}
		}
	}

	year, month, day := time.Now().Date()
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		loc = time.Local
	}
	today := time.Date(year, month, day, 0, 0, 0, 0, loc).Format("2006-01-02")
	today += "%"
	meansMap := map[string]int{}
	deviationsMap := map[string]int{}
	tickers := []string{}
	o.Using("apollo")
	sql = "SELECT DATE_FORMAT(date, '%Y-%m-%d %H:%i'), mean, deviation FROM `apollo`.`deviations`"
	sql += " WHERE platform = ? AND metric = 1 AND date LIKE ? ORDER BY date ASC"
	num, err = o.Raw(sql, platformName, today).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if (num > 0) && (diff > 600) {
		for _, row := range rows {
			ticker := row["DATE_FORMAT(date, '%Y-%m-%d %H:%i')"].(string)
			mean := 0
			value, err := strconv.Atoi(row["mean"].(string))
			if err != nil {
				setError(err.Error(), result)
			} else {
				mean = value
			}
			deviation := 0
			value, err = strconv.Atoi(row["deviation"].(string))
			if err != nil {
				setError(err.Error(), result)
			} else {
				deviation = value
			}
			meansMap[ticker] = mean
			deviationsMap[ticker] = deviation
			tickers = append(tickers, ticker)
		}
		oneSigma := map[string]int{}
		twoSigma := map[string]int{}
		threeSigma := map[string]int{}
		for _, ticker := range tickers {
			mean := meansMap[ticker]
			deviation := deviationsMap[ticker]
			oneSigma[ticker] = mean + deviation
			twoSigma[ticker] = mean + (deviation * 2)
			threeSigma[ticker] = mean + (deviation * 3)
		}
		oneSigmaSeries := [][]int{}
		twoSigmaSeries := [][]int{}
		threeSigmaSeries := [][]int{}
		for _, ticker := range tickers {
			timestamp := getTimestampFromTicker(ticker)
			oneSigmaValue := []int{
				int(timestamp * 1000),
				oneSigma[ticker],
			}
			twoSigmaValue := []int{
				int(timestamp * 1000),
				twoSigma[ticker],
			}
			threeSigmaValue := []int{
				int(timestamp * 1000),
				threeSigma[ticker],
			}
			oneSigmaSeries = append(oneSigmaSeries, oneSigmaValue)
			twoSigmaSeries = append(twoSigmaSeries, twoSigmaValue)
			threeSigmaSeries = append(threeSigmaSeries, threeSigmaValue)
		}
		item := map[string]interface{}{
			"data":       oneSigmaSeries,
			"hostsCount": len(ISPs["all"]),
			"ISP":        "all",
			"metric":     "sigma.x1",
			"platform":   platformName,
		}
		items = append(items, item)
		item = map[string]interface{}{
			"data":       twoSigmaSeries,
			"hostsCount": len(ISPs["all"]),
			"ISP":        "all",
			"metric":     "sigma.x2",
			"platform":   platformName,
		}
		items = append(items, item)
		item = map[string]interface{}{
			"data":       threeSigmaSeries,
			"hostsCount": len(ISPs["all"]),
			"ISP":        "all",
			"metric":     "sigma.x3",
			"platform":   platformName,
		}
		items = append(items, item)
	}

	result["items"] = items
	nodes["result"] = result
	nodes["count"] = len(items)
	nodes["platform"] = platformName
}

func parsePlatformArguments(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	arguments := strings.Split(req.URL.Path, "/")
	if len(arguments) == 6 && arguments[len(arguments)-2] == "bandwidths" && arguments[len(arguments)-1] == "average" {
		platformName := arguments[len(arguments)-3]
		metricType := arguments[len(arguments)-2]
		nodes = getPlatformBandwidthsFiveMinutesAverage(platformName, metricType, rw)
	} else if len(arguments) == 5 && arguments[len(arguments)-1] == "contact" {
		platformName := arguments[len(arguments)-2]
		getPlatformContact(platformName, nodes)

	} else if arguments[4] == "isp" {
		platformName := arguments[3]
		duration := "1d"
		if len(arguments) == 6 {
			duration = arguments[5]
		}
		getPlatformBandwidthsByISP(platformName, duration, nodes)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		errors := []string{}
		var result = make(map[string]interface{})
		result["error"] = errors
		errorMessage := "Error: wrong API path."
		if strings.Index(req.URL.Path, "/bandwidths/") > -1 {
			errorMessage += " Example: /api/platforms/{platformName}/bandwidths/average"
		} else if strings.Index(req.URL.Path, "/contact") > -1 {
			errorMessage += " Example: /api/platforms/{platformName}/contact"
		} else if strings.Index(req.URL.Path, "/isp") > -1 {
			errorMessage += " Example: /api/platforms/{platformName}/isp"
		}
		setError(errorMessage, result)
		nodes["result"] = result
	}
	setResponse(rw, nodes)
}

func getTicker(timestamp int64, unit int) string {
	now := time.Now().Unix()
	diff := now - timestamp
	if diff <= 0 {
		return time.Now().Format("2006-01-02 15:04")
	}
	if diff <= 600 {
		return time.Unix(timestamp, 0).Format("2006-01-02 15:04")
	}
	ticker := ""
	date := time.Unix(timestamp, 0)
	minute := date.Format("04")
	value, err := strconv.Atoi(minute)
	if err == nil {
		residue := int(math.Mod(float64(value), float64(unit)))
		value -= residue
		minute = strconv.Itoa(value)
		if len(minute) == 1 {
			minute = "0" + minute
		}
		ticker = date.Format("2006-01-02 15:") + minute
	}
	return ticker
}

func getNearestTicker(timestamp float64, tickers []string) string {
	tickerNearest := ""
	index := -1
	min := math.Inf(1)
	for key, ticker := range tickers {
		diff := math.Abs(timestamp - float64(getTimestampFromTicker(ticker)))
		if min > diff {
			min = diff
			index = key
		}
	}
	if index >= 0 {
		tickerNearest = tickers[index]
	}
	return tickerNearest
}

func getTimestampFromTicker(ticker string) int64 {
	timestamp := int64(0)
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		loc = time.Local
	}
	timeFormat := "2006-01-02 15:04"
	date, err := time.ParseInLocation(timeFormat, ticker, loc)
	if err == nil {
		timestamp = date.Unix()
	}
	return timestamp
}

func getGraphLastData(metrics []string, hostnames []string, result map[string]interface{}) []*cmodel.GraphQueryResponse {
	var data []*cmodel.GraphQueryResponse
	for _, metric := range metrics {
		for _, hostname := range hostnames {
			param := cmodel.GraphLastParam{
				Endpoint: hostname,
				Counter:  metric,
			}
			resp, err := graph.Last(param)
			if err != nil {
				log.Errorf(err.Error())
			} else if resp != nil {
				item := cmodel.RRDData{
					Timestamp: resp.Value.Timestamp,
					Value:     resp.Value.Value,
				}
				values := []*cmodel.RRDData{}
				values = append(values, &item)
				datum := cmodel.GraphQueryResponse{
					Endpoint: resp.Endpoint,
					Counter:  resp.Counter,
					Values:   values,
				}
				data = append(data, &datum)
			}
		}
	}
	return data
}

func getGraphQueryData(metrics []string, duration string, hostnames []string, step int, result map[string]interface{}) []*cmodel.GraphQueryResponse {
	var data []*cmodel.GraphQueryResponse
	start, end := convertDurationToPoint(duration, result)
	if (time.Now().Unix() - start) < 200 {
		data = getGraphLastData(metrics, hostnames, result)
	} else {
		data = getGraphQueryResponse(hostnames, metrics, duration, "AVERAGE", step, result)
		diff := time.Now().Unix() - end
		secondsInHalfDay := int64(43200)
		if ((end - start) > secondsInHalfDay) && (diff < 600) && (len(data) > 0) {
			dataRecent := getGraphLastData(metrics, hostnames, result)
			data = addRecentData(data, dataRecent)
		}
	}
	return data
}

func getBandwidthsSum(metricType string, duration string, hostnames []string, filter string, result map[string]interface{}) []interface{} {
	items := []interface{}{}
	sort.Strings(hostnames)
	metrics := getMetricsByMetricType(metricType)
	metricMap := map[string]interface{}{}
	valuesMap := map[string]map[string]float64{}
	timestamps := []int64{}
	tickers := []string{}
	tickersMap := map[string]float64{}
	if len(metrics) > 0 && len(hostnames) > 0 {
		data := getGraphQueryData(metrics, duration, hostnames, 1200, result)
		index := -1
		max := 0
		for key, item := range data {
			if len(item.Values) > max {
				max = len(item.Values)
				index = key
			}
		}
		if index >= 0 {
			for _, rrdObj := range data[index].Values {
				ticker := getTicker(rrdObj.Timestamp, 5)
				if _, ok := tickersMap[ticker]; !ok {
					if len(ticker) > 0 {
						tickersMap[ticker] = float64(0)
						tickers = append(tickers, ticker)
					}
				}
				timestamps = append(timestamps, rrdObj.Timestamp)
			}
		}
		if len(tickers) > 0 {
			for _, metric := range metrics {
				tickerMap := map[string]float64{}
				for _, ticker := range tickers {
					tickerMap[ticker] = float64(0)
				}
				valuesMap[metric] = tickerMap
			}
			for _, series := range data {
				metric := series.Counter
				tickerMap := valuesMap[metric]
				for _, rrdObj := range series.Values {
					if !math.IsNaN(float64(rrdObj.Value)) {
						ticker := getNearestTicker(float64(rrdObj.Timestamp), tickers)
						tickerMap[ticker] += float64(rrdObj.Value)
					}
				}
				metricMap[metric] = tickerMap
			}
		}
	}
	if len(tickers) > 0 {
		for _, metric := range metrics {
			tickerMap := metricMap[metric].(map[string]float64)
			data := [][]float64{}
			for _, ticker := range tickers {
				timestamp := getTimestampFromTicker(ticker)
				value := tickerMap[ticker]
				datum := []float64{
					float64(timestamp * 1000),
					value,
				}
				data = append(data, datum)
			}
			item := map[string]interface{}{
				"host":   strings.Join(hostnames, ","),
				"metric": metric,
				"data":   data,
			}
			items = append(items, item)
		}
		if len(filter) > 0 && strings.Index(filter, ",") == -1 {
			queryIDCsBandwidths(filter, result)
			if len(result["error"].([]string)) > 0 {
				result["error"] = []string{}
			} else {
				upperLimit := result["items"].(map[string]interface{})["upperLimitMB"].(float64) * 1000 * 1000
				data := []interface{}{}
				for _, timestamp := range timestamps {
					datum := []interface{}{
						timestamp * 1000,
						upperLimit,
					}
					data = append(data, datum)
				}
				item := map[string]interface{}{
					"host":   strings.Join(hostnames, ","),
					"metric": "net.if.upper.limit.bits",
					"data":   data,
				}
				items = append(items, item)
			}
		}
	}
	return items
}

func getNICOutSpeed(hostname string, result map[string]interface{}) int {
	NICOutSpeed := 0
	metrics := []string{
		"nic.default.out.speed",
	}
	var param cmodel.GraphLastParam
	var params []cmodel.GraphLastParam
	param.Endpoint = hostname
	for _, metric := range metrics {
		param.Counter = metric
		params = append(params, param)
	}

	var data []cmodel.GraphLastResp
	proc.LastRequestCnt.Incr()
	for _, param := range params {
		last, err := graph.Last(param)
		if err != nil {
			setError("graph.last fail, err: "+err.Error(), result)
			return NICOutSpeed
		}
		if last == nil {
			continue
		}
		data = append(data, *last)
	}
	proc.LastRequestItemCnt.IncrBy(int64(len(data)))
	if len(data) > 0 {
		if data[0].Value.Value > 0 {
			NICOutSpeed = int(data[0].Value.Value)
		} else {
			for _, item := range data {
				if NICOutSpeed < int(item.Value.Value) {
					NICOutSpeed = int(item.Value.Value)
				}
			}
		}
	}
	return NICOutSpeed
}

func getHostsBandwidths(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	items := []interface{}{}
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	arguments := strings.Split(req.URL.Path, "/")
	hostnames := []string{}
	metricType := ""
	method := ""
	duration := ""
	if len(arguments) == 7 && arguments[4] == "bandwidths" {
		hostnames = strings.Split(arguments[3], ",")
		metricType = arguments[4]
		method = arguments[5]
		duration = arguments[6]
	} else if len(arguments) == 5 && arguments[2] == "hosts" {
		hostnames = strings.Split(arguments[3], ",")
		method = arguments[4]
	}
	if method == "average" {
		items = getBandwidthsAverage(metricType, duration, hostnames, result)
	} else if method == "sum" {
		filter := req.URL.Query().Get("filter")
		items = getBandwidthsSum(metricType, duration, hostnames, filter, result)
	} else if method == "nic-out-speed" {
		for _, hostname := range hostnames {
			if strings.Index(hostname, "-") > -1 {
				NICOutSpeed := getNICOutSpeed(hostname, result)
				item := map[string]interface{}{
					"hostname":         hostname,
					"nic.out.speed.MB": NICOutSpeed,
				}
				items = append(items, item)
			}
		}
	}
	result["items"] = items
	nodes["result"] = result
	nodes["count"] = len(items)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getIDCsHosts(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	idcsMap := map[string]interface{}{}
	idcIDs := []string{}
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	getPlatformJSON(nodes, result)
	hosts := map[string]interface{}{}
	hostnames := []string{}
	hostnamesMap := map[string]int{}
	if nodes["status"] != nil && int(nodes["status"].(float64)) == 1 {
		hostname := ""
		for _, platform := range nodes["result"].([]interface{}) {
			for _, device := range platform.(map[string]interface{})["ip_list"].([]interface{}) {
				hostname = device.(map[string]interface{})["hostname"].(string)
				if _, ok := hostnamesMap[hostname]; !ok {
					ip := device.(map[string]interface{})["ip"].(string)
					if len(ip) > 0 && ip == getIPFromHostname(hostname, result) {
						hostnames = append(hostnames, hostname)
						idcID := device.(map[string]interface{})["pop_id"].(string)
						host := map[string]interface{}{
							"activate": device.(map[string]interface{})["ip_status"].(string),
							"hostname": hostname,
							"idcID":    idcID,
							"ip":       ip,
						}
						hostnamesMap[hostname] = 1
						hosts[hostname] = host
						idcIDs = appendUniqueString(idcIDs, idcID)
					}
				}
			}
		}
		sort.Strings(hostnames)
		sort.Strings(idcIDs)
		for _, hostname := range hostnames {
			host := hosts[hostname].(map[string]interface{})
			idcID := host["idcID"].(string)
			if _, ok := idcsMap[idcID]; ok {
				idcsMap[idcID] = append(idcsMap[idcID].([]map[string]interface{}), host)
			} else {
				idcsMap[idcID] = []map[string]interface{}{
					host,
				}
			}
		}
		IDCNamesMap := map[string]string{}
		IDCNames := []string{}
		o := orm.NewOrm()
		o.Using("grafana")
		var idcs []*Idc
		sqlcommand := "SELECT pop_id, name FROM `grafana`.`idc` ORDER BY pop_id ASC"
		_, err := o.Raw(sqlcommand).QueryRows(&idcs)
		if err != nil {
			setError(err.Error(), result)
		} else {
			for _, idc := range idcs {
				IDCNamesMap[idc.Name] = strconv.Itoa(idc.Pop_id)
				IDCNames = appendUniqueString(IDCNames, idc.Name)
			}
		}
		sort.Strings(IDCNames)
		for _, IDCName := range IDCNames {
			idcID := IDCNamesMap[IDCName]
			if _, ok := idcsMap[idcID]; ok {
				idc := idcsMap[idcID]
				idcsMap[IDCName] = idc
				delete(idcsMap, idcID)
			}
		}
	}
	if _, ok := nodes["info"]; ok {
		delete(nodes, "info")
	}
	if _, ok := nodes["status"]; ok {
		delete(nodes, "status")
	}
	result["items"] = idcsMap
	nodes["result"] = result
	nodes["count"] = len(idcIDs)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func queryIDCsBandwidths(IDCName string, result map[string]interface{}) {
	var nodes = make(map[string]interface{})
	upperLimitSum := float64(0)
	fcname := g.Config().Api.Name
	fctoken := getFctoken()
	url := g.Config().Api.Uplink
	params := map[string]string{
		"fcname":   fcname,
		"fctoken":  fctoken,
		"pop_name": IDCName,
	}
	s, err := json.Marshal(params)
	if err != nil {
		setError(err.Error(), result)
	}
	reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
	if err != nil {
		setError(err.Error(), result)
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		setError(err.Error(), result)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &nodes)
		if err != nil {
			setError(err.Error(), result)
		}
		if nodes["status"] != nil && int(nodes["status"].(float64)) == 1 {
			if len(nodes["result"].([]interface{})) == 0 {
				errorMessage := "IDC name not found: " + IDCName
				setError(errorMessage, result)
			} else {
				for _, uplink := range nodes["result"].([]interface{}) {
					if upperLimit, ok := uplink.(map[string]interface{})["all_uplink_top"].(float64); ok {
						upperLimitSum += upperLimit
					}
				}
			}
		} else {
			setError("Error occurs", result)
		}
	}
	items := map[string]interface{}{
		"IDCName":      IDCName,
		"upperLimitMB": upperLimitSum,
	}
	if _, ok := nodes["info"]; ok {
		delete(nodes, "info")
	}
	if _, ok := nodes["status"]; ok {
		delete(nodes, "status")
	}
	result["items"] = items
}

func getIDCsBandwidthsUpperLimit(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	arguments := strings.Split(req.URL.Path, "/")
	IDCName := ""
	if len(arguments) == 6 && arguments[len(arguments)-2] == "bandwidths" && arguments[len(arguments)-1] == "limit" {
		IDCName = arguments[len(arguments)-3]
		queryIDCsBandwidths(IDCName, result)
	} else {
		errorMessage := "Error: wrong API path."
		errorMessage += " Example: /api/idcs/{IDCName}/bandwidths/limit"
		setError(errorMessage, result)
	}
	nodes["result"] = result
	setResponse(rw, nodes)
}

func getHostsList(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}
	getPlatformJSON(nodes, result)
	hosts := map[string]interface{}{}
	hostnames := []string{}
	hostnamesMap := map[string]int{}
	if int(nodes["status"].(float64)) == 1 {
		hostname := ""
		for _, platform := range nodes["result"].([]interface{}) {
			platformName := platform.(map[string]interface{})["platform"].(string)
			for _, device := range platform.(map[string]interface{})["ip_list"].([]interface{}) {
				hostname = device.(map[string]interface{})["hostname"].(string)
				ip := device.(map[string]interface{})["ip"].(string)
				if len(ip) > 0 && ip == getIPFromHostname(hostname, result) {
					if _, ok := hostnamesMap[hostname]; !ok {
						ip := device.(map[string]interface{})["ip"].(string)
						if len(ip) > 0 && ip == getIPFromHostname(hostname, result) {
							hostnames = append(hostnames, hostname)
							idcID := device.(map[string]interface{})["pop_id"].(string)
							host := map[string]interface{}{
								"activate":     device.(map[string]interface{})["ip_status"].(string),
								"hostname":     hostname,
								"idcID":        idcID,
								"ip":           ip,
								"platform":     platformName,
								"isp":          strings.Split(hostname, "-")[0],
								"provinceCode": strings.Split(hostname, "-")[1],
							}
							hostnamesMap[hostname] = 1
							hosts[hostname] = host
						}
					} else {
						host := hosts[hostname].(map[string]interface{})
						platforms := strings.Split(host["platform"].(string), ", ")
						platforms = appendUniqueString(platforms, platformName)
						host["platform"] = strings.Join(platforms, ", ")
						hosts[hostname] = host
					}
				}
			}
		}
		sort.Strings(hostnames)
		idcIDsMap := map[string]interface{}{}
		idcNames := []string{}
		o := orm.NewOrm()
		o.Using("grafana")
		var idcs []*Idc
		sqlcommand := "SELECT pop_id, province, city, name FROM `grafana`.`idc` ORDER BY pop_id ASC"
		_, err := o.Raw(sqlcommand).QueryRows(&idcs)
		if err != nil {
			setError(err.Error(), result)
		} else {
			for _, idc := range idcs {
				item := map[string]string{
					"name":     idc.Name,
					"province": idc.Province,
					"city":     idc.City,
				}
				idcIDsMap[strconv.Itoa(idc.Pop_id)] = item
				idcNames = appendUniqueString(idcNames, idc.Name)
			}
		}
		sort.Strings(idcNames)
		for _, hostname := range hostnames {
			host := hosts[hostname].(map[string]interface{})
			idcID := host["idcID"].(string)
			if _, ok := idcIDsMap[idcID]; ok {
				item := idcIDsMap[idcID]
				host["idc"] = item.(map[string]string)["name"]
				host["province"] = item.(map[string]string)["province"]
				host["city"] = item.(map[string]string)["city"]
				delete(host, "idcID")
				items = append(items, host)
			}
		}
	}
	if _, ok := nodes["info"]; ok {
		delete(nodes, "info")
	}
	if _, ok := nodes["status"]; ok {
		delete(nodes, "status")
	}
	result["items"] = items
	nodes["count"] = len(items)
	nodes["result"] = result
	setResponse(rw, nodes)
}

func getHostgroups(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	hosts := queryHostsData(result)
	result["count"] = len(hosts)
	result["items"] = hosts
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func configAPIRoutes() {
	http.HandleFunc("/api/info", queryInfo)
	http.HandleFunc("/api/history", queryHistory)
	http.HandleFunc("/api/endpoints", dashboardEndpoints)
	http.HandleFunc("/api/counters", dashboardCounters)
	http.HandleFunc("/api/chart", dashboardChart)
	http.HandleFunc("/api/alive", getAlive)
	http.HandleFunc("/api/tags/update", setStrategyTags)
	http.HandleFunc("/api/templates/", getTemplateStrategies)
	http.HandleFunc("/api/alive/platforms", getAgentAlivePlatforms)
	http.HandleFunc("/api/metrics.health/", getHostMetricValues)
	http.HandleFunc("/api/apollo/filters", getApolloFilters)
	http.HandleFunc("/api/apollo/charts/", getApolloCharts)
	http.HandleFunc("/api/apollo/remarks/add", addApolloRemark)
	http.HandleFunc("/api/apollo/remarks/", getApolloRemark)
	http.HandleFunc("/api/platforms", getPlatforms)
	http.HandleFunc("/api/platforms/", parsePlatformArguments)
	http.HandleFunc("/api/hosts/", getHostsBandwidths)
	http.HandleFunc("/api/idcs/hosts", getIDCsHosts)
	http.HandleFunc("/api/idcs/", getIDCsBandwidthsUpperLimit)
	http.HandleFunc("/api/hosts", getHostsList)
	http.HandleFunc("/api/hostgroups", getHostgroups)
}
