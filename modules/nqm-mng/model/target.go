package model

import (
	"fmt"
	"strconv"
	"strings"
)

// Represents the data of target used by NQM agent
type NqmTarget struct {
	ID   int32  `gorm:"column:apl_tg_id" json:"id"`
	Host string `gorm:"column:tg_host" json:"host"`

	IspID   int16  `gorm:"column:isp_id" json:"isp_id"`
	IspName string `gorm:"column:isp_name" json:"isp_name"`

	ProvinceID   int16  `gorm:"column:pv_id" json:"province_id"`
	ProvinceName string `gorm:"column:pv_name" json:"province_name"`

	CityID   int16  `gorm:"column:ct_id" json:"ct_id"`
	CityName string `gorm:"column:ct_name" json:"ct_name"`

	NameTagID int16  `gorm:"column:nt_id" json:"nt_id"`
	NameTag   string `gorm:"column:nt_value" json:"nt_value"`

	GroupTags   string `gorm:"column:gts" json:"-"`
	GroupTagIDs []int8 `json:"gt_ids"`
}

func (NqmTarget) TableName() string {
	return "nqm_cache_agent_ping_list_log"
}

func (t *NqmTarget) AfterLoad() {
	if t.GroupTags == "" {
		return
	}
	strs := strings.Split(t.GroupTags, ",")
	for i, s := range strs {
		u, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(fmt.Errorf("Cannot parse value in array to Uint. Index: [%d]. Value: [%v]", i, s))
		}
		t.GroupTagIDs = append(t.GroupTagIDs, int8(u))
	}
}

func (t *NqmTarget) String() string {
	return fmt.Sprintf(
		"ID: [%d] Host: [%s] Isp: \"%s\"(%d) Province: \"%s\"(%d), City: \"%s\"[%d], Name tag: [%s](%d), Group Tags: %v",
		t.ID, t.Host,
		t.IspName, t.IspID,
		t.ProvinceName, t.ProvinceID,
		t.CityName, t.CityID,
		t.NameTag, t.NameTagID,
		t.GroupTagIDs,
	)
}
