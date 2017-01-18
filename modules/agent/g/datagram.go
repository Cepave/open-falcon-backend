package g

import (
	"fmt"
)

type MetricValueExtend struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
	Fields    string      `json:"fields"`
}

func (this *MetricValueExtend) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Type:%s, Tags:%s, Fields:%s, Step:%d, Time:%d, Value:%v>",
		this.Endpoint,
		this.Metric,
		this.Type,
		this.Tags,
		this.Fields,
		this.Step,
		this.Timestamp,
		this.Value,
	)
}
