package nqm

import (
	"fmt"
	"time"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	json "github.com/bitly/go-simplejson"
	"github.com/Cepave/open-falcon-backend/common/utils"
)

type TargetForAdding struct {
	Id int32 `json:"-"`
	Name string `json:"name" conform:"trim" validate:"min=1"`
	Host string `json:"host" conform:"trim" validate:"min=1"`
	ProbedByAll bool `json:"probed_by_all"`
	Status bool `json:"status"`
	Comment string `json:"comment" conform:"trim"`

	IspId int16 `json:"isp_id" validate:"nonZeroId"`
	ProvinceId int16 `json:"province_id" validate:"nonZeroId"`
	CityId int16 `json:"city_id" validate:"nonZeroId"`

	NameTagId int16 `json:"-"`
	NameTagValue string `json:"name_tag" conform:"trim"`
	GroupTags []string `json:"group_tags" conform:"trim"`
}
func (target *TargetForAdding) AreGroupTagsSame(anotherTarget *TargetForAdding) bool {
	return utils.AreArrayOfStringsSame(
		target.GroupTags, anotherTarget.GroupTags,
	)
}
func (target *TargetForAdding) UniqueGroupTags() {
	target.GroupTags = utils.UniqueArrayOfStrings(target.GroupTags)
}

func NewTargetForAdding() *TargetForAdding {
	return &TargetForAdding {
		Status: true,
		ProbedByAll: false,
		IspId: -1,
		ProvinceId: -1,
		CityId: -1,
		NameTagId: -1,
	}
}

type Target struct {
	Id int32 `gorm:"primary_key:true;column:tg_id"`
	Name string `gorm:"column:tg_name"`
	Host string `gorm:"column:tg_host"`
	ProbedByAll bool `gorm:"column:tg_probed_by_all"`
	Status bool `gorm:"column:tg_status"`
	Available bool `gorm:"column:tg_available"`
	Comment string `gorm:"column:tg_comment"`
	CreationTime *time.Time `gorm:"column:tg_created_ts"`

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
func (Target) TableName() string {
	return "nqm_target"
}
func (target *Target) AfterLoad() {
	target.GroupTags = owlModel.SplitToArrayOfGroupTags(
		target.IdsOfGroupTags, ",",
		target.NamesOfGroupTags, "\000",
	)
}
func (target *Target) MarshalJSON() ([]byte, error) {
	jsonObject := json.New()

	jsonObject.Set("id", target.Id)
	jsonObject.Set("name", target.Name)
	jsonObject.Set("host", target.Host)
	jsonObject.Set("probed_by_all", target.ProbedByAll)
	jsonObject.Set("status", target.Status)
	jsonObject.Set("available", target.Available)
	jsonObject.Set("creation_time", target.CreationTime.Unix())

	if target.Comment != "" {
		jsonObject.Set("comment", target.Comment)
	} else {
		jsonObject.Set("comment", nil)
	}

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
func (target *Target) String() string {
	return fmt.Sprintf(
		"Id[%d] Name: [%s]. Host: [%s]. Status: [%v]. Available: [%v]. Probed By All: [%v]",
		target.Id, target.Name, target.Host,
		target.Status, target.Available, target.ProbedByAll,
	)
}
func (target *Target) ToTargetForAdding() *TargetForAdding {
	return &TargetForAdding {
		Id: target.Id,
		Name: target.Name,
		Host: target.Host,
		ProbedByAll: target.ProbedByAll,
		Status: target.Status,
		Comment: target.Comment,

		IspId: target.IspId,
		ProvinceId: target.ProvinceId,
		CityId: target.CityId,

		NameTagId: target.NameTagId,
		GroupTags: owlModel.GroupTags(target.GroupTags).ToNames(),
	}
}

type SimpleTarget1 struct {
	Id   int32  `json:"id" db:"tg_id"`
	Name string `json:"name" db:"tg_name"`
	Host string `json:"host" db:"tg_host"`

	IspId int16 `json:"isp_id" db:"tg_isp_id"`
	ProvinceId int16 `json:"province_id" db:"tg_pv_id"`
	CityId int16 `json:"city_id" db:"tg_ct_id"`
	NameTagId int16 `json:"name_tag_id" db:"tg_nt_id"`
}

type SimpleTarget1s []*SimpleTarget1
func (as SimpleTarget1s) GetInt32s() []int32 {
	result := make([]int32, 0)

	for _, a := range as {
		result = append(result, a.Id)
	}

	return result
}
