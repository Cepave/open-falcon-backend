package portal

import (
	"sync"

	"strconv"
	"strings"
	"time"

	event "github.com/Cepave/open-falcon-backend/modules/fe/model/falcon_portal"
	log "github.com/Sirupsen/logrus"
)

type UpdatedEvents struct {
	Events []event.EventCases
	Enotes []event.EventNote
}
type UpdatedObjMutex struct {
	sync.Mutex
	UpdatedEvents
}

var (
	storeEvents UpdatedObjMutex
)

func (this *PortalController) OnTimeFeeding() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		log.Error(err.Error())
		this.ResposeError(baseResp, err.Error())
		return
	}
	username := this.GetString("cName", "")
	isadmin, tplids, err := event.GetCasePermission(username)
	if tplids == "" {
		tplids = "-1"
	}
	events := []event.EventCases{}
	events = storeEvents.Events
	//not admin & no any related alarms
	if !isadmin && tplids == "-1" {
		events = []event.EventCases{}
	} else if !isadmin && tplids != "-1" {
		log.Debugf("this user is not admin, will use %s for query", tplids)
		eventsTmp := []event.EventCases{}
		//make a templete mapping for check the permission of this user
		tplidsArr := strings.Split(tplids, ",")
		log.Debugf("tplarr: %v", tplidsArr)
		tplidMap := map[int]int{}
		for _, tid := range tplidsArr {
			tidint, err := strconv.Atoi(tid)
			if err != nil {
				log.Error(err.Error())
			}
			tplidMap[tidint] = tidint
		}
		for _, e := range events {
			if _, ok := tplidMap[e.TemplateId]; ok {
				eventsTmp = append(eventsTmp, e)
			}
		}
		events = eventsTmp
	}
	notes := storeEvents.Enotes
	anyNew := false
	if len(events) != 0 || len(notes) != 0 {
		anyNew = true
	}
	baseResp.Data["events"] = events
	baseResp.Data["notes"] = notes
	baseResp.Data["any_new"] = anyNew
	baseResp.Data["admin"] = isadmin
	this.ServeApiJson(baseResp)
	return
}

func CronForQuery(updatTo chan UpdatedEvents, pid chan string) {
	defer func() {
		if r := recover(); r != nil {
			time.Sleep(time.Minute * 1)
			pid <- "CronForQuery"
			return
		}
	}()

	for {
		currentTime := time.Now().Unix()
		startTime := (currentTime - int64(60*1))
		log.Debugf("cron job for on time query: range -> %v ~ %v", startTime, currentTime)
		events, err := event.GetEventCases(false, startTime, currentTime, -1, "ALL", "ALL", 500, 0, "root", "ALL", "")
		if err != nil {
			log.Errorf("crond get evnetCase go err: %v", err.Error())
		}
		enotes, err := event.GetNotes("", 1000, startTime, currentTime, false)
		if err != nil {
			log.Errorf("crond get GetNotes go err: %v", err.Error())
		}
		updatTo <- UpdatedEvents{events, enotes}
		time.Sleep(time.Minute * 1)
	}
}

func CronReciveUpdate(updatTo chan UpdatedEvents, pid chan string) {
	defer func() {
		if r := recover(); r != nil {
			time.Sleep(time.Minute * 1)
			pid <- "CronReciveUpdate"
			return
		}
	}()
	for {
		select {
		case updated := <-updatTo:
			storeEvents.Lock()
			storeEvents.UpdatedEvents = updated
			storeEvents.Unlock()
		}
	}
}

func CornDaemonStart() {
	log.Info("on time query cron job start!")
	defer log.Error("cron stoped , on time feeding function will broke, please restart this application")
	updateChann := make(chan UpdatedEvents)
	supervisorChn := make(chan string)
	go CronForQuery(updateChann, supervisorChn)
	go CronReciveUpdate(updateChann, supervisorChn)

	for {
		select {
		case sup := <-supervisorChn:
			if sup == "CronForQuery" {
				log.Errorf("%s dead will unknown reason, will restart the this rotuine", sup)
				go CronForQuery(updateChann, supervisorChn)
			} else if sup == "CornReceiveUpdate" {
				log.Errorf("%s dead will unknown reason, will restart the this rotuine", sup)
				go CronReciveUpdate(updateChann, supervisorChn)
			} else {
				log.Errorf("got worng params of supervisorChn -> %v .", sup)
			}
		}
	}
}
