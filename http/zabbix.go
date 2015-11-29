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
	"reflect"
	"strconv"
	"time"
)

type Endpoint struct {
	Id       int
	Endpoint string
	Ts       int64
	T_create string
	T_modify string
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
 * @function name:   func hostCreate(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function gets host data for database insertion.
 * @related issues:  OWL-093, OWL-086, OWL-085
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/11/2015
 * @last modified:   10/22/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostCreate(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostCreate()")
	params := nodes["params"].(map[string]interface{})
	var result = make(map[string]interface{})

	if _, ok := params["host"]; ok {
		host := params["host"].(string)
		t := time.Now()
		timestamp := t.Unix()
		log.Println(timestamp)
		now := getNow()
		args := map[string]string {
			"host": host,
		}

		o := orm.NewOrm()
		endpoint := Endpoint{
			Endpoint: host,
			Ts: timestamp,
			T_create: now,
			T_modify: now,
		}

		log.Println("endpoint =", endpoint)
		hostId, err := o.Insert(&endpoint)
		if err != nil {
			log.Println("Error =", err.Error())
			result["error"] = [1]string{string(err.Error())}
		} else {
			hostid := strconv.Itoa(int(hostId))
			hostids := [1]string{string(hostid)}
			result["hostids"] = hostids
			if _, ok := params["groups"]; ok {
				database := "falcon_portal"
				o.Using(database)
				groups := params["groups"].([]interface{})
				groupId := ""
				for i, group := range groups {
					groupId = group.(map[string]interface{})["groupid"].(string)
					log.Println("hostId:", hostId)
					log.Println("groupId:", groupId)
					if i == 0 {
						args["groupId"] = groupId
					}
					grp_id, err := strconv.Atoi(groupId)

					sqlcmd := "SELECT COUNT(*) FROM falcon_portal.grp_host WHERE host_id=? AND grp_id=?"
					log.Println("sqlcmd:", sqlcmd)
					res, err := o.Raw(sqlcmd, hostId, grp_id).Exec()
					if err != nil {
						log.Println("Error =", err.Error())
						result["error"] = [1]string{string(err.Error())}
					} else {
						num, _ := res.RowsAffected()
						log.Println("num:", num)
						if num > 0 {
							log.Println("Record existed. count:", num)
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

		if _, ok := params["interfaces"]; ok {
			interfaces := params["interfaces"].([]interface{})
			for i, arg := range interfaces {
				if i == 0 {
					ip := arg.(map[string]interface{})["ip"].(string)
					port := arg.(map[string]interface{})["port"].(string)
					args["ip"] = ip
					args["port"] = port
				}
			}
		}

		if _, ok := params["templates"]; ok {
			templates := params["templates"].([]interface{})
			for i, template := range templates {
				if i == 0 {
					templateId := template.(map[string]interface{})["templateid"].(string)
					args["templateId"] = templateId
				}
			}
		}

		if _, ok := params["inventory"]; ok {
			inventory := params["inventory"].(map[string]interface{})
			macAddr := inventory["macaddress_a"].(string) + inventory["macaddress_b"].(string)
			args["macAddr"] = macAddr
		}
		log.Println("args =", args)
	} else {
		result["error"] = [1]string{"host name can not be null."}
	}

	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
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
					log.Println("RowsDeleted:", num)
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
		log.Println("mysql row affected nums:", num)
	}
	result["hostids"] = hostids
	resp["result"] = result
	RenderJson(rw, resp)
}

/**
 * @function name:   func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter)
 * @description:     This function updates host data.
 * @related issues:  OWL-093, OWL-086
 * @param:           nodes map[string]interface{}
 * @param:           rw http.ResponseWriter
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/23/2015
 * @last modified:   10/21/2015
 * @called by:       func apiParser(rw http.ResponseWriter, req *http.Request)
 */
func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostUpdate()")
	params := nodes["params"].(map[string]interface{})
	var result = make(map[string]interface{})
	hostId, err := strconv.Atoi(params["hostid"].(string))
	if err != nil {
		log.Println("Error =", err.Error())
		result["error"] = [1]string{string(err.Error())}
	}
	now := getNow()
	log.Println("now:", now)

	o := orm.NewOrm()

	if _, ok := params["host"]; ok {
		hostName := params["host"].(string)
		log.Println("hostName:", hostName)

		if hostName != "" {
			endpoint := Endpoint{Id: hostId}
			log.Println("endpoint:", endpoint)
			err := o.Read(&endpoint)
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				log.Println("endpoint:", endpoint)
				endpoint.Endpoint = hostName
				endpoint.T_modify = now
				num, err := o.Update(&endpoint)
				if err != nil {
					log.Println("Error =", err.Error())
					result["error"] = [1]string{string(err.Error())}
				} else {
					if num > 0 {
						hostids := [1]string{strconv.Itoa(hostId)}
						result["hostids"] = hostids
						log.Println("update hostId:", hostId)
						log.Println("mysql row affected nums:", num)
					}
				}
			}
		}
	}

	if _, ok := params["groups"]; ok {
		groups := params["groups"].([]interface{})
		log.Println("groups:", groups)

		count := 0
		for _, group := range groups {
			log.Println("group:", group)
			count += 1
		}
		log.Println("count:", count)

		if count > 0 {
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
					log.Println("mysql row affected nums:", num)
				}
			}

			for _, group := range groups {
				log.Println("group:", group)
				groupId, err := strconv.Atoi(group.(map[string]interface{})["groupid"].(string))
				log.Println("groupId:", groupId)
				grp_host := Grp_host{Grp_id: groupId, Host_id: hostId}
				log.Println("grp_host:", grp_host)

				database := "falcon_portal"
				o.Using(database)
				_, err = o.Insert(&grp_host)
				if err != nil {
					log.Println("Error =", err.Error())
					result["error"] = [1]string{string(err.Error())}
				} else {
					hostids := [1]string{strconv.Itoa(hostId)}
					result["hostids"] = hostids
					log.Println("update hostId:", hostId)
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
		log.Println("TypeOf(err):", reflect.TypeOf(err))					// *mysql.MySQLError
		log.Println("TypeOf(err.Error()):", reflect.TypeOf(err.Error()))	// string
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
					log.Println("delete hostgroup id:", hostgroupId)
					log.Println("mysql row affected nums:", num)
				}
			}
		}
	}
	result["groupids"] = groupids
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
		log.Println("hostgroupName:", hostgroupName)

		if hostgroupName != "" {
			grp := Grp{Id: hostgroupId}
			log.Println("grp:", grp)
			err := o.Read(&grp)
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				log.Println("grp:", grp)
				grp.Grp_name = hostgroupName
				log.Println("grp:", grp)
				num, err := o.Update(&grp)
				if err != nil {
					log.Println("Error =", err.Error())
					result["error"] = [1]string{string(err.Error())}
				} else {
					if num > 0 {
						groupids := [1]string{strconv.Itoa(hostgroupId)}
						result["groupids"] = groupids
						log.Println("update groupid:", hostgroupId)
						log.Println("mysql row affected nums:", num)
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
					log.Println("delete template id:", templateId)
					log.Println("mysql row affected nums:", num)
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
		log.Println("templateName:", templateName)

		if templateName != "" {
			tpl := Tpl{Id: templateId}
			log.Println("tpl:", tpl)
			err := o.Read(&tpl)
			if err != nil {
				log.Println("Error =", err.Error())
				result["error"] = [1]string{string(err.Error())}
			} else {
				log.Println("tpl:", tpl)
				tpl.Tpl_name = templateName
				log.Println("tpl:", tpl)
				num, err := o.Update(&tpl)
				if err != nil {
					log.Println("Error =", err.Error())
					result["error"] = [1]string{string(err.Error())}
				} else {
					if num > 0 {
						templateids := [1]string{strconv.Itoa(templateId)}
						result["templateids"] = templateids
						log.Println("update template id:", templateId)
						log.Println("mysql row affected nums:", num)
					}
				}
			}
		}
	}

	if _, ok := params["groups"]; ok {
		groups := params["groups"].([]interface{})
		log.Println("groups:", groups)

		count := 0
		for _, group := range groups {
			log.Println("group:", group)
			count += 1
		}
		log.Println("count:", count)

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
					log.Println("mysql row affected nums:", num)
				}
			}

			for _, group := range groups {
				log.Println("group:", group)
				groupId, err := strconv.Atoi(group.(map[string]interface{})["groupid"].(string))
				log.Println("groupId:", groupId)
				grp_tpl := Grp_tpl{Grp_id: groupId, Tpl_id: templateId, Bind_user: user}
				log.Println("grp_tpl:", grp_tpl)

				_, err = o.Insert(&grp_tpl)
				if err != nil {
					log.Println("Error =", err.Error())
					result["error"] = [1]string{string(err.Error())}
				} else {
					templateids := [1]string{strconv.Itoa(templateId)}
					result["templateids"] = templateids
					log.Println("update template id:", templateId)
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

	log.Println("response Status:", resp.Status)	// 200 OK   TypeOf(resp.Status): string
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

/**
 * @function name:   func apiParser(rw http.ResponseWriter, req *http.Request)
 * @description:     This function parses the method of API request.
 * @related issues:  OWL-085
 * @param:           rw http.ResponseWriter
 * @param:           req *http.Request
 * @return:          void
 * @author:          Don Hsieh
 * @since:           09/11/2015
 * @last modified:   10/23/2015
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
		} else if method == "host.update" {
			hostUpdate(nodes, rw)
		} else if method == "host.exists" {
			// hostExist(params)
		} else if method == "hostgroup.create" {
			hostgroupCreate(nodes, rw)
		} else if method == "hostgroup.delete" {
			hostgroupDelete(nodes, rw)
		} else if method == "hostgroup.update" {
			hostgroupUpdate(nodes, rw)
		} else if method == "hostgroup.exists" {
			// hostgroupExist(params)
		} else if method == "template.create" {
			templateCreate(nodes, rw)
		} else if method == "template.delete" {
			templateDelete(nodes, rw)
		} else if method == "template.update" {
			templateUpdate(nodes, rw)
		} else if method == "template.exists" {
			// templateExist(params)
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
