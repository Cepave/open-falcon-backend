package metric_parser

import (
	//"github.com/Cepave/open-falcon-backend/common/utils"
	//"log"
)

type metricType byte

const (
	MetricMax metricType = 1
	MetricMin metricType = 2
	MetricAvg metricType = 3
	MetricMed metricType = 4
	MetricMdev metricType = 5
	MetricLoss metricType = 6
	MetricCount metricType = 7
	MetricPckSent metricType = 8
	MetricPckReceived metricType = 9
	MetricNumAgent metricType = 10
	MetricNumTarget metricType = 11
)

var mapOfMetric = map[string]metricType {
	"max": MetricMax,
	"min": MetricMin,
	"avg": MetricAvg,
	"med": MetricMed,
	"mdev": MetricMdev,
	"loss": MetricLoss,
	"count": MetricCount,
	"pck_sent": MetricPckSent,
	"pck_received": MetricPckReceived,
	"num_agent": MetricNumAgent,
	"num_target": MetricNumTarget,
}
