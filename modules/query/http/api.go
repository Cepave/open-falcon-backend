package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	"github.com/Cepave/open-falcon-backend/modules/query/graph"
	"github.com/Cepave/open-falcon-backend/modules/query/proc"
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
				status = "warm"
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
	url += "/show_active/yes/hostname/yes/pop_id/yes/ip/yes.json"

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

func setGraphQueries(hostnames []string, hostnamesExisted []string, versions map[string]string, result map[string]interface{}) []*cmodel.GraphLastParam {
	var queries []*cmodel.GraphLastParam
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
			var query cmodel.GraphLastParam
			if !strings.Contains(host.Hostname, ".") && strings.Contains(host.Hostname, "-") {
				hostnamesExisted = append(hostnamesExisted, host.Hostname)
				versions[host.Hostname] = host.Agent_version
				query.Endpoint = host.Hostname
				query.Counter = "agent.alive"
				queries = append(queries, &query)
			}
		}
	}
	return queries
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

func classifyAgentAliveResponse(data []cmodel.GraphLastResp, hostnamesExisted []string, versions map[string]string, result map[string]interface{}) {
	name := ""
	version := ""
	status := ""
	alive := 0
	var diff int64
	var timestamp int64
	items := map[string]interface{}{}
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
		version = versions[name]
		if alive > 0 {
			if diff > 3600 {
				status = "warm"
			} else {
				status = "normal"
			}
		}
		item := map[string]interface{}{}
		item["version"] = version
		item["status"] = status
		items[name] = item
	}
	result["items"] = items
}

func getAnomalies(errorHosts []interface{}, result map[string]interface{}) map[string]interface{} {
	anomalies := map[string]interface{}{}
	pop_ids := map[string]string{}
	for _, errorHost := range errorHosts {
		pop_id := errorHost.(map[string]string)["pop_id"]
		pop_ids[pop_id] = pop_id
	}
	arr := []string{}
	for _, pop_id := range pop_ids {
		arr = append(arr, pop_id)
	}
	sort.Strings(arr)

	sqlcmd := "SELECT pop_id, name, province, city FROM grafana.idc WHERE pop_id IN ('"
	sqlcmd += strings.Join(arr, "','") + "')"
	idcs := map[string]interface{}{}
	var rows []orm.Params
	o := orm.NewOrm()
	o.Using("grafana")
	_, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else {
		for _, row := range rows {
			idc := map[string]string{
				"idc":      row["name"].(string),
				"province": row["province"].(string),
				"city":     row["city"].(string),
			}
			idcs[row["pop_id"].(string)] = idc
		}
	}

	anomalies2 := map[string]interface{}{}
	for _, errorHost := range errorHosts {
		pop_id := errorHost.(map[string]string)["pop_id"]
		if idc, ok := idcs[pop_id]; ok {
			errorHost.(map[string]string)["idc"] = idc.(map[string]string)["idc"]
			errorHost.(map[string]string)["city"] = idc.(map[string]string)["city"]
			provinceName := idc.(map[string]string)["province"]
			if provinceName == "特区" {
				provinceName = idc.(map[string]string)["city"]
			}
			errorHost.(map[string]string)["province"] = provinceName
			delete(errorHost.(map[string]string), "pop_id")

			if province, ok := anomalies2[provinceName]; ok {
				province = append(province.([]map[string]string), errorHost.(map[string]string))
				anomalies2[provinceName] = province
			} else {
				anomalies2[provinceName] = []map[string]string{
					errorHost.(map[string]string),
				}
			}
		} else {
			errorHost.(map[string]string)["idc"] = ""
			errorHost.(map[string]string)["city"] = ""
			errorHost.(map[string]string)["province"] = ""
		}
	}

	for provinceName, hosts := range anomalies2 {
		count := len(hosts.([]map[string]string))
		anomalies[provinceName] = map[string]interface{}{
			"count": count,
			"hosts": hosts,
		}
	}
	return anomalies
}

func completeAgentAliveData(groups map[string]interface{}, groupNames []string, result map[string]interface{}) {
	errorHosts := []interface{}{}
	platforms := []interface{}{}
	count := map[string]int{}
	countOfNormalSum := 0
	countOfWarnSum := 0
	countOfErrorSum := 0
	countOfMissSum := 0
	countOfDeactivatedSum := 0
	hostId := 1
	name := ""
	activate := ""
	version := ""
	pop_id := ""
	status := ""
	items := result["items"].(map[string]interface{})
	for _, groupName := range groupNames {
		platform := map[string]interface{}{}
		hosts := []interface{}{}
		count := map[string]int{}
		countOfNormal := 0
		countOfWarn := 0
		countOfError := 0
		countOfMiss := 0
		countOfDeactivated := 0
		group := groups[groupName].([]interface{})
		for _, agent := range group {
			name = agent.(map[string]interface{})["name"].(string)
			activate = agent.(map[string]interface{})["activate"].(string)
			pop_id = agent.(map[string]interface{})["pop_id"].(string)
			status = ""
			version = ""
			if activate == "1" {
				if item, ok := items[name]; ok {
					status = item.(map[string]interface{})["status"].(string)
					version = item.(map[string]interface{})["version"].(string)
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
			} else if status == "warm" {
				countOfWarn++
			} else if status == "error" {
				countOfError++
			}
			host := map[string]string{
				"id":       strconv.Itoa(hostId),
				"name":     name,
				"platform": groupName,
				"pop_id":   pop_id,
				"status":   status,
				"version":  version,
			}
			if host["status"] == "error" {
				errorHosts = append(errorHosts, host)
			} else {
				delete(host, "pop_id")
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

func getPlatforms(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	getPlatformJSON(nodes, result)
	groups := map[string]interface{}{}
	groupNames := []string{}
	hostnames := []string{}
	hostnamesMap := map[string]int{}
	if nodes["status"] != nil && int(nodes["status"].(float64)) == 1 {
		hostname := ""
		for _, platform := range nodes["result"].([]interface{}) {
			groupName := platform.(map[string]interface{})["platform"].(string)
			groupNames = append(groupNames, groupName)
			group := []interface{}{}
			for _, device := range platform.(map[string]interface{})["ip_list"].([]interface{}) {
				hostname = device.(map[string]interface{})["hostname"].(string)
				ip := device.(map[string]interface{})["ip"].(string)
				if len(ip) > 0 && ip == getIPFromHostname(hostname, result) {
					if _, ok := hostnamesMap[hostname]; !ok {
						hostnames = append(hostnames, hostname)
						host := map[string]interface{}{
							"name":     hostname,
							"activate": device.(map[string]interface{})["ip_status"].(string),
							"pop_id":   device.(map[string]interface{})["pop_id"].(string),
						}
						group = append(group, host)
						hostnamesMap[hostname] = 1
					}
				}
			}
			groups[groupName] = group
		}
		sort.Strings(hostnames)
		sort.Strings(groupNames)

		hostnamesExisted := []string{}
		var versions = make(map[string]string)
		queries := setGraphQueries(hostnames, hostnamesExisted, versions, result)
		data := queryAgentAlive(queries, req.Host, result)
		classifyAgentAliveResponse(data, hostnamesExisted, versions, result)
		completeAgentAliveData(groups, groupNames, result)
	}
	if _, ok := nodes["info"]; ok {
		delete(nodes, "info")
	}
	if _, ok := nodes["status"]; ok {
		delete(nodes, "status")
	}
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

func convertDurationToPoint(duration string, result map[string]interface{}) (timestampFrom int64, timestampTo int64) {
	if strings.Index(duration, ",") > -1 {
		var e error
		timestampFrom, e = strconv.ParseInt(strings.Split(duration, ",")[0], 10, 64)
		if e != nil {
			setError(e.Error(), result)
		}
		timestampTo, e = strconv.ParseInt(strings.Split(duration, ",")[1], 10, 64)
		if e != nil {
			setError(e.Error(), result)
		}
		if timestampFrom >= timestampTo {
			setError("Value of timestampFrom should be less than value of timestampTo.", result)
		}
		if timestampTo > time.Now().Unix() {
			setError("Value of timestampTo should be equal to or less than value of now.", result)
		}
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
		timestampTo = now
	}
	return
}

func getGraphQueryResponse(metrics []string, duration string, hostnames []string, result map[string]interface{}) []*cmodel.GraphQueryResponse {
	data := []*cmodel.GraphQueryResponse{}
	start, end := convertDurationToPoint(duration, result)

	proc.HistoryRequestCnt.Incr()
	for _, hostname := range hostnames {
		for _, metric := range metrics {
			request := cmodel.GraphQueryParam{
				Start:     start,
				End:       end,
				ConsolFun: "AVERAGE",
				Endpoint:  hostname,
				Counter:   metric,
			}
			response, err := graph.QueryOne(request)
			if err != nil {
				setError("graph.queryOne fail, "+err.Error(), result)
				return data
			}
			if result == nil {
				continue
			}
			data = append(data, response)
		}
	}

	proc.HistoryResponseCounterCnt.IncrBy(int64(len(data)))
	for _, item := range data {
		proc.HistoryResponseItemCnt.IncrBy(int64(len(item.Values)))
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
		data := getGraphQueryResponse(metrics, duration, []string{hostname}, result)

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
	sliceStr := strings.Join(slice, "','")
	if !strings.Contains(sliceStr, s) {
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

func getExistedHosts(hosts []interface{}, hostnamesExisted []string, result map[string]interface{}) map[string]interface{} {
	hostsExisted := map[string]interface{}{}
	for key, hostname := range hostnamesExisted {
		host := map[string]interface{}{
			"id":   key + 1,
			"name": hostname,
		}
		hostsExisted[hostname] = host
	}
	idcMap := map[string]string{}
	o := orm.NewOrm()
	var idcs []*Idc
	sqlcommand := "SELECT pop_id, name FROM grafana.idc ORDER BY pop_id ASC"
	_, err := o.Raw(sqlcommand).QueryRows(&idcs)
	if err != nil {
		setError(err.Error(), result)
	} else {
		for _, idc := range idcs {
			idcMap[strconv.Itoa(idc.Pop_id)] = idc.Name
		}
	}
	for _, host := range hosts {
		hostname := host.(map[string]interface{})["name"].(string)
		if _, ok := hostsExisted[hostname]; ok {
			hostExisted := hostsExisted[hostname]
			isp := strings.Split(hostname, "-")[0]
			province := strings.Split(hostname, "-")[1]
			popID := host.(map[string]interface{})["popID"].(string)
			idc := idcMap[popID]
			platform := host.(map[string]interface{})["platform"].(string)
			hostExisted.(map[string]interface{})["isp"] = isp
			hostExisted.(map[string]interface{})["province"] = province
			hostExisted.(map[string]interface{})["idc"] = idc
			hostExisted.(map[string]interface{})["platform"] = platform
			hostsExisted[hostname] = hostExisted
		}
	}
	return hostsExisted
}

func completeApolloFiltersData(hostsExisted map[string]interface{}, result map[string]interface{}) {
	hosts := map[string]interface{}{}
	keywords := map[string]interface{}{}
	for _, host := range hostsExisted {
		id := host.(map[string]interface{})["id"].(int)
		platform := host.(map[string]interface{})["platform"].(string)
		tags := []string{}
		tags = appendUniqueString(tags, platform)
		if _, ok := keywords[platform]; ok {
			keywords[platform] = appendUnique(keywords[platform].([]int), id)
		} else {
			keywords[platform] = []int{id}
		}
		isp := host.(map[string]interface{})["isp"].(string)
		tags = appendUniqueString(tags, isp)
		if _, ok := keywords[isp]; ok {
			keywords[isp] = appendUnique(keywords[isp].([]int), id)
		} else {
			keywords[isp] = []int{id}
		}

		province := host.(map[string]interface{})["province"].(string)
		tags = appendUniqueString(tags, province)
		if _, ok := keywords[province]; ok {
			keywords[province] = appendUnique(keywords[province].([]int), id)
		} else {
			keywords[province] = []int{id}
		}

		name := host.(map[string]interface{})["name"].(string)
		fragments := strings.Split(name, "-")
		if len(fragments) == 6 {
			fragments := fragments[2:]
			for _, fragment := range fragments {
				tags = appendUniqueString(tags, fragment)
				if _, ok := keywords[fragment]; ok {
					keywords[fragment] = appendUnique(keywords[fragment].([]int), id)
				} else {
					keywords[fragment] = []int{id}
				}
			}
		}
		ip := getIPFromHostname(name, result)
		keywords[ip] = []int{id}
		tags = appendUniqueString(tags, ip)

		idc := host.(map[string]interface{})["idc"].(string)
		tags = appendUniqueString(tags, idc)
		if _, ok := keywords[idc]; ok {
			keywords[idc] = appendUnique(keywords[idc].([]int), id)
		} else {
			keywords[idc] = []int{id}
		}
		host.(map[string]interface{})["tag"] = strings.Join(tags, ",")
		delete(host.(map[string]interface{}), "id")
		delete(host.(map[string]interface{}), "isp")
		delete(host.(map[string]interface{}), "province")
		hosts[strconv.Itoa(id)] = host
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
	getPlatformJSON(nodes, result)
	hosts := []interface{}{}
	hostnames := []string{}
	if nodes["status"] != nil && int(nodes["status"].(float64)) == 1 {
		hostname := ""
		for _, platform := range nodes["result"].([]interface{}) {
			groupName := platform.(map[string]interface{})["platform"].(string)
			for _, device := range platform.(map[string]interface{})["ip_list"].([]interface{}) {
				hostname = device.(map[string]interface{})["hostname"].(string)
				ip := device.(map[string]interface{})["ip"].(string)
				if len(ip) > 0 && ip == getIPFromHostname(hostname, result) {
					popID := device.(map[string]interface{})["pop_id"].(string)
					hostnames = append(hostnames, hostname)
					host := map[string]interface{}{
						"name":     hostname,
						"platform": groupName,
						"popID":    popID,
					}
					hosts = append(hosts, host)
					hostnames = append(hostnames, hostname)
				}
			}
		}
		sort.Strings(hostnames)
		hostnamesExisted := getExistedHostnames(hostnames, result)
		sort.Strings(hostnamesExisted)
		hostsExisted := getExistedHosts(hosts, hostnamesExisted, result)
		count = len(hostsExisted)
		completeApolloFiltersData(hostsExisted, result)
	}
	if _, ok := nodes["info"]; ok {
		delete(nodes, "info")
	}
	if _, ok := nodes["status"]; ok {
		delete(nodes, "status")
	}
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

func getApolloCharts(rw http.ResponseWriter, req *http.Request) {
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}
	arguments := strings.Split(req.URL.Path, "/")
	metricType := arguments[4]
	hostnames := strings.Split(arguments[5], ",")
	metrics := getMetricsByMetricType(metricType)
	if metricType == "bandwidths" {
		metrics = addNetworkSpeedAndBondMode(metrics, hostnames, result)
	}
	duration := "1d"
	if len(arguments) > 6 {
		duration = arguments[6]
	}
	data := getGraphQueryResponse(metrics, duration, hostnames, result)
	for _, series := range data {
		metric := series.Counter
		if strings.Index(metric, "nic.default.out.speed/device=") > -1 {
			if len(series.Values) > 0 && series.Values[0].Value > 0 {
				series.Counter = "net.transmission.limit.80%"
				limit := series.Values[0].Value
				for _, item := range series.Values {
					value := item.Value
					if value > limit {
						limit = value
						break
					}
				}
				limit *= 1024 * 1024 * 0.8
				for key, _ := range series.Values {
					series.Values[key].Value = limit
				}
			} else {
				series.Counter = ""
			}
		}
		if series.Counter != "" {
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

func getIPFromHostname(hostname string, result map[string]interface{}) string {
	ip := ""
	fragments := strings.Split(hostname, "-")
	slice := []string{}
	if len(fragments) == 6 {
		fragments := fragments[2:]
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
	hostMap := map[string]interface{}{}
	if len(metrics) > 0 && len(hostnames) > 0 {
		data := getGraphQueryResponse(metrics, duration, hostnames, result)
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
			if value, ok := hostMap[series.Endpoint]; ok {
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
			hostMap[series.Endpoint] = item
		}
	}
	for _, hostname := range hostnames {
		host := hostMap[hostname]
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
			for _, name := range strings.Split(platformName, ",") {
				if platform, ok := nodes["result"].(map[string]interface{})[name].([]interface{}); ok {
					items := []interface{}{}
					for _, contact := range platform {
						item := map[string]interface{}{
							"name":  contact.(map[string]interface{})["realname"].(string),
							"phone": contact.(map[string]interface{})["cell"].(string),
							"email": contact.(map[string]interface{})["email"].(string),
						}
						items = append(items, item)
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
	} else {
		errors := []string{}
		var result = make(map[string]interface{})
		result["error"] = errors
		errorMessage := "Error: wrong API path."
		if strings.Index(req.URL.Path, "/bandwidths/") > -1 {
			errorMessage += " Example: /api/platforms/{platformName}/bandwidths/average"
		} else if strings.Index(req.URL.Path, "/contact") > -1 {
			errorMessage += " Example: /api/platforms/{platformName}/contact"
		}
		setError(errorMessage, result)
		nodes["result"] = result
	}
	setResponse(rw, nodes)
}

func getBandwidthsSum(metricType string, duration string, hostnames []string, filter string, result map[string]interface{}) []interface{} {
	items := []interface{}{}
	sort.Strings(hostnames)
	metrics := getMetricsByMetricType(metricType)
	metricMap := map[string]interface{}{}
	timestamps := []float64{}
	if len(metrics) > 0 && len(hostnames) > 0 {
		data := getGraphQueryResponse(metrics, duration, hostnames, result)
		for _, rrdObj := range data[0].Values {
			timestamps = append(timestamps, float64(rrdObj.Timestamp))
		}
		if len(timestamps) > 0 {
			for _, metric := range metrics {
				values := []float64{}
				for _, _ = range timestamps {
					values = append(values, float64(0))
				}
				metricMap[metric] = values
			}
			for _, series := range data {
				values := metricMap[series.Counter].([]float64)
				for key, rrdObj := range series.Values {
					if !math.IsNaN(float64(rrdObj.Value)) {
						values[key] += float64(rrdObj.Value)
					}
				}
				metricMap[series.Counter] = values
			}
		}
	}
	if len(timestamps) > 0 {
		for _, metric := range metrics {
			slice := metricMap[metric].([]float64)
			data := []interface{}{}
			for key, value := range slice {
				datum := []interface{}{
					timestamps[key] * 1000,
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
				upperLimit := result["items"].(map[string]interface{})["upperLimitMB"].(float64) * 1024 * 1024
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
	if len(arguments) == 7 && arguments[len(arguments)-3] == "bandwidths" {
		hostnames = strings.Split(arguments[len(arguments)-4], ",")
		metricType = arguments[len(arguments)-3]
		method = arguments[len(arguments)-2]
		duration = arguments[len(arguments)-1]
	}

	if method == "average" {
		items = getBandwidthsAverage(metricType, duration, hostnames, result)
	} else if method == "sum" {
		filter := req.URL.Query().Get("filter")
		items = getBandwidthsSum(metricType, duration, hostnames, filter, result)
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
		var idcs []*Idc
		sqlcommand := "SELECT pop_id, name FROM grafana.idc ORDER BY pop_id ASC"
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
				setError("IDC name not found", result)
			} else {
				for _, uplink := range nodes["result"].([]interface{}) {
					upperLimit := uplink.(map[string]interface{})["all_uplink_top"].(float64)
					upperLimitSum += upperLimit
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
		var idcs []*Idc
		sqlcommand := "SELECT pop_id, province, city, name FROM grafana.idc ORDER BY pop_id ASC"
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

func configAPIRoutes() {
	http.HandleFunc("/api/info", queryInfo)
	http.HandleFunc("/api/history", queryHistory)
	http.HandleFunc("/api/endpoints", dashboardEndpoints)
	http.HandleFunc("/api/counters", dashboardCounters)
	http.HandleFunc("/api/chart", dashboardChart)
	http.HandleFunc("/api/alive", getAlive)
	http.HandleFunc("/api/tags/update", setStrategyTags)
	http.HandleFunc("/api/templates/", getTemplateStrategies)
	http.HandleFunc("/api/alive/platforms", getPlatforms)
	http.HandleFunc("/api/metrics.health/", getHostMetricValues)
	http.HandleFunc("/api/apollo/filters", getApolloFilters)
	http.HandleFunc("/api/apollo/charts/", getApolloCharts)
	http.HandleFunc("/api/platforms/", parsePlatformArguments)
	http.HandleFunc("/api/hosts/", getHostsBandwidths)
	http.HandleFunc("/api/idcs/hosts", getIDCsHosts)
	http.HandleFunc("/api/idcs/", getIDCsBandwidthsUpperLimit)
	http.HandleFunc("/api/hosts", getHostsList)
}
