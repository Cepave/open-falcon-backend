package nqm

import (
	"encoding/json"
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
	GroupingName_tag = "name_tag"

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
	GroupingName_tag: true,
}

var supportingTargetGrouping = map[string]bool {
	TargetGroupingName: true,
	TargetGroupingHost: true,
	GroupingProvince: true,
	GroupingCity: true,
	GroupingIsp: true,
	GroupingName_tag: true,
}

// The main object of query for compound report
type CompoundQuery struct {
	Filters *CompoundQueryFilter `json:"filters" digest:"1"`

	// Grouping content of report
	Grouping *QueryGrouping `json:"grouping" digest:"11"`

	// Output content of report
	Output *QueryOutput `json:"output" digest:"21"`

	jsonObject *sjson.Json
}
func (q *CompoundQuery) GetDigestValue() []byte {
	return digest.DigestStruct(q, digest.Md5SumFunc)
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

	timeRangeType byte
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

	panic("Unknown type of time filter for digesting")
}

func NewTimeFilter() *TimeFilter {
	return &TimeFilter {
		timeRangeType: 0,
		ToNow: &TimeWithUnit { "", 0 },
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

// Converts this query object to compressed query
func (q *CompoundQuery) GetCompressedQuery() []byte {
	json, jsonErr := json.Marshal(q)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return flateCompressor.MustCompressString(string(json))
}

func (q *CompoundQuery) UnmarshalFromCompressedQuery(compressedQuery []byte) {
	json := flateCompressor.MustDecompressToString(compressedQuery)
	q.UnmarshalJSON([]byte(json))
}

// Converts this query object to compressed query
func (q *CompoundQuery) GetQueryDigest() []byte {
	return nil
}

// Processes the source of JSON to initialize query object
func (query *CompoundQuery) UnmarshalJSON(jsonSource []byte) error {
	json, err := sjson.NewJson(jsonSource)
	if err != nil {
		return err
	}

	query.jsonObject = json

	// Loads "filters" property
	if err = query.loadFilters(); err != nil {
		return err
	}
	// Loads "grouping" property
	if err = query.loadGrouping(); err != nil {
		return err
	}
	// Loads "output" property
	if err = query.loadOutput(); err != nil {
		return err
	}

	return nil
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
			GroupingProvince, GroupingCity, GroupingIsp, GroupingName_tag,
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

func (query *CompoundQuery) loadFilters() (err error) {
	query.Filters.Metrics = purifyStringOfJson(
		query.jsonObject.GetPath("filters", "metrics"),
	)

	if err = query.loadFiltersOfTime(); err != nil {
		return
	}
	if err = query.loadFiltersOfAgent(); err != nil {
		return
	}
	if err = query.loadFiltersOfTarget(); err != nil {
		return
	}

	return
}
func (query *CompoundQuery) loadFiltersOfTime() error {
	jsonObject := query.jsonObject
	timeFilter := query.Filters.Time

	jsonTime := jsonObject.GetPath("filters", "time")

	if startTime, endTime := query.loadAbsoluteTimeRange(jsonTime)
		!startTime.IsZero() && !endTime.IsZero() {
		timeFilter.timeRangeType = TimeRangeAbsolute
		timeFilter.StartTime = ojson.JsonTime(startTime)
		timeFilter.EndTime = ojson.JsonTime(endTime)
	}

	if timeWithUnit := query.loadRelativeTimeRange(jsonTime)
		timeWithUnit != nil {
		timeFilter.timeRangeType = TimeRangeRelative
		timeFilter.ToNow = timeWithUnit
	}

	return nil
}
func (query *CompoundQuery) loadRelativeTimeRange(jsonTime *sjson.Json) *TimeWithUnit {
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
func (query *CompoundQuery) loadAbsoluteTimeRange(jsonTime *sjson.Json) (time.Time, time.Time) {
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
func (query *CompoundQuery) loadFiltersOfAgent() error {
	agentFilter := query.Filters.Agent
	jsonObject := query.jsonObject

	agentFilter.Name = purifyStringArrayOfJsonForValues(
		jsonObject.GetPath("filters", "agent", "name"),
	)
	agentFilter.Hostname = purifyStringArrayOfJsonForValues(
		jsonObject.GetPath("filters", "agent", "hostname"),
	)
	agentFilter.IpAddress = purifyStringArrayOfJsonForValues(
		jsonObject.GetPath("filters", "agent", "ip_address"),
	)
	agentFilter.ConnectionId = purifyStringArrayOfJsonForValues(
		jsonObject.GetPath("filters", "agent", "connection_id"),
	)

	agentFilter.IspIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.GetPath("filters", "agent", "isp_ids"),
	)
	agentFilter.ProvinceIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.GetPath("filters", "agent", "province_ids"),
	)
	agentFilter.CityIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.GetPath("filters", "agent", "city_ids"),
	)
	agentFilter.NameTagIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.GetPath("filters", "agent", "name_tag_ids"),
	)
	agentFilter.GroupTagIds = purifyNumberArrayOfJsonToInt32(
		jsonObject.GetPath("filters", "agent", "group_tag_ids"),
	)

	return nil
}
func (query *CompoundQuery) loadFiltersOfTarget() error {
	targetFilter := query.Filters.Target
	jsonObject := query.jsonObject

	targetFilter.Name = purifyStringArrayOfJsonForValues(
		jsonObject.GetPath("filters", "target", "name"),
	)
	targetFilter.Host = purifyStringArrayOfJsonForValues(
		jsonObject.GetPath("filters", "target", "host"),
	)

	targetFilter.IspIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.GetPath("filters", "target", "isp_ids"),
	)
	targetFilter.ProvinceIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.GetPath("filters", "target", "province_ids"),
	)
	targetFilter.CityIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.GetPath("filters", "target", "city_ids"),
	)
	targetFilter.NameTagIds = purifyNumberArrayOfJsonToInt16(
		jsonObject.GetPath("filters", "target", "name_tag_ids"),
	)
	targetFilter.GroupTagIds = purifyNumberArrayOfJsonToInt32(
		jsonObject.GetPath("filters", "target", "group_tag_ids"),
	)

	return nil
}
func (query *CompoundQuery) loadGrouping() error {
	query.Grouping.Agent = purifyStringArrayOfJsonForDomain(
		query.jsonObject.GetPath("grouping", "agent"),
		supportingAgentGrouping,
	)

	query.Grouping.Target = purifyStringArrayOfJsonForDomain(
		query.jsonObject.GetPath("grouping", "target"),
		supportingTargetGrouping,
	)

	return nil
}
func (query *CompoundQuery) loadOutput() error {
	query.Output.Metrics = purifyStringArrayOfJsonForDomain(
		query.jsonObject.GetPath("output", "metrics"),
		supportingOutput,
	)

	return nil
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
		utils.MakeAbstractArray([]int{}).GetArrayAsType(targetType)
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
