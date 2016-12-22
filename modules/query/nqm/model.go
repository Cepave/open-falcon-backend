package nqm

import (
	dsl "github.com/Cepave/open-falcon-backend/modules/query/dsl/nqm_parser" // As NQM intermediate representation
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

/**
 * Macro-struct re-used by various data
 */
type Metrics struct {
	Max                     int16   `json:"max"`
	Min                     int16   `json:"min"`
	Avg                     float64 `json:"avg"`
	Med                     int16   `json:"med"`
	Mdev                    float64 `json:"mdev"`
	Loss                    float64 `json:"loss"`
	Count                   int32   `json:"count"`
	NumberOfSentPackets     uint64  `json:"number_of_sent_packets"`
	NumberOfReceivedPackets uint64  `json:"number_of_received_packets"`
	NumberOfAgents          int32   `json:"number_of_agents"`
	NumberOfTargets         int32   `json:"number_of_targets"`
}

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

	IdsOfTargets []int32 `json:"ids_of_targets"`
	IdsOfTargetProvinces []int16 `json:"ids_of_target_provinces"`
	IdsOfTargetIsps []int16 `json:"ids_of_target_isps"`
	IdsOfTargetCities []int16 `json:"ids_of_target_cities"`

	IspRelation dsl.HostRelation `json:"isp_relation"`
	ProvinceRelation dsl.HostRelation `json:"province_relation"`
	CityRelation dsl.HostRelation `json:"city_relation"`
}

// The data used for reporting of ICMP statistics(grouping by provinces of agents)
type ProvinceMetric struct {
	Province *owlModel.Province `json:"province"`
	Metrics *Metrics `json:"metrics"`
}

// The data used for reporting of ICMP statistics, which contains detail of target node(grouping by city)
type CityMetric struct {
	City *owlModel.City2 `json:"city"`
	Metrics *Metrics `json:"metrics"`
	Targets []TargetMetric `json:"targets"`
}

// The data used for reporting of ICMP statistics target node
type TargetMetric struct {
	Id int32 `json:"id"`
	Host string `json:"host"`
	Isp *owlModel.Isp `json:"isp"`
	Metrics *Metrics `json:"metrics"`
}
