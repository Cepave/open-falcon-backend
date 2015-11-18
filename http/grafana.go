package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/**
 * @function name:  func getHosts(rw http.ResponseWriter, req *http.Request, hostKeyword string)
 * @description:    This function returns available hosts for Grafana query editor.
 * @related issues: OWL-151
 * @param:          rw http.ResponseWriter
 * @param:          req *http.Request
 * @param:          hostKeyword string
 * @return:         void
 * @author:         Don Hsieh
 * @since:          11/17/2015
 * @last modified:  11/18/2015
 * @called by:      func setQueryEditor(rw http.ResponseWriter, req *http.Request, hostKeyword string)
 */
func getHosts(rw http.ResponseWriter, req *http.Request, hostKeyword string) {
	hostKeyword = strings.Replace(hostKeyword, "*", "", -1)
	if len(hostKeyword) == 0 {
		hostKeyword = "%"
	}
	rand.Seed(time.Now().UTC().UnixNano())
	random64 := rand.Float64()
	_r := strconv.FormatFloat(random64, 'f', -1, 32)
	target := req.URL.Query().Get("target")
	target = strings.Replace(target, "/api/grafana", "/api/endpoints", -1)
	url := target + "?q=" + hostKeyword + "&tags&limit=500&_r=" + _r
	
	reqGet, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error =", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(reqGet)
	if err != nil {
		log.Println("Error =", err.Error())
	}
	defer resp.Body.Close()

	result := []interface{}{}
	chart := map[string]interface{} {
		"text": "chart",
		"expandable": true,
	}
	result = append(result, chart)

	if resp.Status == "200 OK" {
		body, _ := ioutil.ReadAll(resp.Body)
		var nodes = make(map[string]interface{})
		if err := json.Unmarshal(body, &nodes); err != nil {
			log.Println(err.Error())
		}
		for _, host := range nodes["data"].([]interface {}) {
			item := map[string]interface{} {
				"text": host,
				"expandable": true,
			}
			result = append(result, item)
		}
		RenderJson(rw, result)
	} else {
		RenderJson(rw, result)
	}
}

/**
 * @function name:  func getNextCounterSegment(metric string, counter string) string
 * @description:    This function returns next segment of a counter.
 * @related issues: OWL-151
 * @param:          metric string
 * @param:          counter string
 * @return:         segment string
 * @author:         Don Hsieh
 * @since:          11/18/2015
 * @last modified:  11/18/2015
 * @called by:      func getMetrics(rw http.ResponseWriter, req *http.Request, query string)
 */
func getNextCounterSegment(metric string, counter string) string {
	if len(metric) > 0 {
		metric += "."
	}
	counter = strings.Replace(counter, metric, "", 1)
	segment := strings.Split(counter, ".")[0]
	return segment
}

/**
 * @function name:  func checkSegmentExpandable(segment string, counter string) bool
 * @description:    This function checks if a segment is expandable.
 * @related issues: OWL-151
 * @param:          segment string
 * @param:          counter string
 * @return:         expandable bool
 * @author:         Don Hsieh
 * @since:          11/18/2015
 * @last modified:  11/18/2015
 * @called by:      func getMetrics(rw http.ResponseWriter, req *http.Request, query string)
 */
func checkSegmentExpandable(segment string, counter string) bool {
	segments := strings.Split(counter, ".")
	expandable := !(segment == segments[len(segments)-1])
	return expandable
}

/**
 * @function name:  func getMetrics(rw http.ResponseWriter, req *http.Request, query string)
 * @description:    This function returns available segments of a metric for Grafana query editor.
 * @related issues: OWL-151
 * @param:          rw http.ResponseWriter
 * @param:          req *http.Request
 * @param:          query string
 * @return:         void
 * @author:         Don Hsieh
 * @since:          11/17/2015
 * @last modified:  11/18/2015
 * @called by:      func setQueryEditor(rw http.ResponseWriter, req *http.Request, hostKeyword string)
 */
func getMetrics(rw http.ResponseWriter, req *http.Request, query string) {
	result := []interface{}{}
	
	query = strings.Replace(query, ".*", "", -1)
	arrQuery := strings.Split(query, ".")
	host, arrMetric := arrQuery[0], arrQuery[1:]

	if host == "chart" {
		chartBar := map[string]interface {} {
			"text": "bar",
			"expandable": false,
		}
		result = append(result, chartBar)

		chartMap := map[string]interface {} {
			"text": "map",
			"expandable": false,
		}
		result = append(result, chartMap)

		chartPie := map[string]interface {} {
			"text": "pie",
			"expandable": false,
		}
		result = append(result, chartPie)
		RenderJson(rw, result)
	} else {
		metric := strings.Join(arrMetric, ".")
		endpoints := "[\"" + host + "\"]"

		rand.Seed(time.Now().UTC().UnixNano())
		random64 := rand.Float64()
		_r := strconv.FormatFloat(random64, 'f', -1, 32)

		form := url.Values{}
		form.Set("endpoints", endpoints)
		form.Add("q", metric)
		form.Add("limit", "")
		form.Add("_r", _r)
		
		target := req.URL.Query().Get("target")
		target = strings.Replace(target, "/api/grafana", "/api/counters", -1)
		reqPost, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
		if err != nil {
			log.Println("Error =", err.Error())
		}
		reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(reqPost)
		if err != nil {
			log.Println("Error =", err.Error())
		}
		defer resp.Body.Close()

		if resp.Status == "200 OK" {
			body, _ := ioutil.ReadAll(resp.Body)
			var nodes = make(map[string]interface{})
			if err := json.Unmarshal(body, &nodes); err != nil {
				log.Println(err.Error())
			}
			var segmentPool = make(map[string]int)
			for _, data := range nodes["data"].([]interface {}) {
				counter := data.([]interface {})[0].(string)
				segment := getNextCounterSegment(metric, counter)
				expandable := checkSegmentExpandable(segment, counter)
				if _, ok := segmentPool[segment]; !ok {
					item := map[string]interface {} {
						"text": segment,
						"expandable": expandable,
					}
					result = append(result, item)
					segmentPool[segment] = 1
				}
			}
			RenderJson(rw, result)
		} else {
			RenderJson(rw, result)
		}
	}
}

/**
 * @function name:  func setQueryEditor(rw http.ResponseWriter, req *http.Request)
 * @description:    This function returns arranges data for Grafana query editor.
 * @related issues: OWL-151
 * @param:          rw http.ResponseWriter
 * @param:          req *http.Request
 * @return:         void
 * @author:         Don Hsieh
 * @since:          11/17/2015
 * @last modified:  11/18/2015
 * @called by:      func GrafanaApiParser(rw http.ResponseWriter, req *http.Request)
 */
func setQueryEditor(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("query")
	query = strings.Replace(query, ".%", "", -1)
	query = strings.Replace(query, ".undefined", "", -1)
	query = strings.Replace(query, ".select metric", "", -1)
	if !strings.Contains(query, ".") {
		getHosts(rw, req, query)
	} else {
		getMetrics(rw, req, query)
	}
}

/**
 * @function name:  func getValues(rw http.ResponseWriter, req *http.Request)
 * @description:    This function returns metric values for Grafana to draw graph.
 * @related issues: OWL-151
 * @param:          rw http.ResponseWriter
 * @param:          req *http.Request
 * @return:         void
 * @author:         Don Hsieh
 * @since:          11/16/2015
 * @last modified:  11/17/2015
 * @called by:      func GrafanaApiParser(rw http.ResponseWriter, req *http.Request)
 */
func getValues(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	endpoint_counters := []interface{}{}
	for _, target := range req.PostForm["target"] {
		if !strings.Contains(target, ".select metric") {
			targets := strings.Split(target, ".")
			host, targets := targets[0], targets[1:]
			metric := strings.Join(targets, ".")
			if strings.Contains(host, "{") {
				host = strings.Replace(host, "{", "", -1)
				host = strings.Replace(host, "}", "", -1)
				hosts := strings.Split(host, ",")
				for _, host := range hosts {	// Templating metrics request
												// host:"{host1,host2}"
					item := map[string]string {
						"endpoint": host,
						"counter": metric,
					}
					endpoint_counters = append(endpoint_counters, item)
				}
			} else {
				item := map[string]string {
					"endpoint": host,
					"counter": metric,
				}
				endpoint_counters = append(endpoint_counters, item)
			}
		}
	}
	result := []interface{}{}
	if len(endpoint_counters) > 0 {
		from, err := strconv.ParseInt(req.PostForm["from"][0], 10, 64)
		until, err := strconv.ParseInt(req.PostForm["until"][0], 10, 64)
		param := req.URL.Query()
		url := strings.Join(param["target"], "")
		url = strings.Replace(url, "/grafana", "/history", -1)

		args := map[string]interface{} {
			"start": from,
			"end": until,
			"cf": "AVERAGE",
			"endpoint_counters": endpoint_counters,
		}
		bs, err := json.Marshal(args)
		if err != nil {
			log.Println("Error =", err.Error())
		}

		reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bs)))
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

		if resp.Status == "200 OK" {
			body, _ := ioutil.ReadAll(resp.Body)
			nodes := []interface{}{}
			if err := json.Unmarshal(body, &nodes); err != nil {
				log.Println(err.Error())
			}

			for _, node := range nodes {
				if _, ok := node.(map[string]interface{})["Values"]; ok {
					result = append(result, node)
				}
			}
			RenderJson(rw, result)
		} else {
			RenderJson(rw, result)
		}
	} else {
		RenderJson(rw, result)
	}
}

/**
 * @function name:  func GrafanaApiParser(rw http.ResponseWriter, req *http.Request)
 * @description:    This function parses the method of API request.
 * @related issues: OWL-151
 * @param:          rw http.ResponseWriter
 * @param:          req *http.Request
 * @return:         void
 * @author:         Don Hsieh
 * @since:          11/16/2015
 * @last modified:  11/16/2015
 * @called by:      func configGrafanaRoutes()
 */
func GrafanaApiParser(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		setQueryEditor(rw, req)
	} else if req.Method == "POST" {
		getValues(rw, req)
	}
}

/**
 * @function name:  func configGrafanaRoutes()
 * @description:    This function handles API requests.
 * @related issues: OWL-151
 * @param:          void
 * @return:         void
 * @author:         Don Hsieh
 * @since:          11/16/2015
 * @last modified:  11/16/2015
 * @called by:      func Start()
 *                   in http/http.go
 */
func configGrafanaRoutes() {
	http.HandleFunc("/api/grafana/", GrafanaApiParser)
}
