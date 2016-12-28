package nqm

import (
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	sjson "github.com/bitly/go-simplejson"
)

type DynamicRecord struct {
	Agent *DynamicAgentProps `json:"agent,omitempty"`
	Target *DynamicTargetProps `json:"target,omitempty"`
	Metrics *DynamicMetrics `json:"metrics"`
}

type DynamicAgentProps struct {
	Id int32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Isp *owlModel.Isp `json:"isp,omitempty"`
	Province *owlModel.Province `json:"province,omitempty"`
	City *owlModel.City2 `json:"city,omitempty"`
	NameTag *owlModel.NameTag `json:"name_tag,omitempty"`
}
type DynamicTargetProps struct {
	Id int32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Host string `json:"host,omitempty"`
	Isp *owlModel.Isp `json:"isp,omitempty"`
	Province *owlModel.Province `json:"province,omitempty"`
	City *owlModel.City2 `json:"city,omitempty"`
	NameTag *owlModel.NameTag `json:"name_tag,omitempty"`
}

type DynamicMetrics struct {
	Metrics *Metrics
	Output *[]string
}

func (m *DynamicMetrics) MarshalJSON() ([]byte, error) {
	jsonObj := sjson.New()

	metricsHolder := m.Metrics
	for _, column := range *m.Output {
		switch column {
		case MetricMax:
			jsonObj.Set("max", metricsHolder.Max)
		case MetricMin:
			jsonObj.Set("min", metricsHolder.Min)
		case MetricAvg:
			jsonObj.Set("avg", metricsHolder.Avg)
		case MetricMed:
			jsonObj.Set("med", metricsHolder.Med)
		case MetricMdev:
			jsonObj.Set("mdev", metricsHolder.Mdev)
		case MetricLoss:
			jsonObj.Set("loss", metricsHolder.Loss)
		case MetricCount:
			jsonObj.Set("count", metricsHolder.Count)
		case MetricPckSent:
			jsonObj.Set("pck_sent", metricsHolder.NumberOfSentPackets)
		case MetricPckReceived:
			jsonObj.Set("pck_received", metricsHolder.NumberOfReceivedPackets)
		case MetricNumAgent:
			jsonObj.Set("num_agent", metricsHolder.NumberOfAgents)
		case MetricNumTarget:
			jsonObj.Set("num_target", metricsHolder.NumberOfTargets)
		}
	}

	return jsonObj.MarshalJSON()
}
