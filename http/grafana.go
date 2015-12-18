package http

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/Cepave/query/g"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Idc struct {
	Id         int
	Pop_id     int
	Name       string
	Count      int
	Area   string
	Province   string
	City       string
	Updated_at string
}

type Province struct {
	Id         int
	Province   string
	Count      int
	Updated_at string
}

type City struct {
	Id         int
	City       string
	Province   string
	Count      int
	Updated_at string
}

/**
 * @function name:   func getHosts(rw http.ResponseWriter, req *http.Request, hostKeyword string)
 * @description:     This function returns available hosts for Grafana query editor.
 * @related issues:  OWL-221, OWL-151
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @param:           hostKeyword string
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/17/2015
 * @last modified:   12/18/2015
 * @called by:       func setQueryEditor(rw http.ResponseWriter, req *http.Request, hostKeyword string)
 */
func getHosts(rw http.ResponseWriter, req *http.Request, hostKeyword string) {
	hostKeyword = strings.Replace(hostKeyword, "*", "", -1)
	if len(hostKeyword) == 0 {
		hostKeyword = "%"
	}
	rand.Seed(time.Now().UTC().UnixNano())
	random64 := rand.Float64()
	_r := strconv.FormatFloat(random64, 'f', -1, 32)
	url := "/api/endpoints" + "?q=" + hostKeyword + "&tags&limit=500&_r=" + _r
	log.Println("req.Host =", req.Host)
	log.Println("g.Config().Api.Query =", g.Config().Api.Query)
	if strings.Index(g.Config().Api.Query, req.Host) >= 0 {
		url = "http://localhost:9966" + url
	} else {
		url = g.Config().Api.Query + url
	}
	log.Println("url =", url)

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
 * @function name:   func getNextCounterSegment(metric string, counter string) string
 * @description:     This function returns next segment of a counter.
 * @related issues:  OWL-151
 * @param:           metric string
 * @param:           counter string
 * @return:          segment string
 * @author:          Don Hsieh
 * @since:           11/18/2015
 * @last modified:   11/18/2015
 * @called by:       func getMetrics(rw http.ResponseWriter, req *http.Request, query string)
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
 * @function name:   func checkSegmentExpandable(segment string, counter string) bool
 * @description:     This function checks if a segment is expandable.
 * @related issues:  OWL-151
 * @param:           segment string
 * @param:           counter string
 * @return:          expandable bool
 * @author:          Don Hsieh
 * @since:           11/18/2015
 * @last modified:   11/18/2015
 * @called by:       func getMetrics(rw http.ResponseWriter, req *http.Request, query string)
 */
func checkSegmentExpandable(segment string, counter string) bool {
	segments := strings.Split(counter, ".")
	expandable := !(segment == segments[len(segments)-1])
	return expandable
}

/**
 * @function name:   func getMetrics(rw http.ResponseWriter, req *http.Request, query string)
 * @description:     This function returns available segments of a metric for Grafana query editor.
 * @related issues:  OWL-221, OWL-151
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @param:           query string
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/17/2015
 * @last modified:   12/18/2015
 * @called by:       func setQueryEditor(rw http.ResponseWriter, req *http.Request, hostKeyword string)
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

		chartRose := map[string]interface {} {
			"text": "rose",
			"expandable": false,
		}
		result = append(result, chartRose)
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

		log.Println("req.Host =", req.Host)
		log.Println("g.Config().Api.Query =", g.Config().Api.Query)
		target := "/api/counters"
		if strings.Index(g.Config().Api.Query, req.Host) >= 0 {
			target = "http://localhost:9966" + target
		} else {
			target = g.Config().Api.Query + target
		}
		log.Println("target =", target)

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
 * @function name:   func setQueryEditor(rw http.ResponseWriter, req *http.Request)
 * @description:     This function returns arranges data for Grafana query editor.
 * @related issues:  OWL-151
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/17/2015
 * @last modified:   11/18/2015
 * @called by:       func GrafanaApiParser(rw http.ResponseWriter, req *http.Request)
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
 * @function name:   func getLocation(pop_id int) map[string] string
 * @description:     This function gets location of a server room given by pop_id (ID of server room).
 * @related issues:  OWL-159
 * @param:           void
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/25/2015
 * @last modified:   11/25/2015
 * @called by:       func updateMapData()
 */
func getLocation(pop_id int) map[string] string {
	location := map[string] string {
		"area": "",
		"province": "",
		"city": "",
	}
	fcname := g.Config().Api.Name
	fctoken := getFctoken()
	url := g.Config().Api.Geo

	args := map[string]interface{} {
		"fcname": fcname,
		"fctoken": fctoken,
		"pop_id": pop_id,
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
		nodes := map[string]interface{}{}
		if err := json.Unmarshal(body, &nodes); err != nil {
			log.Println(err.Error())
		}
		status := int(nodes["status"].(float64))
		if status == 1 {
			result := nodes["result"]
			location["area"] = result.(map[string]interface {})["area"].(string)
			location["province"] = result.(map[string]interface {})["province"].(string)
			location["city"] = result.(map[string]interface {})["city"].(string)
		}
	}
	return location
}

/**
 * @function name:   func updateCities()
 * @description:     This function updates "city" table in "grafana" database.
 * @related issues:  OWL-159
 * @param:           void
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/25/2015
 * @last modified:   11/25/2015
 * @called by:       func updateMapData()
 */
func updateCities() {
	var rows []orm.Params
	o := orm.NewOrm()
	o.Using("grafana")
	sqlcmd := "SELECT area, province, city, count FROM grafana.idc WHERE 1"
	_, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		log.Println("Error =", err.Error())
	} else {
		keys := map[string]int{}
		for _, row := range rows {
			area := row["area"].(string)
			province := row["province"].(string)
			city := row["city"].(string)
			count, err := strconv.Atoi(row["count"].(string))
			if err != nil {
				log.Println("Error =", err.Error())
			}
			key := province + "_" + city
			if province == city {
				key = province
			} else {
				key = strings.Replace(key, "特区_", "", -1)
				key = strings.Replace(key, "_其他", "", -1)
			}
			key = area + "_" + key
			if _, ok := keys[key]; ok {
				keys[key] += count
			} else {
				keys[key] = count
			}
		}

		provinces := map[string]int{}
		cities := map[string]interface{}{}
		provinceNames := []string{}
		cityNames := []string{}
		for key, count := range keys {
			arr := strings.Split(key, "_")
			area, arr := arr[0], arr[1:]
			key = strings.Join(arr, "_")
			provinceName := ""
			if !strings.Contains(key, "_") {
				provinceName = key
			} else {
				provinceName = strings.Split(key, "_")[0]
			}
			provinceIndex := area + "_" + provinceName
			if _, ok := provinces[provinceIndex]; ok {
				provinces[provinceIndex] += count
			} else {
				provinces[provinceIndex] = count
				provinceNames = append(provinceNames, provinceIndex)
			}
			if strings.Contains(key, "_") {
				cityName := strings.Split(key, "_")[1]
				cityIndex := area + "_" + provinceName + "_" + cityName
				if _, ok := cities[cityIndex]; ok {
					city := cities[cityIndex]
					city.(map[string]int)["count"] += count
				} else {
					city := map[string]interface{}{
						"city": cityName,
						"province": provinceName,
						"count": count,
					}
					cities[cityIndex] = city
					cityNames = append(cityNames, cityIndex)
				}
			}
		}
		sort.Strings(provinceNames)
		sort.Strings(cityNames)
		for _, provinceIndex := range provinceNames {
			count := provinces[provinceIndex]
			provinceName := strings.Split(provinceIndex, "_")[1]

			var rows []orm.Params
			sqlcmd := "SELECT id, province FROM grafana.province WHERE province=?"
			num, err := o.Raw(sqlcmd, provinceName).Values(&rows)
			if err != nil {
				log.Println("Error =", err.Error())
			} else {
				province := Province{
					Province: provinceName,
					Count: count,
					Updated_at: getNow(),
				}
				if num > 0 {	// existed. update data.
					id, err := strconv.Atoi(rows[0]["id"].(string))
					if err != nil {
						log.Println("Error =", err.Error())
					}
					province.Id = id
					num, err := o.Update(&province)
					if err != nil {
						log.Println("Error =", err.Error())
					} else {
						if num > 0 {
							log.Println("update provinceId:", id)
							log.Println("mysql row affected nums:", num)
						}
					}
				} else {		// not existed. insert data.
					provinceId, err := o.Insert(&province)
					if err != nil {
						log.Println("Error =", err.Error())
					} else {
						log.Println("Insert provinceId =", provinceId)
					}
				}
			}
		}
		for _, cityIndex := range cityNames {
			item := cities[cityIndex]
			provinceName := item.(map[string]interface {})["province"].(string)
			cityName := item.(map[string]interface {})["city"].(string)
			count := item.(map[string]interface {})["count"].(int)
			var rows []orm.Params
			sqlcmd := "SELECT id, city FROM grafana.city WHERE city=?"
			num, err := o.Raw(sqlcmd, cityName).Values(&rows)
			if err != nil {
				log.Println("Error =", err.Error())
			} else {
				city := City{
					Province: provinceName,
					City: cityName,
					Count: count,
					Updated_at: getNow(),
				}
				if num > 0 {	// existed. update data.
					id, err := strconv.Atoi(rows[0]["id"].(string))
					if err != nil {
						log.Println("Error =", err.Error())
					}
					city.Id = id
					num, err := o.Update(&city)
					if err != nil {
						log.Println("Error =", err.Error())
					} else {
						if num > 0 {
							log.Println("update cityId:", id)
							log.Println("mysql row affected nums:", num)
						}
					}
				} else {		// not existed. insert data.
					cityId, err := o.Insert(&city)
					if err != nil {
						log.Println("Error =", err.Error())
					} else {
						log.Println("Insert cityId =", cityId)
					}
				}
			}
		}
	}
}

/**
 * @function name:   func updateMapData()
 * @description:     This function updates map data by API return values.
 * @related issues:  OWL-159
 * @param:           void
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/25/2015
 * @last modified:   11/25/2015
 * @called by:       func getMapValues(chartType string) map[string]interface{}
 */
func updateMapData() {
	fcname := g.Config().Api.Name
	fctoken := getFctoken()
	url := g.Config().Api.Map + "/fcname/" + fcname + "/fctoken/" + fctoken
	url += "/eqt/yes/hostname/yes/pop/yes/pop_id/yes/show_active/yes/show_isp/yes.json"
	log.Println("url =", url)

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
	var nodes = make(map[string]interface{})
	if err := json.Unmarshal(body, &nodes); err != nil {
		log.Println("Error =", err.Error())
	}
	result := map[string]int{}
	items := map[string]interface {}{}
	names := []string{}
	if int(nodes["status"].(float64)) == 1 {
		countOfPlatform := 0
		countOfDevice := 0
		for _, platform := range nodes["result"].([]interface {}) {
			for _, device := range platform.(map[string]interface {})["ip_list"].([]interface {}) {
				countOfDevice++
				id, err := strconv.Atoi(device.(map[string]interface {})["pop_id"].(string))
				if err != nil {
					log.Println("Error =", err.Error())
				}
				name := device.(map[string]interface {})["pop"].(string)
				if _, ok := result[name]; ok {
					result[name]++
					item := items[name]
					count := item.(map[string]interface {})["count"].(int)
					count++
					item.(map[string]interface {})["count"] = count
					items[name] = item
				} else {
					result[name] = 1
					item := map[string]interface{} {
						"id": id,
						"name": name,
						"count": 1,
					}
					items[name] = item
					names = append(names, name)
				}
			}
			countOfPlatform++
		}
		log.Println("countOfPlatform =", countOfPlatform)
		log.Println("countOfDevice =", countOfDevice)
		sort.Strings(names)

		o := orm.NewOrm()
		o.Using("grafana")
		for _, name := range names {
			log.Println("item =", items[name])
			item := items[name]
			pop_id := item.(map[string]interface {})["id"].(int)
			name := item.(map[string]interface {})["name"].(string)
			count := item.(map[string]interface {})["count"].(int)
			now := getNow()
			idc := Idc{
				Pop_id: pop_id,
			}
			location := getLocation(pop_id)
			log.Println("location =", location)
			area := location["area"]
			province := location["province"]
			city := location["city"]

			var rows []orm.Params
			sqlcmd := "SELECT id, pop_id FROM grafana.idc WHERE pop_id=?"
			num, err := o.Raw(sqlcmd, pop_id).Values(&rows)
			if err != nil {
				log.Println("Error =", err.Error())
			} else {
				idc.Name = name
				idc.Count = count
				idc.Area = area
				idc.Province = province
				idc.City = city
				idc.Updated_at = now
				if num > 0 {	// existed. update data.
					id, err := strconv.Atoi(rows[0]["id"].(string))
					if err != nil {
						log.Println("Error =", err.Error())
					}
					idc.Id = id
					num, err := o.Update(&idc)
					if err != nil {
						log.Println("Error =", err.Error())
					} else {
						if num > 0 {
							log.Println("update idcId:", id)
							log.Println("mysql row affected nums:", num)
						}
					}
				} else {		// not existed. insert data.
					idcId, err := o.Insert(&idc)
					if err != nil {
						log.Println("Error =", err.Error())
					} else {
						log.Println("Insert idcId =", idcId)
					}
				}
			}
		}
		updateCities()
	}
}

/**
 * @function name:   func getMapValues(chartType string) map[string]interface{}
 * @description:     This function returns map values for Grafana to draw graph.
 * @related issues:  OWL-159
 * @param:           chartType string
 * @return:          hosts map[string]interface{}
 * @author:          Don Hsieh
 * @since:           11/19/2015
 * @last modified:   11/26/2015
 * @called by:       func getValues(rw http.ResponseWriter, req *http.Request)
 */
func getMapValues(chartType string) map[string]interface{} {
	hosts := map[string]interface{}{}
	provinces := []interface{}{}
	sqlcmd := "SELECT province, count FROM grafana.province WHERE 1"
	o := orm.NewOrm()
	var rows []orm.Params
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		log.Println("Error =", err.Error())
	} else if num > 0 {
		for _, row := range rows {
			name := row["province"]
			count := row["count"]
			item := map[string]interface{} {
				"name": name,
				"value": count,
			}
			provinces = append(provinces, item)
		}
	}

	citiesInProvince := []interface{}{}
	sqlcmd = "SELECT city, count FROM grafana.city WHERE 1"
	num, err = o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		log.Println("Error =", err.Error())
	} else if num > 0 {
		for _, row := range rows {
			name := row["city"]
			count := row["count"]
			item := map[string]interface{} {
				"name": name,
				"value": count,
			}
			citiesInProvince = append(citiesInProvince, item)
		}
	}
	hosts["chartType"] = chartType
	hosts["provinces"] = provinces
	hosts["citiesInProvince"] = citiesInProvince
	return hosts
}

/**
 * @function name:   func getBarChartOptions(chartType string) map[string]interface{}
 * @description:     This function returns bar chart options for Grafana to draw graph.
 * @related issues:  OWL-159
 * @param:           chartType string
 * @param:           provinces []interface{}
 * @return:          chart map[string]interface{}
 * @author:          Don Hsieh
 * @since:           11/26/2015
 * @last modified:   11/27/2015
 * @called by:       func getChartOptions(chartType string) map[string]interface{}
 */
func getBarChartOptions(chartType string, provinces []interface{}) map[string]interface{} {
	chart := map[string]interface{}{}
	names := []string{}
	values := []int{}
	for _, item := range provinces {
		name := item.(map[string]interface{})["name"].(string)
		value := item.(map[string]interface{})["value"].(int)
		if value > 40 {
			names = append(names, name)
			values = append(values, value)
		}
	}

	grid := map[string]int{
		"borderWidth": 0,
		"y": 80,
		"y2": 60,
	}

	color := []string {
		"#09aa3c",
	}
	label := map[string]interface{} {
		"show": true,
		"position": "top",
		"formatter": "{b}\n{c}",
	}
	normal := map[string]interface{} {
		"color": color,
		"label": label,
	}
	itemStyle := map[string]interface{} {
		"normal": normal,
	}

	xAxis := map[string]interface{} {
		"type": "category",
		"show": false,
		"data": names,
	}

	yAxis := map[string]interface{} {
		"type": "value",
		"show": false,
	}

	series := map[string]interface{} {
		"data": values,
		"itemStyle": itemStyle,
		"name": "servers",
		"type": "bar",
	}

	tooltip := map[string]string {
		"trigger": "item",
		"formatter": "{a} <br/>{b} : {c} ({d}%)",
	}

	option := map[string]interface{} {
		"calculable" : true,
		"grid" : grid,
		"series": []interface{} {series},
		"xAxis": []interface{} {xAxis},
		"yAxis": []interface{} {yAxis},
		"tooltip": tooltip,
	}

	chart["chartType"] = chartType
	chart["option"] = option
	return chart
}

/**
 * @function name:   func getRoseChartOptions(chartType string) map[string]interface{}
 * @description:     This function returns rose chart options for Grafana to draw graph.
 * @related issues:  OWL-159
 * @param:           chartType string
 * @param:           provinces []interface{}
 * @return:          chart map[string]interface{}
 * @author:          Don Hsieh
 * @since:           11/26/2015
 * @last modified:   11/26/2015
 * @called by:       func getChartOptions(chartType string) map[string]interface{}
 */
func getRoseChartOptions(chartType string, provinces []interface{}) map[string]interface{} {
	chart := map[string]interface{}{}

	showTrue := map[string]bool{
		"show": true,
	}
	showFalse := map[string]bool{
		"show": false,
	}
	emphasis := map[string]interface{} {
		"label": showTrue,
		"labelLine": showTrue,
	}
	normal := map[string]interface{} {
		"label": showFalse,
		"labelLine": showFalse,
	}
	itemStyle := map[string]interface{} {
		"emphasis": emphasis,
		"normal": normal,
	}

	series := map[string]interface{} {
		"center": []string{"50%", "50%"},
		"data": provinces,
		"itemStyle": itemStyle,
		"name": "servers",
		"radius": []int{20, 75},
		"roseType": "radius",
		"type": "pie",
	}

	tooltip := map[string]string {
		"trigger": "item",
		"formatter": "{a} <br/>{b} : {c} ({d}%)",
	}

	option := map[string]interface{} {
		"calculable" : true,
		"series": []interface{} {series},
		"tooltip": tooltip,
	}

	chart["chartType"] = chartType
	chart["option"] = option
	return chart
}

/**
 * @function name:   func getChartOptions(chartType string) map[string]interface{}
 * @description:     This function returns chart options for Grafana to draw graph.
 * @related issues:  OWL-159
 * @param:           chartType string
 * @return:          chart map[string]interface{}
 * @author:          Don Hsieh
 * @since:           11/26/2015
 * @last modified:   11/27/2015
 * @called by:       func getValues(rw http.ResponseWriter, req *http.Request)
 */
func getChartOptions(chartType string) map[string]interface{} {
	chart := map[string]interface{}{}
	provinces := []interface{}{}
	max := 0
	sqlcmd := "SELECT province, count, updated_at FROM grafana.province WHERE 1"
	o := orm.NewOrm()
	var rows []orm.Params
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		log.Println("Error =", err.Error())
	}
	if num > 0 {
		for _, row := range rows {
			name := row["province"]
			count, err := strconv.Atoi(row["count"].(string))
			if err != nil {
				log.Println("Error =", err.Error())
			} else if max < count {
				max = count
			}
			item := map[string]interface{} {
				"name": name,
				"value": count,
			}
			provinces = append(provinces, item)
		}
		updatedAt := rows[0]["updated_at"]
		date, err := time.Parse("2006-01-02 15:04:05", updatedAt.(string))
		if err != nil {
			log.Println("Error =", err.Error())
		}
		hours := time.Since(date).Hours() + 8
		log.Println("hours =", hours)
		if hours > 24 {
			updateMapData()
		}
	} else {
		updateMapData()
	}

	if (chartType == "rose") {
		return getRoseChartOptions(chartType, provinces)
	} else if (chartType == "bar") {
		return getBarChartOptions(chartType, provinces)
	}

	citiesInProvince := []interface{}{}
	sqlcmd = "SELECT city, count FROM grafana.city WHERE 1"
	num, err = o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		log.Println("Error =", err.Error())
	} else if num > 0 {
		for _, row := range rows {
			name := row["city"].(string)
			if !strings.Contains(name, "地区") {
				name += "市"
			}
			count := row["count"].(string)
			item := map[string]interface{} {
				"name": name,
				"value": count,
			}
			citiesInProvince = append(citiesInProvince, item)
		}
	}

	color := []string {
		"maroon",
		"purple",
		"red",
		"orange",
		"yellow",
		"lightgreen",
	}

	dataRange := map[string]interface{} {
		"x": "right",
		"min": 0,
		"max": max,
		"calculable": true,
		"color": color,
	}

	title := "servers"
	legend := map[string]interface{} {
		"orient": "vertical",
		"data": []string {title},
	}

	tooltip := map[string]string {
		"trigger": "item",
	}

	label := map[string]bool{
		"show": true,
	}
	emphasis := map[string]interface{} {
		"label": label,
	}
	itemStyleForProvinces := map[string]interface{} {
		"emphasis": emphasis,
	}
	seriesForProvinces := map[string]interface{} {
		"data": provinces,
		"itemStyle": itemStyleForProvinces,
		"mapLocation": map[string]string{},
		"mapType": "china",
		"name": "server",
		"roam": true,
		"selectedMode": "single",
		"type": "map",
	}

	itemStyleForCities := map[string]interface{} {
		"normal": emphasis,
		"emphasis": emphasis,
	}
	seriesForCities := map[string]interface{} {
		"data": citiesInProvince,
		"itemStyle": itemStyleForCities,
		"mapLocation": map[string]string{},
		"mapType": "china",
		"name": title,
		"roam": true,
		"type": "map",
	}
	series := []interface{} {
		seriesForProvinces,
		seriesForCities,
	}

	option := map[string]interface{} {
		"dataRange": dataRange,
		"legend": legend,
		"series": series,
		"tooltip": tooltip,
	}

	chart["chartType"] = chartType
	chart["option"] = option
	return chart
}

/**
 * @function name:   func getMetricValues(req *http.Request, host string, targets []string, result []interface{}) []interface{}
 * @description:     This function returns map values for Grafana to draw graph.
 * @related issues:  OWL-159
 * @param:           req *http.Request
 * @param:           host string
 * @param:           targets []string
 * @return:          result []interface{}
 * @author:          Don Hsieh
 * @since:           11/23/2015
 * @last modified:   11/24/2015
 * @called by:       func getValues(rw http.ResponseWriter, req *http.Request)
 */
func getMetricValues(req *http.Request, host string, targets []string, result []interface{}) []interface{} {
	endpoint_counters := []interface{}{}
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
		}
	}
	return result
}

/**
 * @function name:   func getValues(rw http.ResponseWriter, req *http.Request)
 * @description:     This function returns metric values for Grafana to draw graph.
 * @related issues:  OWL-159, OWL-151
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/16/2015
 * @last modified:   11/24/2015
 * @called by:       func GrafanaApiParser(rw http.ResponseWriter, req *http.Request)
 */
func getValues(rw http.ResponseWriter, req *http.Request) {
	result := []interface{}{}
	req.ParseForm()
	for _, target := range req.PostForm["target"] {
		if !strings.Contains(target, ".select metric") {
			targets := strings.Split(target, ".")
			host, targets := targets[0], targets[1:]
			if host == "chart" {
				chartType := targets[len(targets)-1]
				chartValues := getChartOptions(chartType)
				result = append(result, chartValues)
			} else {
				result = getMetricValues(req, host, targets, result)
			}
		}
	}
	RenderJson(rw, result)
}

/**
 * @function name:   func GrafanaApiParser(rw http.ResponseWriter, req *http.Request)
 * @description:     This function parses the method of API request.
 * @related issues:  OWL-151
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/16/2015
 * @last modified:   11/16/2015
 * @called by:       func configGrafanaRoutes()
 */
func GrafanaApiParser(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		setQueryEditor(rw, req)
	} else if req.Method == "POST" {
		getValues(rw, req)
	}
}

/**
 * @function name:   func configGrafanaRoutes()
 * @description:     This function handles API requests.
 * @related issues:  OWL-151
 * @param:           void
 * @return:          void
 * @author:          Don Hsieh
 * @since:           11/16/2015
 * @last modified:   11/16/2015
 * @called by:       func Start()
 *                    in http/http.go
 */
func configGrafanaRoutes() {
	http.HandleFunc("/api/grafana/", GrafanaApiParser)
}
