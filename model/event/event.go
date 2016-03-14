package event

import "time"

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

type Tpl struct {
	Id         int    `json:"id" orm:"pk"`
	TplName    string `json:"tpl_name"`
	ParentId   string `json:"parent_id "`
	ActionId   string `json:"action_id"`
	CreateUser string `json:"create_user"`
	CreateAt   string `json:"create_at"`
}
