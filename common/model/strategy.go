package model

import (
	"fmt"
	"strings"

	"github.com/Cepave/open-falcon-backend/common/utils"
)

type Strategy struct {
	Id         int               `json:"id"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	Func       string            `json:"func"`       // e.g. max(#3) all(#3)
	Operator   string            `json:"operator"`   // e.g. < !=
	RightValue float64           `json:"rightValue"` // critical value
	MaxStep    int               `json:"maxStep"`
	Priority   int               `json:"priority"`
	Note       string            `json:"note"`
	Tpl        *Template         `json:"tpl"`
}

func (this *Strategy) String() string {
	return fmt.Sprintf(
		"<Id:%d, Metric:%s, Tags:%v, %s%s%s MaxStep:%d, P%d, %s, %v>",
		this.Id,
		this.Metric,
		this.Tags,
		this.Func,
		this.Operator,
		utils.ReadableFloat(this.RightValue),
		this.MaxStep,
		this.Priority,
		this.Note,
		this.Tpl,
	)
}

type HostStrategy struct {
	Hostname   string     `json:"hostname"`
	Strategies []Strategy `json:"strategies"`
}

func (this *HostStrategy) String() string {
	return fmt.Sprintf(
		"<Hostname:%v, Strategies:%v>",
		this.Hostname,
		this.Strategies,
	)
}

type StrategiesResponse struct {
	HostStrategies []*HostStrategy `json:"hostStrategies"`
}

type NewStrategy struct {
	ID         int    `json:"id"`
	Metric     string `json:"metric"`
	Tags       map[string]string
	TagsStr    string       `json:"tags"`
	Func       string       `json:"func"`               // e.g. max(#3) all(#3)
	Operator   string       `json:"operator"`           // e.g. < !=
	RightValue float64      `json:"right_value,string"` // critical value
	MaxStep    int          `json:"max_step"`
	Priority   int          `json:"priority"`
	Note       string       `json:"note"`
	Tpl        *NewTemplate `json:"tpl"`
}

func (s *NewStrategy) AfterLoad() {
	if s.Tags == nil {
		s.Tags = make(map[string]string)
	}
	if s.TagsStr == "" {
		return
	}
	kvStrs := strings.Split(s.TagsStr, ",")
	for _, kv := range kvStrs {
		kvArr := strings.Split(kv, "=")
		if len(kvArr) != 2 {
			continue
		}
		s.Tags[kvArr[0]] = kvArr[1]
	}
}

func (this *NewStrategy) String() string {
	return fmt.Sprintf(
		"<ID:%d, Metric:%s, Tags:%v, %s%s%s MaxStep:%d, P%d, %s, %v>",
		this.ID,
		this.Metric,
		this.Tags,
		this.Func,
		this.Operator,
		utils.ReadableFloat(this.RightValue),
		this.MaxStep,
		this.Priority,
		this.Note,
		this.Tpl,
	)
}

type NewHostStrategy struct {
	Hostname   string         `json:"hostname"`
	Strategies []*NewStrategy `json:"strategies"`
}
