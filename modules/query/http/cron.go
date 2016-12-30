package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/query/g"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	"github.com/jasonlvhit/gocron"
)

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
	Hostname string
	Platform string
	Updated  string
}

type Platforms struct {
	Id        int
	Platform  string
	Contacts  string
	Principal string
	Deputy    string
	Upgrader  string
	Count     int
	Updated   string
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
		<-gocron.Start()
	}
}

func getIDCMap() map[string]interface{} {
	idcMap := map[string]interface{}{}
	o := orm.NewOrm()
	var idcs []Idc
	sqlcommand := "SELECT pop_id, name, province, city FROM grafana.idc ORDER BY pop_id ASC"
	_, err := o.Raw(sqlcommand).QueryRows(&idcs)
	if err != nil {
		log.Errorf(err.Error())
	}
	for _, idc := range idcs {
		idcMap[strconv.Itoa(idc.Pop_id)] = idc
	}
	return idcMap
}

func queryIDCsHostsCount(IDCName string) int64 {
	count := int64(0)
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT COUNT(*) FROM boss.hosts WHERE idc = ?"
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
	sql := "SELECT updated FROM boss.idcs ORDER BY updated DESC LIMIT 1"
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

func syncHostsTable() {
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM boss.ips WHERE exist = 1 ORDER BY updated DESC LIMIT 1"
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
	platformsMap := map[string]map[string]interface{}{}
	hostnames := []string{}
	IPs := []string{}
	IPKeys := []string{}
	IPsMap := map[string]map[string]string{}
	hostsMap := map[string]map[string]string{}
	idcIDs := []string{}
	hostname := ""
	for _, platform := range nodes["result"].([]interface{}) {
		platformName := platform.(map[string]interface{})["platform"].(string)
		platformNames = appendUniqueString(platformNames, platformName)
		for _, device := range platform.(map[string]interface{})["ip_list"].([]interface{}) {
			hostname = device.(map[string]interface{})["hostname"].(string)
			IP := device.(map[string]interface{})["ip"].(string)
			status := device.(map[string]interface{})["ip_status"].(string)
			item := map[string]string{
				"IP":       IP,
				"status":   status,
				"hostname": hostname,
				"platform": platformName,
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
		platformsMap[platformName] = map[string]interface{}{
			"platformName": platformName,
		}
	}
	sort.Strings(IPs)
	sort.Strings(IPKeys)
	sort.Strings(hostnames)
	sort.Strings(platformNames)
	log.Debugf("platformNames =", platformNames)
	updateIPsTable(IPKeys, IPsMap)
	updateHostsTable(hostnames, hostsMap)
	updatePlatformsTable(platformNames, platformsMap)
	muteFalconHostTable(hostnames, hostsMap)
}

func syncContactsTable() {
	log.Debugf("func syncContactsTable()")
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM boss.contacts ORDER BY updated DESC LIMIT 1"
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
					if (role == "principal") {
						platform.Principal = user["name"]
					} else if (role == "deputy") {
						platform.Deputy = user["name"]
					} else if (role == "upgrader") {
						platform.Upgrader = user["name"]
					}
				}
				if (len(platform.Principal) > 0) {
					contacts = append(contacts, platform.Principal)
				}
				if (len(platform.Deputy) > 0) {
					contacts = append(contacts, platform.Deputy)
				}
				if (len(platform.Upgrader) > 0) {
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
			sql := "INSERT INTO boss.contacts(name, phone, email, updated) VALUES(?, ?, ?, ?)"
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
			sql := "INSERT INTO boss.idcs(popid, idc, bandwidth, count, area, province, city, updated) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
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
	sql := "SELECT updated FROM boss.ips WHERE exist = 1 ORDER BY updated DESC LIMIT 1"
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
			sql += "ip, exist, status, hostname, platform, updated) "
			sql += "VALUES(?, ?, ?, ?, ?, ?)"
			_, err := o.Raw(sql, item["IP"], 1, status, item["hostname"], item["platform"], now).Exec()
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
			sql += " SET ip = ?, exist = ?, status = ?,"
			sql += " hostname = ?, platform = ?, updated = ?"
			sql += " WHERE id = ?"
			_, err := o.Raw(sql, item["IP"], 1, status, item["hostname"], item["platform"], now, ID).Exec()
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
	now := getNow()
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
			host["IDC"] = idc.(Idc).Name
			host["province"] = idc.(Idc).Province
			host["city"] = idc.(Idc).City
		}
		hosts = append(hosts, host)
	}

	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM boss.hosts WHERE exist = 1 ORDER BY updated DESC LIMIT 1"
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

	for _, host := range hosts {
		sql = "SELECT id FROM boss.hosts WHERE hostname = ? LIMIT 1"
		num, err := o.Raw(sql, host["hostname"]).Values(&rows)
		if num == 0 {
			sql := "INSERT INTO boss.hosts("
			sql += "hostname, exist, activate, platform, platforms, idc, ip, "
			sql += "isp, province, city, updated) "
			sql += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			_, err := o.Raw(sql, host["hostname"], 1, host["activate"], host["platform"], host["platforms"], host["IDC"], host["IP"], host["ISP"], host["province"], host["city"], now).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		} else if err != nil {
			log.Errorf(err.Error())
		} else if num > 0 {
			row := rows[0]
			ID := row["id"]
			sql = "UPDATE boss.hosts"
			sql += " SET exist = ?, activate = ?, platform = ?,"
			sql += " platforms = ?, idc = ?, ip = ?, isp = ?,"
			sql += " province = ?, city = ?, updated = ?"
			sql += " WHERE id = ?"
			_, err := o.Raw(sql, 1, host["activate"], host["platform"], host["platforms"], host["IDC"], host["IP"], host["ISP"], host["province"], host["city"], now, ID).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}

	sql = "SELECT id FROM boss.hosts WHERE exist = ?"
	sql += " AND updated <= DATE_SUB(CONVERT_TZ(NOW(),@@session.time_zone,'+08:00'),"
	sql += " INTERVAL 10 MINUTE) LIMIT 30"
	num, err = o.Raw(sql, 1).Values(&rows)
	if err != nil {
		log.Errorf(err.Error())
	} else if num > 0 {
		for _, row := range rows {
			ID := row["id"]
			sql = "UPDATE boss.hosts"
			sql += " SET exist = ?"
			sql += " WHERE id = ?"
			_, err := o.Raw(sql, 0, ID).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}

func muteFalconHostTable(hostnames []string, hostsMap map[string]map[string]string) {
	log.Debugf("func muteFalconHostTable()")
	o := orm.NewOrm()
	var rows []orm.Params
	now := getNow()
	for _, hostname := range hostnames {
		host := hostsMap[hostname]
		sql := "SELECT id FROM falcon_portal.host WHERE hostname = ? LIMIT 1"
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

func updatePlatformsTable(platformNames []string, platformsMap map[string]map[string]interface{}) {
	log.Debugf("func updatePlatformsTable()")
	now := getNow()
	o := orm.NewOrm()
	o.Using("boss")
	var platform Platforms
	for _, platformName := range platformNames {
		count, err := o.QueryTable("ips").Filter("platform", platformName).Filter("exist", 1).Filter("status", 1).Exclude("hostname__isnull", true).Count()
		if err != nil {
			count = 0
		}
		group := platformsMap[platformName]
		err = o.QueryTable("platforms").Filter("platform", group["platformName"]).One(&platform)
		if err == orm.ErrNoRows {
			sql := "INSERT INTO boss.platforms(platform, count, updated) VALUES(?, ?, ?)"
			_, err := o.Raw(sql, group["platformName"], count, now).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		} else if err != nil {
			log.Errorf(err.Error())
		} else {
			platform.Platform = group["platformName"].(string)
			platform.Count = int(count)
			platform.Updated = now
			_, err := o.Update(&platform)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}
