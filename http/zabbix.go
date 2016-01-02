package http

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/bitly/go-simplejson"
	"github.com/Cepave/query/g"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Endpoint struct {
	Id       int
	Endpoint string
	Ts       int64
	T_create string
	T_modify string
	Ipv4     string
	Port     string
}

type Grp struct {
	Id          int
	Grp_name    string
	Create_user string
	Create_at   string
	Come_from   int
}

type Tpl struct {
	Id          int
	Tpl_name    string
	Parent_id   int
	Action_id   int
	Create_user string
	Create_at   string
}

type Grp_tpl struct {
	Id        int
	Grp_id    int
	Tpl_id    int
	Bind_user string
}

type Grp_host struct {
	Id      int
	Grp_id  int
	Host_id int
}

/**
 * @function name:   func getNow() string
 * @description:     This function gets string of current time.
 * @related issues:  OWL-093
 * @param:           void
 * @return:          now sting
 * @author:          Don Hsieh
 * @since:           10/21/2015
 * @last modified:   10/21/2015
 * @called by:       func hostCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 *                   func hostgroupCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 *                   func templateCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 *                   func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 */
func getNow() string {
	t := time.Now()
	now := t.Format("2006-01-02 15:04:05")
	return now
}

/**
 * @function name:   func getHostId(params map[string]interface{}) string
 * @description:     This function gets host ID.
 * @related issues:  OWL-240
 * @param:           params map[string]interface{}
 * @return:          hostId string
 * @author:          Don Hsieh
 * @since:           12/16/2015
 * @last modified:   12/16/2015
 * @called by:       func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 */
func getHostId(params map[string]interface{}) string {
	hostId := ""
	if val, ok := params["hostid"]; ok {
		if val != nil {
			hostId = val.(string)
		}
	}
	return hostId
}

/**
 * @function name:   func getHostName(params map[string]interface{}) string
 * @description:     This function gets host name.
 * @related issues:  OWL-240
 * @param:           params map[string]interface{}
 * @return:          hostName string
 * @author:          Don Hsieh
 * @since:           12/16/2015
 * @last modified:   12/16/2015
 * @called by:       func hostCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 *                   func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 */
func getHostName(params map[string]interface{}) string {
	hostName := ""
	if val, ok := params["host"]; ok {
		if val != nil {
			hostName = val.(string)
		} else if val, ok = params["name"]; ok {
			if val != nil {
				hostName = val.(string)
			}
		}
	}
	return hostName
}

/**
 * @function name:   func checkHostExist(hostId string, hostName string)
 * @description:     This function checks if a host existed.
 * @related issues:  OWL-240
 * @param:           hostId string
 * @param:           hostName string
 * @return:          endpoint Endpoint
 * @author:          Don Hsieh
 * @since:           12/16/2015
 * @last modified:   12/16/2015
 * @called by:       func hostCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 *                   func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 */
func checkHostExist(hostId string, hostName string) Endpoint {
	var endpoint Endpoint
	o := orm.NewOrm()
	if hostId != "" {
		hostIdint, err := strconv.Atoi(hostId)
		if err != nil {
			log.Println("Error =", err.Error())
		} else {
			endpoint := Endpoint{Id: hostIdint}
			err := o.Read(&endpoint)
			if err != nil {
				log.Println("Error =", err.Error())
			}
		}
	} else {
		err := o.QueryTable("endpoint").Filter("endpoint", hostName).One(&endpoint)
		if err == orm.ErrMultiRows {
			// Have multiple records
			log.Printf("Returned multiple rows")
		} else if err == orm.ErrNoRows {
			// No result
			log.Printf("Not found")
		}
	}
	return endpoint
}

/**
 * @function name:   bindGroup(hostId int64, params map[string]interface{}, args map[string]string, result map[string]interface{})
 * @description:     This function binds a host to a host group.
 * @related issues:  OWL-240
 * @param:           hostId int64
 * @param:           params map[string]interface{}
 * @param:           args map[string]string
 * @param:           result map[string]interface{}
 * @return:          void
 * @author:          Don Hsieh
 * @since:           12/15/2015
 * @last modified:   12/15/2015
 * @called by:       func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 *                   func addHost(hostName string, params map[string]interface{}, args map[string]string, result map[string]interface{})
 */
func bindGroup(hostId int64, params map[string]interface{}, args map[string]string, result map[string]interface{}) {
	if _, ok := params["groups"]; ok {
		o := orm.NewOrm()
		database := "falcon_portal"
		o.Using(database)

		sqlcmd := "DELETE FROM falcon_portal.grp_host WHERE host_id=?"
		res, err := o.Raw(sqlcmd, hostId).Exec()
		if err != nil {
			log.Println("Error =", err.Error())
			result["error"] = [1]string{string(err.Error())}
		} else {
			num, _ := res.RowsAffected()
			if num > 0 {
				log.Println("mysql row affected nums =", num)
			}
		}

		groups := params["groups"].([]interface{})
		groupId := ""
		for _, group := range groups {
			groupId = group.(map[string]interface{})["groupid"].(string)
			args["groupId"] = groupId
			grp_id, err := strconv.Atoi(groupId)
			sqlcmd := "SELECT COUNT(*) FROM falcon_portal.grp_host WHERE host_id=? AND grp_id=?"
			res, err := o.Raw(sqlcmd, hostId, grp_id).Exec()
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				num, _ := res.RowsAffected()
				log.Println("num =", num)
				if num > 0 {
					log.Println("Record existed. count =", num)
				} else {	// Record not existed. Insert new one.
					grp_host := Grp_host{
						Grp_id: grp_id,
						Host_id: int(hostId),
					}
					log.Println("grp_host =", grp_host)

					_, err = o.Insert(&grp_host)
					if err != nil {
						log.Println("Error =", err.Error())
						result["error"] = [1]string{string(err.Error())}
					}
				}
			}
		}
	}
}

/**
 * @function name:   bindTemplate(params map[string]interface{}, args map[string]string, result map[string]interface{})
 * @description:     This function binds a host to a template.
 * @related issues:  OWL-240
 * @param:           params map[string]interface{}
 * @param:           args map[string]string
 * @param:           result map[string]interface{}
 * @return:          void
 * @author:          Don Hsieh
 * @since:           12/15/2015
 * @last modified:   12/15/2015
 * @called by:       func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 *                   func addHost(hostName string, params map[string]interface{}, args map[string]string, result map[string]interface{})
 */
func bindTemplate(params map[string]interface{}, args map[string]string, result map[string]interface{}) {
	if _, ok := params["templates"]; ok {
		o := orm.NewOrm()
		database := "falcon_portal"
		o.Using(database)
		groupId := args["groupId"]
		grp_id, _ := strconv.Atoi(groupId)
		templates := params["templates"].([]interface{})
		for _, template := range templates {
			templateId := template.(map[string]interface{})["templateid"].(string)
			tpl_id, err := strconv.Atoi(templateId)
			args["templateId"] = templateId

			sqlcmd := "SELECT COUNT(*) FROM falcon_portal.grp_tpl WHERE grp_id=? AND tpl_id=?"
			res, err := o.Raw(sqlcmd, grp_id, tpl_id).Exec()
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				num, _ := res.RowsAffected()
				log.Println("num =", num)
				if num > 0 {
					log.Println("Record existed. count =", num)
				} else {	// Record not existed. Insert new one.
					grp_tpl := Grp_tpl{
						Grp_id: grp_id,
						Tpl_id: tpl_id,
						Bind_user: "zabbix",
					}
					log.Println("grp_tpl =", grp_tpl)

					_, err = o.Insert(&grp_tpl)
					if err != nil {
						log.Println("Error =", err.Error())
						result["error"] = [1]string{string(err.Error())}
					}
				}
			}
		}
	}
}

/**
 * @function name:   func addHost(hostName string, params map[string]interface{}, args map[string]string, result map[string]interface{})
 * @description:     This function inserts a host to "endpoint" table and binds the host to its group and template.
 * @related issues:  OWL-240
 * @param:           hostName string
 * @param:           params map[string]interface{}
 * @param:           args map[string]string
 * @param:           result map[string]interface{}
 * @return:          void
 * @author:          Don Hsieh
 * @since:           12/21/2015
 * @last modified:   12/28/2015
 * @called by:       func hostCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 *                   func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 */
func addHost(hostName string, params map[string]interface{}, args map[string]string, result map[string]interface{}) {
	ip := ""
	port := ""
	if _, ok := params["interfaces"]; ok {
		interfaces := params["interfaces"].([]interface{})
		for i, arg := range interfaces {
			if i == 0 {
				ip = arg.(map[string]interface{})["ip"].(string)
				port = arg.(map[string]interface{})["port"].(string)
				args["ip"] = ip
				args["port"] = port
			}
		}
	}

	t := time.Now()
	timestamp := t.Unix()
	log.Println(timestamp)
	now := getNow()

	endpoint := Endpoint{
		Endpoint: hostName,
		Ts: timestamp,
		T_create: now,
		T_modify: now,
		Ipv4: ip,
	}
	if len(port) > 0 {
		endpoint.Port = port
	}
	log.Println("endpoint =", endpoint)

	o := orm.NewOrm()
	hostId, err := o.Insert(&endpoint)
	if err != nil {
		log.Println("Error =", err.Error())
		result["error"] = [1]string{string(err.Error())}
	} else {
		bindGroup(hostId, params, args, result)
		hostid := strconv.Itoa(int(hostId))
		hostids := [1]string{string(hostid)}
		result["hostids"] = hostids
		bindTemplate(params, args, result)
	}
}

/**
 * @function name:   func hostCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function gets host data for database insertion.
 * @related issues:  OWL-240, OWL-093, OWL-086, OWL-085
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/11/2015
 * @last modified:   12/28/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostCreate(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostCreate()")
	params := nodes["params"].(map[string]interface{})
	var result = make(map[string]interface{})

	hostName := getHostName(params)
	hostId := ""
	endpoint := checkHostExist(hostId, hostName)
	log.Println("returned endpoint =", endpoint)
	log.Println("endpoint.Id =", endpoint.Id)
	if endpoint.Id > 0 {
		result["error"] = [1]string{"host name existed."}
	} else {
		log.Println("host not existed")
		if len(hostName) > 0 {
			args := map[string]string {
				"host": hostName,
			}
			addHost(hostName, params, args, result)
			if _, ok := params["inventory"]; ok {
				inventory := params["inventory"].(map[string]interface{})
				macAddr := inventory["macaddress_a"].(string) + inventory["macaddress_b"].(string)
				args["macAddr"] = macAddr
			}
			log.Println("args =", args)
		} else {
			result["error"] = [1]string{"host name can not be null."}
		}
	}
	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func unbindGroup(hostId string, result map[string]interface{})
 * @description:     This function unbinds a host to a host group.
 * @related issues:  OWL-241
 * @param:           hostId string
 * @param:           result map[string]interface{}
 * @return:          void
 * @author:          Don Hsieh
 * @since:           01/01/2016
 * @last modified:   01/01/2016
 * @called by:       func removeHost(hostIds []string, result map[string]interface{})
 */
func unbindGroup(hostId string, result map[string]interface{}) {
	o := orm.NewOrm()
	o.Using("falcon_portal")
	sql := "DELETE FROM grp_host WHERE host_id = ?"
	res, err := o.Raw(sql, hostId).Exec()
	if err != nil {
		log.Println("Error =", err.Error())
		result["error"] = [1]string{string(err.Error())}
	}
	num, _ := res.RowsAffected()
	log.Println("mysql row affected nums =", num)
}

/**
 * @function name:   func hostDelete(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function deletes host from "endpoint" table.
 * @related issues:  OWL-093, OWL-086, OWL-085
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/11/2015
 * @last modified:   10/21/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostDelete(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].([]interface {})
	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})

	o := orm.NewOrm()
	hostids := []string{}
	for _, hostId := range params {
		if id, err := strconv.Atoi(hostId.(string)); err == nil {
			num, err := o.Delete(&Endpoint{Id: id})
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				if num > 0 {
					hostids = append(hostids, hostId.(string))
					log.Println("RowsDeleted =", num)
				}
			}
		}
	}

	database := "falcon_portal"
	o.Using(database)
	for _, hostId := range params {
		sql := "DELETE FROM grp_host WHERE host_id = ?"
		res, err := o.Raw(sql, hostId).Exec()
		if err != nil {
			log.Println("Error =", err.Error())
		}
		num, _ := res.RowsAffected()
		log.Println("mysql row affected nums =", num)
	}
	result["hostids"] = hostids
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func hostGet(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function gets existed host data.
 * @related issues:  OWL-254
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           12/29/2015
 * @last modified:   12/30/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostGet(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostGet()")
	params := nodes["params"].(map[string]interface{})
	result := []interface{}{}
	hostNames := []string{}
	queryAll := false
	if val, ok := params["filter"]; ok {
		filter := val.(map[string]interface{})
		if val, ok = filter["host"]; ok {
			for _, hostName := range val.([]interface{}) {
				if hostName.(string) == "_all_" {
					queryAll = true
				} else {
					hostNames = append(hostNames, hostName.(string))
				}
			}
		}
	}
	o := orm.NewOrm()
	if queryAll {
		var endpoints []*Endpoint
		num, err := o.QueryTable("endpoint").All(&endpoints)
		if err != nil {
			log.Println("Error =", err.Error())
		} else {
			log.Println("num =", num)
			for _, endpoint := range endpoints {
				log.Println("endpoint =", endpoint)
				item := map[string]string {}
				var grp_id int
				o.Raw("SELECT grp_id FROM falcon_portal.grp_host WHERE host_id=?", endpoint.Id).QueryRow(&grp_id)
				item["hostid"] = strconv.Itoa(endpoint.Id)
				item["hostname"] = endpoint.Endpoint
				item["ip"] = endpoint.Ipv4
				item["groupid"] = strconv.Itoa(grp_id)
				result = append(result, item)
			}
		}
	} else {
		ip := ""
		hostId := ""
		groupId := ""
		var endpoint Endpoint
		for _, hostName := range hostNames {
			item := map[string]string {}
			ip = ""
			hostId = ""
			groupId = ""
			err := o.QueryTable("endpoint").Filter("endpoint", hostName).One(&endpoint)
			if err == orm.ErrMultiRows {
				log.Printf("Returned multiple rows")
			} else if err == orm.ErrNoRows {
				log.Printf("Not found")
			} else if endpoint.Id > 0 {
				ip = endpoint.Ipv4
				var grp_id int
				o.Raw("SELECT grp_id FROM falcon_portal.grp_host WHERE host_id=?", endpoint.Id).QueryRow(&grp_id)
				log.Println("grp_id =", grp_id)
				hostId = strconv.Itoa(endpoint.Id)
				groupId = strconv.Itoa(grp_id)
			}
			item["hostid"] = hostId
			item["hostname"] = hostName
			item["ip"] = ip
			item["groupid"] = groupId
			result = append(result, item)
		}
	}
	log.Println("result =", result)
	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function updates host data.
 * @related issues:  OWL-240, OWL-093, OWL-086
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/23/2015
 * @last modified:   12/28/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostUpdate()")
	params := nodes["params"].(map[string]interface{})
	var result = make(map[string]interface{})
	hostId := getHostId(params)
	hostName := getHostName(params)
	args := map[string]string {
		"host": hostName,
	}

	endpoint := checkHostExist(hostId, hostName)
	log.Println("returned endpoint =", endpoint)
	log.Println("endpoint.Id =", endpoint.Id)
	if endpoint.Id > 0 {
		log.Println("host existed")
		hostId := endpoint.Id
		now := getNow()
		log.Println("now =", now)
		endpoint.T_modify = now

		o := orm.NewOrm()
		num, err := o.Update(&endpoint)
		if err != nil {
			log.Println("Error =", err.Error())
			result["error"] = [1]string{string(err.Error())}
		} else {
			log.Println("update hostId =", hostId)
			log.Println("mysql row affected nums =", num)
			bindGroup(int64(endpoint.Id), params, args, result)
			hostid := strconv.Itoa(endpoint.Id)
			hostids := [1]string{string(hostid)}
			result["hostids"] = hostids
			bindTemplate(params, args, result)
		}
	} else {
		log.Println("host not existed")
		addHost(hostName, params, args, result)
	}
	log.Println("args =", args)

	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func hostgroupCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function gets hostgroup data for database insertion.
 * @related issues:  OWL-093, OWL-086
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/21/2015
 * @last modified:   10/21/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostgroupCreate(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostgroupCreate()")
	params := nodes["params"].(map[string]interface{})
	hostgroupName := params["name"].(string)
	user := "zabbix"
	now := getNow()

	database := "falcon_portal"
	o := orm.NewOrm()
	o.Using(database)

	grp := Grp{
		Grp_name: hostgroupName,
		Create_user: user,
		Create_at: now,
	}
	log.Println("grp =", grp)

	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})

	id, err := o.Insert(&grp)
	if err != nil {
		log.Println("Error =", err.Error())
		result["error"] = [1]string{string(err.Error())}
	} else {
		groupid := strconv.Itoa(int(id))
		groupids := [1]string{string(groupid)}
		result["groupids"] = groupids
	}
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func hostgroupDelete(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function gets hostgroup data for database insertion.
 * @related issues:  OWL-093, OWL-086
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/21/2015
 * @last modified:   10/21/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostgroupDelete(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostgroupDelete()")
	params := nodes["params"].([]interface {})

	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})

	o := orm.NewOrm()
	database := "falcon_portal"
	o.Using(database)

	args := []interface{}{}
	args = append(args, "DELETE FROM falcon_portal.grp WHERE id=?")
	args = append(args, "DELETE FROM falcon_portal.grp_host WHERE grp_id=?")
	args = append(args, "DELETE FROM falcon_portal.grp_tpl WHERE grp_id=?")
	args = append(args, "DELETE FROM falcon_portal.plugin_dir WHERE grp_id=?")
	log.Println("args =", args)

	groupids := []string{}
	for _, sqlcmd := range args {
		for _, hostgroupId := range params {
			res, err := o.Raw(sqlcmd.(string), hostgroupId).Exec()
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				num, _ := res.RowsAffected()
				if num > 0 && sqlcmd == "DELETE FROM falcon_portal.grp WHERE id=?" {
					groupids = append(groupids, hostgroupId.(string))
					log.Println("delete hostgroup id =", hostgroupId)
					log.Println("mysql row affected nums =", num)
				}
			}
		}
	}
	result["groupids"] = groupids
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func hostgroupGet(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function gets existed hostgroup data.
 * @related issues:  OWL-254
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           12/29/2015
 * @last modified:   12/29/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostgroupGet(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostgroupGet()")
	params := nodes["params"].(map[string]interface{})
	result := []interface{}{}
	groupNames := []string{}
	queryAll := false
	if val, ok := params["filter"]; ok {
		filter := val.(map[string]interface{})
		if val, ok = filter["name"]; ok {
			for _, groupName := range val.([]interface{}) {
				if groupName.(string) == "_all_" {
					queryAll = true
				} else {
					groupNames = append(groupNames, groupName.(string))
				}
			}
		}
	}
	groupId := ""
	o := orm.NewOrm()
	o.Using("falcon_portal")
	if queryAll {
		var grps []*Grp
		_, err := o.QueryTable("grp").All(&grps)
		if err != nil {
			log.Println("Error =", err.Error())
		} else {
			for _, grp := range grps {
				item := map[string]string {}
				item["groupid"] = strconv.Itoa(grp.Id)
				item["groupname"] = grp.Grp_name
				result = append(result, item)
			}
		}
	} else {
		var grp Grp
		for _, groupName := range groupNames {
			item := map[string]string {}
			groupId = ""
			err := o.QueryTable("grp").Filter("grp_name", groupName).One(&grp)
			if err == orm.ErrMultiRows {
				log.Printf("Returned multiple rows")
			} else if err == orm.ErrNoRows {
				log.Printf("Not found")
			} else if grp.Id > 0 {
				groupId = strconv.Itoa(grp.Id)
			}
			item["groupid"] = groupId
			item["groupname"] = groupName
			result = append(result, item)
		}
	}
	log.Println("result =", result)
	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func hostgroupUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function gets hostgroup data for database insertion.
 * @related issues:  OWL-093, OWL-086
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/21/2015
 * @last modified:   10/21/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostgroupUpdate(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostgroupUpdate()")
	params := nodes["params"].(map[string]interface{})
	var result = make(map[string]interface{})
	hostgroupId, err := strconv.Atoi(params["groupid"].(string))
	if err != nil {
		log.Println("Error =", err.Error())
		result["error"] = [1]string{string(err.Error())}
	}
	o := orm.NewOrm()
	database := "falcon_portal"
	o.Using(database)

	if _, ok := params["name"]; ok {
		hostgroupName := params["name"].(string)
		log.Println("hostgroupName =", hostgroupName)

		if hostgroupName != "" {
			grp := Grp{Id: hostgroupId}
			log.Println("grp =", grp)
			err := o.Read(&grp)
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				log.Println("grp =", grp)
				grp.Grp_name = hostgroupName
				log.Println("grp =", grp)
				num, err := o.Update(&grp)
				if err != nil {
					log.Println("Error =", err.Error())
					result["error"] = [1]string{string(err.Error())}
				} else {
					if num > 0 {
						groupids := [1]string{strconv.Itoa(hostgroupId)}
						result["groupids"] = groupids
						log.Println("update groupid =", hostgroupId)
						log.Println("mysql row affected nums =", num)
					}
				}
			}
		}
	}
	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func templateCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function gets template data for database insertion.
 * @related issues:  OWL-093, OWL-086
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/22/2015
 * @last modified:   10/21/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func templateCreate(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func templateCreate()")
	params := nodes["params"].(map[string]interface{})
	templateName := params["host"].(string)
	user := "zabbix"
	groups := params["groups"]
	groupid := groups.(map[string]interface{})["groupid"].(json.Number)
	hostgroupId := string(groupid)
	now := getNow()

	database := "falcon_portal"
	o := orm.NewOrm()
	o.Using(database)

	tpl := Tpl{
		Tpl_name: templateName,
		Create_user: user,
		Create_at: now,
	}
	log.Println("tpl =", tpl)

	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})

	id, err := o.Insert(&tpl)
	if err != nil {
		log.Println("Error =", err.Error())
		result["error"] = [1]string{string(err.Error())}
	} else {
		templateId := strconv.Itoa(int(id))
		templateids := [1]string{string(templateId)}
		result["templateids"] = templateids

		groupId, err := strconv.Atoi(hostgroupId)
		grp_tpl := Grp_tpl{
			Grp_id: groupId,
			Tpl_id: int(id),
			Bind_user: user,
		}
		log.Println("grp_tpl =", grp_tpl)

		_, err = o.Insert(&grp_tpl)
		if err != nil {
			log.Println("Error =", err.Error())
			result["error"] = [1]string{string(err.Error())}
		}
	}
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func templateDelete(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function deletes template data.
 * @related issues:  OWL-093, OWL-086
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/22/2015
 * @last modified:   10/21/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func templateDelete(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func templateDelete()")
	params := nodes["params"].([]interface {})
	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})
	o := orm.NewOrm()
	args := []interface{}{}
	args = append(args, "DELETE FROM falcon_portal.tpl WHERE id=?")
	args = append(args, "DELETE FROM falcon_portal.grp_tpl WHERE tpl_id=?")
	log.Println("args =", args)

	templateids := []string{}
	for _, sqlcmd := range args {
		log.Println(sqlcmd)
		for _, templateId := range params {
			log.Println("templateId =", templateId)
			res, err := o.Raw(sqlcmd.(string), templateId).Exec()
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				num, _ := res.RowsAffected()
				if num > 0 && sqlcmd == "DELETE FROM falcon_portal.tpl WHERE id=?" {
					templateids = append(templateids, templateId.(string))
					log.Println("delete template id =", templateId)
					log.Println("mysql row affected nums =", num)
				}
			}
		}
	}
	result["templateids"] = templateids
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func templateUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function gets hostgroup data for database insertion.
 * @related issues:  OWL-093, OWL-086
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/22/2015
 * @last modified:   10/23/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func templateUpdate(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].(map[string]interface{})
	var result = make(map[string]interface{})
	templateId, err := strconv.Atoi(params["templateid"].(string))
	if err != nil {
		log.Println("Error =", err.Error())
		result["error"] = [1]string{string(err.Error())}
	}
	o := orm.NewOrm()
	database := "falcon_portal"
	o.Using(database)

	if _, ok := params["name"]; ok {
		templateName := params["name"].(string)
		log.Println("templateName =", templateName)

		if templateName != "" {
			tpl := Tpl{Id: templateId}
			log.Println("tpl =", tpl)
			err := o.Read(&tpl)
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				log.Println("tpl =", tpl)
				tpl.Tpl_name = templateName
				log.Println("tpl =", tpl)
				num, err := o.Update(&tpl)
				if err != nil {
					log.Println("Error =", err.Error())
					result["error"] = [1]string{string(err.Error())}
				} else {
					if num > 0 {
						templateids := [1]string{strconv.Itoa(templateId)}
						result["templateids"] = templateids
						log.Println("update template id =", templateId)
						log.Println("mysql row affected nums =", num)
					}
				}
			}
		}
	}

	if _, ok := params["groups"]; ok {
		groups := params["groups"].([]interface{})
		log.Println("groups =", groups)

		count := 0
		for _, group := range groups {
			log.Println("group =", group)
			count += 1
		}
		log.Println("count =", count)

		if count > 0 {
			user := "zabbix"
			sqlcmd := "DELETE FROM falcon_portal.grp_tpl WHERE tpl_id=?"
			res, err := o.Raw(sqlcmd, templateId).Exec()
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				num, _ := res.RowsAffected()
				if num > 0 {
					log.Println("mysql row affected nums =", num)
				}
			}

			for _, group := range groups {
				log.Println("group =", group)
				groupId, err := strconv.Atoi(group.(map[string]interface{})["groupid"].(string))
				log.Println("groupId =", groupId)
				grp_tpl := Grp_tpl{Grp_id: groupId, Tpl_id: templateId, Bind_user: user}
				log.Println("grp_tpl =", grp_tpl)

				_, err = o.Insert(&grp_tpl)
				if err != nil {
					log.Println("Error =", err.Error())
					result["error"] = [1]string{string(err.Error())}
				} else {
					templateids := [1]string{strconv.Itoa(templateId)}
					result["templateids"] = templateids
					log.Println("update template id =", templateId)
				}
			}
		}
	}
	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func getFctoken() fctoken string
 * @description:     This function returns fctoken for API request.
 * @related issues:  OWL-159
 * @param:           void
 * @return:          fctoken string
 * @author:          Don Hsieh
 * @since:           11/24/2015
 * @last modified:   11/24/2015
 * @called by:       func apiAlert(rw http.ResponseWriter, req *http.Request)
 *                    in query/http/zabbix.go
 *                   func getMapValues(chartType string) map[string]interface{}
 *                    in query/http/grafana.go
 */
func getFctoken() string {
	hasher := md5.New()
	io.WriteString(hasher, g.Config().Api.Token)
	s := hex.EncodeToString(hasher.Sum(nil))

	t := time.Now()
	now := t.Format("20060102")
	s = now + s

	hasher = md5.New()
	io.WriteString(hasher, s)
	fctoken := hex.EncodeToString(hasher.Sum(nil))
	return fctoken
}

/**
 * @function name:   func apiAlert(rw http.ResponseWriter, req *http.Request)
 * @description:     This function handles alarm API request.
 * @related issues:  OWL-159, OWL-093
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/29/2015
 * @last modified:   11/24/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func apiAlert(rw http.ResponseWriter, req *http.Request) {
	fcname := g.Config().Api.Name
	fctoken := getFctoken()
	param := req.URL.Query()
	log.Println("param =", param)
	arr := param["endpoint"]
	hostname := arr[0]
	arr = param["time"]
	datetime := arr[0]

	arr = param["stra_id"]
	trigger_id, err := strconv.Atoi(arr[0])
	if err != nil {
		log.Println(err.Error())
	}
	arr = param["metric"]
	metric := arr[0]
	arr = param["step"]
	step := arr[0]
	arr = param["tpl_id"]
	tpl_id := arr[0]
	arr = param["status"]
	zabbix_status := arr[0]
	arr = param["priority"]
	zabbix_level := arr[0]
	summary := "[OWL] " + metric + "_" + step + "_" + zabbix_level

	args := map[string]interface{} {
		"summary": summary,
		"zabbix_status": zabbix_status,		// "PROBLEM",
		"zabbix_level": "Information",		// "Information" or "High"
		"trigger_id": trigger_id,
		"host_ip": "",
		"hostname": hostname,
		"event_id": tpl_id,
		"template_name": "Template Server Basic Monitor",
		"datetime": datetime,
		"fcname": fcname,
		"fctoken": fctoken,
	}

	log.Println("args =", args)
	bs, err := json.Marshal(args)
	if err != nil {
		log.Println("Error =", err.Error())
	}

	url := g.Config().Api.Event
	log.Println("url =", url)

	reqAlert, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bs)))
	if err != nil {
		log.Println("Error =", err.Error())
	}
	reqAlert.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqAlert)
	if err != nil {
		log.Println("Error =", err.Error())
	}
	defer resp.Body.Close()

	log.Println("response Status =", resp.Status)	// 200 OK   TypeOf(resp.Status): string
	log.Println("response Headers =", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body =", string(body))
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

/**
 * @function name:   func apiParser(rw http.ResponseWriter, req *http.Request)
 * @description:     This function parses the method of API request.
 * @related issues:  OWL-254, OWL-085
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/11/2015
 * @last modified:   12/29/2015
 * @called by:       http.HandleFunc("/api", apiParser)
 *                    in func main()
 */
func apiParser(rw http.ResponseWriter, req *http.Request) {
	log.Println("func apiParser(rw http.ResponseWriter, req *http.Request)")
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	log.Println(buf.Len())
	if buf.Len() == 0 {
		apiAlert(rw, req)
	} else {
		s := buf.String() // Does a complete copy of the bytes in the buffer.
		log.Println("s =", s)
		json, err := simplejson.NewJson(buf.Bytes())
		if err != nil {
			log.Println(err.Error())
		}

		var nodes = make(map[string]interface{})
		nodes, _ = json.Map()

		method := nodes["method"]
		log.Println(method)
		delete(nodes, "method")
		delete(nodes, "auth")

		if method == "host.create" {
			hostCreate(nodes, rw)
		} else if method == "host.delete" {
			hostDelete(nodes, rw)
		} else if method == "host.get" {
			hostGet(nodes, rw)
		} else if method == "host.update" {
			hostUpdate(nodes, rw)
		} else if method == "hostgroup.create" {
			hostgroupCreate(nodes, rw)
		} else if method == "hostgroup.delete" {
			hostgroupDelete(nodes, rw)
		} else if method == "hostgroup.get" {
			hostgroupGet(nodes, rw)
		} else if method == "hostgroup.update" {
			hostgroupUpdate(nodes, rw)
		} else if method == "template.create" {
			templateCreate(nodes, rw)
		} else if method == "template.delete" {
			templateDelete(nodes, rw)
		} else if method == "template.update" {
			templateUpdate(nodes, rw)
		}
	}
}

/**
 * @function name:   func configZabbixRoutes()
 * @description:     This function handles API requests.
 * @related issues:  OWL-093, OWL-085
 * @param:           void
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/09/2015
 * @last modified:   10/21/2015
 * @called by:       func Start()
 *                    in http/http.go
 */
func configZabbixRoutes() {
	http.HandleFunc("/api", apiParser)
}
