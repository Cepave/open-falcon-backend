package http

import (
	"github.com/astaxie/beego/orm"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	log "github.com/Sirupsen/logrus"
)

type Nqm_node struct {
	Id         int
	Nid        string
	Pid        string
	Pname      string
	Pname_abbr string
	Cid        string
	Cname      string
	Cname_abbr string
	Iid        int
	Iname      string
	Status     int
	Note       string
	Addtime    int
}

func getNQMNodes(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}
	o := orm.NewOrm()
	o.Using("gz_nqm")
	var NQMNodes []*Nqm_node
	_, err := o.Raw("SELECT nid, pname, cname, status, note FROM `gz_nqm`.`nqm_node` ORDER BY nid ASC").QueryRows(&NQMNodes)
	if err != nil {
		log.Debugf("Error = %v", err.Error())
	} else {
		for _, node := range NQMNodes {
			idc := map[string]string{
				"idc":      node.Nid,
				"status":   strconv.Itoa(node.Status),
				"province": node.Pname,
				"city":     node.Cname,
				"note":     node.Note,
			}
			items = append(items, idc)
		}
	}
	result["items"] = items
	nodes["result"] = result
	nodes["count"] = len(items)
	setResponse(rw, nodes)
}

func getLatestTimestamp(tableName string, result map[string]interface{}) int64 {
	timestamp := int64(0)
	o := orm.NewOrm()
	o.Using("gz_nqm")
	sqlcmd := "SELECT mtime FROM `gz_nqm`.`" + tableName + "` ORDER BY mtime DESC LIMIT 1"
	var rows []orm.Params
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		log.Debugf("Error = %v", err.Error())
	} else if num > 0 {
		mtime, err := strconv.Atoi(rows[0]["mtime"].(string))
		if err == nil {
			timestamp = int64(mtime)
		}
	}
	return timestamp
}

func getNearestTimestamp(tableName string, bound int64, result map[string]interface{}) int64 {
	timestamp := int64(0)
	o := orm.NewOrm()
	o.Using("gz_nqm")
	sqlcmd := "SELECT mtime FROM `gz_nqm`.`" + tableName
	sqlcmd += "` WHERE mtime <= ? ORDER BY mtime DESC LIMIT 1"
	var rows []orm.Params
	num, err := o.Raw(sqlcmd, bound).Values(&rows)
	if err != nil {
		log.Debugf("Error = %v", err.Error())
	} else if num > 0 {
		mtime, err := strconv.Atoi(rows[0]["mtime"].(string))
		if err == nil {
			timestamp = int64(mtime)
		}
	}
	return timestamp
}

func getSum(slice []float64) float64 {
	sum := float64(0)
	for _, number := range slice {
		sum += number
	}
	return sum
}

func getPacketLossAndAveragePingTime(nodeName string, timestamps []int64) []map[string]interface{} {
	result := []map[string]interface{}{}
	tableName := "nqm_log_" + strings.Replace(nodeName, "-", "_", -1)
	for _, timestamp := range timestamps {
		sends := []float64{}
		receives := []float64{}
		averages := []float64{}
		o := orm.NewOrm()
		o.Using("gz_nqm")
		sqlcmd := "SELECT send, receive, avg FROM `gz_nqm`.`" + tableName + "` WHERE mtime = ?"
		var rows []orm.Params
		num, err := o.Raw(sqlcmd, strconv.Itoa(int(timestamp))).Values(&rows)
		if err != nil {
			log.Debugf("Error = %v", err.Error())
		} else if num > 0 {
			for _, row := range rows {
				send, err := strconv.ParseFloat(row["send"].(string), 64)
				if err != nil {
					log.Debugf("Error = %v", err.Error())
				} else {
					sends = append(sends, send)
				}
				receive, err := strconv.ParseFloat(row["receive"].(string), 64)
				if err != nil {
					log.Debugf("Error = %v", err.Error())
				} else {
					receives = append(receives, receive)
				}
				avg, err := strconv.ParseFloat(row["avg"].(string), 64)
				if err != nil {
					log.Debugf("Error = %v", err.Error())
				} else {
					averages = append(averages, avg)
				}
			}
		}
		item := map[string]interface{}{
			"node":                        nodeName,
			"packetLossRate":              "",
			"averagePingTimeMilliseconds": "",
			"time": time.Unix(timestamp, 0).Format("2006-01-02 15:04"),
		}
		divider := float64(len(sends))
		if divider > 0 {
			packetLossRate := 1 - (getSum(receives) / getSum(sends))
			averagePingTimeMilliseconds := getSum(averages) / divider
			item["packetLossRate"] = packetLossRate
			item["averagePingTimeMilliseconds"] = averagePingTimeMilliseconds
		}
		result = append(result, item)
	}
	return result
}

func getNQMPacketLoss(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}

	tablesMap := map[string]interface{}{}
	o := orm.NewOrm()
	o.Using("gz_nqm")
	sqlcmd := "SELECT table_name FROM information_schema.tables WHERE table_name LIKE 'nqm_log_prober_%'"
	var rows []orm.Params
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		log.Debugf("Error = %v", err.Error())
	} else if num > 0 {
		for _, row := range rows {
			tablesMap[row["table_name"].(string)] = ""
		}
	}

	nids := []string{}
	pidMap := map[string]interface{}{}
	var NQMNodes []*Nqm_node
	_, err = o.Raw("SELECT nid, pid, status FROM `gz_nqm`.`nqm_node` ORDER BY nid ASC").QueryRows(&NQMNodes)
	if err != nil {
		log.Debugf("Error = %v", err.Error())
	} else {
		for _, node := range NQMNodes {
			if node.Status > 0 {
				nids = append(nids, node.Nid)
				pidMap[node.Nid] = node.Pid
			}
		}
	}
	idcNames := []string{}
	for _, nid := range nids {
		if _, ok := tablesMap["nqm_log_prober_"+strings.Replace(nid, "-", "_", -1)]; ok {
			idcNames = append(idcNames, nid)
		}
	}
	sort.Strings(idcNames)
	for _, idcName := range idcNames {
		tableName := "nqm_log_" + strings.Replace(idcName, "-", "_", -1)
		timestamp := getLatestTimestamp(tableName, result)
		if timestamp > 0 {
			result := getPacketLossAndAveragePingTime(idcName, []int64{timestamp})
			idc := result[0]
			idc["nodeName"] = idcName
			idc["pid"] = pidMap[idcName]
			items = append(items, idc)
		}
	}
	result["items"] = items
	nodes["result"] = result
	nodes["count"] = len(items)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getJaguar(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []interface{}{}
	timestamp := int64(0)
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		loc = time.Local
	}
	timeFormat := "2006-01-02 15:04"
	timeInput := req.URL.Query().Get("time")
	date, err := time.ParseInLocation(timeFormat, timeInput, loc)
	if err == nil {
		timestamp = date.Unix()
	}
	nodeMap := map[string]map[string]string{}
	nodeNames := []string{}
	o := orm.NewOrm()
	o.Using("gz_nqm")
	var NQMNodes []*Nqm_node
	_, err = o.Raw("SELECT nid, pname, cname, iname FROM `gz_nqm`.`nqm_node` ORDER BY nid ASC").QueryRows(&NQMNodes)
	if err != nil {
		log.Debugf("Error = %v", err.Error())
	} else {
		for _, node := range NQMNodes {
			nodeName := node.Nid
			node := map[string]string{
				"node":     nodeName,
				"province": node.Pname,
				"city":     node.Cname,
				"isp":      node.Iname,
			}
			nodeMap[nodeName] = node
			nodeNames = append(nodeNames, nodeName)
		}
	}
	sqlcmd := "SELECT nid, ip, note FROM `gz_nqm`.`nqm_dev` WHERE nid IN ('"
	sqlcmd += strings.Join(nodeNames, "','") + "')"
	var rows []orm.Params
	_, err = o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		log.Debugf("Error = %v", err.Error())
	} else {
		for _, row := range rows {
			nodeName := row["nid"].(string)
			nodeMap[nodeName]["ip"] = row["ip"].(string)
			nodeMap[nodeName]["idc"] = row["note"].(string)
		}
	}
	for _, nodeName := range nodeNames {
		node := nodeMap[nodeName]
		item := map[string]interface{}{
			"node":     node["node"],
			"province": node["province"],
			"city":     node["city"],
			"isp":      node["isp"],
			"ip":       node["ip"],
			"idc":      node["idc"],
			"loss":     nil,
			"ping.ms":  nil,
			"time":     nil,
		}
		tableName := "nqm_log_" + strings.Replace(nodeName, "-", "_", -1)
		timestampLatest := getLatestTimestamp(tableName, result)
		if timestampLatest > 0 {
			timestampNearest := timestampLatest
			if timestamp > 0 {
				timestampNearest = getNearestTimestamp(tableName, timestamp, result)
			}
			result := getPacketLossAndAveragePingTime(nodeName, []int64{timestampNearest})
			resp := result[0]
			if _, ok := resp["packetLossRate"]; ok {
				item["loss"] = resp["packetLossRate"]
				item["ping.ms"] = resp["averagePingTimeMilliseconds"]
				item["time"] = resp["time"]
			}
		}
		items = append(items, item)
	}
	result["items"] = items
	nodes["result"] = result
	nodes["count"] = len(items)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getTimestamps(tableName string, timestampFrom int64, timestampTo int64) []int64 {
	timestamps := []int64{}
	o := orm.NewOrm()
	o.Using("gz_nqm")
	sqlcmd := "SELECT DISTINCT mtime FROM `gz_nqm`.`" + tableName
	sqlcmd += "` WHERE mtime BETWEEN ? AND ? ORDER BY mtime ASC"
	var rows []orm.Params
	num, err := o.Raw(sqlcmd, timestampFrom, timestampTo).Values(&rows)
	if err != nil {
		log.Debugf("Error = %v", err.Error())
	} else if num > 0 {
		for _, row := range rows {
			timestamp, err := strconv.ParseInt(row["mtime"].(string), 10, 64)
			if err != nil {
				log.Debugf("Error = %v", err.Error())
			} else {
				timestamps = append(timestamps, timestamp)
			}
		}
	}
	return timestamps
}

func getSnorlax(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	items := []map[string]interface{}{}
	countOfTimestamps := 0
	nodeName := req.URL.Query().Get("node")
	from := req.URL.Query().Get("from")
	to := req.URL.Query().Get("to")
	page := 1
	if len(req.URL.Query().Get("page")) > 0 {
		pageInput, err := strconv.Atoi(req.URL.Query().Get("page"))
		if err == nil && pageInput > 0 {
			page = pageInput
		}
	}
	o := orm.NewOrm()
	o.Using("gz_nqm")
	tableName := "nqm_log_" + strings.Replace(nodeName, "-", "_", -1)
	timestampFrom := int64(0)
	timestampTo := int64(0)
	timestampLatest := getLatestTimestamp(tableName, result)
	if timestampLatest > 0 {
		timestampTo = timestampLatest
		loc, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			loc = time.Local
		}
		timeFormat := "2006-01-02 15:04"
		date, err := time.ParseInLocation(timeFormat, to, loc)
		if err == nil {
			timestampTo = date.Unix()
		}
		if timestampTo > 0 {
			timestampTo = getNearestTimestamp(tableName, timestampTo, result)
		}

		date, err = time.ParseInLocation(timeFormat, from, loc)
		if err == nil {
			timestampFrom = date.Unix()
		}
		if timestampFrom > 0 {
			timestampFrom = getNearestTimestamp(tableName, timestampFrom, result)
		}
		timestamps := getTimestamps(tableName, timestampFrom, timestampTo)
		countOfTimestamps = len(timestamps)
		rowsPerPage := 20
		begin := 0
		end := 19
		begin = (page - 1) * rowsPerPage
		end =  begin + rowsPerPage
		lastIndex := len(timestamps)
		if (begin > lastIndex) {
			begin = lastIndex
		}
		if (end > lastIndex) {
			end = lastIndex
		}
		timestamps = timestamps[begin:end]
		items = getPacketLossAndAveragePingTime(nodeName, timestamps)
	}
	result["items"] = items
	nodes["result"] = result
	nodes["count"] = len(items)
	nodes["countOfTimestamps"] = countOfTimestamps
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getPhoenix(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	node := map[string]string{}
	items := []map[string]string{}
	nodeName := req.URL.Query().Get("node")

	timestamp := int64(0)
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		loc = time.Local
	}
	timeFormat := "2006-01-02 15:04"
	timeInput := req.URL.Query().Get("time")
	date, err := time.ParseInLocation(timeFormat, timeInput, loc)
	if err == nil {
		timestamp = date.Unix()
	}
	tableName := "nqm_log_" + strings.Replace(nodeName, "-", "_", -1)
	log.Debugf("tableName = %v", tableName)
	timestampLatest := getLatestTimestamp(tableName, result)
	timestampNearest := int64(0)
	if timestampLatest > 0 {
		timestampNearest = timestampLatest
		if timestamp > 0 {
			timestampNearest = getNearestTimestamp(tableName, timestamp, result)
		}
	}

	if timestampNearest > 0 {
		log.Debugf("timestampNearest = %v", timestampNearest)
		o := orm.NewOrm()
		o.Using("gz_nqm")
		sqlcmd := "SELECT ip, dest_ip, dest_id, loss, max, min, avg FROM `gz_nqm`.`" + tableName + "` WHERE mtime = ?"
		log.Debugf("sqlcmd = %v", sqlcmd)
		var rows []orm.Params
		num, err := o.Raw(sqlcmd, strconv.Itoa(int(timestampNearest))).Values(&rows)
		if err != nil {
			log.Debugf("Error = %v", err.Error())
		} else if num > 0 {
			row := rows[0]
			node = map[string]string{
				"node": nodeName,
				"IP": row["ip"].(string),
			}
			for _, row := range rows {
				IP := row["dest_ip"].(string)
				destination := map[string]string{
					"IDC": row["dest_id"].(string),
					"IP": IP,
					"max": row["max"].(string),
					"min": row["min"].(string),
					"avg": row["avg"].(string),
					"loss": row["loss"].(string),
					"time": time.Unix(timestampNearest, 0).Format("2006-01-02 15:04"),
				}
				items = append(items, destination)
			}
		}
	}
	result["items"] = items
	nodes["result"] = result
	nodes["count"] = len(items)
	nodes["node"] = node
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func configNQMRoutes() {
	http.HandleFunc("/api/nqm/nodes", getNQMNodes)
	http.HandleFunc("/api/nqm/loss", getNQMPacketLoss)
	http.HandleFunc("/api/nqm/jaguar", getJaguar)
	http.HandleFunc("/api/snorlax", getSnorlax)
	http.HandleFunc("/api/phoenix", getPhoenix)
}
