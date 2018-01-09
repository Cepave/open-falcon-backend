package model

import (
	"fmt"
	"time"

	"github.com/Cepave/open-falcon-backend/common/utils"
)

type JudgeItem struct {
	Endpoint        string            `json:"endpoint"`
	Metric          string            `json:"metric"`
	Value           float64           `json:"value"`
	JudgeType       string            `json:"judgeType"`
	Tags            map[string]string `json:"tags"`
	Timestamp       int64             `json:"timestamp"`
	SourceTimestamp int64             `json:"source_timestamp"`
}

func (this *JudgeItem) String() string {
	alignTime := time.Unix(this.Timestamp, 0)
	sourceTime := time.Unix(this.SourceTimestamp, 0)

	return fmt.Sprintf(
		"<Endpoint[%s], Metric[%s], Value[%f], Timestamp(align)[%s], Source Timestamp[%s], JudgeType[%s] Tags[%v]>",
		this.Endpoint, this.Metric, this.Value,
		alignTime, sourceTime,
		this.JudgeType, this.Tags,
	)
}

func (this *JudgeItem) PrimaryKey() string {
	return utils.Md5(utils.PK(this.Endpoint, this.Metric, this.Tags))
}

type HistoryData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}
