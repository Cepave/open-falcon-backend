package nqm

import (
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	sjson "github.com/bitly/go-simplejson"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"net"
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

// Defines the function used for utils.Comparison of two dynamic records
type CompareDynamicRecord func(*DynamicRecord, *DynamicRecord, byte) int

const (
	Larger = 1
	Equal = 0
	Lesser = -1
)

var CompareFunctions = map[string]CompareDynamicRecord {
	"agent_name": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Agent, right.Agent, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Agent.Name, right.Agent.Name, direction)
	},
	"agent_hostname": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Agent, right.Agent, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Agent.Hostname, right.Agent.Hostname, direction)
	},
	"agent_ip_address": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Agent, right.Agent, direction); hasNil {
			return r
		}

		leftIp, rightIp := net.ParseIP(left.Agent.IpAddress), net.ParseIP(right.Agent.IpAddress)
		return utils.CompareIpAddress(leftIp, rightIp, direction)
	},
	"agent_isp": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Agent, right.Agent, direction); hasNil {
			return r
		}
		if r, hasNil := utils.CompareNil(left.Agent.Isp, right.Agent.Isp, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Agent.Isp.Name, right.Agent.Isp.Name, direction)
	},
	"agent_province": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Agent, right.Agent, direction); hasNil {
			return r
		}
		if r, hasNil := utils.CompareNil(left.Agent.Province, right.Agent.Province, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Agent.Province.Name, right.Agent.Province.Name, direction)
	},
	"agent_city": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Agent, right.Agent, direction); hasNil {
			return r
		}
		if r, hasNil := utils.CompareNil(left.Agent.City, right.Agent.City, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Agent.City.Name, right.Agent.City.Name, direction)
	},
	"agent_name_tag": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Agent, right.Agent, direction); hasNil {
			return r
		}
		if r, hasNil := utils.CompareNil(left.Agent.NameTag, right.Agent.NameTag, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Agent.NameTag.Value, right.Agent.NameTag.Value, direction)
	},

	"target_name": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Target, right.Target, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Target.Name, right.Target.Name, direction)
	},
	"target_host": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Target, right.Target, direction); hasNil {
			return r
		}

		leftIfIp, rightIfIp := net.ParseIP(left.Target.Host), net.ParseIP(right.Target.Host)
		if leftIfIp != nil && rightIfIp != nil {
			return utils.CompareIpAddress(leftIfIp, rightIfIp, direction)
		}

		return utils.CompareString(left.Target.Host, right.Target.Host, direction)
	},
	"target_isp": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Target, right.Target, direction); hasNil {
			return r
		}
		if r, hasNil := utils.CompareNil(left.Target.Isp, right.Target.Isp, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Target.Isp.Name, right.Target.Isp.Name, direction)
	},
	"target_province": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Target, right.Target, direction); hasNil {
			return r
		}
		if r, hasNil := utils.CompareNil(left.Target.Province, right.Target.Province, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Target.Province.Name, right.Target.Province.Name, direction)
	},
	"target_city": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Target, right.Target, direction); hasNil {
			return r
		}
		if r, hasNil := utils.CompareNil(left.Target.City, right.Target.City, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Target.City.Name, right.Target.City.Name, direction)
	},
	"target_name_tag": func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		if r, hasNil := utils.CompareNil(left.Target, right.Target, direction); hasNil {
			return r
		}
		if r, hasNil := utils.CompareNil(left.Target.NameTag, right.Target.NameTag, direction); hasNil {
			return r
		}

		return utils.CompareString(left.Target.NameTag.Value, right.Target.NameTag.Value, direction)
	},

	MetricMax: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareInt(
			int64(left.Metrics.Metrics.Max),
			int64(right.Metrics.Metrics.Max),
			direction,
		)
	},
	MetricMin: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareInt(
			int64(left.Metrics.Metrics.Min),
			int64(right.Metrics.Metrics.Min),
			direction,
		)
	},
	MetricAvg: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareFloat(
			left.Metrics.Metrics.Avg,
			right.Metrics.Metrics.Avg,
			direction,
		)
	},
	MetricMed: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareInt(
			int64(left.Metrics.Metrics.Med),
			int64(right.Metrics.Metrics.Med),
			direction,
		)
	},
	MetricMdev: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareFloat(
			left.Metrics.Metrics.Mdev,
			right.Metrics.Metrics.Mdev,
			direction,
		)
	},
	MetricLoss: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareFloat(
			left.Metrics.Metrics.Loss,
			right.Metrics.Metrics.Loss,
			direction,
		)
	},
	MetricCount: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareInt(
			int64(left.Metrics.Metrics.Count),
			int64(right.Metrics.Metrics.Count),
			direction,
		)
	},
	MetricPckSent: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareUint(
			left.Metrics.Metrics.NumberOfSentPackets,
			right.Metrics.Metrics.NumberOfSentPackets,
			direction,
		)
	},
	MetricPckReceived: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareUint(
			left.Metrics.Metrics.NumberOfReceivedPackets,
			right.Metrics.Metrics.NumberOfReceivedPackets,
			direction,
		)
	},
	MetricNumAgent: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareInt(
			int64(left.Metrics.Metrics.NumberOfAgents),
			int64(right.Metrics.Metrics.NumberOfAgents),
			direction,
		)
	},
	MetricNumTarget: func(left *DynamicRecord, right *DynamicRecord, direction byte) int {
		return utils.CompareInt(
			int64(left.Metrics.Metrics.NumberOfTargets),
			int64(right.Metrics.Metrics.NumberOfTargets),
			direction,
		)
	},
}
