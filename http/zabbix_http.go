package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/toolkits/file"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type DatabaseConfig struct {
	Addr     string `json:"addr"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	Hostname string          `json:"hostname"`
	IP       string          `json:"ip"`
	Database *DatabaseConfig `json:"database"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func initDb() {
	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/graph?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}
}

func readDb(endpointId int) {
	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/graph?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	stmtOut, err := db.Prepare("SELECT endpoint FROM graph.endpoint WHERE id = ?")
	if err != nil {
		log.Println(err.Error())
	}
	defer stmtOut.Close()

	var endpoint string // we "scan" the result in here

	err = stmtOut.QueryRow(endpointId).Scan(&endpoint) // WHERE id = endpointId
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("The endpoint name of %d is: %s", endpointId, endpoint)
}

func writeDb(sqlcmd string, args []interface{}) {
	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/graph?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	stmtIns, err := db.Prepare(sqlcmd)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	if result, err := stmtIns.Exec(args); err == nil {
		if id, err := result.LastInsertId(); err == nil {
			log.Println("insert id :", id)
		} else {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func hostCreate(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].(map[string]interface{})
	host := params["host"].(string)
	interfaces := params["interfaces"].([]interface{})
	ip := ""
	port := ""
	for i, arg := range interfaces {
		if i == 0 {
			ip = arg.(map[string]interface{})["ip"].(string)
			port = arg.(map[string]interface{})["port"].(string)
		}
	}
	groups := params["groups"].([]interface{})
	groupId := ""
	for i, group := range groups {
		if i == 0 {
			groupId = group.(map[string]interface{})["groupid"].(string)
		}
	}

	templates := params["templates"].([]interface{})
	templateId := ""
	for i, template := range templates {
		if i == 0 {
			templateId = template.(map[string]interface{})["templateid"].(string)
		}
	}

	inventory := params["inventory"].(map[string]interface{})
	macAddr := inventory["macaddress_a"].(string) + inventory["macaddress_b"].(string)

	args2 := map[string]string{
		"host":       host,
		"ip":         ip,
		"port":       port,
		"groupId":    groupId,
		"templateId": templateId,
		"macAddr":    macAddr,
	}
	t := time.Now()
	timestamp := t.Unix()
	now := t.Format("2006-01-02 15:04:05")
	args := []interface{}{}
	args = append(args, host)
	args = append(args, timestamp)
	args = append(args, now)
	args = append(args, now)

	sqlcmd := "INSERT INTO graph.endpoint (endpoint,ts,t_create,t_modify) VALUES(?, ?, ?, ?)"
	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/graph?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	stmtIns, err := db.Prepare(sqlcmd)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})
	if sqlResult, err := stmtIns.Exec(host, timestamp, now, now); err == nil {
		if id, err := sqlResult.LastInsertId(); err == nil {
			hostid := strconv.Itoa(int(id))
			hostids := [1]string{string(hostid)}
			result["hostids"] = hostids
		} else {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
	resp["result"] = result
	RenderJson(rw, resp)
}

func hostDelete(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].([]interface{})

	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/graph?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	sqlcmd := "DELETE FROM graph.endpoint WHERE id=?"
	stmtIns, err := db.Prepare(sqlcmd)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	hostids := []string{}
	for _, hostId := range params {
		if result, err := stmtIns.Exec(hostId); err == nil {
			if RowsAffected, err := result.RowsAffected(); err == nil {
				if RowsAffected > 0 {
					hostids = append(hostids, hostId.(string))
				}
			} else {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}
	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})
	result["hostids"] = hostids
	resp["result"] = result
	RenderJson(rw, resp)
}

func hostUpdate(nodes map[string]interface{}, rw http.ResponseWriter) {
	log.Println("func hostUpdate()")
	params := nodes["params"].(map[string]interface{})
	hostName := params["host"].(string)
	hostId := params["hostid"].(string)
	now := time.Now().Format("2006-01-02 15:04:05")
	sqlcmd := "UPDATE graph.endpoint SET endpoint = ?, t_modify = ? WHERE id = ?"

	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/graph?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	stmtIns, err := db.Prepare(sqlcmd)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	var result = make(map[string]interface{})
	if sqlResult, err := stmtIns.Exec(hostName, now, hostId); err == nil {
		if RowsAffected, err := sqlResult.RowsAffected(); err == nil {
			if RowsAffected > 0 {
				hostids := [1]string{hostId}
				result["hostids"] = hostids
			}
		} else {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

func hostgroupCreate(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].(map[string]interface{})
	hostgroupName := params["name"].(string)

	sqlcmd := "INSERT INTO falcon_portal.grp (grp_name) VALUES(?)"
	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/falcon_portal?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	stmtIns, err := db.Prepare(sqlcmd)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	var result = make(map[string]interface{})
	if sqlResult, err := stmtIns.Exec(hostgroupName); err == nil {
		if id, err := sqlResult.LastInsertId(); err == nil {
			groupid := strconv.Itoa(int(id))
			groupids := [1]string{string(groupid)}
			result["groupids"] = groupids
		} else {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

func hostgroupDelete(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].([]interface{})

	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/falcon_portal?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	args := []interface{}{}
	args = append(args, "DELETE FROM falcon_portal.grp WHERE id=?")
	args = append(args, "DELETE FROM falcon_portal.grp_host WHERE grp_id=?")
	args = append(args, "DELETE FROM falcon_portal.grp_tpl WHERE grp_id=?")
	args = append(args, "DELETE FROM falcon_portal.plugin_dir WHERE grp_id=?")

	groupids := []string{}
	for _, sqlcmd := range args {
		stmtIns, err := db.Prepare(sqlcmd.(string))
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer stmtIns.Close()

		for _, hostgroupId := range params {
			if result, err := stmtIns.Exec(hostgroupId); err == nil {
				if RowsAffected, err := result.RowsAffected(); err == nil {
					if RowsAffected > 0 && sqlcmd == "DELETE FROM falcon_portal.grp WHERE id=?" {
						groupids = append(groupids, hostgroupId.(string))
						log.Println("delete hostgroup id:", hostgroupId)
					}
				} else {
					log.Println(err.Error())
				}
			} else {
				log.Println(err.Error())
			}
		}
	}
	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})
	result["groupids"] = groupids
	resp["result"] = result
	RenderJson(rw, resp)
}

func hostgroupUpdate(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].(map[string]interface{})
	hostgroupId := params["groupid"].(string)
	hostgroupName := params["name"].(string)
	sqlcmd := "UPDATE falcon_portal.grp SET grp_name = ? WHERE id = ?"

	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/falcon_portal?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	stmtIns, err := db.Prepare(sqlcmd)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	var result = make(map[string]interface{})
	if sqlResult, err := stmtIns.Exec(hostgroupName, hostgroupId); err == nil {
		if RowsAffected, err := sqlResult.RowsAffected(); err == nil {
			if RowsAffected > 0 {
				groupids := [1]string{hostgroupId}
				result["groupids"] = groupids
				log.Println("update groupid : ", hostgroupId)
			}
		} else {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

func templateCreate(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].(map[string]interface{})
	templateName := params["host"].(string)
	user := "root"
	groups := params["groups"]
	groupid := groups.(map[string]interface{})["groupid"].(json.Number)
	hostgroupId := string(groupid)

	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/falcon_portal?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	sqlcmd := "INSERT INTO falcon_portal.tpl (tpl_name, create_user) VALUES(?, ?)"

	stmtIns, err := db.Prepare(sqlcmd)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	var result = make(map[string]interface{})
	if sqlResult, err := stmtIns.Exec(templateName, user); err == nil {
		if id, err := sqlResult.LastInsertId(); err == nil {
			templateId := strconv.Itoa(int(id))
			templateids := [1]string{string(templateId)}
			result["templateids"] = templateids

			sqlcmd = "INSERT INTO falcon_portal.grp_tpl (grp_id, tpl_id, bind_user) VALUES(?, ?, ?)"
			stmtIns, err = db.Prepare(sqlcmd)
			if err != nil {
				log.Println(err.Error())
				return
			}
			defer stmtIns.Close()

			if result, err := stmtIns.Exec(hostgroupId, templateId, user); err == nil {
				if id, err := result.LastInsertId(); err == nil {
					log.Println("insert id :", id)
				} else {
					log.Println(err.Error())
				}
			} else {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
	resp := nodes
	delete(resp, "params")
	resp["result"] = result
	RenderJson(rw, resp)
}

func templateDelete(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].([]interface{})

	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/falcon_portal?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	args := []interface{}{}
	args = append(args, "DELETE FROM falcon_portal.tpl WHERE id=?")
	args = append(args, "DELETE FROM falcon_portal.grp_tpl WHERE tpl_id=?")

	templateids := []string{}
	for _, sqlcmd := range args {
		log.Println(sqlcmd)
		stmtIns, err := db.Prepare(sqlcmd.(string))
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer stmtIns.Close()

		for _, templateId := range params {
			if result, err := stmtIns.Exec(templateId); err == nil {
				if RowsAffected, err := result.RowsAffected(); err == nil {
					if RowsAffected > 0 {
						templateids = append(templateids, templateId.(string))
					}
					log.Println("delete id:", templateId)
				} else {
					log.Println(err.Error())
				}
			} else {
				log.Println(err.Error())
			}
		}
	}
	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})
	result["templateids"] = templateids
	resp["result"] = result
	RenderJson(rw, resp)
}

func templateUpdate(nodes map[string]interface{}, rw http.ResponseWriter) {
	params := nodes["params"].(map[string]interface{})
	templateId := params["templateid"].(string)
	templateName := params["name"].(string)
	sqlcmd := "UPDATE falcon_portal.tpl SET tpl_name = ? WHERE id = ?"

	str := config.Database.Account + ":" + config.Database.Password + "@tcp(" + config.Database.Addr + ")/falcon_portal?charset=utf8"
	db, err := sql.Open("mysql", str)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	if err != nil {
		log.Println("Oh noez, could not connect to database")
		return
	}

	stmtIns, err := db.Prepare(sqlcmd)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer stmtIns.Close()

	templateids := []string{}
	if result, err := stmtIns.Exec(templateName, templateId); err == nil {
		if RowsAffected, err := result.RowsAffected(); err == nil {
			if RowsAffected > 0 {
				templateids = append(templateids, templateId)
				log.Println("update groupid : ", templateId)
			}
		} else {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
	resp := nodes
	delete(resp, "params")
	var result = make(map[string]interface{})
	result["templateids"] = templateids
	resp["result"] = result
	RenderJson(rw, resp)
}

func apiParser(rw http.ResponseWriter, req *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	s := buf.String() // Does a complete copy of the bytes in the buffer.

	json, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		log.Println(err.Error())
	}

	f, err := os.OpenFile("falcon_api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	var nodes = make(map[string]interface{})
	nodes, _ = json.Map()

	method := nodes["method"]
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

func parseConfig(cfg string) {
	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}
	ConfigFile = cfg
	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
		return
	}
	lock.Lock()
	defer lock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
}

func configZabbixRoutes() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()
	parseConfig(*cfg)
	http.HandleFunc("/api", apiParser)
}
