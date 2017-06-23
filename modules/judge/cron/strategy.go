package cron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/judge/g"
	log "github.com/sirupsen/logrus"
)

func SyncStrategies(pid chan string) {

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("run time panic: %v", r)
			pid <- "SyncStrategies"
			return
		}
	}()

	duration := time.Duration(g.Config().Hbs.Interval) * time.Second
	for {
		syncStrategies()
		syncExpression()
		dumpAllJudgedEvents()
		time.Sleep(duration)
	}
}

//dump current events into file, can solve the problem of huge alarm triggered when restart judge
func dumpAllJudgedEvents() {
	event := g.LastEvents.GetAll()
	if len(event) == 0 {
		return
	}
	var tmpStore []*model.Event
	for _, value := range event {
		tmpStore = append(tmpStore, value)
	}
	data, err := json.Marshal(tmpStore)
	if err != nil {
		log.Error(err.Error())
	}

	filepath := g.Config().Alarm.EventsStoreFilePath

	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug("dumpAllJudgedEvents into local file.")
}

func syncStrategies() {
	var strategiesResponse model.StrategiesResponse
	err := g.HbsClient.Call("Hbs.GetStrategies", model.NullRpcRequest{}, &strategiesResponse)
	if err != nil {
		log.Println("[ERROR] Hbs.GetStrategies:", err)
		return
	}

	rebuildStrategyMap(&strategiesResponse)
}

func rebuildStrategyMap(strategiesResponse *model.StrategiesResponse) {
	// endpoint:metric => [strategy1, strategy2 ...]
	m := make(map[string][]model.Strategy)
	for _, hs := range strategiesResponse.HostStrategies {
		hostname := hs.Hostname
		if g.Config().Debug && hostname == g.Config().DebugHost {
			log.Println(hostname, "strategies:")
			bs, _ := json.Marshal(hs.Strategies)
			fmt.Println(string(bs))
		}
		for _, strategy := range hs.Strategies {
			key := fmt.Sprintf("%s/%s", hostname, strategy.Metric)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], strategy)
			} else {
				m[key] = []model.Strategy{strategy}
			}
		}
	}

	g.StrategyMap.ReInit(m)
}

func syncExpression() {
	var expressionResponse model.ExpressionResponse
	err := g.HbsClient.Call("Hbs.GetExpressions", model.NullRpcRequest{}, &expressionResponse)
	if err != nil {
		log.Println("[ERROR] Hbs.GetExpressions:", err)
		return
	}

	rebuildExpressionMap(&expressionResponse)
}

func rebuildExpressionMap(expressionResponse *model.ExpressionResponse) {
	m := make(map[string][]*model.Expression)
	for _, exp := range expressionResponse.Expressions {
		for k, v := range exp.Tags {
			key := fmt.Sprintf("%s/%s=%s", exp.Metric, k, v)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], exp)
			} else {
				m[key] = []*model.Expression{exp}
			}
		}
	}

	g.ExpressionMap.ReInit(m)
}
