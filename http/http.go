package http

import (
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

func Start() {
	if !g.Config().Http.Enable {
		log.Println("http.Start warning, not enable")
		return
	}

	// config http routes
	configCommonRoutes()
	configProcHttpRoutes()
	configGraphRoutes()
	configZabbixRoutes()

	// start mysql database
	account := g.Config().Database.Account
	password := g.Config().Database.Password
	ip := g.Config().Database.Ip
	port := g.Config().Database.Port
	database := "graph"
	strConn := account + ":" + password + "@tcp(" + ip + ":" + port + ")/" + database + "?charset=utf8"
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	maxIdle := 30
	maxConn := 30
	orm.RegisterDataBase("default", "mysql", strConn, maxIdle, maxConn)

	database = "falcon_portal"
	strConn = account + ":" + password + "@tcp(" + ip + ":" + port + ")/" + database + "?charset=utf8"
	orm.RegisterDataBase("falcon_portal", "mysql", strConn, maxIdle, maxConn)
	orm.RegisterModel(new(Endpoint))
	orm.RegisterModel(new(Grp))
	orm.RegisterModel(new(Grp_host))
	orm.RegisterModel(new(Grp_tpl))
	orm.RegisterModel(new(Tpl))

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
