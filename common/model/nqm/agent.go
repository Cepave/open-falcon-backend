package nqm

import (
	"fmt"
	"net"
	"time"

	owlGin "github.com/Cepave/open-falcon-backend/common/gin"
	owlJson "github.com/Cepave/open-falcon-backend/common/json"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	"github.com/Cepave/open-falcon-backend/common/utils"
	json "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
)

type AgentForAdding struct {
	Id           int32   `json:"-"`
	Name         *string `json:"name" conform:"trimToNil"`
	Comment      *string `json:"comment" conform:"trimToNil"`
	ConnectionId string  `json:"connection_id" conform:"trim" validate:"min=1"`
	Status       bool    `json:"status"`

	Hostname  string `json:"hostname" conform:"trim" validate:"min=1"`
	IpAddress net.IP `json:"-"`

	IspId      int16 `json:"isp_id" validate:"nonZeroId"`
	ProvinceId int16 `json:"province_id" validate:"nonZeroId"`
	CityId     int16 `json:"city_id" validate:"nonZeroId"`

	NameTagId    int16   `json:"-"`
	NameTagValue *string `json:"name_tag" conform:"trim"`

	GroupTags []string `json:"group_tags" conform:"trim"`
}

func NewAgentForAdding() *AgentForAdding {
	return &AgentForAdding{
		Status:     true,
		Hostname:   "",
		IpAddress:  net.ParseIP("0.0.0.0").To4(),
		IspId:      -1,
		ProvinceId: -1,
		CityId:     -1,
		NameTagId:  -1,
	}
}

func (agent *AgentForAdding) AreGroupTagsSame(anotherAgent *AgentForAdding) bool {
	return utils.AreArrayOfStringsSame(
		agent.GroupTags, anotherAgent.GroupTags,
	)
}

func (agent *AgentForAdding) UniqueGroupTags() {
	agent.GroupTags = utils.UniqueArrayOfStrings(agent.GroupTags)
}

func (agent *AgentForAdding) GetIpAddressAsBytes() []byte {
	return ([]byte)(agent.IpAddress.To4())
}

func (agent *AgentForAdding) GetIpAddressAsString() string {
	return agent.IpAddress.String()
}

type Agent struct {
	Id                    int32   `gorm:"primary_key:true;column:ag_id"`
	Name                  *string `gorm:"column:ag_name"`
	ConnectionId          string  `gorm:"column:ag_connection_id"`
	Hostname              string  `gorm:"column:ag_hostname"`
	NumOfEnabledPingtasks int32   `gorm:"column:ag_num_of_enabled_pingtasks"`

	IpAddress net.IP `gorm:"column:ag_ip_address"`

	Status        bool             `gorm:"column:ag_status"`
	Comment       *string          `gorm:"column:ag_comment"`
	LastHeartBeat owlJson.JsonTime `gorm:"column:ag_last_heartbeat"`

	IspId   int16  `gorm:"column:isp_id"`
	IspName string `gorm:"column:isp_name"`

	ProvinceId   int16  `gorm:"column:pv_id"`
	ProvinceName string `gorm:"column:pv_name"`

	CityId   int16  `gorm:"column:ct_id"`
	CityName string `gorm:"column:ct_name"`

	NameTagId    int16  `gorm:"column:nt_id"`
	NameTagValue string `gorm:"column:nt_value"`

	IdsOfGroupTags   string `gorm:"column:gt_ids"`
	NamesOfGroupTags string `gorm:"column:gt_names"`
	GroupTags        []*owlModel.GroupTag
}

func (Agent) TableName() string {
	return "nqm_agent"
}

func (agentView *Agent) ToSimpleJson() *json.Json {
	jsonObject := json.New()

	jsonObject.Set("id", agentView.Id)
	jsonObject.Set("connection_id", agentView.ConnectionId)
	jsonObject.Set("hostname", agentView.Hostname)
	jsonObject.Set("ip_address", agentView.IpAddress)
	jsonObject.Set("status", agentView.Status)

	jsonObject.Set("last_heartbeat_time", owlJson.JsonTime(agentView.LastHeartBeat))
	jsonObject.Set("name", agentView.Name)
	jsonObject.Set("comment", agentView.Comment)
	jsonObject.Set("num_of_enabled_pingtasks", agentView.NumOfEnabledPingtasks)

	jsonIsp := json.New()
	jsonIsp.Set("id", agentView.IspId)
	jsonIsp.Set("name", agentView.IspName)
	jsonObject.Set("isp", jsonIsp)

	jsonProvince := json.New()
	jsonProvince.Set("id", agentView.ProvinceId)
	jsonProvince.Set("name", agentView.ProvinceName)
	jsonObject.Set("province", jsonProvince)

	jsonCity := json.New()
	jsonCity.Set("id", agentView.CityId)
	jsonCity.Set("name", agentView.CityName)
	jsonObject.Set("city", jsonCity)

	jsonNameTag := json.New()
	jsonNameTag.Set("id", agentView.NameTagId)
	jsonNameTag.Set("value", agentView.NameTagValue)
	jsonObject.Set("name_tag", jsonNameTag)

	jsonGroupTags := owlModel.GroupTags(agentView.GroupTags).ToJson()
	jsonObject.Set("group_tags", jsonGroupTags)

	return jsonObject
}

func (agentView *Agent) ToAgentForAdding() *AgentForAdding {
	groupTags := make([]string, 0, len(agentView.GroupTags))
	for _, groupTag := range agentView.GroupTags {
		groupTags = append(groupTags, groupTag.Name)
	}

	nameTagValue := agentView.NameTagValue

	return &AgentForAdding{
		Id:      agentView.Id,
		Name:    agentView.Name,
		Comment: agentView.Comment,
		Status:  agentView.Status,

		ConnectionId: agentView.ConnectionId,
		Hostname:     agentView.Hostname,
		IpAddress:    agentView.IpAddress,

		IspId:      agentView.IspId,
		ProvinceId: agentView.ProvinceId,
		CityId:     agentView.CityId,

		NameTagId:    agentView.NameTagId,
		NameTagValue: &nameTagValue,
		GroupTags:    owlModel.GroupTags(agentView.GroupTags).ToNames(),
	}
}

func (agentView *Agent) MarshalJSON() ([]byte, error) {
	return agentView.ToSimpleJson().MarshalJSON()
}

func (agentView *Agent) AfterLoad() {
	agentView.GroupTags = owlModel.SplitToArrayOfGroupTags(
		agentView.IdsOfGroupTags, ",",
		agentView.NamesOfGroupTags, "\000",
	)
}

func (agent *Agent) String() string {
	return fmt.Sprintf(
		"Id: [%d]. Name: [%s]. Connection Id: [%s]. IpAddress: [%s]. Hostname: [%s]. Status: [%v]",
		agent.Id, agent.Name, agent.ConnectionId, agent.IpAddress, agent.Hostname, agent.Status,
	)
}

type AgentWithPingTask struct {
	Agent
	ApplyingPingTask bool `gorm:"column:applying_ping_task"`
}

func (a *AgentWithPingTask) MarshalJSON() ([]byte, error) {
	json := a.ToSimpleJson()
	json.Set("applying_ping_task", a.ApplyingPingTask)

	return json.MarshalJSON()
}

type SimpleAgent1 struct {
	Id        int32   `json:"id" db:"ag_id"`
	Name      *string `json:"name" db:"ag_name"`
	Hostname  string  `json:"hostname" db:"ag_hostname"`
	IpAddress net.IP  `json:"ip_address" db:"ag_ip_address"`

	IspId      int16 `json:"isp_id" db:"ag_isp_id"`
	ProvinceId int16 `json:"province_id" db:"ag_pv_id"`
	CityId     int16 `json:"city_id" db:"ag_ct_id"`
	NameTagId  int16 `json:"name_tag_id" db:"ag_nt_id"`
}

type SimpleAgent1s []*SimpleAgent1

func (as SimpleAgent1s) GetInt32s() []int32 {
	result := make([]int32, 0)

	for _, a := range as {
		result = append(result, a.Id)
	}

	return result
}

type SimpleAgent1InCity struct {
	City   *owlModel.City2 `json:"city"`
	Agents []*SimpleAgent1 `json:"agents"`
}

type PingListLog struct {
	NumberOfTargets int32     `db:"apll_number_of_targets"`
	AccessTime      time.Time `db:"apll_time_access"`
	RefreshTime     time.Time `db:"apll_time_refresh"`
}

type CacheAgentPingListLog struct {
	ApllAgId          int32     `gorm:"primary_key:true;column:apll_ag_id"`
	CachedRefreshTime time.Time `gorm:"column:apll_time_refresh"`
}

func (CacheAgentPingListLog) TableName() string {
	return "nqm_cache_agent_ping_list_log"
}

type AgentPingListTarget struct {
	Target
	AgentPingListTimeAccess time.Time        `gorm:"column:tg_apl_time_access"`
	ProbedTime              owlJson.JsonTime //`json:"probed_time"`
}

func (target *AgentPingListTarget) MarshalJSON() ([]byte, error) {
	jsonObject := json.New()

	jsonObject.Set("id", target.Id)
	jsonObject.Set("name", target.Name)
	jsonObject.Set("host", target.Host)
	jsonObject.Set("probed_by_all", target.ProbedByAll)
	jsonObject.Set("status", target.Status)
	jsonObject.Set("available", target.Available)
	jsonObject.Set("creation_time", target.CreationTime.Unix())
	jsonObject.Set("comment", target.Comment)
	jsonObject.Set("probed_time", target.ProbedTime)

	jsonIsp := json.New()
	jsonIsp.Set("id", target.IspId)
	jsonIsp.Set("name", target.IspName)
	jsonObject.Set("isp", jsonIsp)

	jsonProvince := json.New()
	jsonProvince.Set("id", target.ProvinceId)
	jsonProvince.Set("name", target.ProvinceName)
	jsonObject.Set("province", jsonProvince)

	jsonCity := json.New()
	jsonCity.Set("id", target.CityId)
	jsonCity.Set("name", target.CityName)
	jsonObject.Set("city", jsonCity)

	jsonNameTag := json.New()
	jsonNameTag.Set("id", target.NameTagId)
	jsonNameTag.Set("value", target.NameTagValue)
	jsonObject.Set("name_tag", jsonNameTag)

	jsonGroupTags := owlModel.GroupTags(target.GroupTags).ToJson()
	jsonObject.Set("group_tags", jsonGroupTags)

	return jsonObject.MarshalJSON()
}

type TargetsOfAgent struct {
	CacheRefreshTime owlJson.JsonTime       `json:"cache_refresh_time"`
	CacheLifeTime    int                    `json:"cache_lifetime"`
	Targets          []*AgentPingListTarget `json:"targets"`
}

type ClearCacheView struct {
	RowsAffected int8 `json:"rows_affected"`
}

func (l *PingListLog) GetDurationOfLastAccess(checkedTime time.Time) int64 {
	return int64(checkedTime.Sub(l.AccessTime) / time.Minute)
}

type HeartbeatRequest struct {
	// The connection id of the NQM agent
	ConnectionId string `json:"connection_id" binding:"required"`
	// The hostname of the machine running NQM agent
	Hostname string `json:"hostname" binding:"required"`
	// The IP address of the NQM agent. It supports both IPv4 and IPv6 format
	IpAddress owlJson.IP `json:"ip_address" binding:"required"`
	// The timestamp generated by hbs
	Timestamp owlJson.JsonTime `json:"timestamp" binding:"required"`
}

func (r *HeartbeatRequest) Bind(c *gin.Context) {
	owlGin.BindJson(c, r)
}

func (r HeartbeatRequest) String() string {
	return fmt.Sprintf(
		"HeartbeatRequest{ConnectionId=%v Hostname=%v IpAddress=%v Timestamp=%v}",
		r.ConnectionId, r.Hostname, r.IpAddress, r.Timestamp,
	)
}

type AgentView struct {
	Id                    int32            `json:"id"`
	Name                  string           `json:"name"`
	ConnectionId          string           `json:"connection_id"`
	Hostname              string           `json:"hostname"`
	IpAddress             net.IP           `json:"ip_address"`
	Status                bool             `json:"status"`
	Comment               *string          `json:"comment"`
	LastHeartBeat         owlJson.JsonTime `json:"last_heartbeat_time"`
	NumOfEnabledPingtasks int32            `json:"num_of_enabled_pingtasks"`

	ISP       *ISP      `json:"isp"`
	Province  *Province `json:"province"`
	City      *City     `json:"city"`
	NameTag   *NameTag  `json:"name_tag"`
	GroupTags []int32   `json:"group_tags"`
}

func (a AgentView) String() string {
	return fmt.Sprintf(
		"Id: [%d]. Name: [%s]. Connection Id: [%s]. IpAddress: [%s]. Hostname: [%s]. Status: [%v]",
		a.Id, a.Name, a.ConnectionId, a.IpAddress, a.Hostname, a.Status,
	)
}

type ISP struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type Province struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type NameTag struct {
	ID    int16  `json:"id"`
	Value string `json:"value"`
}
