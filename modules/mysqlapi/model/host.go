package model

import (
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
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
