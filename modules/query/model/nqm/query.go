package nqm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
	"github.com/Cepave/open-falcon-backend/common/utils"
	sjson "github.com/bitly/go-simplejson"
)

const (
	METRICS_MAX = "max"
    METRICS_MIN = "min"
    METRICS_AVG = "avg"
    METRICS_MED = "med"
    METRICS_MDEV = "mdev"
    METRICS_LOSS = "loss"
    METRICS_COUNT = "count"
    METRICS_PCK_SENT = "pck_sent"
    METRICS_PCK_RECEIVED = "pck_received"
    METRICS_NUM_AGENT = "num_agent"
    METRICS_NUM_TARGET = "num_target"

	AGENT_GP_NAME = "name"
	AGENT_GP_IP_ADDRESS = "ip_address"
	AGENT_GP_HOSTNAME = "hostname"

	TARGET_GP_NAME = "name"
	TARGET_GP_HOST = "host"

	GROUPING_PROVINCE = "province"
	GROUPING_CITY = "city"
	GROUPING_ISP = "isp"
	GROUPING_NAME_TAG = "name_tag"

	TIME_RANGE_ABSOLUTE byte = 1
	TIME_RANGE_RELATIVE byte = 2

	TU_YEAR byte = 6
	TU_MONTH byte = 5
	TU_WEEK byte = 4
	TU_DAY byte = 3
	TU_HOUR byte = 2
	TU_MINUTE byte = 1
)

type MetricsFilterParseError struct {
}

func (e *MetricsFilterParseError) Error() string {
	return ""
}

var supportingTimeUnit = map[string]byte {
	"y": TU_YEAR,
	"m": TU_MONTH,
	"w": TU_WEEK,
	"d": TU_DAY,
	"h": TU_HOUR,
	"n": TU_MINUTE,
}

var supportingOutput = map[string]bool {
	METRICS_MAX: true,
    METRICS_MIN: true,
    METRICS_AVG: true,
    METRICS_MED: true,
    METRICS_MDEV: true,
    METRICS_LOSS: true,
    METRICS_COUNT: true,
    METRICS_PCK_SENT: true,
    METRICS_PCK_RECEIVED: true,
    METRICS_NUM_AGENT: true,
    METRICS_NUM_TARGET: true,
}

var supportingAgentGrouping = map[string]bool {
	AGENT_GP_NAME: true,
	AGENT_GP_IP_ADDRESS: true,
	AGENT_GP_HOSTNAME: true,
	GROUPING_PROVINCE: true,
	GROUPING_CITY: true,
	GROUPING_ISP: true,
	GROUPING_NAME_TAG: true,
}

var supportingTargetGrouping = map[string]bool {
	TARGET_GP_NAME: true,
	TARGET_GP_HOST: true,
	GROUPING_PROVINCE: true,
	GROUPING_CITY: true,
	GROUPING_ISP: true,
	GROUPING_NAME_TAG: true,
}

// The main object of query for compound report
type CompoundQuery struct {
	// Output content of report
	Output *struct {
		Metrics []string
	}

	// Grouping content of report
	Grouping *struct {
		Agent []string
		Target []string
	}

	Filters *CompoundQueryFilter

	jsonObject *sjson.Json
}

type CompoundQueryFilter struct {
	Time *TimeFilter
	Agent *AgentFilter
	Target *TargetFilter
	Metrics string
}

type TimeFilter struct {
	TimeRangeType byte
	StartTime time.Time
	EndTime time.Time
	ToNow *TimeWithUnit
}

func NewTimeFilter() *TimeFilter {
	return &TimeFilter {
		TimeRangeType: 0,
		StartTime: time.Unix(0, 0),
		EndTime: time.Unix(0, 0),
		ToNow: &TimeWithUnit{},
	}
}

type TimeWithUnit struct {
	Unit byte
	Value int
}

type AgentFilter struct {
	Name []string
	Hostname []string
	IpAddress []string
	ConnectionId []string
	IspIds []int16
	ProvinceIds []int16
	CityIds []int16
	NameTagIds []int16
	GroupTagIds []int32
}

type TargetFilter struct {
	Name []string
	Host []string
	IspIds []int16
	ProvinceIds []int16
	CityIds []int16
	NameTagIds []int16
	GroupTagIds []int32
}

func NewCompoundQuery() *CompoundQuery {
	return &CompoundQuery {
		Output: &struct {
			Metrics []string
		} {},
		Grouping: &struct {
			Agent []string
			Target []string
		} {},
		Filters: &CompoundQueryFilter {
			Time: NewTimeFilter(),
			Agent: &AgentFilter{},
			Target: &TargetFilter{},
		},
	}
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

	jsonStartTime, hasStartTime := jsonTime.CheckGet("start_time")
	jsonEndTime, hasEndTime := jsonTime.CheckGet("end_time")
	jsonToNow, hasToNow := jsonTime.CheckGet("to_now")

	/**
	 * Parse time range of absolute
	 */
	if hasStartTime && hasEndTime {
		timeFilter.TimeRangeType = TIME_RANGE_ABSOLUTE
		timeFilter.StartTime = time.Unix(jsonStartTime.MustInt64(), 0)
		timeFilter.EndTime = time.Unix(jsonEndTime.MustInt64(), 0)
		return nil
	}
	// :~)

	if !hasToNow {
		return nil
	}

	/**
	 * Parse tiem range of relative
	 */
	timeFilter.TimeRangeType = TIME_RANGE_RELATIVE
	timeFilter.ToNow.Value = jsonToNow.Get("value").MustInt()

	stringOfTimeUnit := strings.ToLower(jsonToNow.GetPath("unit").MustString())
	valueOfTimeUnit, ok := supportingTimeUnit[stringOfTimeUnit]
	if !ok {
		return fmt.Errorf("Unknown time unit: [%s](\"filters.time.to_now.unit\").", stringOfTimeUnit)
	}

	timeFilter.ToNow.Unit = valueOfTimeUnit
	// :~)

	return nil
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
		return nil
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
		return nil
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
		return nil
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
