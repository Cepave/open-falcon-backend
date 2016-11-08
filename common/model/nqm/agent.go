package nqm

import (
	"fmt"
	"net"
	"time"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	json "github.com/bitly/go-simplejson"
	"strings"
)

type AgentForAdding struct {
	Id int32 `json:"-"`
	Name string `json:"name" conform:"trim"`
	ConnectionId string `json:"connection_id" conform:"trim" validate:"min=1"`
	Comment string `json:"comment" conform:"trim"`
	Status bool `json:"status"`

	Hostname string `json:"hostname" conform:"trim" validate:"min=1"`
	IpAddress net.IP `json:"ip_address"`

	IspId int16 `json:"isp_id" validate:"nonZeroId"`
	ProvinceId int16 `json:"province_id" validate:"nonZeroId"`
	CityId int16 `json:"city_id" validate:"nonZeroId"`

	NameTagId int16 `json:"-"`
	NameTagValue string `json:"name_tag" conform:"trim"`

	GroupTags []string `json:"group_tags" conform:"trim"`
}
func NewAgentForAdding() *AgentForAdding {
	return &AgentForAdding {
		Status: true,
		Hostname: "0.0.0.0",
		IpAddress: net.ParseIP("0.0.0.0").To4(),
		IspId: -1,
		ProvinceId: -1,
		CityId: -1,
		NameTagId: -1,
	}
}
func (agent *AgentForAdding) GetIpAddressAsBytes() []byte {
	return ([]byte)(agent.IpAddress.To4())
}
func (agent *AgentForAdding) GetIpAddressAsString() string {
	return agent.IpAddress.String()
}

type Agent struct {
	Id int32 `gorm:"primary_key:true;column:ag_id"`
	Name string `gorm:"column:ag_name"`
	ConnectionId string `gorm:"column:ag_connection_id""`
	Hostname string `gorm:"column:ag_hostname"`

	IpAddress net.IP `gorm:"column:ag_ip_address"`

	Status bool `gorm:"column:ag_status"`
	Comment string `gorm:"column:ag_comment"`
	LastHeartBeat time.Time `gorm:"column:ag_last_heartbeat"`

	IspId int16 `gorm:"column:isp_id"`
	IspName string `gorm:"column:isp_name"`

	ProvinceId int16 `gorm:"column:pv_id"`
	ProvinceName string `gorm:"column:pv_name"`

	CityId int16 `gorm:"column:ct_id"`
	CityName string `gorm:"column:ct_name"`

	NameTagId int16 `gorm:"column:nt_id"`
	NameTagValue string `gorm:"column:nt_value"`

	IdsOfGroupTags string `gorm:"column:gt_ids"`
	NamesOfGroupTags string `gorm:"column:gt_names"`
	GroupTags []*owlModel.GroupTag
}
func (Agent) TableName() string {
	return "nqm_agent"
}
func (agentView *Agent) MarshalJSON() ([]byte, error) {
	jsonObject := json.New()

	jsonObject.Set("id", agentView.Id)
	jsonObject.Set("connection_id", agentView.ConnectionId)
	jsonObject.Set("hostname", agentView.Hostname)
	jsonObject.Set("ip_address", agentView.IpAddress)
	jsonObject.Set("status", agentView.Status)
	if agentView.Name != "" {
		jsonObject.Set("name", agentView.Name)
	} else {
		jsonObject.Set("name", nil)
	}
	if agentView.Comment != "" {
		jsonObject.Set("comment", agentView.Comment)
	} else {
		jsonObject.Set("comment", nil)
	}

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

	groupTags := make([]*json.Json, 0, len(agentView.GroupTags))
	for _, groupTag := range agentView.GroupTags {
		jsonGroupTag := json.New()
		jsonGroupTag.Set("id", groupTag.Id)
		jsonGroupTag.Set("name", groupTag.Name)

		groupTags = append(groupTags, jsonGroupTag)
	}
	jsonObject.Set("group_tags", groupTags)

	return jsonObject.MarshalJSON()
}
func (agentView *Agent) AfterLoad() {
	if agentView.IdsOfGroupTags == "" {
		return
	}

	allIds := commonDb.GroupedPlainStringToUintArray(agentView.IdsOfGroupTags, ",")
	allNames := strings.Split(agentView.NamesOfGroupTags, "\000")

	for i, groupTagId := range allIds {
		agentView.GroupTags = append(
			agentView.GroupTags,
			&owlModel.GroupTag {
				Id: int32(groupTagId),
				Name: allNames[i],
			},
		)
	}
}
func (agent *Agent) String() string {
	return fmt.Sprintf(
		"Id: [%d]. Name: [%s]. Connection Id: [%s]. IpAddress: [%s]. Hostname: [%s]. Status: [%v]",
		agent.Id, agent.Name, agent.ConnectionId, agent.IpAddress, agent.Hostname, agent.Status,
	)
}
