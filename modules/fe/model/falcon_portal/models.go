package falconPortal

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
	MaxStep       int       `json:"max_step"`
	CurrentStep   int       `json:"current_step"`
	Priority      int       `json:"priority"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"start_at"`
	UpdateAt      time.Time `json:"update_at"`
	ProcessNote   int       `json:"process_note"`
	ProcessStatus string    `json:"process_status"`
	TplCreator    string    `json:"tpl_creator"`
	ExpressionId  int       `json:"expression_id"`
	StrategyId    int       `json:"strategy_id"`
	TemplateId    int       `json:"template_id"`
	Events        []*Events `json:"evevnts" orm:"reverse(many)"`
}

type Events struct {
	Id          int         `json:"id" orm:"pk"`
	Step        int         `json:"step"`
	Cond        string      `json:"cond"`
	Status      int         `json:"status"`
	Timestamp   time.Time   `json:"timestamp"`
	EventCaseId *EventCases `json:"event_caseId" orm:"rel(fk)"`
}

type EventsRsp struct {
	Id         int       `json:"id"`
	Step       int       `json:"step"`
	Cond       string    `json:"cond"`
	Timestamp  time.Time `json:"timestamp"`
	Eid        string    `json:"event_caseId" orm:"column(eid)"`
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

type Action struct {
	Id                 int    `json:"id"`
	Uic                string `json:"uic"`
	Url                string `json:"url"`
	Callback           int    `json:"callback"`
	BeforeCallbackSms  int    `json:"before_callback_sms"`
	BeforeCallbackMail int    `json:"before_callback_mail"`
	AfterCallbackSms   int    `json:"after_callback_sms"`
	AfterCallbackMail  int    `json:"after_callback_mail"`
}

type EventNote struct {
	Id          int       `json:"id" orm:"pk"`
	EventCaseId string    `json:"event_caseId" orm:"column(event_caseId)"`
	Note        string    `json:"note"`
	CaseId      string    `json:"case_id"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	UserId      int64     `json:"-"`
	UserName    string    `json:"user_name" orm:"user_name"`
}

type Host struct {
	Id            int    `json:"id" orm:"pk"`
	Hostname      string `json:"hostname" orm:"hostname"`
	Ip            string
	AgentVersion  string
	PluginVersion string
	MaintainBegin int64
	MaintainEnd   int64
	UpdateAt      time.Time
}
