package event

import (
	"fmt"
	coommonModel "github.com/Cepave/common/model"
	"github.com/Cepave/common/utils"
	"github.com/astaxie/beego/orm"
	"log"
	"time"
)

type Event struct {
	// uniuq
	Id       string `json:"id" orm:"pk"`
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Func     string `json:"func"`
	//leftValue + operator + rightValue
	Cond         string    `json:"cond"`
	Note         string    `json:"note"`
	MaxStep      int       `json:"max_step"`
	CurrentStep  int       `json:"current_step"`
	Priority     int       `json:"priority"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
	UpdateAt     time.Time `json:"update_at"`
	ClosedAt     time.Time `json:"closed_at"`
	UserModified int       `json:"user_modified"`
	ExpressionId int       `json:"expression_id"`
	StrategyId   int       `json:"strategy_id"`
	TemplateId   int       `json:"template_id"`
}

func InsertEvent(eve *coommonModel.Event) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	var event []Event
	q.Raw("select * from event where id = ?", eve.Id).QueryRows(&event)
	if len(event) == 0 {

		sqltemplete := `INSERT INTO event (
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
					expression_id,
					strategy_id,
					template_id
					) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		res, err := q.Raw(
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
			//timestamp
			time.Unix(eve.EventTime, 0),
			//update_at
			time.Unix(eve.EventTime, 0),
			eve.ExpressionId(),
			eve.StrategyId(),
			//template_id
			eve.TplId()).Exec()
		log.Printf("%v, %v", res, err)

	} else {

		res, err := q.Raw(
			"UPDATE event SET update_at = ?, current_step = ?, cond = ?  WHERE id = ?",
			time.Unix(eve.EventTime, 0),
			eve.CurrentStep,
			fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			eve.Id).Exec()
		log.Printf("%v, %v", res, err)

	}
}

func counterGen(metric string, tags string) (mycounter string) {
	mycounter = metric
	if tags != "" {
		mycounter = fmt.Sprintf("%s/%s", metric, tags)
	}
	return
}
