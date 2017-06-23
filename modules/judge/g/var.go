package g

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/Cepave/open-falcon-backend/common/model"
	log "github.com/sirupsen/logrus"
)

type SafeStrategyMap struct {
	sync.RWMutex
	// endpoint:metric => [strategy1, strategy2 ...]
	M map[string][]model.Strategy
}

type SafeExpressionMap struct {
	sync.RWMutex
	// metric:tag1 => [exp1, exp2 ...]
	// metric:tag2 => [exp1, exp2 ...]
	M map[string][]*model.Expression
}

type SafeEventMap struct {
	sync.RWMutex
	M map[string]*model.Event
}

var (
	HbsClient     *SingleConnRpcClient
	StrategyMap   = &SafeStrategyMap{M: make(map[string][]model.Strategy)}
	ExpressionMap = &SafeExpressionMap{M: make(map[string][]*model.Expression)}
	LastEvents    = &SafeEventMap{M: make(map[string]*model.Event)}
)

func InitHbsClient() {
	HbsClient = &SingleConnRpcClient{
		RpcServers: Config().Hbs.Servers,
		Timeout:    time.Duration(Config().Hbs.Timeout) * time.Millisecond,
	}
}

func InitLastEvents() {

	if !Config().Alarm.StoreEventToFile {
		return
	}

	filepath := Config().Alarm.EventsStoreFilePath
	_, err := os.Stat(filepath)
	if err == nil {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Error(err.Error())
		}
		var v []*model.Event
		err = json.Unmarshal(data, &v)
		if err != nil {
			log.Error(err.Error())
		}
		for _, event := range v {
			LastEvents.Set(event.Id, event)
		}
	} else {
		log.Error(err.Error())
	}

	if err != nil {
		log.Error(err.Error())
	} else if len(LastEvents.GetAll()) != 0 {
		log.Infof("init lastEvent of %s , %d events is inserted", filepath, len(LastEvents.GetAll()))
	}
}

func (this *SafeStrategyMap) ReInit(m map[string][]model.Strategy) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeStrategyMap) Get() map[string][]model.Strategy {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeExpressionMap) ReInit(m map[string][]*model.Expression) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeExpressionMap) Get() map[string][]*model.Expression {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeEventMap) GetAll() map[string]*model.Event {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeEventMap) Get(key string) (*model.Event, bool) {
	this.RLock()
	defer this.RUnlock()
	event, exists := this.M[key]
	return event, exists
}

func (this *SafeEventMap) Set(key string, event *model.Event) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = event
}

func (this *SafeEventMap) Delete(key string) bool {
	this.Lock()
	defer this.Unlock()
	delete(this.M, key)
	_, exist := this.M[key]
	return exist
}
