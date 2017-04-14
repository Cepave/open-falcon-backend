package nqm

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
)

// The query conditions of agent
type AgentQuery struct {
	Name         string `mvc:"query[name]" conform:"trim"`
	ConnectionId string `mvc:"query[connection_id]" conform:"trim"`
	Hostname     string `mvc:"query[hostname]" conform:"trim"`
	IpAddress    string `mvc:"query[ip_address]" conform:"trim"`

	IspId         int16 `mvc:"query[isp_id]"`
	HasIspIdParam bool  `mvc:"query[?isp_id]"`

	Status         bool `mvc:"query[status]"`
	HasStatusParam bool `mvc:"query[?status]"`
}

// Gets the []byte used to perform like in MySql
func (query *AgentQuery) GetIpForLikeCondition() []byte {
	bytes, err := commonDb.IpV4ToBytesForLike(query.IpAddress)
	if err != nil {
		panic(err)
	}

	return bytes
}

type AgentQueryWithPingTask struct {
	AgentQuery
	PingTaskId int32  `mvc:"param[pingtask_id]"`
	Applied    bool   `mvc:"query[applied]"`
	HasApplied string `mvc:"query[applied] default[!N!]"`
}

func (query *AgentQueryWithPingTask) HasAppliedCondition() bool {
	return query.HasApplied != "!N!"
}

type TargetsOfAgentQuery struct {
	AgentID     int32 `mvc:"param[agent_id]"`
	TargetQuery *TargetQuery
}

// The query conditions of target
type TargetQuery struct {
	Name string `mvc:"query[name]"`
	Host string `mvc:"query[host]"`

	IspId         int16 `mvc:"query[isp_id]"`
	HasIspIdParam bool  `mvc:"query[?isp_id]"`

	Status         bool `mvc:"query[status]"`
	HasStatusParam bool `mvc:"query[?status]"`
}

type AgentFilter struct {
	Name         []string `json:"name" digest:"1"`
	Hostname     []string `json:"hostname" digest:"2"`
	IpAddress    []string `json:"ip_address" digest:"3"`
	ConnectionId []string `json:"connection_id" digest:"4"`
	IspIds       []int16  `json:"isp_ids" digest:"21"`
	ProvinceIds  []int16  `json:"province_ids" digest:"22"`
	CityIds      []int16  `json:"city_ids" digest:"23"`
	NameTagIds   []int16  `json:"name_tag_ids" digest:"24"`
	GroupTagIds  []int32  `json:"group_tag_ids" digest:"25"`
}

func (f *AgentFilter) HasAgentDescriptive() bool {
	return len(f.Name)+len(f.Hostname)+
		len(f.IpAddress)+len(f.ConnectionId) > 0
}

type TargetFilter struct {
	Name        []string `json:"name" digest:"1"`
	Host        []string `json:"host" digest:"2"`
	IspIds      []int16  `json:"isp_ids" digest:"21"`
	ProvinceIds []int16  `json:"province_ids" digest:"22"`
	CityIds     []int16  `json:"city_ids" digest:"23"`
	NameTagIds  []int16  `json:"name_tag_ids" digest:"24"`
	GroupTagIds []int32  `json:"group_tag_ids" digest:"25"`
}

func (f *TargetFilter) HasTargetDescriptive() bool {
	return len(f.Name)+len(f.Host) > 0
}

// The query parameters filtering pingtasks
type PingtaskQuery struct {
	Period             string `mvc:"query[period]"`
	Name               string `mvc:"query[name]"`
	Enable             string `mvc:"query[enable]"`
	Comment            string `mvc:"query[comment]"`
	NumOfEnabledAgents string `mvc:"query[num_of_enabled_agents]"`
}
