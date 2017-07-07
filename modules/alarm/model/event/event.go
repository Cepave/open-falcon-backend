package event

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

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
	// add on 2017/07
	AlarmTypeId  int    `json:"alarm_type_id" orm:"column(alarm_type_id)"`
	IP           string `json:"ip" orm:"column(ip)"`
	IDC          string `json:"idc" orm:"column(idc)"`
	Platform     string `json:"platform" orm:"column(platform)"`
	Contact      string `json:"contact" orm:"column(contact)"`
	ExtendedBLOB string `json:"extended_blob" orm:"column(extended_blob)"`
}

type Events struct {
	Id          int         `json:"id" orm:"pk"`
	Step        int         `json:"step"`
	Cond        string      `json:"cond"`
	Status      int         `json:"status"`
	Timestamp   time.Time   `json:"timestamp"`
	EventCaseId *EventCases `json:"event_caseId" orm:"rel(fk)"`
}

type AlarmType struct {
	Id           int       `json:"id" orm:"pk"`
	Name         string    `json:"name" orm:"column(name)"`
	InternalData int       `json:"internal_data" orm:"column(internal_data)"`
	Description  string    `json:"description" orm:"column(description)"`
	Created      time.Time `json:"created" orm:"column(created)"`
}

type ExternalEvent struct {
	AlarmType          string            `json:"alarm_type"`
	Status             int               `json:"status"`
	Target             string            `json:"target"`
	Metric             string            `json:"metric"`
	CurrentStep        int               `json:"current_step"`
	EventTime          int64             `json:"event_time"`
	Priority           int               `json:"priority"`
	TriggerId          int               `json:"trigger_id"`
	TriggerDescription string            `json:"trigger_description"`
	TriggerCondition   string            `json:"trigger_condition"`
	Note               string            `json:"note"`
	PushedTags         map[string]string `json:"pushed_tags"`
	ExtendedBlob       map[string]string `json:"extended_blob"`
}

func (mine *ExternalEvent) CheckFormating() error {
	var errorStrings []string
	if mine.Status != 0 && mine.Status != 1 {
		errorStrings = append(errorStrings, "status not vaild")
	}
	if mine.Target == "" {
		errorStrings = append(errorStrings, "target is empty")
	}
	if mine.Metric == "" {
		errorStrings = append(errorStrings, "metric is empty")
	}
	if mine.CurrentStep <= 0 {
		errorStrings = append(errorStrings, "current step not vaild")
	}
	if mine.EventTime <= 0 {
		errorStrings = append(errorStrings, "event time can not be empty")
	}
	if mine.Priority < 0 || mine.Priority > 6 {
		errorStrings = append(errorStrings, "priority is not vaild")
	}
	if mine.TriggerId <= 0 {
		errorStrings = append(errorStrings, "trigger id can not be empty")
	}
	if mine.TriggerCondition == "" {
		errorStrings = append(errorStrings, "trigger condistion can not be empty")
	}
	if len(errorStrings) == 0 {
		return nil
	} else {
		eTmp := fmt.Sprintf("[[\"%v\"]", errorStrings[0])
		for n := 1; n < len(errorStrings); n++ {
			eTmp = fmt.Sprintf("%v, [\"%v\"]", eTmp, errorStrings[n])
		}
		eTmp = fmt.Sprintf("%v]", eTmp)
		return errors.New(eTmp)
	}
}

func (mine *ExternalEvent) Id() string {
	hasher := sha1.New()
	bv := fmt.Sprintf("%v_%v_%v_%v", mine.Target, mine.Metric, mine.PushedTags, mine.TriggerId)
	hasher.Write([]byte(bv))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	id := fmt.Sprintf("%v_%v", mine.AlarmType, sha)
	return id
}

func (mine *ExternalEvent) StatusStr() string {
	switch mine.Status {
	case 0:
		return "PROBLEM"
	case 1:
		return "OK"
	default:
		return "UNKNOWN"
	}
}

// set 99 is defult display value.
func (mine *ExternalEvent) MaxStep() int {
	return 99
}

func (mine *ExternalEvent) GetKey(key string) string {
	if val, ok := mine.ExtendedBlob[key]; ok {
		return val
	} else {
		return ""
	}
}

func (mine *ExternalEvent) ExpressionId() int {
	return 0
}

func (mine *ExternalEvent) StrategyId() int {
	return 0
}

func (mine *ExternalEvent) TplId() int {
	return 0
}

func (mine *ExternalEvent) ExtendedBlobStr() string {
	blobStr, err := json.Marshal(mine.ExtendedBlob)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(blobStr)
}

func (mine ExternalEvent) ForceFixStepWhenStatusOk() ExternalEvent {
	if mine.StatusStr() == "OK" {
		mine.CurrentStep = 1
	}
	return mine
}
