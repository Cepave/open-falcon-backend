package nqm_parser

import (
	"fmt"
	"time"
)

// The parameters for query
type QueryParams struct {
	StartTime time.Time
	EndTime time.Time
	AgentFilter NodeFilter
	TargetFilter NodeFilter
}

// The filter of node
type NodeFilter struct {
	MatchProvinces []string
	MatchCities []string
	MatchIsps []string
}

// Checks the paramters
//
// 1. provinces and cities cannot be assigned at the same time
// 2. duplicated value would be eliminated
const FORMAT_ERROR_LOCATION_FILTER = "%v filter for provinces:%v and cities:%v are both set"
func (p* QueryParams) checkParams() (err error) {
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

func (p *QueryParams) addIspOfAgent(values ...string) {
	p.AgentFilter.MatchIsps = append(p.AgentFilter.MatchIsps, values...)
}
func (p *QueryParams) addProvinceOfAgent(values ...string) {
	p.AgentFilter.MatchProvinces = append(p.AgentFilter.MatchProvinces, values...)
}
func (p *QueryParams) addCityOfAgent(values ...string) {
	p.AgentFilter.MatchCities = append(p.AgentFilter.MatchCities, values...)
}
func (p *QueryParams) addIspOfTarget(values ...string) {
	p.TargetFilter.MatchIsps = append(p.TargetFilter.MatchIsps, values...)
}
func (p *QueryParams) addProvinceOfTarget(values ...string) {
	p.TargetFilter.MatchProvinces = append(p.TargetFilter.MatchProvinces, values...)
}
func (p *QueryParams) addCityOfTarget(values ...string) {
	p.TargetFilter.MatchCities = append(p.TargetFilter.MatchCities, values...)
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
	for k, _ := range mapOfValues {
		resultValues = append(resultValues, k)
	}

	return resultValues
}
