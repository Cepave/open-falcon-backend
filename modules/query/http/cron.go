package http

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/astaxie/beego/orm"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
)

type IDCMapItem struct {
	Popid    int
	Idc      string
	Province string
	City     string
}

type Contacts struct {
	Id      int
	Name    string
	Phone   string
	Email   string
	Updated string
}

type Hosts struct {
	Id        int
	Hostname  string
	Exist     int
	Activate  int
	Platform  string
	Platforms string
	Idc       string
	Ip        string
	Isp       string
	Province  string
	City      string
	Status    string
	Bonding   int
	Speed     int
	Remark    string
	Updated   string
}

type Idcs struct {
	Id        int
	Popid     int
	Idc       string
	Bandwidth int
	Count     int
	Area      string
	Province  string
	City      string
	Updated   string
}

type Ips struct {
	Id       int
	Ip       string
	Exist    int
	Status   int
	Type     string
	Hostname string
	Platform string
	Updated  string
}

type Platforms struct {
	Id          int
	Platform    string
	Type        string
	Visible     int
	Contacts    string
	Principal   string
	Deputy      string
	Upgrader    string
	Count       int
	Department  string
	Team        string
	Description string
	Updated     string
}

func SyncHostsAndContactsTable() {
	if g.Config().Hosts.Enabled || g.Config().Contacts.Enabled {
		if g.Config().Hosts.Enabled {
			syncIDCsTable()
			syncHostsTable()
			intervalToSyncHostsTable := uint64(g.Config().Hosts.Interval)
			gocron.Every(intervalToSyncHostsTable).Seconds().Do(syncHostsTable)
			intervalToSyncContactsTable := uint64(g.Config().Contacts.Interval)
			gocron.Every(intervalToSyncContactsTable).Seconds().Do(syncIDCsTable)
		}
		if g.Config().Contacts.Enabled {
			syncContactsTable()
			intervalToSyncContactsTable := uint64(g.Config().Contacts.Interval)
			gocron.Every(intervalToSyncContactsTable).Seconds().Do(syncContactsTable)
		}
		if g.Config().Net.Enabled {
			syncNetTable()
			gocron.Every(1).Day().At(g.Config().Net.Time).Do(syncNetTable)
		}
		if g.Config().Deviations.Enabled {
			syncDeviationsTable()
			gocron.Every(1).Day().At(g.Config().Deviations.Time).Do(syncDeviationsTable)
		}
		if g.Config().Speed.Enabled {
			addBondingAndSpeedToHostsTable()
			gocron.Every(1).Day().At(g.Config().Speed.Time).Do(addBondingAndSpeedToHostsTable)
		}
		<-gocron.Start()
	}
}

func getIDCMap() map[string]interface{} {
	idcMap := map[string]interface{}{}
	o := orm.NewOrm()
	o.Using("boss")
	var idcs []IDCMapItem
	sqlcommand := "SELECT `popid`, `idc`, `province`, `city` FROM `idcs` ORDER BY popid ASC"
	_, err := o.Raw(sqlcommand).QueryRows(&idcs)
	if err != nil {
		log.Errorf(err.Error())
	}
	for _, idc := range idcs {
		idcMap[strconv.Itoa(idc.Popid)] = idc
	}
	return idcMap
}

func queryIDCsHostsCount(IDCName string) int64 {
	count := int64(0)
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT COUNT(*) FROM `boss`.`hosts` WHERE idc = ?"
	num, err := o.Raw(sql, IDCName).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
	} else if num > 0 {
		row := rows[0]
		countInt, err := strconv.Atoi(row["COUNT(*)"].(string))
		if err == nil {
			count = int64(countInt)
		}
	}
	return count
}

func syncIDCsTable() {
	log.Debugf("func syncIDCsTable()")
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM `boss`.`idcs` ORDER BY updated DESC LIMIT 1"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		format := "2006-01-02 15:04:05"
		updatedTime, _ := time.Parse(format, rows[0]["updated"].(string))
		currentTime, _ := time.Parse(format, getNow())
		diff := currentTime.Unix() - updatedTime.Unix()
		if int(diff) < g.Config().Contacts.Interval {
			return
		}
	}
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	fcname := g.Config().Api.Name
	fctoken := getFctoken()
	url := g.Config().Api.Map + "/fcname/" + fcname + "/fctoken/" + fctoken
	url += "/pop/yes/pop_id/yes.json"
	log.Debugf("url = %v", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Error = %v", err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error = %v", err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var nodes = make(map[string]interface{})
	if err := json.Unmarshal(body, &nodes); err != nil {
		log.Errorf("Error = %v", err.Error())
		return
	}
	if nodes["status"] == nil {
		return
	} else if int(nodes["status"].(float64)) != 1 {
		return
	}
	IDCsMap := map[string]map[string]string{}
	IDCNames := []string{}
	for _, platform := range nodes["result"].([]interface{}) {
		for _, device := range platform.(map[string]interface{})["ip_list"].([]interface{}) {
			IDCName := device.(map[string]interface{})["pop"].(string)
			if _, ok := IDCsMap[IDCName]; !ok {
				popID := device.(map[string]interface{})["pop_id"].(string)
				item := map[string]string{
					"idc":   IDCName,
					"popid": popID,
				}
				IDCsMap[IDCName] = item
				IDCNames = appendUniqueString(IDCNames, IDCName)
			}
		}
	}
	for _, IDCName := range IDCNames {
		idc := IDCsMap[IDCName]
		idcID, err := strconv.Atoi(idc["popid"])
		if err == nil {
			location := getLocation(idcID)
			log.Debugf("location = %v", location)
			idc["area"] = location["area"]
			idc["province"] = location["province"]
			idc["city"] = location["city"]
		}
		queryIDCsBandwidths(IDCName, result)
		idc["bandwidth"] = "0"
		if val, ok := result["items"].(map[string]interface{})["upperLimitMB"].(float64); ok {
			bandwidth := int(val)
			idc["bandwidth"] = strconv.Itoa(bandwidth)
		}
		count := int(queryIDCsHostsCount(IDCName))
		idc["count"] = strconv.Itoa(count)
		IDCsMap[IDCName] = idc
	}
	updateIDCsTable(IDCNames, IDCsMap)
}

func getHostsBondingAndSpeed(hostname string) map[string]int {
	item := map[string]int{}
	param := cmodel.GraphLastParam{
		Endpoint: hostname,
	}
	param.Counter = "nic.bond.mode"
	resp, err := graph.Last(param)
	if err != nil {
		log.Errorf(err.Error())
	} else if resp != nil {
		value := int(resp.Value.Value)
		if value >= 0 {
			item["bonding"] = value
		}
	}
	param.Counter = "nic.default.out.speed"
	resp, err = graph.Last(param)
	if err != nil {
		log.Errorf(err.Error())
	} else if resp != nil {
		value := int(resp.Value.Value)
		if value > 0 {
			item["speed"] = value
		}
	}
	return item
}

func addBondingAndSpeedToHostsTable() {
	log.Debugf("func addBondingAndSpeedToHostsTable()")
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT id, hostname FROM `boss`.`hosts` WHERE exist = 1"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
	} else if num > 0 {
		var host Hosts
		for _, row := range rows {
			hostname := row["hostname"].(string)
			item := getHostsBondingAndSpeed(hostname)
			o.Using("boss")
			err = o.QueryTable("hosts").Filter("hostname", hostname).One(&host)
			if err != nil {
				log.Errorf(err.Error())
			} else {
				if _, ok := item["bonding"]; ok {
					host.Bonding = item["bonding"]
				}
				if _, ok := item["speed"]; ok {
					host.Speed = item["speed"]
				}
				host.Updated = getNow()
				_, err = o.Update(&host)
				if err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}
}

func getPlatformsType(nodes map[string]interface{}, result map[string]interface{}, platformsMap map[string]map[string]string) map[string]map[string]string {
	fcname := g.Config().Api.Name
	fctoken := getFctoken()
	url := g.Config().Api.Platform
	params := map[string]string{
		"fcname":  fcname,
		"fctoken": fctoken,
	}
	s, err := json.Marshal(params)
	if err != nil {
		log.Errorf(err.Error())
		return platformsMap
	}
	reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
	if err != nil {
		log.Errorf(err.Error())
		return platformsMap
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		log.Errorf(err.Error())
		return platformsMap
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &nodes)
		if err != nil {
			log.Errorf(err.Error())
			return platformsMap
		}
		if nodes["status"] != nil && int(nodes["status"].(float64)) == 1 {
			if len(nodes["result"].([]interface{})) == 0 {
				errorMessage := "No platforms returned"
				setError(errorMessage, result)
				return platformsMap
			} else {
				re_inside_whiteSpaces := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
				for _, platform := range nodes["result"].([]interface{}) {
					platformName := ""
					if platform.(map[string]interface{})["platform"] != nil {
						platformName = platform.(map[string]interface{})["platform"].(string)
					}
					platformType := ""
					if platform.(map[string]interface{})["platform_type"] != nil {
						platformType = platform.(map[string]interface{})["platform_type"].(string)
					}
					department := ""
					if platform.(map[string]interface{})["department"] != nil {
						department = platform.(map[string]interface{})["department"].(string)
					}
					team := ""
					if platform.(map[string]interface{})["team"] != nil {
						team = platform.(map[string]interface{})["team"].(string)
					}
					visible := ""
					if platform.(map[string]interface{})["visible"] != nil {
						visible = platform.(map[string]interface{})["visible"].(string)
					}
					description := platform.(map[string]interface{})["description"].(string)
					if len(description) > 0 {
						description = strings.Replace(description, "\r", " ", -1)
						description = strings.Replace(description, "\n", " ", -1)
						description = strings.Replace(description, "\t", " ", -1)
						description = strings.TrimSpace(description)
						description = re_inside_whiteSpaces.ReplaceAllString(description, " ")
						if len(description) > 200 {
							description = string([]rune(description)[0:100])
						}
					}
					if value, ok := platformsMap[platformName]; ok {
						value["type"] = platformType
						value["visible"] = visible
						value["department"] = department
						value["team"] = team
						value["description"] = description
						platformsMap[platformName] = value
					}
				}
			}
		} else {
			setError("Error occurs", result)
		}
	}
	return platformsMap
}

func getDurationForNetTableQuery(offset int) (int64, int64) {
	year, month, day := time.Now().Date()
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		loc = time.Local
	}
	timestampFrom := time.Date(year, month, day-offset, 0, 0, 0, 0, loc).Unix() - 300
	timestampTo := time.Date(year, month, day-offset, 23, 59, 59, 0, loc).Unix()
	return timestampFrom, timestampTo
}

func getPlatformsDailyTrafficData(platformName string, offset int) (map[string]map[string]int, string, map[string]int) {
	data := map[string]map[string]int{
		"in":  {},
		"out": {},
	}
	date := ""
	counts := map[string]int{
		"in":  0,
		"out": 0,
	}
	hostnames := []string{}
	var rows []orm.Params
	o := orm.NewOrm()
	o.Using("boss")
	sql := "SELECT DISTINCT hostname FROM `boss`.`ips`"
	sql += " WHERE platform = ? AND exist = 1 ORDER BY hostname ASC"
	num, err := o.Raw(sql, platformName).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
	} else if num > 0 {
		for _, row := range rows {
			hostnames = append(hostnames, row["hostname"].(string))
		}
	}
	metrics := getMetricsByMetricType("bandwidths")
	timestampFrom, timestampTo := getDurationForNetTableQuery(offset)
	responses := []*cmodel.GraphQueryResponse{}
	for _, hostname := range hostnames {
		for _, metric := range metrics {
			request := cmodel.GraphQueryParam{
				Endpoint:  hostname,
				Counter:   metric,
				Start:     timestampFrom,
				End:       timestampTo,
				ConsolFun: "AVERAGE",
				Step:      1200,
			}
			response, err := graph.QueryOne(request)
			if err != nil {
				log.Debugf("graph.queryOne fail = %v", err.Error())
			} else {
				responses = append(responses, response)
			}
		}
	}
	dataRaw := map[string]map[string]float64{
		"in":  {},
		"out": {},
	}
	tickers := []string{}
	if len(responses) > 0 {
		index := -1
		max := 0
		for key, item := range responses {
			if max < len(item.Values) {
				max = len(item.Values)
				index = key
			}
		}
		if index == -1 {
			date = time.Unix(timestampTo, 0).Format("2006-01-02")
			return data, date, counts
		}
		unit := 20
		tickersMap := map[string]float64{}
		for _, rrdObj := range responses[index].Values {
			ticker := getTicker(rrdObj.Timestamp, unit)
			if _, ok := tickersMap[ticker]; !ok {
				if len(ticker) > 0 {
					tickersMap[ticker] = float64(0)
					tickers = append(tickers, ticker)
				}
			}
		}
		for _, series := range responses {
			metric := strings.Replace(series.Counter, "net.if.", "", -1)
			metric = strings.Replace(metric, ".bits/iface=eth_all", "", -1)
			for _, rrdObj := range series.Values {
				value := float64(rrdObj.Value)
				if !math.IsNaN(value) {
					timestamp := rrdObj.Timestamp
					ticker := getNearestTicker(float64(timestamp), tickers)
					if len(ticker) > 0 {
						if _, ok := dataRaw[metric][ticker]; ok {
							dataRaw[metric][ticker] += value
						} else {
							dataRaw[metric][ticker] = value
						}
					}
				}
			}
			counts[metric]++
		}
	}
	for metric, series := range dataRaw {
		for _, ticker := range tickers {
			value := int(math.Floor(series[ticker]))
			date = strings.Split(ticker, " ")[0]
			ticker = strings.Split(ticker, " ")[1]
			data[metric][ticker] = value
		}
	}
	return data, date, counts
}

func getMean(values []int) int {
	mean := 0
	if len(values) == 0 {
		return mean
	}
	sum := 0
	for _, value := range values {
		sum += value
	}
	mean = sum / len(values)
	return mean
}

func getStandardDeviation(values []int) int {
	deviation := 0
	if len(values) == 0 {
		return deviation
	}
	total := 0
	mean := getMean(values)
	for _, value := range values {
		total += (value - mean) * (value - mean)
	}
	variance := float64(total) / float64(len(values))
	deviation = int(math.Sqrt(variance))
	return deviation
}

func getMinMaxAvg(values []int) (int, int, int) {
	avg := 0
	min := 0
	max := 0
	if len(values) > 0 {
		sum := 0
		for _, value := range values {
			sum += value
		}
		avg = sum / len(values)
		sort.Ints(values)
		min = values[0]
		max = values[len(values)-1]
	}
	return min, max, avg
}

func writeToDeviationsTable(platformName string, hour int, minute int, date string, ticker string) {
	o := orm.NewOrm()
	o.Using("apollo")
	var rows []orm.Params
	dateFull := date + " " + ticker + ":00"
	sql := "SELECT metric, COUNT(DISTINCT date), AVG(bits), STD(bits) "
	sql += "FROM `apollo`.`net` WHERE platform = ? AND hour = ? AND minute = ? "
	sql += "AND date >= DATE_SUB(?, INTERVAL 7 DAY) "
	sql += "AND date < DATE_SUB(?, INTERVAL 1 DAY) GROUP BY metric"
	num, err := o.Raw(sql, platformName, hour, minute, dateFull, dateFull).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		for _, row := range rows {
			samples := 0
			value, err := strconv.Atoi(row["COUNT(DISTINCT date)"].(string))
			if err != nil {
				log.Errorf(err.Error())
			} else {
				samples = value
			}
			if samples >= 3 {
				metricKey := 0
				value, err := strconv.Atoi(row["metric"].(string))
				if err != nil {
					log.Errorf(err.Error())
				} else {
					metricKey = value
				}
				mean := 0
				val, err := strconv.ParseFloat(row["AVG(bits)"].(string), 64)
				if err != nil {
					log.Errorf(err.Error())
				} else {
					mean = int(math.Floor(val))
				}
				deviation := 0
				val, err = strconv.ParseFloat(row["STD(bits)"].(string), 64)
				if err != nil {
					log.Errorf(err.Error())
				} else {
					deviation = int(math.Floor(val))
				}
				sql = "SELECT id FROM `apollo`.`deviations` WHERE date = ? AND platform = ? AND metric = ? LIMIT 1"
				num, err = o.Raw(sql, date+" "+ticker, platformName, metricKey).Values(&rows)
				if err != nil {
					log.Errorf(err.Error())
				} else if num == 0 {
					sql = "INSERT INTO `apollo`.`deviations`(`date`, `platform`, `metric`,"
					sql += "`samples`, `mean`, `deviation`, `updated`) VALUES("
					sql += "?, ?, ?, ?, ?, ?, ?)"
					_, err := o.Raw(sql, date+" "+ticker, platformName, metricKey, samples, mean, deviation,
						getNow()).Exec()
					if err != nil {
						log.Errorf(err.Error())
					}
				}
			}
		}
	}
}

func syncDeviationsTable() {
	platformNames := []string{}
	platformsMap := map[string]map[string]string{}
	o := orm.NewOrm()
	o.Using("apollo")
	bo := orm.NewOrm()
	bo.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM `apollo`.`deviations` ORDER BY updated DESC LIMIT 1"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		format := "2006-01-02 15:04:05"
		updatedTime, _ := time.Parse(format, rows[0]["updated"].(string))
		currentTime, _ := time.Parse(format, getNow())
		diff := currentTime.Unix() - updatedTime.Unix()
		if int(diff) < g.Config().Contacts.Interval {
			return
		}
	}
	sql = "SELECT platform, principal FROM `boss`.`platforms` WHERE type LIKE '%业务' AND visible = 1 AND count > 0 ORDER BY platform ASC"
	num, err = bo.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		for _, row := range rows {
			platformName := row["platform"].(string)
			platformsMap[platformName] = map[string]string{
				"contact": row["principal"].(string),
			}
			platformNames = append(platformNames, platformName)
		}
	}
	hours := []int{}
	for hour := 0; hour < 24; hour++ {
		hours = append(hours, hour)
	}
	minutes := []int{0, 20, 40}
	for _, platformName := range platformNames {
		for i := 0; i < 30; i++ {
			offset := i * (-1)
			date := time.Now().AddDate(0, 0, offset).Format("2006-01-02")
			dateFull := date + " 00:00:00"
			sql = "SELECT DISTINCT date FROM `apollo`.`net` "
			sql += "WHERE platform = ? AND hour = ? AND minute = ? "
			sql += "AND date >= DATE_SUB(?, INTERVAL 7 DAY) "
			sql += "AND date < DATE_SUB(?, INTERVAL 1 DAY) ORDER BY date DESC"
			num, err = o.Raw(sql, platformName, 0, 0, dateFull, dateFull).Values(&rows)
			if err != nil {
				log.Errorf(err.Error())
				break
			} else if num > 1 {
				for _, hour := range hours {
					for _, minute := range minutes {
						ticker := strconv.Itoa(hour) + ":"
						if hour < 10 {
							ticker = "0" + ticker
						}
						if minute == 0 {
							ticker += "00"
						} else {
							ticker += strconv.Itoa(minute)
						}
						dateQuery := date + " " + ticker + "%"
						sql = "SELECT date FROM `apollo`.`deviations` WHERE platform = ? AND date LIKE ? LIMIT 1"
						num, err = o.Raw(sql, platformName, dateQuery).Values(&rows)
						if err != nil {
							log.Errorf(err.Error())
						} else if num == 0 {
							writeToDeviationsTable(platformName, hour, minute, date, ticker)
						}
					}
				}
			} else {
				break
			}
		}
	}
}

func writeToNetTable(platformName string, offset int) {
	hours := []int{}
	for hour := 0; hour < 24; hour++ {
		hours = append(hours, hour)
	}
	minutes := []int{0, 20, 40}
	o := orm.NewOrm()
	o.Using("apollo")
	var rows []orm.Params
	data, date, counts := getPlatformsDailyTrafficData(platformName, offset)
	metrics := []string{
		"in",
		"out",
	}
	for metricKey, metric := range metrics {
		for _, hour := range hours {
			for _, minute := range minutes {
				ticker := strconv.Itoa(hour) + ":"
				if hour < 10 {
					ticker = "0" + ticker
				}
				if minute == 0 {
					ticker += "00"
				} else {
					ticker += strconv.Itoa(minute)
				}
				bits := 0
				if val, ok := data[metric][ticker]; ok {
					bits = val
				}
				sql := "SELECT id, date, hour, minute, platform, metric, count FROM `apollo`.`net` "
				sql += "WHERE date = ? AND hour = ? AND minute = ? AND platform = ? AND metric = ? LIMIT 1"
				num, err := o.Raw(sql, date, hour, minute, platformName, metricKey).Values(&rows)
				if err != nil {
					log.Errorf(err.Error())
				} else if num == 0 {
					sql = "INSERT INTO `apollo`.`net`(`date`, `hour`, `minute`,"
					sql += "`platform`, `metric`, `count`, `bits`,"
					sql += "`updated`) VALUES("
					sql += "?, ?, ?, ?, ?, ?, ?, ?)"
					_, err := o.Raw(sql, date, hour, minute, platformName,
						metricKey, counts[metric], bits,
						getNow()).Exec()
					if err != nil {
						log.Errorf(err.Error())
					}
				} else if num > 0 {
					count, _ := strconv.Atoi(rows[0]["count"].(string))
					if count < counts[metric] {
						ID := rows[0]["id"]
						sql := "UPDATE `apollo`.`net`"
						sql += " SET `date` = ?, `hour` = ?, `minute` = ?,"
						sql += " `platform` = ?, `metric` = ?, `count` = ?,"
						sql += " `bits` = ?, `updated` = ?"
						sql += " WHERE id = ?"
						_, err := o.Raw(sql, date, hour, minute, platformName,
							metricKey, counts[metric], bits,
							getNow(), ID).Exec()
						if err != nil {
							log.Errorf(err.Error())
						}
					}
				}
			}
		}
	}
}

func syncNetTable() {
	o := orm.NewOrm()
	o.Using("apollo")
	bo := orm.NewOrm()
	bo.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM `apollo`.`net` ORDER BY updated DESC LIMIT 1"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		format := "2006-01-02 15:04:05"
		updatedTime, _ := time.Parse(format, rows[0]["updated"].(string))
		currentTime, _ := time.Parse(format, getNow())
		diff := currentTime.Unix() - updatedTime.Unix()
		if int(diff) < g.Config().Contacts.Interval {
			return
		}
	}
	platformNames := []string{}
	sql = "SELECT platform, count FROM `boss`.`platforms` WHERE type LIKE '%业务' AND visible = 1 AND count > 0 ORDER BY platform ASC"
	num, err = bo.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		for _, row := range rows {
			platformName := row["platform"].(string)
			platformNames = append(platformNames, platformName)
		}
	}
	for _, platformName := range platformNames {
		for i := 1; i < 7; i++ {
			hostCountOfData := 0
			offset := i * (-1)
			date := time.Now().AddDate(0, 0, offset).Format("2006-01-02")
			sql = "SELECT MIN(count) FROM `apollo`.`net` "
			sql += "WHERE platform = ? AND date LIKE ?"
			num, err = o.Raw(sql, platformName, date+"%").Values(&rows)
			if err != nil {
				log.Errorf(err.Error())
			} else if num > 0 {
				if val, ok := rows[0]["MIN(count)"]; ok {
					if val != nil {
						value, err := strconv.Atoi(val.(string))
						if err == nil {
							hostCountOfData = value
						}
					}
				}
			}
			if hostCountOfData == 0 {
				writeToNetTable(platformName, i)
			} else {
				hostCountOfPlatform := 0
				sql = "SELECT DISTINCT hostname FROM `boss`.`ips` "
				sql += "WHERE platform = ? AND exist = 1 ORDER BY hostname ASC"
				num, err = bo.Raw(sql, platformName).Values(&rows)
				if err != nil {
					log.Errorf(err.Error())
				} else {
					hostCountOfPlatform = int(num)
				}
				if (hostCountOfData < hostCountOfPlatform) && (hostCountOfPlatform > 0) {
					writeToNetTable(platformName, i)
				}
			}
		}
	}
}

func syncHostsTable() {
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM `boss`.`ips` WHERE exist = 1 ORDER BY updated DESC LIMIT 1"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		format := "2006-01-02 15:04:05"
		updatedTime, _ := time.Parse(format, rows[0]["updated"].(string))
		currentTime, _ := time.Parse(format, getNow())
		diff := currentTime.Unix() - updatedTime.Unix()
		if int(diff) < g.Config().Hosts.Interval {
			return
		}
	}
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	getPlatformJSON(nodes, result)
	if nodes["status"] == nil {
		return
	} else if int(nodes["status"].(float64)) != 1 {
		return
	}
	platformNames := []string{}
	platformsMap := map[string]map[string]string{}
	hostname := ""
	hostnames := []string{}
	hostsMap := map[string]map[string]string{}
	IPs := []string{}
	IPKeys := []string{}
	IPsMap := map[string]map[string]string{}
	idcIDs := []string{}
	for _, platform := range nodes["result"].([]interface{}) {
		platformName := platform.(map[string]interface{})["platform"].(string)
		platformNames = appendUniqueString(platformNames, platformName)
		for _, device := range platform.(map[string]interface{})["ip_list"].([]interface{}) {
			hostname = device.(map[string]interface{})["hostname"].(string)
			IP := device.(map[string]interface{})["ip"].(string)
			status := device.(map[string]interface{})["ip_status"].(string)
			IPType := device.(map[string]interface{})["ip_type"].(string)
			item := map[string]string{
				"IP":       IP,
				"status":   status,
				"hostname": hostname,
				"platform": platformName,
				"type":     strings.ToLower(IPType),
			}
			IPs = append(IPs, IP)
			IPKey := platformName + "_" + IP
			IPKeys = append(IPKeys, IPKey)
			if _, ok := IPsMap[IP]; !ok {
				IPsMap[IPKey] = item
			}
			if len(hostname) > 0 {
				if host, ok := hostsMap[hostname]; !ok {
					hostnames = append(hostnames, hostname)
					idcID := device.(map[string]interface{})["pop_id"].(string)
					host := map[string]string{
						"hostname":  hostname,
						"activate":  "0",
						"platforms": "",
						"idcID":     idcID,
						"IP":        IP,
					}
					if len(IP) > 0 && IP == getIPFromHostname(hostname, result) {
						host["IP"] = IP
						host["platform"] = platformName
						platforms := []string{}
						if len(host["platforms"]) > 0 {
							platforms = strings.Split(host["platforms"], ",")
						}
						platforms = appendUniqueString(platforms, platformName)
						host["platforms"] = strings.Join(platforms, ",")
					}
					if status == "1" {
						host["activate"] = "1"
					}
					hostsMap[hostname] = host
					idcIDs = appendUniqueString(idcIDs, idcID)
				} else {
					if len(IP) > 0 && IP == getIPFromHostname(hostname, result) {
						host["IP"] = IP
						host["platform"] = platformName
						platforms := []string{}
						if len(host["platforms"]) > 0 {
							platforms = strings.Split(host["platforms"], ",")
						}
						platforms = appendUniqueString(platforms, platformName)
						host["platforms"] = strings.Join(platforms, ",")
					}
					if status == "1" {
						host["activate"] = "1"
					}
					hostsMap[hostname] = host
				}
			}
		}
		platformsMap[platformName] = map[string]string{
			"platformName": platformName,
			"type":         "",
			"visible":      "",
			"department":   "",
			"team":         "",
			"description":  "",
		}
	}
	sort.Strings(IPs)
	sort.Strings(IPKeys)
	sort.Strings(hostnames)
	sort.Strings(platformNames)
	log.Debugf("platformNames =", platformNames)
	updateIPsTable(IPKeys, IPsMap)
	updateHostsTable(hostnames, hostsMap)
	platformsMap = getPlatformsType(nodes, result, platformsMap)
	updatePlatformsTable(platformNames, platformsMap)
	muteFalconHostTable(hostnames, hostsMap)
}

func syncContactsTable() {
	log.Debugf("func syncContactsTable()")
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM `boss`.`contacts` ORDER BY updated DESC LIMIT 1"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		format := "2006-01-02 15:04:05"
		updatedTime, _ := time.Parse(format, rows[0]["updated"].(string))
		currentTime, _ := time.Parse(format, getNow())
		diff := currentTime.Unix() - updatedTime.Unix()
		if int(diff) < g.Config().Contacts.Interval {
			return
		}
	}
	platformNames := []string{}
	sql = "SELECT DISTINCT platform FROM boss.platforms ORDER BY platform ASC"
	num, err = o.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		for _, row := range rows {
			platformNames = append(platformNames, row["platform"].(string))
		}
	}

	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	getPlatformContact(strings.Join(platformNames, ","), nodes)
	contactNames := []string{}
	contactsMap := map[string]map[string]string{}
	contacts := nodes["result"].(map[string]interface{})["items"].(map[string]interface{})
	for _, platformName := range platformNames {
		if items, ok := contacts[platformName]; ok {
			for _, user := range items.(map[string]map[string]string) {
				contactName := user["name"]
				if _, ok := contactsMap[contactName]; !ok {
					contactsMap[contactName] = user
					contactNames = append(contactNames, contactName)
				}
			}
		}
	}
	sort.Strings(contactNames)
	updateContactsTable(contactNames, contactsMap)
	addContactsToPlatformsTable(contacts)
}

func addContactsToPlatformsTable(contacts map[string]interface{}) {
	log.Debugf("func addContactsToPlatformsTable()")
	now := getNow()
	o := orm.NewOrm()
	o.Using("boss")
	var platforms []Platforms
	_, err := o.QueryTable("platforms").All(&platforms)
	if err != nil {
		log.Errorf(err.Error())
	} else {
		for _, platform := range platforms {
			platformName := platform.Platform
			if items, ok := contacts[platformName]; ok {
				contacts := []string{}
				for role, user := range items.(map[string]map[string]string) {
					if role == "principal" {
						platform.Principal = user["name"]
					} else if role == "deputy" {
						platform.Deputy = user["name"]
					} else if role == "upgrader" {
						platform.Upgrader = user["name"]
					}
				}
				if len(platform.Principal) > 0 {
					contacts = append(contacts, platform.Principal)
				}
				if len(platform.Deputy) > 0 {
					contacts = append(contacts, platform.Deputy)
				}
				if len(platform.Upgrader) > 0 {
					contacts = append(contacts, platform.Upgrader)
				}
				platform.Contacts = strings.Join(contacts, ",")
			}
			platform.Updated = now
			_, err := o.Update(&platform)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}

func updateContactsTable(contactNames []string, contactsMap map[string]map[string]string) {
	log.Debugf("func updateContactsTable()")
	o := orm.NewOrm()
	o.Using("boss")
	var contact Contacts
	for _, contactName := range contactNames {
		user := contactsMap[contactName]
		err := o.QueryTable("contacts").Filter("name", user["name"]).One(&contact)
		if err == orm.ErrNoRows {
			sql := "INSERT INTO `boss`.`contacts`(name, phone, email, updated) VALUES(?, ?, ?, ?)"
			_, err := o.Raw(sql, user["name"], user["phone"], user["email"], getNow()).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		} else if err != nil {
			log.Errorf(err.Error())
		} else {
			contact.Email = user["email"]
			contact.Phone = user["phone"]
			contact.Updated = getNow()
			_, err := o.Update(&contact)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}

func updateIDCsTable(IDCNames []string, IDCsMap map[string]map[string]string) {
	log.Debugf("func updateIDCsTable()")
	now := getNow()
	o := orm.NewOrm()
	o.Using("boss")
	var idc Idcs
	for _, IDCName := range IDCNames {
		item := IDCsMap[IDCName]
		err := o.QueryTable("idcs").Filter("idc", IDCName).One(&idc)
		if err == orm.ErrNoRows {
			sql := "INSERT INTO `boss`.`idcs`(popid, idc, bandwidth, count, area, province, city, updated) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
			_, err := o.Raw(sql, item["popid"], item["idc"], item["bandwidth"], item["count"], item["area"], item["province"], item["city"], now).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		} else if err != nil {
			log.Errorf(err.Error())
		} else {
			popID, _ := strconv.Atoi(item["popid"])
			bandwidth, _ := strconv.Atoi(item["bandwidth"])
			count, _ := strconv.Atoi(item["count"])
			idc.Popid = popID
			idc.Idc = item["idc"]
			idc.Bandwidth = bandwidth
			idc.Count = count
			idc.Area = item["area"]
			idc.Province = item["province"]
			idc.City = item["city"]
			idc.Updated = now
			_, err := o.Update(&idc)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}

func updateIPsTable(IPNames []string, IPsMap map[string]map[string]string) {
	log.Debugf("func updateIPsTable()")
	now := getNow()
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM `boss`.`ips` WHERE exist = 1 ORDER BY updated DESC LIMIT 1"
	num, err := o.Raw(sql).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		format := "2006-01-02 15:04:05"
		updatedTime, _ := time.Parse(format, rows[0]["updated"].(string))
		currentTime, _ := time.Parse(format, getNow())
		diff := currentTime.Unix() - updatedTime.Unix()
		if int(diff) < g.Config().Hosts.Interval {
			return
		}
	}
	for _, IPName := range IPNames {
		item := IPsMap[IPName]
		sql := "SELECT id FROM boss.ips WHERE ip = ? AND platform = ? LIMIT 1"
		num, err := o.Raw(sql, item["IP"], item["platform"]).Values(&rows)
		if num == 0 {
			status, _ := strconv.Atoi(item["status"])
			sql := "INSERT INTO boss.ips("
			sql += "ip, exist, status, type, hostname, platform, updated) "
			sql += "VALUES(?, ?, ?, ?, ?, ?, ?)"
			_, err := o.Raw(sql, item["IP"], 1, status, item["type"], item["hostname"], item["platform"], now).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		} else if err != nil {
			log.Errorf(err.Error())
		} else if num > 0 {
			row := rows[0]
			ID := row["id"]
			status, _ := strconv.Atoi(item["status"])
			sql := "UPDATE boss.ips"
			sql += " SET ip = ?, exist = ?, status = ?, type = ?,"
			sql += " hostname = ?, platform = ?, updated = ?"
			sql += " WHERE id = ?"
			_, err := o.Raw(sql, item["IP"], 1, status, item["type"], item["hostname"], item["platform"], now, ID).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}

	sql = "SELECT id FROM boss.ips WHERE exist = ?"
	sql += " AND updated <= DATE_SUB(CONVERT_TZ(NOW(),@@session.time_zone,'+08:00'),"
	sql += " INTERVAL 10 MINUTE) LIMIT 30"
	num, err = o.Raw(sql, 1).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
	} else if num > 0 {
		for _, row := range rows {
			ID := row["id"]
			sql = "UPDATE boss.ips"
			sql += " SET exist = ?"
			sql += " WHERE id = ?"
			_, err := o.Raw(sql, 0, ID).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}

func updateHostsTable(hostnames []string, hostsMap map[string]map[string]string) {
	log.Debugf("func updateHostsTable()")
	idcMap := getIDCMap()
	hosts := []map[string]string{}
	for _, hostname := range hostnames {
		host := hostsMap[hostname]
		if len(host["platform"]) == 0 {
			host["platform"] = strings.Split(host["platforms"], ",")[0]
		}
		ISP := ""
		str := strings.Replace(host["hostname"], "_", "-", -1)
		slice := strings.Split(str, "-")
		if len(slice) >= 4 {
			ISP = slice[0]
		}
		if len(ISP) > 5 {
			ISP = ""
		}
		host["ISP"] = ISP
		idcID := host["idcID"]
		if idc, ok := idcMap[idcID]; ok {
			host["IDC"] = idc.(IDCMapItem).Idc
			host["province"] = idc.(IDCMapItem).Province
			host["city"] = idc.(IDCMapItem).City
		}
		hosts = append(hosts, host)
	}

	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	num, err := o.Raw("SELECT * FROM hosts limit 1;").Values(&rows)
	if num > 0 { // not empty
		sql := fmt.Sprintf("SELECT * FROM hosts WHERE exist = 1 AND updated <= DATE_SUB(NOW(), INTERVAL %v SECOND) limit 1;", g.Config().Hosts.Interval)
		numStale, err := o.Raw(sql).Values(&rows)
		if err != nil {
			log.Errorf(err.Error())
			return
		} else if numStale == 0 {
			log.Debugln("boss.hosts have no rows out of date.")
			return
		}
	}
	// do real update and insert
	// put data into temporary table
	_, err = o.Raw("DROP TABLE IF EXISTS tempBossHosts;").Exec()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	_, err = o.Raw("CREATE TABLE tempBossHosts LIKE hosts;").Exec()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	now := time.Now().Unix()
	// SQL prepare statement
	// batchSize := 32
	sql := `
		INSERT INTO tempBossHosts(
			hostname, exist, activate, platform, platforms, idc, ip, isp, province, city, updated
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, FROM_UNIXTIME(?))
	`
	o.Begin()
	p, err := o.Raw(sql).Prepare()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	// use transition to batch insert multiple values
	for _, host := range hosts {
		_, err := p.Exec(
			host["hostname"], 1, host["activate"], host["platform"], host["platforms"], host["IDC"],
			host["IP"], host["ISP"], host["province"], host["city"], now)

		if err != nil {
			log.Errorf(err.Error())
			o.Rollback()
			return
		}

	}
	o.Commit()
	// begin transaction to join base table with temporary table
	err = o.Begin()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	// SQL update statement
	sql = `
		UPDATE hosts
			INNER JOIN tempBossHosts t
			ON hosts.hostname = t.hostname
		SET hosts.exist = 1, hosts.activate = t.activate, hosts.platform = t.platform,
			hosts.platforms = t.platforms, hosts.idc = t.idc, hosts.ip = t.ip,
			hosts.isp = t.isp, hosts.province = t.province, hosts.city = t.city,
			hosts.updated = t.updated
	`
	if _, UpdateError := o.Raw(sql).Exec(); UpdateError != nil {
		log.Errorf("Update hosts has error: %v", UpdateError)
		o.Rollback()
		return
	}
	// SQL insert statement
	sql = `
		INSERT INTO hosts(
			hostname, exist, activate, platform, platforms, idc, ip, isp, province, city, updated
		)
		SELECT t.hostname, t.exist, t.activate, t.platform, t.platforms, t.idc,
			t.ip, t.isp, t.province, t.city, t.updated
		FROM tempBossHosts t
			LEFT JOIN hosts
			ON t.hostname = hosts.hostname
		WHERE hosts.hostname IS NULL
	`
	if _, InsertError := o.Raw(sql).Exec(); InsertError != nil {
		log.Errorf("Insert new data into hosts has error: %v", InsertError)
		o.Rollback()
		return
	}
	// SQL update exist property
	sql = fmt.Sprintf("UPDATE hosts SET exist = 0, updated = FROM_UNIXTIME(%d) WHERE exist = 1 AND updated <= DATE_SUB(NOW(), INTERVAL 10 MINUTE);", now)
	if _, StaleError := o.Raw(sql).Exec(); StaleError != nil {
		log.Errorf("Update hosts to non-exist has error: %v", StaleError)
		o.Rollback()
		return
	}

	// Commit or Rollback
	err = o.Commit()
	if err != nil {
		log.Errorf("Commit has error: %v", err)
	}

	_, err = o.Raw("DROP TABLE IF EXISTS tempBossHosts;").Exec()
	if err != nil {
		log.Errorf("DROP TEMP TABLE tempBossHosts has error: %v", err)
	}
}

func muteFalconHostTable(hostnames []string, hostsMap map[string]map[string]string) {
	log.Debugf("func muteFalconHostTable()")
	o := orm.NewOrm()
	o.Using("default")
	var rows []orm.Params
	now := getNow()
	for _, hostname := range hostnames {
		host := hostsMap[hostname]
		sql := "SELECT id FROM `falcon_portal`.`host` WHERE hostname = ? LIMIT 1"
		num, err := o.Raw(sql, host["hostname"]).Values(&rows)
		if err != nil {
			log.Errorf(err.Error())
		} else if num > 0 {
			activate := host["activate"]
			if activate == "0" || activate == "1" {
				begin := int64(0)
				end := int64(0)
				if activate == "0" {
					begin = int64(946684800) // Sat, 01 Jan 2000 00:00:00 GMT
					end = int64(4292329420)  // Thu, 07 Jan 2106 17:43:40 GMT
				}
				row := rows[0]
				ID := row["id"]
				sql = "UPDATE falcon_portal.host"
				sql += " SET maintain_begin = ?, maintain_end = ?, update_at = ?"
				sql += " WHERE id = ?"
				_, err := o.Raw(sql, begin, end, now, ID).Exec()
				if err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}
}

func updatePlatformsTable(platformNames []string, platformsMap map[string]map[string]string) {
	log.Debugf("func updatePlatformsTable()")
	now := getNow()
	o := orm.NewOrm()
	o.Using("boss")
	var platform Platforms
	var rows []orm.Params
	sql := "SELECT DISTINCT hostname FROM `boss`.`ips`"
	sql += " WHERE platform = ? AND exist = 1 ORDER BY hostname ASC"
	sqlInsert := "INSERT INTO `boss`.`platforms`"
	sqlInsert += "(platform, type, visible, count, department, team, description, updated) "
	sqlInsert += "VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	for _, platformName := range platformNames {
		count, err := o.Raw(sql, platformName).Values(&rows)
		if err != nil {
			count = 0
			log.Errorf(err.Error())
		}
		group := platformsMap[platformName]
		err = o.QueryTable("platforms").Filter("platform", group["platformName"]).One(&platform)
		if err == orm.ErrNoRows {
			_, err := o.Raw(sqlInsert, group["platformName"], group["type"], group["visible"], count, group["department"], group["team"], group["description"], now).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		} else if err != nil {
			log.Errorf(err.Error())
		} else {
			platform.Platform = group["platformName"]
			if len(group["type"]) > 0 {
				platform.Type = group["type"]
			}
			if len(group["visible"]) > 0 {
				platform.Visible = 0
				if group["visible"] == "1" {
					platform.Visible = 1
				}
			}
			platform.Count = int(count)
			if len(group["department"]) > 0 {
				platform.Department = group["department"]
			}
			platform.Team = group["team"]
			platform.Description = group["description"]
			platform.Updated = now
			_, err := o.Update(&platform)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}
