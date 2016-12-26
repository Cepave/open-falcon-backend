package nqm

import (
	"fmt"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	sjson "github.com/bitly/go-simplejson"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

type TimeFilterDetail TimeFilter
func (t *TimeFilterDetail) MarshalJSON() ([]byte, error) {
	json := sjson.New()

	switch t.timeRangeType {
	case TimeRangeAbsolute:
		json.Set("start_time", t.StartTime)
		json.Set("end_time", t.EndTime)
	case TimeRangeRelative:
		startTime, endTime := (*TimeFilter)(t).GetNetTimeRange()

		json.Set("start_time", ojson.JsonTime(startTime))
		json.Set("end_time", ojson.JsonTime(endTime))
		json.Set("to_now", t.ToNow)
	default:
		panic(fmt.Sprintf("Unknown type of time filter: %d", t.timeRangeType))
	}

	return json.MarshalJSON()
}

type CompoundQueryDetail struct {
	Time *TimeFilterDetail `json:"time"`
	Metrics ojson.JsonString `json:"metrics"`
	Agent *AgentOfQueryDetail `json:"agent"`
	Target *TargetOfQueryDetail `json:"target"`
	Output *OutputDetail `json:"output"`
}

type AgentOfQueryDetail struct {
	Name []string `json:"name"`
	Hostname []string `json:"hostname"`
	IpAddress []string `json:"ip_address"`
	ConnectionId []string `json:"connection_id"`

	Isps []*owlModel.Isp `json:"isps"`
	Provinces []*owlModel.Province `json:"provinces"`
	Cities []*owlModel.City2 `json:"cities"`

	NameTags []*owlModel.NameTag `json:"name_tags"`
	GroupTags []*owlModel.GroupTag `json:"group_tags"`
}
type TargetOfQueryDetail struct {
	Name []string `json:"name"`
	Host []string `json:"host"`

	Isps []*owlModel.Isp `json:"isps"`
	Provinces []*owlModel.Province `json:"provinces"`
	Cities []*owlModel.City2 `json:"cities"`

	NameTags []*owlModel.NameTag `json:"name_tags"`
	GroupTags []*owlModel.GroupTag `json:"group_tags"`
}

type OutputDetail struct {
	Agent []string `json:"agent"`
	Target []string `json:"target"`
	Metrics []string `json:"metrics"`
}
