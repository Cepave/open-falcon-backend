package nqm

import (
	sjson "github.com/bitly/go-simplejson"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
)

// Defines the interface for process a instance of metrics
type MetricFilter interface {
	// Checks whether or not the metrics matches conditions of filter
	IsMatch(model *Metrics) bool
}

/**
 * Macro-struct re-used by various data
 */
type Metrics struct {
	Max                     int16   `json:"max"`
	Min                     int16   `json:"min"`
	Avg                     float64 `json:"avg"`
	Med                     int16   `json:"med"`
	Mdev                    float64 `json:"mdev"`
	Loss                    float64 `json:"loss"`
	Count                   int32   `json:"count"`
	NumberOfSentPackets     uint64  `json:"number_of_sent_packets"`
	NumberOfReceivedPackets uint64  `json:"number_of_received_packets"`
	NumberOfAgents          int32   `json:"number_of_agents"`
	NumberOfTargets         int32   `json:"number_of_targets"`
}
func (m *Metrics) UnmarshalSimpleJson(jsonObject *sjson.Json) {
	jsonExt := ojson.ToJsonExt(jsonObject)

	m.Max = jsonExt.GetExt("max").MustInt16()
	m.Min = jsonExt.GetExt("min").MustInt16()
	m.Avg = jsonExt.Get("avg").MustFloat64()
	m.Med = jsonExt.GetExt("med").MustInt16()
	m.Mdev = jsonExt.Get("mdev").MustFloat64()
	m.Loss = jsonExt.Get("loss").MustFloat64()
	m.Count = jsonExt.GetExt("count").MustInt32()
	m.NumberOfSentPackets = jsonExt.Get("number_of_sent_packets").MustUint64()
	m.NumberOfReceivedPackets = jsonExt.Get("number_of_received_packets").MustUint64()
	m.NumberOfAgents = jsonExt.GetExt("number_of_agents").MustInt32()
	m.NumberOfTargets = jsonExt.GetExt("number_of_targets").MustInt32()
}
