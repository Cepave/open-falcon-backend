package event

import "time"

type EventCases struct {
	// uniuq
	Id       string `json:"id" orm:"pk"`
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Func     string `json:"func"`
	Cond     string `json:"cond"`
	Note     string `json:"note"`
	//leftValue + operator + rightValue
	MaxStep      int       `json:"max_step"`
	CurrentStep  int       `json:"current_step"`
	Priority     int       `json:"priority"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"start_at"`
	UpdateAt     time.Time `json:"update_at"`
	ClosedAt     time.Time `json:"closed_at"`
	ClosedNote   string    `json:"closed_note"`
	UserModified int       `json:"user_modified"`
	TplCreator   string    `json:"tpl_creator"`
	ExpressionId int       `json:"expression_id"`
	StrategyId   int       `json:"strategy_id"`
	TemplateId   int       `json:"template_id"`
	Events       []*Events `json:"evevnts" orm:"reverse(many)"`
}

type Events struct {
	Id          int         `json:"id" orm:"pk"`
	Step        int         `json:"step"`
	Cond        string      `json:"cond"`
	Timestamp   time.Time   `json:"timestamp"`
	EventCaseId *EventCases `json:"event_caseId" orm:"rel(fk)"`
}

type EventsRsp struct {
	Id         int       `json:"id"`
	Step       int       `json:"step"`
	Cond       string    `json:"cond"`
	Timestamp  time.Time `json:"timestamp"`
	Eid        string    `json:"event_caseId" orm:"eid"`
	TplCreator string    `json:"tpl_creator"`
	Metric     string    `json:"metric"`
	Endpoint   string    `json:"endpoint"`
}

type Tpl struct {
	Id         int    `json:"id" orm:"pk"`
	TplName    string `json:"tpl_name"`
	ParentId   string `json:"parent_id "`
	ActionId   string `json:"action_id"`
	CreateUser string `json:"create_user"`
	CreateAt   string `json:"create_at"`
}
