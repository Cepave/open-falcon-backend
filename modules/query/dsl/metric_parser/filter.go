package metric_parser

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/modules/query/nqm"
	"reflect"
	"strconv"
)

// Defines the interface for process a instance of metrics
type MetricFilter interface {
	// Checks whether or not the metrics matches conditions of filter
	IsMatch(model *nqm.Metrics) bool
}

type boolFilterImpl struct {
	expectedResult bool
	matchers []MetricFilter
}
func (boolFilter *boolFilterImpl) IsMatch(model *nqm.Metrics) bool {
	for _, filter := range boolFilter.matchers {
		if filter.IsMatch(model) == boolFilter.expectedResult {
			return boolFilter.expectedResult
		}
	}

	return !boolFilter.expectedResult
}
func newBoolFilterImpl(checkResult bool, leftTerm interface{}, restTerms interface{}) *boolFilterImpl {
	matchers := make([]MetricFilter, 0)

	matchers = append(matchers, leftTerm.(MetricFilter))
	for _, restFilter := range restTerms.([]interface{}) {
		matchers = append(matchers, restFilter.(MetricFilter))
	}

	return &boolFilterImpl {
		checkResult, matchers,
	}
}

type filterImpl struct {
	leftFactor interface{}
	op string
	rightFactor interface{}
}
func newFilterImpl(leftFactor interface{}, op string, rightFactor interface{}) *filterImpl {
	valueOfLeftFactor := getFactorValue(leftFactor)
	valueOfRightFactor := getFactorValue(rightFactor)

	return &filterImpl {
		valueOfLeftFactor,
		op,
		valueOfRightFactor,
	}
}

func (f *filterImpl) IsMatch(model *nqm.Metrics) bool {
	leftValue := f.getValue(f.leftFactor, model)
	rightValue := f.getValue(f.rightFactor, model)

	switch f.op {
	case ">":
		return leftValue > rightValue
	case "<":
		return leftValue < rightValue
	case "==":
		return leftValue == rightValue
	case ">=":
		return leftValue >= rightValue
	case "<=":
		return leftValue <= rightValue
	case "!=":
		return leftValue != rightValue
	}

	panic(fmt.Errorf("Unsupported operator: [%s]", f.op))
}
func (f *filterImpl) getValue(factor interface{}, model *nqm.Metrics) float64 {
	switch factor.(type) {
	case metricType:
		switch factor.(metricType) {
		case MetricMax:
			return float64(model.Max)
		case MetricMin:
			return float64(model.Min)
		case MetricAvg:
			return float64(model.Avg)
		case MetricMed:
			return float64(model.Med)
		case MetricMdev:
			return float64(model.Mdev)
		case MetricLoss:
			return float64(model.Loss)
		case MetricCount:
			return float64(model.Count)
		case MetricPckSent:
			return float64(model.NumberOfSentPackets)
		case MetricPckReceived:
			return float64(model.NumberOfReceivedPackets)
		case MetricNumAgent:
			return float64(model.NumberOfAgents)
		case MetricNumTarget:
			return float64(model.NumberOfTargets)
		}
	case float64:
		return factor.(float64)
	}

	panic(fmt.Errorf("Unknown type of factor: %v", factor))
}
func (f *filterImpl) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		getFactorString(f.leftFactor),
		f.op,
		getFactorString(f.rightFactor),
	)
}

func getFactorValue(v interface{}) interface{} {
	switch v.(type) {
	case string:
		stringValue := v.(string)
		floatValue, e := strconv.ParseFloat(stringValue, 64)
		if e != nil {
			panic(e)
		}

		return floatValue
	case metricType:
		return v
	}

	panic(fmt.Errorf("Unknown type of factor: [%s]", reflect.TypeOf(v)))
}

func getFactorString(v interface{}) string {
	switch v.(type) {
	case metricType:
		return fmt.Sprintf("Metric(%d)", v)
	case float64:
		return fmt.Sprintf("Float64(%d)", v)
	}

	return fmt.Sprintf("Unknown type[%s](%d)", reflect.TypeOf(v).Name(), v)
}
