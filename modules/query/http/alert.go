package http

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/astaxie/beego/orm"
)

func parsePlatformJSON(result map[string]interface{}) map[string]interface{} {
	var nodes = make(map[string]interface{})
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
							"platform": platformName,
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
		sort.Strings(hostnames)
		idcIDsMap := map[string]string{}
		idcNames := []string{}
		o := orm.NewOrm()
		var idcs []*Idc
		sqlcommand := "SELECT pop_id, name FROM grafana.idc ORDER BY pop_id ASC"
		_, err := o.Raw(sqlcommand).QueryRows(&idcs)
		if err != nil {
			setError(err.Error(), result)
		} else {
			for _, idc := range idcs {
				idcIDsMap[strconv.Itoa(idc.Pop_id)] = idc.Name
				idcNames = appendUniqueString(idcNames, idc.Name)
			}
		}
		sort.Strings(idcNames)
		for _, hostname := range hostnames {
			host := hosts[hostname].(map[string]interface{})
			idcID := host["idcID"].(string)
			if _, ok := idcIDsMap[idcID]; ok {
				idcName := idcIDsMap[idcID]
				host["idc"] = idcName
				delete(host, "idcID")
			}
		}
	}
	if _, ok := nodes["info"]; ok {
		delete(nodes, "info")
	}
	if _, ok := nodes["status"]; ok {
		delete(nodes, "status")
	}
	return hosts
}

func getSeverity(priority string) string {
	severity := "Lower"
	if priority == "0" {
		severity = "High"
	} else if priority == "1" {
		severity = "Medium"
	} else if priority == "2" || priority == "3" {
		severity = "Low"
	}
	return severity
}

func getStatus(statusRaw string) string {
	status := ""
	if statusRaw == "PROBLEM" {
		status = "Triggered"
	} else if statusRaw == "OK" {
		status = "Recovered"
	}
	return status
}

func getProcess(status string) string {
	process := ""
	if status == "Recovered" {
		process = "Resolved"
	} else if status == "Triggered" {
		process = "Unresolved"
	}
	return process
}

func getDuration(timeTriggered string, result map[string]interface{}) string {
	date, err := time.Parse("2006-01-02 15:04", timeTriggered)
	if err != nil {
		setError(err.Error(), result)
	}
	now := time.Now().Unix()
	diff := now - date.Unix()
	if diff <= 60 {
		return "just now"
	}
	if diff <= 120 {
		return "1 minute ago"
	}
	if diff <= 3600 {
		return fmt.Sprintf("%d minutes ago", diff/60)
	}
	if diff <= 7200 {
		return "1 hour ago"
	}
	if diff <= 3600*24 {
		return fmt.Sprintf("%d hours ago", diff/3600)
	}
	if diff <= 3600*24*2 {
		return "1 day ago"
	}
	return fmt.Sprintf("%d days ago", diff/3600/24)
}

func getUserID(username string, result map[string]interface{}) (string, string) {
	userID := ""
	userRole := ""
	o := orm.NewOrm()
	var rows []orm.Params
	sqlcmd := "SELECT id, role FROM uic.user WHERE name = ? LIMIT 1"
	num, err := o.Raw(sqlcmd, username).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			userID = row["id"].(string)
			userRole = row["role"].(string)
		}
	}
	return userID, userRole
}

func checkLoggedIn(userID string, sig string, result map[string]interface{}) bool {
	isLoggedIn := false
	o := orm.NewOrm()
	var rows []orm.Params
	expired := ""
	sqlcmd := "SELECT expired FROM uic.session WHERE uid = ? AND sig = ? LIMIT 1"
	num, err := o.Raw(sqlcmd, userID, sig).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			expired = row["expired"].(string)
		}
	}
	expiredInt, err := strconv.Atoi(expired)
	if err != nil {
		setError(err.Error(), result)
	} else {
		now := time.Now().Unix()
		isLoggedIn = int64(expiredInt) > now
	}
	return isLoggedIn
}

func getReceiverTeamIDs(userID string, result map[string]interface{}) []string {
	receiverTeamIDs := []string{}
	o := orm.NewOrm()
	var rows []orm.Params
	sqlcmd := "SELECT tid FROM uic.rel_team_user WHERE uid = ?"
	num, err := o.Raw(sqlcmd, userID).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			receiverTeamID := row["tid"].(string)
			receiverTeamIDs = append(receiverTeamIDs, receiverTeamID)
		}
	}
	return receiverTeamIDs
}

func getReceiverTeamNames(receiverTeamIDs []string, result map[string]interface{}) []string {
	receiverTeamNames := []string{}
	o := orm.NewOrm()
	var rows []orm.Params
	sqlcmd := "SELECT name FROM uic.team WHERE id IN ("
	sqlcmd += strings.Join(receiverTeamIDs, ",") + ")"
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			receiverTeamName := row["name"].(string)
			receiverTeamNames = append(receiverTeamNames, receiverTeamName)
		}
	}
	return receiverTeamNames
}

func getActionIDs(receiverTeamNames []string, result map[string]interface{}) []string {
	actionIDs := []string{}
	o := orm.NewOrm()
	var rows []orm.Params
	sqlcmd := "SELECT id FROM falcon_portal.action WHERE uic IN ('"
	sqlcmd += strings.Join(receiverTeamNames, "','") + "')"
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			actionID := row["id"].(string)
			actionIDs = append(actionIDs, actionID)
		}
	}
	return actionIDs
}

func getTemplateIDs(username string, sig string, result map[string]interface{}) string {
	templateIDs := []string{}
	userID, userRole := getUserID(username, result)
	if userID == "" {
		setError("User not found", result)
		return ""
	}
	isLoggedIn := checkLoggedIn(userID, sig, result)
	if !isLoggedIn {
		setError("Please log in first", result)
		return ""
	}
	if userRole == "1" || userRole == "2" { // admin user
		return "*"
	}
	receiverTeamIDs := getReceiverTeamIDs(userID, result)
	if len(receiverTeamIDs) == 0 {
		setError("User not subscribe any alerts", result)
		return ""
	}
	receiverTeamNames := getReceiverTeamNames(receiverTeamIDs, result)
	if len(receiverTeamNames) == 0 {
		setError("User not subscribe any alerts", result)
		return ""
	}
	actionIDs := getActionIDs(receiverTeamNames, result)
	if len(actionIDs) == 0 {
		setError("User not subscribe any alerts", result)
		return ""
	}
	o := orm.NewOrm()
	var rows []orm.Params
	sqlcmd := "SELECT id FROM falcon_portal.tpl WHERE action_id IN ("
	sqlcmd += strings.Join(actionIDs, ",") + ")"
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			templateID := row["id"].(string)
			templateIDs = append(templateIDs, templateID)
		}
	}
	return strings.Join(templateIDs, ",")
}

func getUsers(result map[string]interface{}) map[string]string {
	users := map[string]string{}
	o := orm.NewOrm()
	var rows []orm.Params
	sqlcmd := "SELECT id, name FROM uic.user"
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			userID := row["id"].(string)
			users[userID] = row["name"].(string)
		}
	}
	return users
}

func getNote(hash string, timestamp string) []map[string]string {
	o := orm.NewOrm()
	var rows []orm.Params
	queryStr := fmt.Sprintf(`SELECT note.id as id, note.event_caseId as event_caseId, note.note as note, note.case_id as case_id, note.status as status, note.timestamp as timestamp, user.name as name from
	(SELECT * from falcon_portal.event_note WHERE event_caseId = '%s' AND timestamp >= '%s')
	 note LEFT JOIN uic.user as user on note.user_id = user.id;`, hash, timestamp)

	num, err := o.Raw(queryStr).Values(&rows)
	notes := []map[string]string{}
	if err != nil {
		log.Error(err.Error())
	} else if num == 0 {
		return notes
	} else {
		for _, row := range rows {
			hash := row["event_caseId"].(string)
			time := row["timestamp"].(string)
			time = time[:len(time)-3]
			user := row["name"].(string)
			note := map[string]string{
				"note":   row["note"].(string),
				"status": row["status"].(string),
				"user":   user,
				"hash":   hash,
				"time":   time,
			}
			notes = append(notes, note)
		}
	}
	return notes
}

func setSQLQuery(templateIDs string, req *http.Request, result map[string]interface{}) string {
	sqlcmd := "SELECT id, endpoint, metric, func, cond, note, max_step, current_step, priority, status, "
	sqlcmd += "timestamp, update_at, template_id, tpl_creator, process_status "
	sqlcmd += "FROM falcon_portal.event_cases "
	whereConditions := []string{}
	query := req.URL.Query()
	if query.Get("start") != "" && query.Get("end") != "" {
		start := query.Get("start")
		end := query.Get("end")
		if start != "" && end != "" {
			timeCondition := "`update_at` BETWEEN FROM_UNIXTIME(" + start + ") AND FROM_UNIXTIME(" + end + ")"
			whereConditions = append(whereConditions, timeCondition)
		}
	}
	if query.Get("priority") != "" {
		priority := query.Get("priority")
		whereConditions = append(whereConditions, "`priority` = "+priority)
	}
	if query.Get("status") != "" && strings.Index(query.Get("status"), "ALL") == -1 {
		status := query.Get("status") // "PROBLEM", "OK"
		if strings.Index(status, ",") > -1 {
			status = strings.Replace(status, ",", "','", -1)
			whereConditions = append(whereConditions, "`status` IN ('"+status+"')")
		} else {
			whereConditions = append(whereConditions, "`status` = '"+status+"'")
		}
	}
	if query.Get("process") != "" {
		process := query.Get("process") //  "unresolved", "in progress", "resolved", "ignored"
		if strings.Index(process, ",") > -1 {
			process = strings.Replace(process, ",", "','", -1)
			whereConditions = append(whereConditions, "`process_status` IN ('"+process+"')")
		} else {
			whereConditions = append(whereConditions, "`process_status` = '"+process+"'")
		}
	}
	if query.Get("metric") != "" {
		metric := query.Get("metric") // "http.get.time", "net.port.listen/port=80"
		if strings.Index(metric, ",") > -1 {
			metric = strings.Replace(metric, ",", "','", -1)
			whereConditions = append(whereConditions, "`metric` IN ('"+metric+"')")
		} else {
			whereConditions = append(whereConditions, "`metric` = '"+metric+"'")
		}
	}
	limit := "500"
	if query.Get("limit") != "" {
		limit = query.Get("limit")
	}
	if templateIDs != "*" {
		whereConditions = append(whereConditions, "template_id IN ("+templateIDs+")")
	}
	if len(whereConditions) > 0 {
		conditions := strings.Join(whereConditions, " AND ")
		sqlcmd += "WHERE " + conditions
	}
	sqlcmd += " ORDER BY update_at DESC LIMIT " + limit
	return sqlcmd
}

func getEvents(hash string, eventsLimit string, result map[string]interface{}) []interface{} {
	events := []interface{}{}
	o := orm.NewOrm()
	var rows []orm.Params
	sqlcmd := "SELECT * "
	sqlcmd += "FROM falcon_portal.events "
	sqlcmd += "WHERE event_caseId = '" + hash + "' ORDER BY timestamp DESC LIMIT " + eventsLimit
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			event := map[string]interface{}{
				"id":        row["id"].(string),
				"condition": row["cond"].(string),
				"status":    row["status"].(string),
				"step":      row["step"].(string),
				"timestamp": row["timestamp"].(string),
			}
			events = append(events, event)
		}
	}
	return events
}

func queryAlerts(sqlcmd string, req *http.Request, result map[string]interface{}) []interface{} {
	alerts := []interface{}{}
	if sqlcmd == "" {
		return alerts
	}

	query := req.URL.Query()
	eventsLimit := "10"
	if query.Get("elimit") != "" {
		eventsLimit = query.Get("elimit")
	}

	o := orm.NewOrm()
	var rows []orm.Params
	num, err := o.Raw(sqlcmd).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			hash := row["id"].(string)
			hostname := row["endpoint"].(string)
			metric := row["metric"].(string)
			metricType := strings.Split(metric, ".")[0]
			content := row["note"].(string)
			priority := row["priority"].(string)
			statusRaw := row["status"].(string)
			timeStart := row["timestamp"].(string)
			timeStart = timeStart[:len(timeStart)-3]
			timeUpdate := row["update_at"].(string)
			timeUpdate = timeUpdate[:len(timeUpdate)-3]
			process := strings.ToLower(row["process_status"].(string))
			process = strings.Replace(process, process[:1], strings.ToUpper(process[:1]), 1)
			note := getNote(hash, row["timestamp"].(string))
			//this is a work around for auto clean process when the status is expired
			if len(note) == 0 {
				process = "unresolved"
			}
			process = strings.Replace(process, process[:1], strings.ToUpper(process[:1]), 1)
			templateID := row["template_id"].(string)
			author := row["tpl_creator"].(string)
			alert := map[string]interface{}{
				"hash":       hash,
				"hostname":   hostname,
				"metric":     metric,
				"author":     author,
				"templateID": templateID,
				"priority":   priority,
				"severity":   getSeverity(priority),
				"status":     getStatus(statusRaw),
				"statusRaw":  statusRaw,
				"type":       metricType,
				"content":    content,
				"timeStart":  timeStart,
				"timeUpdate": timeUpdate,
				"duration":   getDuration(timeUpdate, result),
				"notes":      note,
				"events":     getEvents(hash, eventsLimit, result),
				"process":    process,
				"function":   row["func"].(string),
				"condition":  row["cond"].(string),
				"stepLimit":  row["max_step"].(string),
				"step":       row["current_step"].(string),
			}
			alerts = append(alerts, alert)
		}
	}
	return alerts
}

func addPlatformToAlerts(alerts []interface{}, result map[string]interface{}, nodes map[string]interface{}, rw http.ResponseWriter) []interface{} {
	items := []interface{}{}
	hostnames := []string{}
	platformNames := []string{}
	hostsMap := parsePlatformJSON(result)
	for _, item := range alerts {
		hostname := item.(map[string]interface{})["hostname"].(string)
		if _, ok := hostsMap[hostname]; ok {
			host := hostsMap[hostname].(map[string]interface{})
			activate := host["activate"].(string)
			if activate == "1" {
				item.(map[string]interface{})["ip"] = host["ip"].(string)
				item.(map[string]interface{})["platform"] = host["platform"].(string)
				item.(map[string]interface{})["idc"] = host["idc"].(string)
				contact := map[string]string{
					"email": "",
					"name":  "",
					"phone": "",
				}
				contacts := []interface{}{
					contact,
				}
				item.(map[string]interface{})["contact"] = contacts
				items = append(items, item)
			} else {
				item.(map[string]interface{})["ip"] = host["ip"].(string) + " (deactivated)"
				item.(map[string]interface{})["platform"] = host["platform"].(string)
				item.(map[string]interface{})["idc"] = host["idc"].(string)
				contact := map[string]string{
					"email": "-",
					"name":  "-",
					"phone": "-",
				}
				contacts := []interface{}{
					contact,
				}
				item.(map[string]interface{})["contact"] = contacts
				items = append(items, item)
				log.Debugf("host deactivated: %v", hostname)
			}
		} else {
			item.(map[string]interface{})["ip"] = "not found"
			item.(map[string]interface{})["platform"] = "not found"
			item.(map[string]interface{})["idc"] = "not found"
			contact := map[string]string{
				"email": "-",
				"name":  "-",
				"phone": "-",
			}
			contacts := []interface{}{
				contact,
			}
			item.(map[string]interface{})["contact"] = contacts
			items = append(items, item)
			log.Debugf("hostname not found: %v", hostname)
		}
	}
	for _, item := range items {
		hostname := item.(map[string]interface{})["hostname"].(string)
		hostnames = appendUniqueString(hostnames, hostname)
	}
	sort.Strings(hostnames)

	hostsTriggeredMap := map[string]string{}
	for _, hostname := range hostnames {
		if _, ok := hostsMap[hostname]; ok {
			host := hostsMap[hostname].(map[string]interface{})
			platformName := host["platform"].(string)
			if strings.Index(platformName, ", ") > -1 {
				for _, name := range strings.Split(platformName, ", ") {
					platformNames = appendUniqueString(platformNames, name)
				}
			} else {
				platformNames = appendUniqueString(platformNames, platformName)
			}
			hostsTriggeredMap[hostname] = platformName
		}
	}
	sort.Strings(platformNames)
	getPlatformContact(strings.Join(platformNames, ","), nodes)
	platforms := nodes["result"].(map[string]interface{})["items"].(map[string]interface{})
	if len(platforms) > 0 {
		for _, item := range items {
			hostname := item.(map[string]interface{})["hostname"].(string)
			if _, ok := hostsMap[hostname]; ok {
				host := hostsMap[hostname].(map[string]interface{})
				platformName := host["platform"].(string)
				if strings.Index(platformName, ", ") > -1 {
					platformName = strings.Split(platformName, ", ")[0]
				}
				if contact, ok := platforms[platformName]; ok {
					item.(map[string]interface{})["contact"] = contact
				} else {
					item.(map[string]interface{})["contact"] = "BOSS 沒有平台負責人資訊"
				}
			}
		}
	}
	return items
}

func getAlertSeverityCounts(items []interface{}) map[string]int {
	count := map[string]int{
		"all":    len(items),
		"high":   0,
		"medium": 0,
		"low":    0,
		"lower":  0,
	}
	for _, item := range items {
		severity := item.(map[string]interface{})["severity"].(string)
		if severity == "High" {
			count["high"]++
		} else if severity == "Medium" {
			count["medium"]++
		} else if severity == "Low" {
			count["low"]++
		} else if severity == "Lower" {
			count["lower"]++
		}
	}
	return count
}

func getAlertProcessCounts(items []interface{}) map[string]int {
	count := map[string]int{
		"unresolved":  0,
		"in progress": 0,
		"resolved":    0,
		"ignored":     0,
	}
	for _, item := range items {
		process := strings.ToLower(item.(map[string]interface{})["process"].(string))
		if process == "unresolved" {
			count["unresolved"]++
		} else if process == "in progress" {
			count["in progress"]++
		} else if process == "resolved" {
			count["resolved"]++
		} else if process == "ignored" {
			count["ignored"]++
		}
	}
	return count
}

func getAlertMetricTypeCounts(items []interface{}) map[string]int {
	count := map[string]int{
		"cpu":          0,
		"disk":         0,
		"memory":       0,
		"net":          0,
		"others":       0,
		"agent":        0,
		"check":        0,
		"chk":          0,
		"dev":          0,
		"fcd":          0,
		"file":         0,
		"fm":           0,
		"http":         0,
		"nic":          0,
		"proc":         0,
		"tags":         0,
		"zabbix-agent": 0,
	}
	for _, item := range items {
		metric := item.(map[string]interface{})["metric"].(string)
		metricType := strings.ToLower(strings.Split(metric, ".")[0])
		if metricType == "cpu" {
			count["cpu"]++
		} else if metricType == "disk" {
			count["disk"]++
		} else if metricType == "memory" {
			count["memory"]++
		} else if metricType == "net" {
			count["net"]++
		} else {
			count["others"]++
			if metricType == "agent" {
				count["agent"]++
			} else if metricType == "check" {
				count["check"]++
			} else if metricType == "chk" {
				count["chk"]++
			} else if metricType == "dev" {
				count["dev"]++
			} else if metricType == "fcd" {
				count["fcd"]++
			} else if metricType == "file" {
				count["file"]++
			} else if metricType == "fm" {
				count["fm"]++
			} else if metricType == "http" {
				count["http"]++
			} else if metricType == "nic" {
				count["nic"]++
			} else if metricType == "proc" {
				count["proc"]++
			} else if metricType == "zabbix-agent" {
				count["zabbix-agent"]++
			} else if strings.Index(metricType, "tags") > -1 {
				count["tags"]++
			}
		}
	}
	return count
}

func getAlerts(rw http.ResponseWriter, req *http.Request) {
	var nodes = make(map[string]interface{})
	errors := []string{}
	var result = make(map[string]interface{})
	result["error"] = errors
	alerts := []interface{}{}
	username := req.URL.Query().Get("user")
	sig := req.URL.Query().Get("sig")

	templateIDs := getTemplateIDs(username, sig, result)
	if templateIDs == "" {
		nodes["result"] = result
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		setResponse(rw, nodes)
	} else {
		items := []interface{}{}
		sqlcmd := setSQLQuery(templateIDs, req, result)
		alerts = queryAlerts(sqlcmd, req, result)
		items = addPlatformToAlerts(alerts, result, nodes, rw)
		countOfSeverity := getAlertSeverityCounts(items)
		countOfProcess := getAlertProcessCounts(items)
		countOfMetricType := getAlertMetricTypeCounts(items)
		result["items"] = items
		nodes["result"] = result
		nodes["count"] = countOfSeverity
		nodes["countOfSeverity"] = countOfSeverity
		nodes["countOfProcess"] = countOfProcess
		nodes["countOfMetricType"] = countOfMetricType
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		setResponse(rw, nodes)
	}
}

func configAlertRoutes() {
	http.HandleFunc("/api/alerts", getAlerts)
}
