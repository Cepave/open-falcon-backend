package cron

import (
	"encoding/json"
	"fmt"
	"github.com/Cepave/alarm/logger"
	"time"

	"github.com/Cepave/alarm/g"
	eventmodel "github.com/Cepave/alarm/model/event"
	"github.com/Cepave/common/model"
	"github.com/garyburd/redigo/redis"
)

func ReadHighEvent() {
	queues := g.Config().Redis.HighQueues
	if len(queues) == 0 {
		return
	}

	for {
		event, err := popEvent(queues)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		consume(event, true)
	}
}

func ReadLowEvent() {
	queues := g.Config().Redis.LowQueues
	if len(queues) == 0 {
		return
	}

	for {
		event, err := popEvent(queues)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		consume(event, false)
	}
}

func popEvent(queues []string) (*model.Event, error) {
	log := logger.Logger()
	count := len(queues)

	params := make([]interface{}, count+1)
	for i := 0; i < count; i++ {
		params[i] = queues[i]
	}
	// set timeout 0
	params[count] = 0

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	reply, err := redis.Strings(rc.Do("BRPOP", params...))
	if err != nil {
		log.Error(fmt.Sprintf("get alarm event from redis fail: %v", err))
		return nil, err
	}

	var event model.Event
	err = json.Unmarshal([]byte(reply[1]), &event)
	if err != nil {
		log.Error(fmt.Sprintf("parse alarm event fail: %v", err))
		return nil, err
	}

	log.Debug(event.String())
	//insert event into database
	eventmodel.InsertEvent(&event)
	// save in memory. display in dashboard
	g.Events.Put(&event)

	return &event, nil
}
