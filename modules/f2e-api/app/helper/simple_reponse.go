package helper

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DataWaper struct {
	Data interface{} `json:"data,omitempty"`
	Page interface{} `json:"page,omitempty"`
	Msg  string      `json:"message,omitempty"`
}

type RespJson struct {
	Error string `json:"error,omitempty"`
	Msg   string `json:"message,omitempty"`
}

// func JSONR(c *gin.Context, wcode int, msg interface{}) (werror error) {
func JSONR(c *gin.Context, arg ...interface{}) (werror error) {
	var (
		wcode int
		msg   interface{}
	)
	if len(arg) == 1 {
		wcode = http.StatusOK
		msg = arg[0]
	} else {
		wcode = arg[0].(int)
		msg = arg[1]
	}
	need_doc := viper.GetBool("gen_doc")
	var body interface{}
	defer func() {
		if need_doc {
			ds, _ := json.Marshal(body)
			bodys := string(ds)
			log.Debugf("body: %v, bodys: %v ", body, bodys)
			c.Set("body_doc", bodys)
		}
	}()
	if wcode == 200 {
		switch msg.(type) {
		case string:
			body = RespJson{Msg: msg.(string)}
			c.JSON(http.StatusOK, body)
		default:
			c.JSON(http.StatusOK, msg)
			body = msg
		}
	} else {
		switch msg.(type) {
		case string:
			msgerr := checkEnv(c, msg.(string))
			body = RespJson{Error: msgerr}
			c.JSON(wcode, body)
		case error:
			msgerr := checkEnv(c, msg.(error).Error())
			body = RespJson{Error: msgerr}
			c.JSON(wcode, body)
		default:
			body = RespJson{Error: "system type error. please ask admin for help"}
			c.JSON(wcode, body)
		}
	}
	return
}

func checkEnv(c *gin.Context, errMsg string) (msg string) {
	msg = errMsg
	// check error message is from sql?
	matched, _ := regexp.MatchString("^\\s*Error\\s\\d+", errMsg)
	if viper.Get("env") == "production" && matched {
		msg = "got error of sql query. please check your api params"
		log.Errorf("sql error: %v, url: %v, form ip: %v", errMsg, c.Request.URL.String(), c.ClientIP())
	}
	return
}
