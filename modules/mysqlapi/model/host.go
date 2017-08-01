package model

import (
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	json "github.com/bitly/go-simplejson"
)

type HostsResult struct {
	ID       int                  `gorm:"primary_key:true;column:id" json:"id"`
	Hostname string               `gorm:"column:hostname" json:"hostname" conform:"trim"`
	Groups   []*owlModel.GroupTag `json:"groups"`

	IdsOfGroups   string `gorm:"column:gt_ids"`
	NamesOfGroups string `gorm:"column:gt_names"`
}

func (HostsResult) TableName() string {
	return "host"
}

func (host *HostsResult) AfterLoad() {
	host.Groups = owlModel.SplitToArrayOfGroupTags(
		host.IdsOfGroups, ",",
		host.NamesOfGroups, "\000",
	)
}

type PluginTag struct {
	ID  int32  `json:"id"`
	Dir string `json:"dir"`
}

func (pluginTag *PluginTag) ToJson() *json.Json {
	jsonPluginTag := json.New()
	jsonPluginTag.Set("id", pluginTag.ID)
	jsonPluginTag.Set("dir", pluginTag.Dir)

	return jsonPluginTag
}

type PluginTags []*PluginTag

func (pluginTags PluginTags) ToJson() []*json.Json {
	jsonPluginTags := make([]*json.Json, len(pluginTags))
	for idx, tag := range pluginTags {
		jsonPluginTags[idx] = tag.ToJson()
	}

	return jsonPluginTags
}

type HostgroupsResult struct {
	ID      int          `gorm:"primary_key:true;column:id" json:"id"`
	Name    string       `gorm:"column:grp_name" json:"name" conform:"trim"`
	Plugins []*PluginTag `json:"plugins"`

	IdsOfGroups   string `gorm:"column:gt_ids"`
	NamesOfGroups string `gorm:"column:gt_names"`
}

func (HostgroupsResult) TableName() string {
	return "grp"
}

func (hg *HostgroupsResult) ToSimpleJson() *json.Json {
	jsonObject := json.New()
	jsonObject.Set("id", hg.ID)
	jsonObject.Set("name", hg.Name)

	jsonPluginTags := PluginTags(hg.Plugins).ToJson()
	jsonObject.Set("plugins", jsonPluginTags)

	return jsonObject
}

func (hg *HostgroupsResult) MarshalJSON() ([]byte, error) {
	return hg.ToSimpleJson().MarshalJSON()
}

func (hg *HostgroupsResult) AfterLoad() {
	groups := owlModel.SplitToArrayOfGroupTags(
		hg.IdsOfGroups, ",",
		hg.NamesOfGroups, "\000",
	)
	hg.Plugins = make([]*PluginTag, len(groups))
	for idx, val := range groups {
		hg.Plugins[idx] = &PluginTag{
			ID:  val.Id,
			Dir: val.Name,
		}
	}
}
