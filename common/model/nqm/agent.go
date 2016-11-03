package nqm

import (
	"net"
	"time"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	json "github.com/bitly/go-simplejson"
	"strings"
)

type Agent struct {
	Id int32 `gorm:"primary_key:true;column:ag_id"`
	Name string `gorm:"column:ag_name"`
	ConnectionId string `gorm:"column:ag_connection_id"`
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
	jsonObject.Set("name", agentView.Name)
	jsonObject.Set("connection_id", agentView.ConnectionId)
	jsonObject.Set("hostname", agentView.Hostname)
	jsonObject.Set("ip_address", agentView.IpAddress)

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
