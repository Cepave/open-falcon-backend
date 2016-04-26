package event

import (
	"fmt"
	"time"

	"github.com/Cepave/fe/g"
	"github.com/Cepave/fe/model/uic"
	"github.com/astaxie/beego/orm"
)

func getUserRole(username string) bool {
	user := uic.ReadUserByName(username)
	if user.Role == 2 {
		return true
	} else {
		return false
	}
}

func GetEventCases(startTime int64, endTime int64, priority int, status string, limit int, username string) (result []EventCases, err error) {
	config := g.Config()
	if limit == 0 || limit > config.FalconPortal.Limit {
		limit = config.FalconPortal.Limit
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
			queryTmp = fmt.Sprintf("%v and status = '%s'", queryTmp, status)
		} else {
			flag = true
			queryTmp = fmt.Sprintf("%v status = '%s'", queryTmp, status)
		}
	}
	isadmin := getUserRole(username)
	if queryTmp != "" && !isadmin {
		_, err = q.Raw(fmt.Sprintf("SELECT * FROM `event_cases` WHERE tpl_creator = '%s' AND %v order by update_at DESC limit %d", username, queryTmp, limit)).QueryRows(&result)
	} else if isadmin {
		_, err = q.Raw(fmt.Sprintf("SELECT * FROM `event_cases` WHERE %v order by update_at DESC limit %d", queryTmp, limit)).QueryRows(&result)
	} else {
		_, err = q.Raw(fmt.Sprintf("SELECT * FROM `event_cases` WHERE tpl_creator = '%s' order by update_at DESC", username)).QueryRows(&result)
	}

	if len(result) == 0 {
		result = []EventCases{}
	} else {
		for indx, event := range result {
			var eventArr []*Events
			q.Raw(fmt.Sprintf("SELECT * FROM `events` WHERE event_caseId = '%s' order by timestamp DESC", event.Id)).QueryRows(&eventArr)
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

func CloseEvent(username string, colsed_note string, id string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	userid := uic.ReadUserIdByName(username)
	_, err = q.Raw("Update event_cases SET user_modified = ?, closed_at = ?, status = ?, closed_note = ? WHERE id = ?", userid, time.Now(), "SOLVED", colsed_note, id).Exec()
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
