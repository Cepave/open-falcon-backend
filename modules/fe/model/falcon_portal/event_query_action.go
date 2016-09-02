package falconPortal

import (
	"fmt"
	"time"

	"strings"

	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/uic"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
)

//generate status filter SQL templete
func genStatusQueryTemplete(status string, feildsName string) (filterTemplete string) {
	statusList := strings.Split(status, ",")
	var filterTempleteArr []string
	for _, s := range statusList {
		filterTempleteArr = append(filterTempleteArr, fmt.Sprintf("%s = '%s'", feildsName, s))
	}
	filterTemplete = strings.Join(filterTempleteArr, " OR ")
	filterTemplete = fmt.Sprintf("( %s )", filterTemplete)
	return
}

func genSqlFilterTemplete(whereConditions []string) string {
	if len(whereConditions) == 0 {
		return ""
	}
	conditions := strings.Join(whereConditions, " AND ")
	return fmt.Sprintf("WHERE %s", conditions)
}

const SkipFilter = "ALL"

func GetEventCases(startTime int64, endTime int64, priority int, status string, progressStatus string, limit int, elimit int, username string, metrics string, caseId string) (result []EventCases, err error) {
	config := g.Config()
	q := orm.NewOrm()
	q.Using("falcon_portal")
	var whereConditions []string
	if limit == 0 || limit > config.FalconPortal.Limit {
		limit = config.FalconPortal.Limit
	}

	isadmin, tplids, err := GetCasePermission(username)
	if tplids == "" {
		tplids = "-1"
	}

	//fot generate sql filter
	if startTime != 0 && endTime != 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("update_at BETWEEN FROM_UNIXTIME(%d) AND FROM_UNIXTIME(%d)", startTime, endTime))
	}
	if priority != -1 {
		whereConditions = append(whereConditions, fmt.Sprintf("priority = %d", priority))
	}
	if status != SkipFilter {
		log.Debug("statis ", status)
		whereConditions = append(whereConditions, genStatusQueryTemplete(status, "status"))
	}
	if progressStatus != SkipFilter {
		whereConditions = append(whereConditions, genStatusQueryTemplete(progressStatus, "process_status"))
	}
	if metrics != SkipFilter {
		whereConditions = append(whereConditions, genStatusQueryTemplete(metrics, "metric"))
	}
	if caseId != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("id = '%s'", caseId))
	}
	//perpare ssql statement
	if !isadmin {
		whereConditions = append(whereConditions, fmt.Sprintf("(tpl_creator = '%s' OR template_id in (%s))", username, tplids))
	}
	_, err = q.Raw(fmt.Sprintf("SELECT * FROM `event_cases` %s limit %d", genSqlFilterTemplete(whereConditions), limit)).QueryRows(&result)

	if len(result) == 0 {
		result = []EventCases{}
		return
	}
	//set default number of event
	var eventLimit int
	if eventLimit = elimit; elimit == 0 {
		eventLimit = 10
	}
	for indx, event := range result {
		var eventArr []*Events
		q.Raw(fmt.Sprintf("SELECT * FROM `events` WHERE event_caseId = '%s' order by timestamp DESC Limit %d", event.Id, eventLimit)).QueryRows(&eventArr)
		if len(eventArr) != 0 {
			event.Events = eventArr
		} else {
			event.Events = []*Events{}
		}
		result[indx] = event
	}
	return
}

func GetEvents(startTime int64, endTime int64, status string, limit int, caseId string) (result []EventsRsp, err error) {

	config := g.Config()
	q := orm.NewOrm()
	q.Using("falcon_portal")

	var whereConditions []string

	if status != SkipFilter {
		if status == "OK" {
			whereConditions = append(whereConditions, fmt.Sprintf("events.status = %d", 1))
		} else if status == "PROBLEM" {
			whereConditions = append(whereConditions, fmt.Sprintf("events.status = %d", 0))
		}
	}

	if startTime != 0 && endTime != 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("events.timestamp BETWEEN FROM_UNIXTIME(%d) AND FROM_UNIXTIME(%d)", startTime, endTime))
	}
	if caseId != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("events.event_caseId = '%s'", caseId))
	}
	if limit == 0 {
		limit = config.FalconPortal.Limit
	}

	_, err = q.Raw(fmt.Sprintf(`SELECT events.id as id,
				events.step as step,
				events.cond as cond,
				events.timestamp as timestamp,
				events.event_caseId as eid,
				event_cases.tpl_creator as tpl_creator,
				event_cases.metric as metric,
				event_cases.endpoint as endpoint
				FROM events LEFT JOIN event_cases on event_cases.id = events.event_caseId
				%s ORDER BY events.timestamp DESC limit %d`, genSqlFilterTemplete(whereConditions), limit)).QueryRows(&result)

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

func GetNotes(eventCaseId string, limit int, startTime int64, endTime int64, filterIgnored bool) (enotes []EventNote, err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	whereConditions := []string{}
	if eventCaseId != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("event_note.event_caseId = '%s' ", eventCaseId))
	}
	switch {
	//allow api only set the startTime and use the currentTime as the endTime
	case startTime != 0 && endTime == 0:
		endTime = time.Now().Unix()
	case startTime != 0 && endTime != 0:
		tempTime := ""
		q.Raw("SELECT timestamp FROM event_cases WHERE id = ?", eventCaseId).QueryRow(&tempTime)
		if tempTime != "" {
			myzone, _ := time.Now().Zone()
			parsedTime, err := time.Parse("2006-01-02 15:04:05 MST", fmt.Sprintf("%s %s", tempTime, myzone))
			log.Debugf("got time: %v , convertedTime: %v, Unix: %v", fmt.Sprintf("%s %s", tempTime, myzone), parsedTime, parsedTime.Unix())
			if err == nil {
				startTime = parsedTime.Unix()
			} else {
				log.Debug(err.Error())
			}
		}
		endTime = time.Now().Unix()
	}
	if startTime > 0 && endTime > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("event_note.timestamp BETWEEN FROM_UNIXTIME(%d) AND FROM_UNIXTIME(%d)", startTime, endTime))
	}
	if filterIgnored {
		whereConditions = append(whereConditions, "event_note.status != 'ignored' ")
	}
	limitTemplete := ""
	if limit != 0 {
		limitTemplete = fmt.Sprintf("LIMIT %d", limit)
	}
	sqlTemplete := fmt.Sprintf(`SELECT event_note.id as id,
		event_note.event_caseId as event_caseId,
		event_note.note as note,
		event_note.case_id as case_id,
		event_note.event_caseId as eid,
		event_note.status as status,
		event_note.timestamp as timestamp,
		user.name as user_name
		FROM falcon_portal.event_note as event_note LEFT JOIN uic.user as user on event_note.user_id = user.id
		%s ORDER BY event_note.timestamp DESC %s`, genSqlFilterTemplete(whereConditions), limitTemplete)
	_, err = q.Raw(sqlTemplete).QueryRows(&enotes)
	if len(enotes) == 0 {
		enotes = []EventNote{}
	}
	return
}

func GetNote(noteId int64) (EventNote, error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	var eventNote EventNote
	err := q.Raw(`SELECT * from event_note WHERE event_note.id = ?`, noteId).QueryRow(&eventNote)
	if err == nil {
		user := uic.ReadUserById(eventNote.UserId)
		eventNote.UserName = user.Name
	}
	return eventNote, err
}
