package service

import (
	"strconv"
	"time"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	cutils "github.com/Cepave/open-falcon-backend/common/utils"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
)

// process new metric values
func RecvMetricValues(sourceMetrics []*cmodel.MetricValue, reply *cmodel.TransferResponse, from string) error {
	startTime := time.Now()

	validCount := int64(0)
	invalidCount := 0

	relayStation := DefaultRelayStationFactory.Build()

	for _, sourceMetric := range sourceMetrics {
		refinedValue, failedCount := checkAndRefineMetric(sourceMetric, startTime)
		if failedCount > 0 {
			invalidCount++
			continue
		}

		relayStation.Dispatch(refinedValue)
		validCount++
	}

	/**
	 * Updates the counter of statistics
	 */
	proc.RecvCnt.IncrBy(validCount)
	if from == "rpc" {
		proc.RpcRecvCnt.IncrBy(validCount)
	} else if from == "http" {
		proc.HttpRecvCnt.IncrBy(validCount)
	}
	// :~)

	/**
	 * Relays the metrics to corresponding queue
	 */
	relayStation.RelayToQueue()
	// :~)

	reply.Message = "ok"
	reply.Invalid = invalidCount
	reply.Total = len(sourceMetrics)
	reply.Latency = (time.Now().UnixNano() - startTime.UnixNano()) / 1000000

	return nil
}

func checkAndRefineMetric(metric *cmodel.MetricValue, startTime time.Time) (*cmodel.MetaData, int) {
	if metric == nil {
		return nil, 1
	}

	/**
	 * 历史遗留问题.
	 * 老版本 agent 上报的 metric=kernel.hostname 的数据,其取值为 string 类型,现在已经不支持了;
	 * 所以,这里硬编码过滤掉
	 */
	if metric.Metric == "kernel.hostname" {
		return nil, 1
	}
	// :~)

	if metric.Metric == "" || metric.Endpoint == "" {
		return nil, 1
	}

	if metric.Type != g.COUNTER && metric.Type != g.GAUGE && metric.Type != g.DERIVE {
		return nil, 1
	}

	if metric.Value == "" {
		return nil, 1
	}

	if metric.Step <= 0 {
		return nil, 1
	}

	if len(metric.Metric)+len(metric.Tags) > 510 {
		return nil, 1
	}

	refinedValue, success := refineValue(metric.Value)
	if !success {
		return nil, 1
	}

	/**
	 * If the timestamp of metric is too early(< 0) or comes from future(2 hours),
	 * following code modify it to the start time of processing this batch of metrics.
	 */
	refinedTimestamp := metric.Timestamp
	startTimeAsUnix := startTime.Unix()
	if metric.Timestamp <= 0 || metric.Timestamp > startTimeAsUnix*1+7200 {
		refinedTimestamp = startTimeAsUnix
	}
	// :~)

	refinedMetric := &cmodel.MetaData{
		Metric:       metric.Metric,
		Endpoint:     metric.Endpoint,
		Timestamp:    refinedTimestamp,
		Step:         metric.Step,
		CounterType:  metric.Type,
		Tags:         cutils.DictedTagstring(metric.Tags), //TODO: tags键值对的个数,要做一下限制
		SourceMetric: metric,
	}

	refinedMetric.Value = refinedValue

	return refinedMetric, 0
}

func refineValue(value interface{}) (float64, bool) {
	switch typedValue := value.(type) {
	case string:
		finalValue, err := strconv.ParseFloat(typedValue, 64)
		if err == nil {
			return finalValue, true
		}
	/**
	 * Improve speed for first version of ordinary code
	 */
	case float64:
		return typedValue, true
	case int64:
		return float64(typedValue), true
	// :~)
	case int:
		return float64(typedValue), true
	case int8:
		return float64(typedValue), true
	case int16:
		return float64(typedValue), true
	case int32:
		return float64(typedValue), true
	}

	return 0, false
}
