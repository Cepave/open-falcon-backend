package http

import (
	"bytes"
	"encoding/json"
	cmodel "github.com/Cepave/common/model"
	"github.com/Cepave/query/g"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Tag struct {
	StrategyId int
	Name       string
	Value      string
	Create_at  string
	Update_at  string
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
func postByJson(rw http.ResponseWriter, req *http.Request, url string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	s := buf.String()
	reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
	if err != nil {
		log.Println("Error =", err.Error())
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		log.Println("Error =", err.Error())
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
	postByJson(rw, req, url)
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
	postByJson(rw, req, url)
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
		log.Println("Error =", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error =", err.Error())
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
		log.Println("Error =", err.Error())
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
				countOfWarn += 1
			} else {
				status = "normal"
				countOfNormal += 1
			}
		} else {
			countOfDead += 1
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

func configApiRoutes() {
	http.HandleFunc("/api/info", queryInfo)
	http.HandleFunc("/api/history", queryHistory)
	http.HandleFunc("/api/endpoints", dashboardEndpoints)
	http.HandleFunc("/api/counters", dashboardCounters)
	http.HandleFunc("/api/chart", dashboardChart)
	http.HandleFunc("/api/alive", getAlive)
}
