package nqm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/common/compress"
	"github.com/Cepave/open-falcon-backend/common/digest"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	sjson "github.com/bitly/go-simplejson"
)

var flateCompressor = compress.NewDefaultFlateCompressor()

// Defines the IR for relation of hosts(between agent and target)
type PropRelation int8

const (
	/**
	 * Input value of realtion
	 */
	RelationSame = -11
	RelationNotSame = -12
	// :~)

	// The relation is unknown
	NoCondition PropRelation = -1
	// Means a property of agent and target must be same
	SameValue PropRelation = 1
	// Means a property of agent and target may not be same
	NotSameValue PropRelation = 2
)

const (
	MetricMax = "max"
    MetricMin = "min"
    MetricAvg = "avg"
    MetricMed = "med"
    MetricMdev = "mdev"
    MetricLoss = "loss"
    MetricCount = "count"
    MetricPckSent = "pck_sent"
    MetricPckReceived = "pck_received"
    MetricNumAgent = "num_agent"
    MetricNumTarget = "num_target"

	AgentGroupingName = "name"
	AgentGroupingIpAddress = "ip_address"
	AgentGroupingHostname = "hostname"

	TargetGroupingName = "name"
	TargetGroupingHost = "host"

	GroupingProvince = "province"
	GroupingCity = "city"
	GroupingIsp = "isp"
	GroupingNameTag = "name_tag"

	TimeRangeAbsolute byte = 1
	TimeRangeRelative byte = 2

	TimeUnitYear = "y"
	TimeUnitMonth = "m"
	TimeUnitWeek = "w"
	TimeUnitDay = "d"
	TimeUnitHour = "h"
	TimeUnitMinute = "n"
)

type MetricsFilterParseError struct {
	error
}

var supportingTimeUnit = map[string]bool {
	TimeUnitYear: false,
	TimeUnitMonth: false,
	TimeUnitWeek: false,
	TimeUnitDay: false,
	TimeUnitHour: false,
	TimeUnitMinute: false,
}

var supportingOutput = map[string]bool {
	MetricMax: true,
    MetricMin: true,
    MetricAvg: true,
    MetricMed: true,
    MetricMdev: true,
    MetricLoss: true,
    MetricCount: true,
    MetricPckSent: true,
    MetricPckReceived: true,
    MetricNumAgent: true,
    MetricNumTarget: true,
}

var supportingAgentGrouping = map[string]bool {
	AgentGroupingName: true,
	AgentGroupingIpAddress: true,
	AgentGroupingHostname: true,
	GroupingProvince: true,
	GroupingCity: true,
	GroupingIsp: true,
	GroupingNameTag: true,
}

var supportingTargetGrouping = map[string]bool {
	TargetGroupingName: true,
	TargetGroupingHost: true,
	GroupingProvince: true,
	GroupingCity: true,
	GroupingIsp: true,
	GroupingNameTag: true,
}

// The main object of query for compound report
type CompoundQuery struct {
	Filters *CompoundQueryFilter `json:"filters" digest:"1"`

	// Grouping content of report
	Grouping *QueryGrouping `json:"grouping" digest:"11"`

	// Output content of report
	Output *QueryOutput `json:"output" digest:"21"`

	metricFilter MetricFilter
}

func NewCompoundQuery() *CompoundQuery {
	return &CompoundQuery {
		Output: &QueryOutput { },
		Grouping: &QueryGrouping { },
		Filters: &CompoundQueryFilter {
			Time: NewTimeFilter(),
			Agent: &nqmModel.AgentFilter{},
			Target: &nqmModel.TargetFilter{},
		},
	}
}

type CompoundQueryFilter struct {
	Time *TimeFilter `json:"time" digest:"1"`
	Agent *nqmModel.AgentFilter `json:"agent" digest:"2"`
	Target *nqmModel.TargetFilter `json:"target" digest:"3"`
	Metrics string `json:"metrics" digest:"4"`
}

type TimeFilter struct {
	StartTime ojson.JsonTime `json:"start_time"`
	EndTime ojson.JsonTime `json:"end_time"`
	ToNow *TimeWithUnit `json:"to_now"`

	netStartTime time.Time
	netEndTime time.Time

	timeRangeType byte
}

func NewTimeFilter() *TimeFilter {
	return &TimeFilter {
		timeRangeType: 0,
	}
}

type TimeWithUnit struct {
	Unit string `json:"unit" digest:"1"`
	Value int `json:"value" digest:"2"`
}

type QueryOutput struct {
	Metrics []string `json:"metrics" digest:"1"`
}

type QueryGrouping struct {
	Agent []string `json:"agent" digest:"1"`
	Target []string `json:"target" digest:"2"`
}

func (q *CompoundQuery) GetDigestValue() []byte {
	return digest.DigestStruct(q, digest.Md5SumFunc)
}
func (q *CompoundQuery) SetMetricFilter(filter MetricFilter) {
	q.metricFilter = filter
}

func (q *CompoundQuery) GetIspRelation() PropRelation {
	filters := q.Filters
	return getRelation(append(filters.Agent.IspIds, filters.Target.IspIds...))
}
func (q *CompoundQuery) GetProvinceRelation() PropRelation {
	filters := q.Filters
	return getRelation(append(filters.Agent.ProvinceIds, filters.Target.ProvinceIds...))
}
func (q *CompoundQuery) GetCityRelation() PropRelation {
	filters := q.Filters
	return getRelation(append(filters.Agent.CityIds, filters.Target.CityIds...))
}
func (q *CompoundQuery) GetNameTagRelation() PropRelation {
	filters := q.Filters
	return getRelation(append(filters.Agent.NameTagIds, filters.Target.NameTagIds...))
}

// Converts this query object to compressed query
func (q *CompoundQuery) GetCompressedQuery() []byte {
	return flateCompressor.MustCompressString(
		ojson.MarshalJSON(q),
	)
}

func (q *CompoundQuery) UnmarshalFromCompressedQuery(compressedQuery []byte) {
	json := flateCompressor.MustDecompressToString(compressedQuery)
	q.UnmarshalJSON([]byte(json))
}

// Processes the source of JSON to initialize query object
func (query *CompoundQuery) UnmarshalJSON(jsonSource []byte) (err error) {
	json, jsonErr := sjson.NewJson(jsonSource)
	if jsonErr != nil {
		return jsonErr
	}

	defer func() {
		r := recover()
		if r != nil {
			switch errValue := r.(type) {
			case error:
				err = errValue
			default:
				err = fmt.Errorf("Umarshal JSON of query has error. %v", errValue)
			}
		}
	}()

	// Loads "filters" property
	err = query.Filters.UnmarshalSimpleJSON(json.Get("filters"))
	if err != nil {
		return
	}

	// Loads "grouping" property
	err = query.Grouping.UnmarshalSimpleJSON(json.Get("grouping"))
	if err != nil {
		return
	}

	// Loads "output" property
	err = query.Output.UnmarshalSimpleJSON(json.Get("output"))
	if err != nil {
		return
	}

	return
}

func (query *CompoundQuery) SetupDefault() {
	if len(query.Output.Metrics) == 0 {
		query.Output.Metrics = []string{
			MetricMax, MetricMin, MetricAvg, MetricLoss, MetricCount,
		}
	}

	if len(query.Grouping.Agent) + len(query.Grouping.Target) == 0 {
		query.Grouping.Agent = []string {
			AgentGroupingName, AgentGroupingIpAddress,
			GroupingProvince, GroupingCity, GroupingIsp, GroupingNameTag,
		}
	}

	timeRange := query.Filters.Time
	if timeRange.timeRangeType == 0 {
		timeRange.timeRangeType = TimeRangeRelative
		timeRange.ToNow = &TimeWithUnit {
			Unit: TimeUnitHour,
			Value: 1,
		}
	}
}

func (f *CompoundQueryFilter) UnmarshalSimpleJSON(jsonObject *sjson.Json) (err error) {
	if err = f.Time.UnmarshalSimpleJSON(jsonObject.Get("time"))
		err != nil {
		return
	}

	f.loadFilterOfAgent(jsonObject.Get("agent"))
	f.loadFilterOfTarget(jsonObject.Get("target"))
	f.loadFilterOfMetrics(jsonObject.Get("metrics"))

	return
}
func (f *CompoundQueryFilter) loadFilterOfAgent(jsonObject *sjson.Json) {
	agentFilter := f.Agent

	agentFilter.Name = purifyStringArrayOfJsonForValues(
		jsonObject.Get("name"),
	)
	agentFilter.Hostname = purifyStringArrayOfJsonForValues(
		jsonObject.Get("hostname"),
	)
	agentFilter.IpAddress = purifyStringArrayOfJsonForValues(
		jsonObject.Get("ip_address"),
	)
	agentFilter.ConnectionId = purifyStringArrayOfJsonForValues(
		jsonObject.Get("connection_id"),
	)

	agentFilter.IspIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.Get("isp_ids"),
	)
	agentFilter.ProvinceIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.Get("province_ids"),
	)
	agentFilter.CityIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.Get("city_ids"),
	)
	agentFilter.NameTagIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.Get("name_tag_ids"),
	)
	agentFilter.GroupTagIds = purifyNumberArrayOfJsonToInt32(
		jsonObject.Get("group_tag_ids"),
	)
}
func (f *CompoundQueryFilter) loadFilterOfTarget(jsonObject *sjson.Json) {
	targetFilter := f.Target

	targetFilter.Name = purifyStringArrayOfJsonForValues(
		jsonObject.Get("name"),
	)
	targetFilter.Host = purifyStringArrayOfJsonForValues(
		jsonObject.Get("host"),
	)

	targetFilter.IspIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.Get("isp_ids"),
	)
	targetFilter.ProvinceIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.Get("province_ids"),
	)
	targetFilter.CityIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.Get("city_ids"),
	)
	targetFilter.NameTagIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.Get("name_tag_ids"),
	)
	targetFilter.GroupTagIds = purifyNumberArrayOfJsonToInt32(
		jsonObject.Get("group_tag_ids"),
	)
}
func (f *CompoundQueryFilter) loadFilterOfMetrics(jsonObject *sjson.Json) {
	f.Metrics = purifyStringOfJson(jsonObject)
}

func (f *TimeFilter) GetDigest() []byte {
	switch f.timeRangeType {
	case TimeRangeAbsolute:
		bytesOfRange := make([]byte, 0)
		bytesOfRange = append(bytesOfRange, digest.DigestableTime(f.StartTime).GetDigest()...)
		bytesOfRange = append(bytesOfRange, digest.DigestableTime(f.EndTime).GetDigest()...)
		return digest.Md5SumFunc(bytesOfRange)
	case TimeRangeRelative:
		return digest.GetBytesGetter(f.ToNow, digest.Md5SumFunc)()
	}

	panic(fmt.Sprintf("Unknown type of time filter for digesting: %#v", f))
}
// Retrieves the start time of net(whether or not the time is absolute or relative)
func (f *TimeFilter) GetNetTimeRange() (time.Time, time.Time) {
	switch f.timeRangeType {
	case TimeRangeAbsolute:
		return time.Time(f.StartTime), time.Time(f.EndTime)
	case TimeRangeRelative:
		if f.netStartTime.IsZero() {
			f.netStartTime, f.netEndTime = f.getRelativeTimeRangeOfNet(time.Now())
		}
		return f.netStartTime, f.netEndTime
	default:
		panic("Unknown time type")
	}
}
func (f *TimeFilter) getRelativeTimeRangeOfNet(baseTime time.Time) (time.Time, time.Time) {
	var startTime, endTime time.Time

	startTimeValue, endTimeValue := f.ToNow.Value, f.ToNow.Value
	if endTimeValue == 0 {
		endTimeValue = 1
	}

	durationStartTime, durationEndTime := time.Duration(startTimeValue), time.Duration(endTimeValue)

	switch f.ToNow.Unit {
	case TimeUnitYear:
		startTime = time.Date(
			baseTime.Year(), 1, 1, 0, 0, 0, 0, baseTime.Location(),
		).AddDate(-startTimeValue, 0, 0)
		endTime = startTime.AddDate(endTimeValue, 0, 0)
	case TimeUnitMonth:
		startTime = time.Date(
			baseTime.Year(), baseTime.Month(), 1, 0, 0, 0, 0, baseTime.Location(),
		).AddDate(0, -startTimeValue, 0)
		endTime = startTime.AddDate(0, endTimeValue, 0)
	case TimeUnitWeek:
		minusDays := -((int(baseTime.Weekday()) + 6) % 7)
		startTime = time.Date(
			baseTime.Year(), baseTime.Month(), baseTime.Day(), 0, 0, 0, 0, baseTime.Location(),
		).AddDate(0, 0, minusDays).AddDate(0, 0, -startTimeValue * 7)
		endTime = startTime.AddDate(0, 0, endTimeValue * 7)
	case TimeUnitDay:
		startTime = time.Date(
			baseTime.Year(), baseTime.Month(), baseTime.Day(), 0, 0, 0, 0, baseTime.Location(),
		).AddDate(0, 0, -startTimeValue)
		endTime = startTime.AddDate(0, 0, endTimeValue)
	case TimeUnitHour:
		startTime = time.Date(
			baseTime.Year(), baseTime.Month(), baseTime.Day(), baseTime.Hour(), 0, 0, 0, baseTime.Location(),
		).Add(-durationStartTime * time.Hour)
		endTime = startTime.Add(durationEndTime * time.Hour)
	case TimeUnitMinute:
		startTime = time.Date(
			baseTime.Year(), baseTime.Month(), baseTime.Day(), baseTime.Hour(), baseTime.Minute(), 0, 0, baseTime.Location(),
		).Add(-durationStartTime * time.Minute)
		endTime = startTime.Add(durationEndTime * time.Minute)
	}

	return startTime, endTime
}
func (f *TimeFilter) UnmarshalSimpleJSON(jsonObject *sjson.Json) error {
	if jsonObject == nil {
		return nil
	}

	if startTime, endTime := f.loadAbsoluteTimeRange(jsonObject)
		!startTime.IsZero() && !endTime.IsZero() {
		f.timeRangeType = TimeRangeAbsolute
		f.StartTime = ojson.JsonTime(startTime)
		f.EndTime = ojson.JsonTime(endTime)
		f.ToNow = nil
	} else if timeWithUnit := f.loadRelativeTimeRange(jsonObject)
		timeWithUnit != nil {
		f.timeRangeType = TimeRangeRelative
		f.ToNow = timeWithUnit
	}

	return nil
}
func (f *TimeFilter) loadRelativeTimeRange(jsonTime *sjson.Json) *TimeWithUnit {
	jsonToNow, hasToNow := jsonTime.CheckGet("to_now")
	if !hasToNow {
		return nil
	}

	stringOfTimeUnit := strings.ToLower(
		strings.TrimSpace(jsonToNow.GetPath("unit").MustString()),
	)
	if _, ok := supportingTimeUnit[stringOfTimeUnit]; !ok {
		return nil
	}

	return &TimeWithUnit{
		Unit: stringOfTimeUnit,
		Value: jsonToNow.Get("value").MustInt(),
	}
}
func (f *TimeFilter) loadAbsoluteTimeRange(jsonTime *sjson.Json) (time.Time, time.Time) {
	jsonStartTime := jsonTime.Get("start_time")
	jsonEndTime := jsonTime.Get("end_time")

	var startTime, endTime time.Time
	if jsonStartTime.Interface() != nil {
		startTime = time.Unix(jsonStartTime.MustInt64(), 0)
	}
	if jsonStartTime.Interface() != nil {
		endTime = time.Unix(jsonEndTime.MustInt64(), 0)
	}

	return startTime, endTime
}

func (o *QueryOutput) UnmarshalSimpleJSON(jsonObject *sjson.Json) error {
	o.Metrics = purifyStringArrayOfJsonForDomain(
		jsonObject.Get("metrics"),
		supportingOutput,
	)

	return nil
}
func (o *QueryOutput) HasMetric(metricName string) bool {
	for _, metric := range o.Metrics {
		if metricName == metric {
			return true
		}
	}

	return false
}

func (tu *TimeWithUnit) String() string {
	return fmt.Sprintf("%d %s", tu.Value, tu.Unit)
}

var eachAgentGrouping = map[string]bool {
	AgentGroupingName: true,
	AgentGroupingIpAddress: true,
	AgentGroupingHostname: true,
}
var eachTargetGrouping = map[string]bool {
	TargetGroupingName: true,
	TargetGroupingHost: true,
}

func (g *QueryGrouping) UnmarshalSimpleJSON(jsonObject *sjson.Json) error {
	g.Agent = purifyStringArrayOfJsonForDomain(
		jsonObject.Get("agent"),
		supportingAgentGrouping,
	)

	g.Target = purifyStringArrayOfJsonForDomain(
		jsonObject.Get("target"),
		supportingTargetGrouping,
	)

	return nil
}
func (g *QueryGrouping) IsForEachAgent() bool {
	for _, agentGroup := range g.Agent {
		if _, ok := eachAgentGrouping[agentGroup]; ok {
			return true
		}
	}

	return false
}
func (g *QueryGrouping) IsForEachTarget() bool {
	for _, targetGroup := range g.Target {
		if _, ok := eachTargetGrouping[targetGroup]; ok {
			return true
		}
	}

	return false
}

func purifyNumberArrayOfJsonToInt16(jsonObject *sjson.Json) []int16 {
	return purifyNumberArrayOfJson(
		jsonObject, utils.TypeOfInt16,
	).([]int16)
}
func purifyNumberArrayOfJsonToInt32(jsonObject *sjson.Json) []int32 {
	return purifyNumberArrayOfJson(
		jsonObject, utils.TypeOfInt32,
	).([]int32)
}
func purifyNumberArrayOfJson(jsonObject *sjson.Json, targetType reflect.Type) interface{} {
	if jsonObject == nil {
		return reflect.MakeSlice(reflect.SliceOf(targetType), 0, 0).Interface()
	}

	uniqueFilter := utils.NewUniqueFilter(utils.TypeOfInt)
	arrayObject := utils.MakeAbstractArray(jsonObject.MustArray()).
		MapTo(
			func (v interface{}) interface{} {
				int64Value, err := v.(json.Number).Int64()
				if err != nil {
					panic(err)
				}

				return int(int64Value)
			},
			utils.TypeOfInt,
		).
		FilterWith(uniqueFilter)

	intArray := arrayObject.GetArray().([]int)
	sort.Ints(intArray)

	return utils.MakeAbstractArray(intArray).GetArrayAsType(targetType)
}

func purifyStringOfJson(jsonObject *sjson.Json) string {
	if jsonObject == nil {
		return ""
	}

	return strings.ToLower(strings.TrimSpace(
		jsonObject.MustString(),
	))
}
func purifyStringArrayOfJsonForDomain(jsonObject *sjson.Json, domain map[string]bool) []string {
	if jsonObject == nil {
		return []string{}
	}

	uniqueFilter := utils.NewUniqueFilter(utils.TypeOfString)
	domainFilter := utils.NewDomainFilter(domain)
	arrayObject := utils.MakeAbstractArray(jsonObject.MustStringArray()).
		MapTo(
			utils.TypedFuncToMapper(
				func (v string) string {
					return strings.ToLower(strings.TrimSpace(v))
				},
			),
			utils.TypeOfString,
		).
		FilterWith(func(v interface{}) bool {
			return utils.EmptyStringFilter(v) &&
				domainFilter(v) &&
				uniqueFilter(v)
		})

	return arrayObject.GetArray().([]string)
}
func purifyStringArrayOfJsonForValues(jsonObject *sjson.Json) []string {
	if jsonObject == nil {
		return []string{}
	}

	uniqueFilter := utils.NewUniqueFilter(utils.TypeOfString)
	arrayObject := utils.MakeAbstractArray(jsonObject.MustStringArray()).
		MapTo(utils.TrimStringMapper, utils.TypeOfString).
		FilterWith(func(v interface{}) bool {
			return utils.EmptyStringFilter(v) &&
				uniqueFilter(v)
		})

	resultArray := arrayObject.GetArray().([]string)
	sort.Strings(resultArray)

	return resultArray
}

func getRelation(arrayOfNumber interface{}) PropRelation {
	newIntArray := utils.MakeAbstractArray(arrayOfNumber).
		GetArrayAsType(utils.TypeOfInt).([]int)

	for _, v := range newIntArray {
		if v == RelationNotSame {
			return NotSameValue
		}

		if v == RelationSame {
			return SameValue
		}
	}

	return NoCondition
}
