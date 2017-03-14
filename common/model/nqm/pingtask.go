package nqm

import (
	"fmt"
	"strings"

	owlGin "github.com/Cepave/open-falcon-backend/common/gin"
	commonOwlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	"github.com/chyeh/cast"
	"gopkg.in/gin-gonic/gin.v1"
)

type AgentPingtask struct {
	AgentID    int32
	PingtaskID int32
}

type pingtaskFilter struct {
	IspFilters      []*commonOwlModel.IspOfPingtaskView      `json:"isps"`
	ProvinceFilters []*commonOwlModel.ProvinceOfPingtaskView `json:"provinces"`
	CityFilters     []*commonOwlModel.CityOfPingtaskView     `json:"cities"`
	NameTagFilters  []*commonOwlModel.NameTagOfPingtaskView  `json:"name_tags"`
	GroupTagFilters []*commonOwlModel.GroupTagOfPingtaskView `json:"group_tags"`
}

type pingtaskModifyFilter struct {
	IspIds      []int8  `json:"ids_of_isp"`
	ProvinceIds []int8  `json:"ids_of_province"`
	CityIds     []int8  `json:"ids_of_city"`
	NameTagIds  []int8  `json:"ids_of_name_tag"`
	GroupTagIds []int16 `json:"ids_of_group_tag"`
}

type PingtaskView struct {
	ID                 int32   `gorm:"primary_key:true;column:pt_id" json:"id"`
	Period             int8    `gorm:"column:pt_period" json:"period"`
	Name               *string `gorm:"column:pt_name" json:"name"`
	Enable             bool    `gorm:"column:pt_enable" json:"enable"`
	Comment            *string `gorm:"column:pt_comment" json:"comment"`
	NumOfEnabledAgents int8    `gorm:"column:pt_num_of_enabled_agents" json:"num_of_enabled_agents"`

	IdsOfIspFilters  string `gorm:"column:pt_isp_filter_ids" json:"-"`
	NamesOfIspFilter string `gorm:"column:pt_isp_filter_names" json:"-"`

	IdsOfProvinceFilters  string `gorm:"column:pt_province_filter_ids" json:"-"`
	NamesOfProvinceFilter string `gorm:"column:pt_province_filter_names" json:"-"`

	IdsOfCityFilters  string `gorm:"column:pt_city_filter_ids" json:"-"`
	NamesOfCityFilter string `gorm:"column:pt_city_filter_names" json:"-"`

	IdsOfNameTagFilters  string `gorm:"column:pt_name_tag_filter_ids" json:"-"`
	NamesOfNameTagFilter string `gorm:"column:pt_name_tag_filter_values" json:"-"`

	IdsOfGroupTagFilters  string `gorm:"column:pt_group_tag_filter_ids" json:"-"`
	NamesOfGroupTagFilter string `gorm:"column:pt_group_tag_filter_names" json:"-"`

	Filter pingtaskFilter `json:"filter"`
}

func (PingtaskView) TableName() string {
	return "nqm_ping_task"
}

func (p *PingtaskView) AfterLoad() {

	var ids []string
	var names []string
	p.Filter.IspFilters = make([]*commonOwlModel.IspOfPingtaskView, 0)
	p.Filter.ProvinceFilters = make([]*commonOwlModel.ProvinceOfPingtaskView, 0)
	p.Filter.CityFilters = make([]*commonOwlModel.CityOfPingtaskView, 0)
	p.Filter.NameTagFilters = make([]*commonOwlModel.NameTagOfPingtaskView, 0)
	p.Filter.GroupTagFilters = make([]*commonOwlModel.GroupTagOfPingtaskView, 0)

	if p.IdsOfIspFilters != "" {
		ids = strings.Split(p.IdsOfIspFilters, ",")
		names = strings.Split(p.NamesOfIspFilter, "\000")
		if len(ids) != len(names) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.IspFilters = append(
				p.Filter.IspFilters,
				&commonOwlModel.IspOfPingtaskView{
					Id:   cast.ToInt(ids[i]),
					Name: names[i],
				},
			)
		}
	}

	if p.IdsOfProvinceFilters != "" {
		ids = strings.Split(p.IdsOfProvinceFilters, ",")
		names = strings.Split(p.NamesOfProvinceFilter, "\000")
		if len(ids) != len(names) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.ProvinceFilters = append(
				p.Filter.ProvinceFilters,
				&commonOwlModel.ProvinceOfPingtaskView{
					Id:   cast.ToInt(ids[i]),
					Name: names[i],
				},
			)
		}
	}

	if p.IdsOfCityFilters != "" {
		ids = strings.Split(p.IdsOfCityFilters, ",")
		names = strings.Split(p.NamesOfCityFilter, "\000")
		if len(ids) != len(names) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.CityFilters = append(
				p.Filter.CityFilters,
				&commonOwlModel.CityOfPingtaskView{
					Id:   cast.ToInt(ids[i]),
					Name: names[i],
				},
			)
		}
	}

	if p.IdsOfNameTagFilters != "" {
		ids = strings.Split(p.IdsOfNameTagFilters, ",")
		values := strings.Split(p.NamesOfNameTagFilter, "\000")
		if len(ids) != len(values) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.NameTagFilters = append(
				p.Filter.NameTagFilters,
				&commonOwlModel.NameTagOfPingtaskView{
					Id:    cast.ToInt(ids[i]),
					Value: values[i],
				},
			)
		}
	}

	if p.IdsOfGroupTagFilters != "" {
		ids = strings.Split(p.IdsOfGroupTagFilters, ",")
		names = strings.Split(p.NamesOfGroupTagFilter, "\000")
		if len(ids) != len(names) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.GroupTagFilters = append(
				p.Filter.GroupTagFilters,
				&commonOwlModel.GroupTagOfPingtaskView{
					Id:   cast.ToInt(ids[i]),
					Name: names[i],
				},
			)
		}
	}
}

type PingtaskModify struct {
	Period  int32                `json:"period"`
	Name    string               `json:"name" conform:"trim"`
	Enable  bool                 `json:"enable"`
	Comment string               `json:"comment" conform:"trim"`
	Filter  pingtaskModifyFilter `json:"filter"`
}

func (p *PingtaskModify) Bind(c *gin.Context) {
	owlGin.BindJson(c, p)
}
