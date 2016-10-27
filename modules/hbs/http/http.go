package http

import (
	"encoding/json"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"gopkg.in/gin-gonic/gin.v1"
	log "github.com/Sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var ginRouter *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	ginRouter = gin.New()

	configCommonRoutes(ginRouter)
	configProcRoutes(ginRouter)
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

func Start() {
	if !g.Config().Http.Enabled {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}
	s := &http.Server{
		Addr: addr,
		Handler: ginRouter,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http listening", addr)
	log.Fatalln(s.ListenAndServe())
}
