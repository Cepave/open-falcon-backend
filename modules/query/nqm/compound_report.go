package nqm

import (
	"fmt"

	"github.com/satori/go.uuid"

	"github.com/Cepave/open-falcon-backend/common/utils"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"

	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
	metricDsl "github.com/Cepave/open-falcon-backend/modules/query/dsl/metric_parser"
)

// Loads data by compound query.
//
// This function do not filter the data with metric DSL.
func LoadIcmpRecordsOfCompoundQuery(q *model.CompoundQuery) []*model.DynamicRecord {
	result := make([]*model.DynamicRecord, 0)

	/**
	 * Set-up DSL by compound query
	 */
	dsl := buildNqmDslByCompoundQuery(q)
	dsl.GroupingColumns = buildGroupingColumnOfDsl(q.Grouping)
	// :~)

	/**
	 * Loads data from data store of ICMP logs
	 */
	icmpLogs, err := getStatisticsOfIcmpByDsl(dsl)
	if err != nil {
		panic(err)
	}
	// :~)

	for _, icmpLog := range icmpLogs {
		newRecord := &model.DynamicRecord{
			Agent: &model.DynamicAgentProps{},
			Target: &model.DynamicTargetProps{},
			Metrics: &model.DynamicMetrics {
				Metrics: icmpLog.metrics,
				Output: &q.Output.Metrics,
			},
		}

		for i, column := range dsl.GroupingColumns {
			currentId := icmpLog.grouping[i]

			switch column {
			/**
			 * Grouping by every single node
			 */
			case "ag_id":
				processAgentGrouping(newRecord.Agent, icmpLog.grouping[i], q)
			case "tg_id":
				processTargetGrouping(newRecord.Target, icmpLog.grouping[i], q)
			// :~)
			/**
			 * Grouping by node's property
			 */
			case "ag_isp_id":
				newRecord.Agent.Isp = ispService.GetIspById(int16(currentId))
			case "ag_pv_id":
				newRecord.Agent.Province = provinceService.GetProvinceById(int16(currentId))
			case "ag_ct_id":
				newRecord.Agent.City = cityService.GetCity2ById(int16(currentId))
			case "ag_nt_id":
				newRecord.Agent.NameTag = nameTagService.GetNameTagById(int16(currentId))
			case "tg_isp_id":
				newRecord.Target.Isp = ispService.GetIspById(int16(currentId))
			case "tg_pv_id":
				newRecord.Target.Province = provinceService.GetProvinceById(int16(currentId))
			case "tg_ct_id":
				newRecord.Target.City = cityService.GetCity2ById(int16(currentId))
			case "tg_nt_id":
				newRecord.Target.NameTag = nameTagService.GetNameTagById(int16(currentId))
			}
			// :~)
		}

		if len(q.Grouping.Agent) == 0 {
			newRecord.Agent = nil
		}
		if len(q.Grouping.Target) == 0 {
			newRecord.Target = nil
		}

		result = append(result, newRecord)
	}

	return result
}

func BuildQuery(q *model.CompoundQuery) *owlModel.Query {
	var digestValue [16]byte
	copy(digestValue[:], q.GetDigestValue())

	queryObject := &owlModel.Query {
		Content: q.GetCompressedQuery(),
		Md5Content: digestValue,
	}

	queryService.CreateOrLoadQuery(queryObject)
	return queryObject
}
func GetCompoundQueryByUuid(uuid uuid.UUID) *model.CompoundQuery {
	queryObject := queryService.LoadQueryByUuid(uuid)
	if queryObject == nil {
		return nil
	}

	compoundQuery := model.NewCompoundQuery()
	compoundQuery.UnmarshalFromCompressedQuery(queryObject.Content)

	filter, err := metricDsl.ParseToMetricFilter(compoundQuery.Filters.Metrics)
	if err != nil {
		panic(fmt.Errorf("Loads query by UUID has error on metric DSL: %v", err))
	}
	compoundQuery.SetMetricFilter(filter)

	return compoundQuery
}
func ToQueryDetail(q *model.CompoundQuery) *model.CompoundQueryDetail {
	agentFilter := q.Filters.Agent
	targetFilter := q.Filters.Target

	return &model.CompoundQueryDetail{
		Time: (*model.TimeFilterDetail)(q.Filters.Time),
		Metrics: ojson.JsonString(q.Filters.Metrics),
		Agent: &model.AgentOfQueryDetail {
			Name: agentFilter.Name,
			Hostname: agentFilter.Hostname,
			IpAddress: agentFilter.IpAddress,
			ConnectionId: agentFilter.ConnectionId,

			Isps: ispService.GetIspsByIds(agentFilter.IspIds...),
			Provinces: provinceService.GetProvincesByIds(agentFilter.ProvinceIds...),
			Cities: cityService.GetCity2sByIds(agentFilter.CityIds...),

			NameTags: nameTagService.GetNameTagsByIds(agentFilter.NameTagIds...),
			GroupTags: groupTagService.GetGroupTagsByIds(agentFilter.GroupTagIds...),
		},
		Target: &model.TargetOfQueryDetail {
			Name: targetFilter.Name,
			Host: targetFilter.Host,

			Isps: ispService.GetIspsByIds(targetFilter.IspIds...),
			Provinces: provinceService.GetProvincesByIds(targetFilter.ProvinceIds...),
			Cities: cityService.GetCity2sByIds(targetFilter.CityIds...),

			NameTags: nameTagService.GetNameTagsByIds(targetFilter.NameTagIds...),
			GroupTags: groupTagService.GetGroupTagsByIds(targetFilter.GroupTagIds...),
		},
		Output: &model.OutputDetail {
			Agent: q.Grouping.Agent,
			Target: q.Grouping.Target,
			Metrics: q.Output.Metrics,
		},
	}
}

func buildNqmDslByCompoundQuery(q *model.CompoundQuery) *NqmDsl {
	filters := q.Filters
	startTime, endTime := filters.Time.GetNetTimeRange()

	loadAgentIdsFunc := func() nqmModel.Int32sGetter {
		return nqmModel.SimpleAgent1s(
			nqmDb.LoadSimpleAgent1sByFilter(filters.Agent),
		)
	}
	loadTargetIdsFunc := func() nqmModel.Int32sGetter {
		return nqmModel.SimpleTarget1s(
			nqmDb.LoadSimpleTarget1sByFilter(filters.Target),
		)
	}

	return &NqmDsl {
		GroupingColumns: buildGroupingColumnOfDsl(q.Grouping),

		StartTime: EpochTime(startTime.Unix()),
		EndTime: EpochTime(endTime.Unix()),

		IdsOfAgents: loadInt32Ids(filters.Agent.HasAgentDescriptive(), loadAgentIdsFunc),
		IdsOfAgentIsps: filterRelationIdsOnInt16(filters.Agent.IspIds),
		IdsOfAgentProvinces: filterRelationIdsOnInt16(filters.Agent.ProvinceIds),
		IdsOfAgentCities: filterRelationIdsOnInt16(filters.Agent.CityIds),
		IdsOfAgentNameTags: filterRelationIdsOnInt16(filters.Agent.NameTagIds),
		IdsOfAgentGroupTags: filterRelationIdsOnInt32(filters.Agent.GroupTagIds),

		IdsOfTargets: loadInt32Ids(filters.Target.HasTargetDescriptive(), loadTargetIdsFunc),
		IdsOfTargetIsps: filterRelationIdsOnInt16(filters.Target.IspIds),
		IdsOfTargetProvinces: filterRelationIdsOnInt16(filters.Target.ProvinceIds),
		IdsOfTargetCities: filterRelationIdsOnInt16(filters.Target.CityIds),
		IdsOfTargetNameTags: filterRelationIdsOnInt16(filters.Target.NameTagIds),
		IdsOfTargetGroupTags: filterRelationIdsOnInt32(filters.Target.GroupTagIds),

		IspRelation: q.GetIspRelation(),
		ProvinceRelation: q.GetProvinceRelation(),
		CityRelation: q.GetCityRelation(),
		NameTagRelation: q.GetNameTagRelation(),
	}
}

func loadInt32Ids(
	hasCondition bool,
	getterFunc func() nqmModel.Int32sGetter,
) []int32 {
	if !hasCondition {
		return []int32{}
	}

	result := getterFunc().GetInt32s()
	if len(result) == 0 {
		return []int32{ -2 }
	}

	return result
}

var groupingMappingOfAgent = map[string]string {
	model.GroupingProvince: "ag_pv_id",
	model.GroupingCity: "ag_ct_id",
	model.GroupingIsp: "ag_isp_id",
	model.GroupingNameTag: "ag_nt_id",
}
var groupingMappingOfTarget = map[string]string {
	model.GroupingProvince: "tg_pv_id",
	model.GroupingCity: "tg_ct_id",
	model.GroupingIsp: "tg_isp_id",
	model.GroupingNameTag: "tg_nt_id",
}
func buildGroupingColumnOfDsl(grouping *model.QueryGrouping) []string {
	groupingColumns := make([]string, 0)

	if grouping.IsForEachAgent() {
		groupingColumns = append(groupingColumns, "ag_id")
	} else {
		for _, groupColumn := range grouping.Agent {
			columnOfDsl, ok := groupingMappingOfAgent[groupColumn]
			if !ok {
				panic(fmt.Sprintf("Unsupported grouping on agent: [%s]", groupColumn))
			}

			groupingColumns = append(groupingColumns, columnOfDsl)
		}
	}

	if grouping.IsForEachTarget() {
		groupingColumns = append(groupingColumns, "tg_id")
	} else {
		for _, groupColumn := range grouping.Target {
			columnOfDsl, ok := groupingMappingOfTarget[groupColumn]
			if !ok {
				panic(fmt.Sprintf("Unsupported grouping on target: [%s]", groupColumn))
			}

			groupingColumns = append(groupingColumns, columnOfDsl)
		}
	}

	return groupingColumns
}

func filterRelationIdsOnInt16(v []int16) []int16 {
	return utils.MakeAbstractArray(v).
		FilterWith(func(v interface{}) bool {
			int16v := v.(int16)
			return int16v != model.RelationSame &&
				int16v != model.RelationNotSame
		}).
		GetArray().([]int16)
}
func filterRelationIdsOnInt32(v []int32) []int32 {
	return utils.MakeAbstractArray(v).
		FilterWith(func(v interface{}) bool {
			int32v := v.(int32)
			return int32v != model.RelationSame &&
				int32v != model.RelationNotSame
		}).
		GetArray().([]int32)
}

func processAgentGrouping(agentProps *model.DynamicAgentProps, agentId int32, query *model.CompoundQuery) {
	agentDetail := agentService.GetSimpleAgent1ById(agentId)
	if agentDetail == nil {
		panic(fmt.Sprintf("Cannot find detail of agent by id: [%d]", agentId))
	}

	for _, selectedPropOfAgent := range query.Grouping.Agent {
		agentProps.Id = agentDetail.Id
		switch selectedPropOfAgent {
			case model.AgentGroupingName:
				agentProps.Name = agentDetail.Name
			case model.AgentGroupingHostname:
				agentProps.Hostname = agentDetail.Hostname
			case model.AgentGroupingIpAddress:
				agentProps.IpAddress = agentDetail.IpAddress.String()
			case model.GroupingIsp:
				agentProps.Isp = ispService.GetIspById(agentDetail.IspId)
			case model.GroupingProvince:
				agentProps.Province = provinceService.GetProvinceById(agentDetail.ProvinceId)
			case model.GroupingCity:
				agentProps.City = cityService.GetCity2ById(agentDetail.CityId)
			case model.GroupingNameTag:
				agentProps.NameTag = nameTagService.GetNameTagById(agentDetail.NameTagId)
			default:
				panic(fmt.Sprintf("Unsupported grouping for agent: [%s]", selectedPropOfAgent))
		}
	}
}
func processTargetGrouping(targetProps *model.DynamicTargetProps, targetId int32, query *model.CompoundQuery) {
	targetDetail := targetService.GetSimpleTarget1ById(targetId)
	if targetDetail == nil {
		panic(fmt.Sprintf("Cannot find detail of target by id: [%d]", targetId))
	}

	for _, selectedPropOfTarget := range query.Grouping.Target {
		targetProps.Id = targetDetail.Id

		switch selectedPropOfTarget {
			case model.TargetGroupingName:
				targetProps.Name = targetDetail.Name
			case model.TargetGroupingHost:
				targetProps.Host = targetDetail.Host
			case model.GroupingIsp:
				targetProps.Isp = ispService.GetIspById(targetDetail.IspId)
			case model.GroupingProvince:
				targetProps.Province = provinceService.GetProvinceById(targetDetail.ProvinceId)
			case model.GroupingCity:
				targetProps.City = cityService.GetCity2ById(targetDetail.CityId)
			case model.GroupingNameTag:
				targetProps.NameTag = nameTagService.GetNameTagById(targetDetail.NameTagId)
			default:
				panic(fmt.Sprintf("Unsupported grouping for target: [%s]", selectedPropOfTarget))
		}
	}
}
