package funcs

import (
	"github.com/Cepave/open-falcon-backend/common/model"
)

func AgentMetrics() []*model.MetricValue {
	return []*model.MetricValue{GaugeValue("agent.alive", 1)}
}
