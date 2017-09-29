package event

import (
	"fmt"
	"time"

	"database/sql"

	coommonModel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/alarm/model/boss"
	"github.com/astaxie/beego/orm"
	log "github.com/sirupsen/logrus"
)

const timeLayout = "2006-01-02 15:04:05"

func insertEvent(q orm.Ormer, eve *coommonModel.Event) (res interface{}, err error) {
	var status int
	if status = 0; eve.Status == "OK" {
		status = 1
	}
	sqltemplete := `INSERT INTO events (
		event_caseId,
		step,
		cond,
		status,
		timestamp
	) VALUES(?,?,?,?,?)`
	res, err = q.Raw(
		sqltemplete,
		eve.Id,
		eve.CurrentStep,
		fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
		status,
		time.Unix(eve.EventTime, 0).Format(timeLayout),
	).Exec()
	return
}

func getBossInfoByEndpointName(name string) boss.BossInfo {
	q := orm.NewOrm()
	bossInfo := []boss.BossInfo{}
	q.Using("boss")
	q.Raw(`select h.hostname as hostname, h.exist as exist, h.activate as activate, h.platform as platform, h.platforms as platforms,
		 		 h.idc as idc, h.ip as ip, h.isp as isp, h.province as province, pt.contacts as contacts
				 from hosts as h left join platforms as pt on h.platform = pt.platform where h.hostname = ? and h.exist != ?`, name, 0).QueryRows(&bossInfo)
	if len(bossInfo) == 0 {
		b := boss.BossInfo{}
		return b.New()
	} else {
		return bossInfo[0]
	}
}

func InsertEvent(eve *coommonModel.Event, alarmType string) error {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	var event []EventCases
	var sqlLog sql.Result
	var errRes error
	var aType AlarmType
	errRes, aType = findAlarmTypeByName(alarmType)
	if errRes != nil {
		return errRes
	}
	q.Raw("select * from event_cases where id = ? limit 1", eve.Id).QueryRows(&event)
	log.Debugf("events: %v", eve)
	log.Debugf("express is null: %v", eve.Expression == nil)
	bossInfo := getBossInfoByEndpointName(eve.Endpoint)
	if len(event) == 0 {
		//create cases
		sqltemplete := `INSERT INTO event_cases (
					id,
					endpoint,
					metric,
					func,
					cond,
					note,
					max_step,
					current_step,
					priority,
					status,
					timestamp,
					update_at,
					tpl_creator,
					expression_id,
					strategy_id,
					template_id,
					alarm_type_id,
					ip,
					idc,
					platform,
					contact) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		tpl_creator := ""
		if eve.Tpl() != nil {
			tpl_creator = eve.Tpl().Creator
		}
		sqlLog, errRes = q.Raw(
			sqltemplete,
			eve.Id,
			eve.Endpoint,
			counterGen(eve.Metric(), utils.SortedTags(eve.PushedTags)),
			eve.Func(),
			//cond
			fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			eve.Note(),
			eve.MaxStep(),
			eve.CurrentStep,
			eve.Priority(),
			eve.Status,
			//start_at
			time.Unix(eve.EventTime, 0).Format(timeLayout),
			//update_at
			time.Unix(eve.EventTime, 0).Format(timeLayout),
			tpl_creator,
			eve.ExpressionId(),
			eve.StrategyId(),
			//template_id
			eve.TplId(),
			aType.Id,
			bossInfo.IP,
			bossInfo.Idc,
			bossInfo.Platform,
			bossInfo.Contact(),
		).Exec()
	} else {
		sqltemplete := `UPDATE event_cases SET
				update_at = ?,
				max_step = ?,
				current_step = ?,
				note = ?,
				cond = ?,
				status = ?,
				func = ?,
				priority = ?,
				tpl_creator = ?,
				expression_id = ?,
				strategy_id = ?,
				template_id = ?,
				alarm_type_id = ?,
				ip = ?,
				idc = ?,
				platform = ?,
				contact = ?`
		//reopen case
		if event[0].ProcessStatus == "resolved" || event[0].ProcessStatus == "ignored" {
			sqltemplete = fmt.Sprintf("%v ,process_status = '%s', process_note = %d", sqltemplete, "unresolved", 0)
		}

		tpl_creator := ""
		if eve.Tpl() != nil {
			tpl_creator = eve.Tpl().Creator
		}
		if eve.CurrentStep == 1 && eve.Status != "OK" {
			//update start time of cases
			sqltemplete = fmt.Sprintf("%v ,timestamp = ? WHERE id = ?", sqltemplete)
			sqlLog, errRes = q.Raw(
				sqltemplete,
				time.Unix(eve.EventTime, 0).Format(timeLayout),
				eve.MaxStep(),
				eve.CurrentStep,
				eve.Note(),
				fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
				eve.Status,
				eve.Func(),
				eve.Priority(),
				tpl_creator,
				eve.ExpressionId(),
				eve.StrategyId(),
				eve.TplId(),
				aType.Id,
				bossInfo.IP,
				bossInfo.Idc,
				bossInfo.Platform,
				bossInfo.Contact(),
				time.Unix(eve.EventTime, 0).Format(timeLayout),
				eve.Id,
			).Exec()
		} else {
			sqltemplete = fmt.Sprintf("%v WHERE id = ?", sqltemplete)
			sqlLog, errRes = q.Raw(
				sqltemplete,
				time.Unix(eve.EventTime, 0).Format(timeLayout),
				eve.MaxStep(),
				eve.CurrentStep,
				eve.Note(),
				fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
				eve.Status,
				eve.Func(),
				eve.Priority(),
				tpl_creator,
				eve.ExpressionId(),
				eve.StrategyId(),
				eve.TplId(),
				aType.Id,
				bossInfo.IP,
				bossInfo.Idc,
				bossInfo.Platform,
				bossInfo.Contact(),
				eve.Id,
			).Exec()
		}
	}
	log.Debug(fmt.Sprintf("%v, %v", sqlLog, errRes))
	//insert case
	res2, err := insertEvent(q, eve)
	log.Debug(fmt.Sprintf("%v, %v", res2, err))
	return err
}

func counterGen(metric string, tags string) (mycounter string) {
	mycounter = metric
	if tags != "" {
		mycounter = fmt.Sprintf("%s/%s", metric, tags)
	}
	return
}

func findAlarmTypeByName(name string) (err error, aType AlarmType) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	aType = AlarmType{}
	err = q.Raw("select * from alarm_types where name = ?", name).QueryRow(&aType)
	if err != nil {
		err = fmt.Errorf("[findAlarmTypeByName]: %v - name: %s", err.Error(), name)
	}
	return
}

func insertExternalEvent(q orm.Ormer, exevent ExternalEvent) (res interface{}, err error) {
	var status int
	status = 0
	switch exevent.StatusStr() {
	case "OK":
		status = 1
	case "PROBLEM":
		status = 0
	case "UNKNOWN":
		// if case == unknow mean, this event's foramt is borken. just skip!
		return
	}
	sqltemplete := `INSERT INTO events (
		event_caseId,
		step,
		cond,
		status,
		timestamp
	) VALUES(?,?,?,?,?)`
	res, err = q.Raw(
		sqltemplete,
		exevent.Id(),
		exevent.CurrentStep,
		exevent.TriggerCondition,
		status,
		time.Unix(exevent.EventTime, 0).Format(timeLayout),
	).Exec()
	return
}

func InsertExternalEvent(exevent ExternalEvent) error {
	q := orm.NewOrm()
	err := q.Using("falcon_portal")
	if err != nil {
		log.Errorf("q.Using(falcon_portal) got error: %v", err.Error())
		return err
	}
	var event []EventCases
	var sqlLog sql.Result
	var errRes error
	var aType AlarmType
	errRes, aType = findAlarmTypeByName(exevent.AlarmType)
	if errRes != nil {
		return errRes
	}
	q.Raw("select * from event_cases where id = ? limit 1", exevent.Id()).QueryRows(&event)
	q.Begin()
	if len(event) == 0 {
		//create cases
		sqltemplete := `INSERT INTO event_cases (
					id,
					endpoint,
					metric,
					func,
					cond,
					note,
					max_step,
					current_step,
					priority,
					status,
					timestamp,
					update_at,
					tpl_creator,
					expression_id,
					strategy_id,
					template_id,
					alarm_type_id,
					ip,
					idc,
					platform,
					contact, extended_blob) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		tpl_creator := ""
		sqlLog, errRes = q.Raw(
			sqltemplete,
			exevent.Id(),
			exevent.Target,
			counterGen(exevent.Metric, utils.SortedTags(exevent.PushedTags)),
			exevent.TriggerDescription,
			//cond
			exevent.TriggerCondition,
			exevent.Note,
			exevent.MaxStep(),
			exevent.CurrentStep,
			exevent.Priority,
			exevent.StatusStr(),
			//start_at
			time.Unix(exevent.EventTime, 0).Format(timeLayout),
			//update_at
			time.Unix(exevent.EventTime, 0).Format(timeLayout),
			tpl_creator,
			exevent.ExpressionId(),
			exevent.StrategyId(),
			//template_id
			exevent.TplId(),
			aType.Id,
			exevent.GetKey("ip"),
			exevent.GetKey("idc"),
			exevent.GetKey("platform"),
			exevent.GetKey("contact"),
			exevent.ExtendedBlobStr(),
		).Exec()
	} else {
		sqltemplete := `UPDATE event_cases SET
				update_at = ?,
				max_step = ?,
				current_step = ?,
				note = ?,
				cond = ?,
				status = ?,
				func = ?,
				priority = ?,
				tpl_creator = ?,
				expression_id = ?,
				strategy_id = ?,
				template_id = ?,
				alarm_type_id = ?,
				ip = ?,
				idc = ?,
				platform = ?,
				contact = ?,
				extended_blob = ?`
		//reopen case
		if event[0].ProcessStatus == "resolved" || event[0].ProcessStatus == "ignored" {
			sqltemplete = fmt.Sprintf("%v ,process_status = '%s', process_note = %d", sqltemplete, "unresolved", 0)
		}

		tpl_creator := ""
		if exevent.CurrentStep == 1 {
			//update start time of cases
			sqltemplete = fmt.Sprintf("%v ,timestamp = ? WHERE id = ?", sqltemplete)
			sqlLog, errRes = q.Raw(
				sqltemplete,
				time.Unix(exevent.EventTime, 0).Format(timeLayout),
				exevent.MaxStep(),
				exevent.CurrentStep,
				exevent.Note,
				exevent.TriggerCondition,
				exevent.StatusStr(),
				exevent.TriggerDescription,
				exevent.Priority,
				tpl_creator,
				exevent.ExpressionId(),
				exevent.StrategyId(),
				exevent.TplId(),
				aType.Id,
				exevent.GetKey("ip"),
				exevent.GetKey("idc"),
				exevent.GetKey("platform"),
				exevent.GetKey("contact"),
				exevent.ExtendedBlobStr(),
				time.Unix(exevent.EventTime, 0).Format(timeLayout),
				exevent.Id(),
			).Exec()
		} else {
			sqltemplete = fmt.Sprintf("%v WHERE id = ?", sqltemplete)
			sqlLog, errRes = q.Raw(
				sqltemplete,
				time.Unix(exevent.EventTime, 0).Format(timeLayout),
				exevent.MaxStep(),
				exevent.CurrentStep,
				exevent.Note,
				exevent.TriggerCondition,
				exevent.StatusStr(),
				exevent.TriggerDescription,
				exevent.Priority,
				tpl_creator,
				exevent.ExpressionId(),
				exevent.StrategyId(),
				exevent.TplId(),
				aType.Id,
				exevent.GetKey("ip"),
				exevent.GetKey("idc"),
				exevent.GetKey("platform"),
				exevent.GetKey("contact"),
				exevent.ExtendedBlobStr(),
				exevent.Id(),
			).Exec()
		}
	}
	log.Debug(fmt.Sprintf("%v, %v", sqlLog, errRes))
	if errRes != nil {
		q.Rollback()
		return errRes
	}
	//insert case
	res2, err := insertExternalEvent(q, exevent)
	log.Debug(fmt.Sprintf("%v, %v", res2, err))
	if err != nil {
		q.Rollback()
	}
	q.Commit()
	return err
}
