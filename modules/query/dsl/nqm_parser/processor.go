package nqm_parser

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	_ = iota

	parse_iso8601_minute = "2006-01-02T15:04Z07:00"
)

var currentTimeZone = time.Now().Format("Z07:00")

type paramSetter func(*QueryParams)
var empty_param_setter paramSetter = func(queryParams *QueryParams) {}

func buildErrorForInvalidParam(paramName interface{}, assignedValue interface{}) error {
	if assignedValue == nil {
		return fmt.Errorf("\"%v\" missed \"=<value>\" ?", paramName);
	}

	if assignedValue != nil {
		paramValue := assignedValue.([]interface{})[1]

		if paramValue != nil {
			return fmt.Errorf("\"%v\" cannot accept \"%v\"", paramName, paramValue)
		}

		return fmt.Errorf("\"%v=\" need set value", paramName)
	}

	return fmt.Errorf("Unknown error")
}

// Converts []interface{} to []paramSetter
func toSetters(finalResult interface{}) []paramSetter {
	if finalResult == nil {
		return make([]paramSetter, 0)
	}

	finalResultAsArray := finalResult.([]interface{})
	resultAsParamSetters := make([]paramSetter, len(finalResultAsArray))

	for i, untypedParamSetter := range finalResultAsArray {
		resultAsParamSetters[i] = untypedParamSetter.(paramSetter)
	}

	return resultAsParamSetters
}

// Builds paramSetter by content of param
func buildSetterFunc(paramName interface{}, srcParamValue interface{}) (newSetter paramSetter, err error) {
	paramNameAsString := paramName.(string)

	switch srcParamValue.(type) {
	case time.Time:
		newSetter, err = buildSetterFuncForTime(paramNameAsString, srcParamValue.(time.Time))
	case []string:
		newSetter, err = buildSetterFuncForStringArray(paramNameAsString, srcParamValue.([]string))
	case HostRelation:
		newSetter, err = buildSetterFuncForHostRelation(paramNameAsString, srcParamValue.(HostRelation))
	default:
		err = fmt.Errorf("Unsupported for type of param[%v]. Param name: [%v]", reflect.TypeOf(srcParamValue), paramNameAsString)
	}

	return
}

func buildSetterFuncForTime(paramName string, timeValue time.Time) (newSetter paramSetter, err error) {
	newSetter = nil

	switch paramName {
		case "starttime":
			newSetter = func (p *QueryParams) {
				p.StartTime = timeValue
			}
		case "endtime":
			newSetter = func (p *QueryParams) {
				p.EndTime = timeValue
			}
	}

	if newSetter == nil {
		err = fmt.Errorf("Unsupported time value for property: [%v]", paramName)
	}

	return
}

func buildSetterFuncForStringArray(paramName string, values []string) (newSetter paramSetter, err error) {
	newSetter = nil

	switch paramName {
	case "agent.isp":
		newSetter = func (p *QueryParams) {
			appendForString(&p.AgentFilter.MatchIsps, values)
		}
	case "agent.province":
		newSetter = func (p *QueryParams) {
			appendForString(&p.AgentFilter.MatchProvinces, values)
		}
	case "agent.city":
		newSetter = func (p *QueryParams) {
			appendForString(&p.AgentFilter.MatchCities, values)
		}
	case "target.isp":
		newSetter = func (p *QueryParams) {
			appendForString(&p.TargetFilter.MatchIsps, values)
		}
	case "target.province":
		newSetter = func (p *QueryParams) {
			appendForString(&p.TargetFilter.MatchProvinces, values)
		}
	case "target.city":
		newSetter = func (p *QueryParams) {
			appendForString(&p.TargetFilter.MatchCities, values)
		}
	}

	if newSetter == nil {
		err = fmt.Errorf("Unsupported []string value for property: [%v]", paramName)
	}

	return
}

// If multiple %<AUTO_COND>% has set on the same properties of agent/target,
// this building would apply the last one
func buildSetterFuncForHostRelation(paramName string, relationValue HostRelation) (newSetter paramSetter, err error) {
	newSetter = nil

	switch paramName {
	case "agent.isp":
		newSetter = func (p *QueryParams) {
			p.IspRelation = relationValue
		}
	case "agent.province":
		newSetter = func (p *QueryParams) {
			p.ProvinceRelation = relationValue
		}
	case "agent.city":
		newSetter = func (p *QueryParams) {
			p.CityRelation = relationValue
		}
	case "target.isp":
		newSetter = func (p *QueryParams) {
			p.IspRelation = relationValue
		}
	case "target.province":
		newSetter = func (p *QueryParams) {
			p.ProvinceRelation = relationValue
		}
	case "target.city":
		newSetter = func (p *QueryParams) {
			p.CityRelation = relationValue
		}
	}

	if newSetter == nil {
		err = fmt.Errorf("Unsupported \"HostRelation\" value for property: [%v]", paramName)
	}

	return
}

func appendForString(valuesHolder *[]string, values []string) {
	*valuesHolder = append(*valuesHolder, values...)
}

func combineStringLiterals(first interface{}, rest interface{}) []string {
	allRests := rest.([]interface{})

	result := make([]string, 0, len(allRests) + 1)
	result = append(result, first.(string))

	for _, v := range allRests {
		result = append(result, v.(string))
	}

	return result
}

// Parses the unix time by string representation of "1213213"(Epoch tiem)
func parseUnixTime(c *current) (time.Time, error) {
	unixTimeInt64, parseErr := strconv.ParseInt(string(c.text), 0, 64)

	if parseErr != nil {
		return time.Unix(0, 0), parseErr
	}

	return time.Unix(unixTimeInt64, 0), nil
}

func parseAutoCondition(autoConditionValue interface{}) (HostRelation, error) {
	stringValue := autoConditionValue.(string)

	switch stringValue {
	case "MATCH_ANOTHER":
		return SAME_VALUE, nil
	case "NOT_MATCH_ANOTHER":
		return NOT_SAME_VALUE, nil
	default:
		return UNKNOWN_RELATION, fmt.Errorf("Unknown auto-condition: %%%v%%", stringValue);
	}
}

// Parses the string representation of ISO-8601 format
//
// For example(assumes local timezone is "+08:00"):
//
// 2010-05-05 -> 2010-05-05T00:00+08:00
// 2010-04-12T02 -> 2010-04-12T02:00+08:00
// 2010-03-05T10:30 -> 2010-03-05T10:30+08:00
// 2012-01-10T10:30+02:00 -> 2012-01-10T10:30+02:00
func parseIso8601(c *current) (time.Time, error) {
	timeStr := string(c.text)

	switch (len(timeStr)) {
	case 10:
		timeStr = fmt.Sprintf("%sT00:00%s", timeStr, currentTimeZone)
	case 13:
		timeStr = fmt.Sprintf("%s:00%s", timeStr, currentTimeZone)
	case 16:
		timeStr = fmt.Sprintf("%s%s", timeStr, currentTimeZone)
	}

	return time.Parse(parse_iso8601_minute, timeStr)
}
