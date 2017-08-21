package model

import (
	"fmt"

	"github.com/Cepave/open-falcon-backend/common/utils"
)

type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (this *Config) String() string {
	return fmt.Sprintf(
		"<Key:%s, Value:%s>",
		this.Key,
		this.Value,
	)
}

type Expression struct {
	Id         int               `json:"id"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	Func       string            `json:"func"`       // e.g. max(#3) all(#3)
	Operator   string            `json:"operator"`   // e.g. < !=
	RightValue float64           `json:"rightValue"` // critical value
	MaxStep    int               `json:"maxStep"`
	Priority   int               `json:"priority"`
	Note       string            `json:"note"`
	ActionId   int               `json:"actionId"`
}

func (this *Expression) String() string {
	return fmt.Sprintf(
		"<Id:%d, Metric:%s, Tags:%v, %s%s%s MaxStep:%d, P%d %s ActionId:%d>",
		this.Id,
		this.Metric,
		this.Tags,
		this.Func,
		this.Operator,
		utils.ReadableFloat(this.RightValue),
		this.MaxStep,
		this.Priority,
		this.Note,
		this.ActionId,
	)
}

type ExpressionResponse struct {
	Expressions []*Expression `json:"expressions"`
}

type NewExpression struct {
	ID         int               `json:"id"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	Func       string            `json:"func"`        // e.g. max(#3) all(#3)
	Operator   string            `json:"operator"`    // e.g. < !=
	RightValue float64           `json:"right_value"` // critical value
	MaxStep    int               `json:"max_step"`
	Priority   int               `json:"priority"`
	Note       string            `json:"note"`
	ActionID   int               `json:"action_id"`
}

func (this NewExpression) String() string {
	return fmt.Sprintf(
		"<ID:%d, Metric:%s, Tags:%v, %s%s%s MaxStep:%d, P%d %s ActionID:%d>",
		this.ID,
		this.Metric,
		this.Tags,
		this.Func,
		this.Operator,
		utils.ReadableFloat(this.RightValue),
		this.MaxStep,
		this.Priority,
		this.Note,
		this.ActionID,
	)
}
