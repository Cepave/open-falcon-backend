package nqm_parser

import (
	"fmt"
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
	"time"
)

// The parameters for query
type QueryParams struct {
	StartTime        time.Time
	EndTime          time.Time
	AgentFilter      NodeFilter
	TargetFilter     NodeFilter
	AgentFilterById  NodeFilterById
	TargetFilterById NodeFilterById
	IspRelation      model.PropRelation
	ProvinceRelation model.PropRelation
	CityRelation     model.PropRelation
}

// The filter of node
type NodeFilter struct {
	MatchProvinces []string
	MatchCities    []string
	MatchIsps      []string
}

// The filter of node
type NodeFilterById struct {
	MatchIds       []int32
	MatchProvinces []int16
	MatchCities    []int16
	MatchIsps      []int16
}

// Initialize query parameters with default values
func NewQueryParams() *QueryParams {
	queryParams := &QueryParams{}

	queryParams.IspRelation = model.NoCondition
	queryParams.ProvinceRelation = model.NoCondition
	queryParams.CityRelation = model.NoCondition

	return queryParams
}

// Sets-up the parameters
func (p *QueryParams) SetUpParams(paramSetters interface{}) {
	for _, setterImpl := range paramSetters.([]paramSetter) {
		setterImpl(p)
	}
}

/**
 * Checks:
 * 1. The end time must be after or equals the start time
 */
func (p *QueryParams) CheckRationalOfParameters() error {
	if !p.EndTime.After(p.StartTime) {
		return fmt.Errorf(
			"Start time is not valid. Start Time: [%s]. End Time: [%s]",
			p.StartTime.Format(time.RFC3339), p.EndTime.Format(time.RFC3339),
		)
	}

	return nil
}

// Checks the paramters
//
// 1. provinces and cities cannot be assigned at the same time(except auto-condition)
// 2. duplicated value would be eliminated
const FORMAT_ERROR_LOCATION_FILTER = "%v filter for provinces:%v and cities:%v are both set"

func (p *QueryParams) checkParams() (err error) {
	err = nil

	if err = buildErrorIfBothAreSet(
		p.AgentFilter.MatchProvinces, p.AgentFilter.MatchCities,
		FORMAT_ERROR_LOCATION_FILTER, "Agent",
	); err != nil {
		return
	}

	if err = buildErrorIfBothAreSet(
		p.TargetFilter.MatchProvinces, p.TargetFilter.MatchCities,
		FORMAT_ERROR_LOCATION_FILTER, "Target",
	); err != nil {
		return
	}

	p.AgentFilter.MatchProvinces = eliminateDuplicatedValues(p.AgentFilter.MatchProvinces)
	p.AgentFilter.MatchCities = eliminateDuplicatedValues(p.AgentFilter.MatchCities)
	p.AgentFilter.MatchIsps = eliminateDuplicatedValues(p.AgentFilter.MatchIsps)
	p.TargetFilter.MatchProvinces = eliminateDuplicatedValues(p.TargetFilter.MatchProvinces)
	p.TargetFilter.MatchCities = eliminateDuplicatedValues(p.TargetFilter.MatchCities)
	p.TargetFilter.MatchIsps = eliminateDuplicatedValues(p.TargetFilter.MatchIsps)

	return
}

func buildErrorIfBothAreSet(leftValues, rightValues []string, format string, title string) error {
	if len(leftValues) > 0 && len(rightValues) > 0 {
		return fmt.Errorf(format, title, leftValues, rightValues)
	}

	return nil
}

func eliminateDuplicatedValues(values []string) []string {
	mapOfValues := map[string]bool{}

	for _, v := range values {
		mapOfValues[v] = true
	}

	resultValues := make([]string, 0, len(mapOfValues))
	for k := range mapOfValues {
		resultValues = append(resultValues, k)
	}

	return resultValues
}
