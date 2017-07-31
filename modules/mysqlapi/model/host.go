package model

import (
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

type HostsResult struct {
	Hostname string               `gorm:"column:hostname" json:"hostname" conform:"trim"`
	ID       int                  `gorm:"primary_key:true;column:id" json:"id"`
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

type HostgroupsResult struct {
	Name   string               `gorm:"column:grp_name" json:"name" conform:"trim"`
	ID     int                  `gorm:"primary_key:true;column:id" json:"id"`
	Groups []*owlModel.GroupTag `json:"plugins"`

	IdsOfGroups   string `gorm:"column:gt_ids"`
	NamesOfGroups string `gorm:"column:gt_names"`
}

func (HostgroupsResult) TableName() string {
	return "grp"
}

func (hg *HostgroupsResult) AfterLoad() {
	hg.Groups = owlModel.SplitToArrayOfGroupTags(
		hg.IdsOfGroups, ",",
		hg.NamesOfGroups, "\000",
	)
}
