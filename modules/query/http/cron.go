package http

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	"github.com/jasonlvhit/gocron"
	log "github.com/Sirupsen/logrus"
)

type Contacts struct {
	Id      int
	Name    string
	Phone   string
	Email   string
	Updated string
}

type Hosts struct {
	Id       int
	Hostname string
	Exist    int
	Activate int
	Platform string
	Idc      string
	Ip       string
	Isp      string
	Province string
	City     string
	Status   string
	Updated  string
}

type Platforms struct {
	Id       int
	Platform string
	Contacts string
	Count    int
	Updated  string
}

func SyncHostsAndContactsTable() {
	if g.Config().Hosts.Enabled || g.Config().Contacts.Enabled {
		if g.Config().Hosts.Enabled {
			syncHostsTable()
			intervalToSyncHostsTable := uint64(10)
			gocron.Every(intervalToSyncHostsTable).Seconds().Do(syncHostsTable)
		}
		if g.Config().Contacts.Enabled {
			syncContactsTable()
			intervalToSyncContactsTable := uint64(g.Config().Contacts.Interval)
			gocron.Every(intervalToSyncContactsTable).Seconds().Do(syncContactsTable)
		}
		<- gocron.Start()
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

func updateHostsTable(hostnames []string, hostsMap map[string]map[string]string) {
	log.Debugf("func updateHostsTable()")
	var hosts []Hosts
	o := orm.NewOrm()
	o.Using("boss")
	_, err := o.QueryTable("hosts").All(&hosts)
	if err != nil {
		log.Errorf(err.Error())
	} else {
		format := "2006-01-02 15:04:05"
		for _, host := range hosts {
			updatedTime, _ := time.Parse(format, host.Updated)
			currentTime, _ := time.Parse(format, getNow())
			diff := currentTime.Unix() - updatedTime.Unix()
			if diff > 600 {
				host.Exist = 0
				_, err := o.Update(&host)
				if err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}

	idcMap := getIDCMap()
	var host Hosts
	for _, hostname := range hostnames {
		item := hostsMap[hostname]
		idcID := item["idcID"]
		if _, ok := idcMap[idcID]; ok {
			idc := idcMap[idcID]
			activate, _ := strconv.Atoi(item["activate"])
			host.Hostname = item["hostname"]
			host.Exist = 1
			host.Activate = activate
			host.Platform = item["platform"]
			host.Idc = idc.(Idc).Name
			host.Province = idc.(Idc).Province
			host.City = idc.(Idc).City
			host.Ip = item["ip"]
			host.Isp = strings.Split(item["hostname"], "-")[0]
			host.Updated = getNow()
			hosts = append(hosts, host)
		}
	}
	for _, item := range hosts {
		err := o.QueryTable("hosts").Filter("hostname", item.Hostname).One(&host)
		if err == orm.ErrNoRows {
			sql := "INSERT INTO boss.hosts("
			sql += "hostname, exist, activate, platform, idc, ip, "
			sql += "isp, province, city, updated) "
			sql += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			_, err := o.Raw(sql, item.Hostname, item.Exist, item.Activate, item.Platform, item.Idc, item.Ip, item.Isp, item.Province, item.City, item.Updated).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		} else if err != nil {
			log.Errorf(err.Error())
		} else {
			item.Id = host.Id
			_, err := o.Update(&item)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}

func updatePlatformsTable(platformNames []string, platformsMap map[string]map[string]interface{}) {
	log.Debugf("func updatePlatformsTable()")
	o := orm.NewOrm()
	o.Using("boss")
	var platform Platforms
	for _, platformName := range platformNames {
		group := platformsMap[platformName]
		err := o.QueryTable("platforms").Filter("platform", group["platformName"]).One(&platform)
		if err == orm.ErrNoRows {
			sql := "INSERT INTO boss.platforms(platform, contacts, count, updated) VALUES(?, ?, ?, ?)"
			_, err := o.Raw(sql, group["platformName"], group["contacts"], group["count"], getNow()).Exec()
			if err != nil {
				log.Errorf(err.Error())
			}
		} else if err != nil {
			log.Errorf(err.Error())
		} else {
			platform.Platform = group["platformName"].(string)
			platform.Count = group["count"].(int)
			platform.Updated = getNow()
			_, err := o.Update(&platform)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}

func updateContactsTable(contactNames []string, contactsMap map[string]map[string]interface{}) {
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
			contact.Email = user["email"].(string)
			contact.Phone = user["phone"].(string)
			contact.Updated = getNow()
			_, err := o.Update(&contact)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}

func addContactsToPlatformsTable(contacts map[string]interface{}) {
	log.Debugf("func addContactsToPlatformsTable()")
	o := orm.NewOrm()
	o.Using("boss")
	var platforms []Platforms
	_, err := o.QueryTable("platforms").All(&platforms)
	if err != nil {
		log.Errorf(err.Error())
	} else {
		for _, platform := range platforms {
			contactsOfPlatform := []string{}
			platformName := platform.Platform
			if users, ok := contacts[platformName]; ok {
				for _, user := range users.([]interface{}) {
					contactName := user.(map[string]interface{})["name"].(string)
					contactsOfPlatform = appendUniqueString(contactsOfPlatform, contactName)
				}
			}
			if len(contactsOfPlatform) > 0 {
				platform.Contacts = strings.Join(contactsOfPlatform, ",")
				platform.Updated = getNow()
				_, err := o.Update(&platform)
				if err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}
}

func syncHostsTable() {
	log.Debugf("func syncHostsTable()")
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM boss.hosts WHERE exist = 1 ORDER BY updated DESC LIMIT 1"
	num, err := o.Raw(sql).Values(&rows)
	diff := int64(0)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		format := "2006-01-02 15:04:05"
		updatedTime, _ := time.Parse(format, rows[0]["updated"].(string))
		currentTime, _ := time.Parse(format, getNow())
		diff = currentTime.Unix() - updatedTime.Unix()
	}
	log.Debugf("diff =", diff)
	if int(diff) < g.Config().Hosts.Interval {
		return
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
	hostsMap := map[string]map[string]string{}
	hostnamesMap := map[string]int{}
	idcIDs := []string{}
	hostname := ""
	for _, platform := range nodes["result"].([]interface{}) {
		countOfHosts := 0
		platformName := platform.(map[string]interface{})["platform"].(string)
		platformNames = appendUniqueString(platformNames, platformName)
		for _, device := range platform.(map[string]interface{})["ip_list"].([]interface{}) {
			hostname = device.(map[string]interface{})["hostname"].(string)
			ip := device.(map[string]interface{})["ip"].(string)
			if len(ip) > 0 && ip == getIPFromHostname(hostname, result) {
				if _, ok := hostnamesMap[hostname]; !ok {
					hostnames = append(hostnames, hostname)
					idcID := device.(map[string]interface{})["pop_id"].(string)
					host := map[string]string{
						"hostname": hostname,
						"activate": device.(map[string]interface{})["ip_status"].(string),
						"platform": platformName,
						"idcID":    idcID,
						"ip":       ip,
					}
					hostsMap[hostname] = host
					idcIDs = appendUniqueString(idcIDs, idcID)
					hostnamesMap[hostname] = 1
					countOfHosts++
				}
			}
		}
		platformsMap[platformName] = map[string]interface{}{
			"platformName": platformName,
			"count": countOfHosts,
			"contacts": "",
		}
	}
	sort.Strings(hostnames)
	sort.Strings(platformNames)
	log.Debugf("platformNames =", platformNames)
	updateHostsTable(hostnames, hostsMap)
	updatePlatformsTable(platformNames, platformsMap)
}

func syncContactsTable() {
	log.Debugf("func syncContactsTable()")
	o := orm.NewOrm()
	o.Using("boss")
	var rows []orm.Params
	sql := "SELECT updated FROM boss.contacts ORDER BY updated DESC LIMIT 1"
	num, err := o.Raw(sql).Values(&rows)
	diff := int64(0)
	if err != nil {
		log.Errorf(err.Error())
		return
	} else if num > 0 {
		format := "2006-01-02 15:04:05"
		updatedTime, _ := time.Parse(format, rows[0]["updated"].(string))
		currentTime, _ := time.Parse(format, getNow())
		diff = currentTime.Unix() - updatedTime.Unix()
	}
	log.Debugf("diff =", diff)
	if int(diff) < g.Config().Contacts.Interval {
		return
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
	contactsMap := map[string]map[string]interface{}{}
	contacts := nodes["result"].(map[string]interface{})["items"].(map[string]interface{})
	for _, platformName := range platformNames {
		if items, ok := contacts[platformName]; ok {
			for _, user := range items.([]interface{}) {
				contactName := user.(map[string]interface{})["name"].(string)
				if _, ok := contactsMap[contactName]; !ok {
					contactsMap[contactName] = user.(map[string]interface{})
					contactNames = append(contactNames, contactName)
				}
			}
		}
	}
	sort.Strings(contactNames)
	updateContactsTable(contactNames, contactsMap)
	addContactsToPlatformsTable(contacts)
}
