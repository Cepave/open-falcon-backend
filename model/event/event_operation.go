package event

import (
	"fmt"
	"time"

	"github.com/Cepave/alarm/logger"
	coommonModel "github.com/Cepave/common/model"
	"github.com/Cepave/common/utils"
	"github.com/astaxie/beego/orm"
)

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
		time.Unix(eve.EventTime, 0),
	).Exec()
	return
}

func InsertEvent(eve *coommonModel.Event) {
	log := logger.Logger()
	q := orm.NewOrm()
	q.Using("falcon_portal")
	var event []EventCases
	q.Raw("select * from event_cases where id = ?", eve.Id).QueryRows(&event)
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
					template_id
					) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		res1, err := q.Raw(
			sqltemplete,
			eve.Id,
			eve.Endpoint,
			counterGen(eve.Metric(), utils.SortedTags(eve.PushedTags)),
			eve.Func(),
			//cond
			fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			eve.Strategy.Note,
			eve.MaxStep(),
			eve.CurrentStep,
			eve.Priority(),
			eve.Status,
			//start_at
			time.Unix(eve.EventTime, 0),
			//update_at
			time.Unix(eve.EventTime, 0),
			eve.Strategy.Tpl.Creator,
			eve.ExpressionId(),
			eve.StrategyId(),
			//template_id
			eve.TplId()).Exec()
		log.Debug(fmt.Sprintf("%v, %v", res1, err))

		//insert case
		res2, err := insertEvent(q, eve)
		log.Debug(fmt.Sprintf("%v, %v", res2, err))

	} else {
		sqltemplete := `UPDATE event_cases SET
				update_at = ?,
				max_step = ?,
				current_step = ?,
				note = ?,
				cond = ?,
				status = ?`
		//reopen case
		if event[0].ProcessStatus == "resolved" || event[0].ProcessStatus == "ignored" {
			sqltemplete = fmt.Sprintf("%v ,process_status = '%s', process_note = %d", sqltemplete, "unresolved", 0)
		}
		if eve.CurrentStep == 1 {

			sqltemplete = fmt.Sprintf("%v ,timestamp = ? WHERE id = ?", sqltemplete)
			//update start time of cases
			res1, err := q.Raw(
				sqltemplete,
				time.Unix(eve.EventTime, 0),
				eve.MaxStep(),
				eve.CurrentStep,
				eve.Strategy.Note,
				fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
				eve.Status,
				time.Unix(eve.EventTime, 0),
				eve.Id).Exec()
			log.Debug(fmt.Sprintf("%v, %v", res1, err))
			//insert case
			res2, err := insertEvent(q, eve)
			log.Debug(fmt.Sprintf("%v, %v", res2, err))
		} else {
			sqltemplete = fmt.Sprintf("%v WHERE id = ?", sqltemplete)
			res1, err := q.Raw(
				sqltemplete,
				time.Unix(eve.EventTime, 0),
				eve.MaxStep(),
				eve.CurrentStep,
				eve.Strategy.Note,
				fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
				eve.Status,
				eve.Id).Exec()
			log.Debug(fmt.Sprintf("%v, %v", res1, err))
			//insert case
			res2, err := insertEvent(q, eve)
			log.Debug(fmt.Sprintf("%v, %v", res2, err))
		}
	}
}

func counterGen(metric string, tags string) (mycounter string) {
	mycounter = metric
	if tags != "" {
		mycounter = fmt.Sprintf("%s/%s", metric, tags)
	}
	return
}
