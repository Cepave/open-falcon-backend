package nqm

import (
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
)

/**
 * Aliases of type for DSL
 */
type EpochTime int64
// :~)

// This value is used to indicate the non-existing id for data(province, city, or ISP)
// Instead of -1(e.x. A agent doesn't has information of ISP), this value is used in query.
const UNKNOWN_ID_FOR_QUERY = -2
const UNKNOWN_NAME_FOR_QUERY = "<UNKNOWN>"

// Represents the DSL for query over Icmp log
type NqmDsl struct {
	GroupingColumns []string `json:"grouping_columns"`

    StartTime EpochTime `json:"start_time"`
	EndTime EpochTime `json:"end_time"`

	IdsOfAgents []int32 `json:"ids_of_agents"`
	IdsOfAgentIsps []int16 `json:"ids_of_agent_isps"`
	IdsOfAgentProvinces []int16 `json:"ids_of_agent_provinces"`
	IdsOfAgentCities []int16 `json:"ids_of_agent_cities"`
	IdsOfAgentNameTags []int16 `json:"ids_of_agent_name_tags"`
	IdsOfAgentGroupTags []int32 `json:"ids_of_agent_group_tags"`

	IdsOfTargets []int32 `json:"ids_of_targets"`
	IdsOfTargetProvinces []int16 `json:"ids_of_target_provinces"`
	IdsOfTargetIsps []int16 `json:"ids_of_target_isps"`
	IdsOfTargetCities []int16 `json:"ids_of_target_cities"`
	IdsOfTargetNameTags []int16 `json:"ids_of_target_name_tags"`
	IdsOfTargetGroupTags []int32 `json:"ids_of_target_group_tags"`

	IspRelation model.PropRelation `json:"isp_relation"`
	ProvinceRelation model.PropRelation `json:"province_relation"`
	CityRelation model.PropRelation `json:"city_relation"`
	NameTagRelation model.PropRelation `json:"name_tag_relation"`
}

// The data used for reporting of ICMP statistics(grouping by provinces of agents)
type ProvinceMetric struct {
	Province *owlModel.Province `json:"province"`
	Metrics *model.Metrics `json:"metrics"`
}

// The data used for reporting of ICMP statistics, which contains detail of target node(grouping by city)
type CityMetric struct {
	City *owlModel.City2 `json:"city"`
	Metrics *model.Metrics `json:"metrics"`
	Targets []TargetMetric `json:"targets"`
}

// The data used for reporting of ICMP statistics target node
type TargetMetric struct {
	Id int32 `json:"id"`
	Host string `json:"host"`
	Isp *owlModel.Isp `json:"isp"`
	Metrics *model.Metrics `json:"metrics"`
}
