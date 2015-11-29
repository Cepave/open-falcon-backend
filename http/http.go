package http

import (
	"strings"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/Cepave/query/g"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func InitDatabase() {
	// set default database
	config := g.Config()
	orm.RegisterDataBase("default", "mysql", config.Db.Addr, config.Db.Idle, config.Db.Max)
	// register model
	orm.RegisterModel(new(Endpoint))

	strConn := strings.Replace(config.Db.Addr, "graph", "falcon_portal", 1)
	orm.RegisterDataBase("falcon_portal", "mysql", strConn, config.Db.Idle, config.Db.Max)
	orm.RegisterModel(new(Grp), new(Grp_host), new(Grp_tpl), new(Tpl))

	strConn = strings.Replace(config.Db.Addr, "graph", "grafana", 1)
	orm.RegisterDataBase("grafana", "mysql", strConn, config.Db.Idle, config.Db.Max)
	orm.RegisterModel(new(Province), new(City), new(Idc))

	if config.Debug == true {
		orm.Debug = true
	}
}

func Start() {
	if !g.Config().Http.Enable {
		log.Println("http.Start warning, not enable")
		return
	}

	// config http routes
	configCommonRoutes()
	configProcHttpRoutes()
	configGraphRoutes()
	configApiRoutes()
	configGrafanaRoutes()
	configZabbixRoutes()

	// start mysql database
	InitDatabase()

	// start http server
	addr := g.Config().Http.Listen
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http.Start ok, listening on", addr)
	log.Fatalln(s.ListenAndServe())
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

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	RenderDataJson(w, data)
}

func StdRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)
		RenderMsgJson(w, err.Error())
		return
	}
	RenderJson(w, data)
}
