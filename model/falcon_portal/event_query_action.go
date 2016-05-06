package falconPortal

import (
	"fmt"

	"github.com/Cepave/fe/g"
	"github.com/astaxie/beego/orm"
)

func GetEventCases(startTime int64, endTime int64, priority int, status string, limit int, elimit int, username string) (result []EventCases, err error) {
	config := g.Config()
	if limit == 0 || limit > config.FalconPortal.Limit {
		limit = config.FalconPortal.Limit
	}

	isadmin, tplids, err := GetCasePermission(username)
	if tplids == "" {
		tplids = "-1"
	}
	q := orm.NewOrm()
	q.Using("falcon_portal")
	flag := false
	queryTmp := ""
	if startTime != 0 && endTime != 0 {
		flag = true
		queryTmp = fmt.Sprintf(" %v update_at >= %d and  update_at <= %d", queryTmp, startTime, endTime)
	}
	if priority != -1 {
		if flag {
			queryTmp = fmt.Sprintf("%v and priority = %d", queryTmp, priority)
		} else {
			flag = true
			queryTmp = fmt.Sprintf("%v priority = %d", queryTmp, priority)
		}
	}
	if status != "ALL" {
		if flag {
			queryTmp = fmt.Sprintf("%v and (status = '%s' or status = 'OK')", queryTmp, status)
		} else {
			flag = true
			queryTmp = fmt.Sprintf("%v (status = '%s' or status = 'OK')", queryTmp, status)
		}
	}
	if queryTmp != "" && !isadmin {
		_, err = q.Raw(fmt.Sprintf("SELECT * FROM `event_cases` WHERE (tpl_creator = '%s' OR template_id in (%s)) AND %v order by update_at DESC limit %d", username, tplids, queryTmp, limit)).QueryRows(&result)
	} else {
		_, err = q.Raw(fmt.Sprintf("SELECT * FROM `event_cases` WHERE %v order by update_at DESC limit %d", queryTmp, limit)).QueryRows(&result)
	}

	if len(result) == 0 {
		result = []EventCases{}
	} else {
		var eventLimit int
		if eventLimit = elimit; elimit == 0 {
			eventLimit = 10
		}
		for indx, event := range result {
			var eventArr []*Events
			q.Raw(fmt.Sprintf("SELECT * FROM `events` WHERE event_caseId = '%s' order by timestamp DESC Limit %d", event.Id, eventLimit)).QueryRows(&eventArr)
			fmt.Sprintf("%v", eventArr)
			if len(eventArr) != 0 {
				event.Events = eventArr
			} else {
				event.Events = []*Events{}
			}
			result[indx] = event
		}
	}
	return
}

func GetEvents(startTime int64, endTime int64, limit int) (result []EventsRsp, err error) {
	config := g.Config()
	if limit == 0 || limit > config.FalconPortal.Limit {
		limit = config.FalconPortal.Limit
	}

	q := orm.NewOrm()
	q.Using("falcon_portal")
	queryTmp := ""
	if startTime != 0 && endTime != 0 {
		queryTmp = fmt.Sprintf(" %v events.timestamp >= %d and  events.timestamp <= %d", queryTmp, startTime, endTime)
	}
	if queryTmp != "" {
		_, err = q.Raw(fmt.Sprintf(`SELECT events.id as id,
					events.step as step,
					events.cond as cond,
					events.timestamp as timestamp,
					events.event_caseId as eid,
					event_cases.tpl_creator as tpl_creator,
					event_cases.metric as metric,
					event_cases.endpoint as endpoint
					FROM events LEFT JOIN event_cases on event_cases.id = events.event_caseId
					WHERE %v ORDER BY events.timestamp DESC limit %d`, queryTmp, limit)).QueryRows(&result)
	} else {
		_, err = q.Raw(fmt.Sprintf(`SELECT
					events.id as id,
					events.step as step,
					events.cond as cond,
					events.timestamp as timestamp,
					events.event_caseId as eid,
					event_cases.tpl_creator as tpl_creator,
					event_cases.metric as metric,
					event_cases.endpoint as endpoint
					FROM events LEFT JOIN event_cases on event_cases.id = events.event_caseId
					ORDER BY events.timestamp DESC limit %d`, limit)).QueryRows(&result)
	}
	if len(result) == 0 {
		result = []EventsRsp{}
	}
	return
}

func CountNumOfTlp() (c int, err error) {
	var h []Tpl
	q := orm.NewOrm()
	q.Using("falcon_portal")
	_, err = q.Raw("select * from `tpl`").QueryRows(&h)
	c = len(h)
	return
}
