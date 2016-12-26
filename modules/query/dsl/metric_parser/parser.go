package metric_parser

import (
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
)

func ParseToMetricFilter(dsl string) (model.MetricFilter, error) {
	if dsl == "" {
		return nil, nil
	}

	filter, err := Parse("Metric", []byte(dsl))
	if err != nil{
		return nil, err
	}

	return filter.(model.MetricFilter), nil
}
