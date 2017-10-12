package cron

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/alarm/api"
	"github.com/Cepave/open-falcon-backend/modules/alarm/g"
	db "github.com/Cepave/open-falcon-backend/modules/alarm/model"
	"testing"
)

func init() {
	g.ParseConfig("../cfg.json")
	g.InitRedisConnPool()
	db.InitDatabase()
}

// for OWL-1980 to check if phone field is empty string
func TestParseUserSms(t *testing.T) {
	e := &model.Event{
		Id: "s_2226_a904da7af7aaf52a466454e8943517f4",
		Strategy: &model.Strategy{0,
			"metric_str",
			map[string]string{},
			"func_str",
			"oper_str",
			90,
			3,
			5,
			"note_str", nil},
		Expression:  nil,
		Status:      "PROBLEM",
		Endpoint:    "aaa",
		LeftValue:   0,
		CurrentStep: 0,
		EventTime:   0,
		PushedTags:  map[string]string{}}
	a := &api.Action{
		Id:                 5,
		Uic:                "ethanhao",
		Url:                "qqqqqq",
		Callback:           5,
		BeforeCallbackSms:  0,
		BeforeCallbackMail: 0,
		AfterCallbackSms:   1,
		AfterCallbackMail:  2}
	ParseUserSms(e, a)
}
