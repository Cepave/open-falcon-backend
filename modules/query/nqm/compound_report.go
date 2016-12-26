package nqm

import (
	"fmt"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	metricDsl "github.com/Cepave/open-falcon-backend/modules/query/dsl/metric_parser"
	"github.com/satori/go.uuid"
)

func BuildQuery(q *model.CompoundQuery) *owlModel.Query {
	var digestValue [16]byte
	copy(digestValue[:], q.GetDigestValue())

	queryObject := &owlModel.Query {
		Content: q.GetCompressedQuery(),
		Md5Content: digestValue,
	}

	queryService.CreateOrLoadQuery(queryObject)
	return queryObject
}
func GetCompoundQueryByUuid(uuid uuid.UUID) *model.CompoundQuery {
	queryObject := queryService.LoadQueryByUuid(uuid)
	if queryObject == nil {
		return nil
	}

	compoundQuery := model.NewCompoundQuery()
	compoundQuery.UnmarshalFromCompressedQuery(queryObject.Content)

	filter, err := metricDsl.ParseToMetricFilter(compoundQuery.Filters.Metrics)
	if err != nil {
		panic(fmt.Errorf("Loads query by UUID has error on metric DSL: %v", err))
	}
	compoundQuery.SetMetricFilter(filter)

	return compoundQuery
}
func ToQueryDetail(q *model.CompoundQuery) *model.CompoundQueryDetail {
	agentFilter := q.Filters.Agent
	targetFilter := q.Filters.Target

	return &model.CompoundQueryDetail{
		Time: (*model.TimeFilterDetail)(q.Filters.Time),
		Metrics: ojson.JsonString(q.Filters.Metrics),
		Agent: &model.AgentOfQueryDetail {
			Name: agentFilter.Name,
			Hostname: agentFilter.Hostname,
			IpAddress: agentFilter.IpAddress,
			ConnectionId: agentFilter.ConnectionId,

			Isps: ispService.GetIspsByIds(agentFilter.IspIds...),
			Provinces: provinceService.GetProvincesByIds(agentFilter.ProvinceIds...),
			Cities: cityService.GetCity2sByIds(agentFilter.CityIds...),

			NameTags: nameTagService.GetNameTagsByIds(agentFilter.NameTagIds...),
			GroupTags: groupTagService.GetGroupTagsByIds(agentFilter.GroupTagIds...),
		},
		Target: &model.TargetOfQueryDetail {
			Name: targetFilter.Name,
			Host: targetFilter.Host,

			Isps: ispService.GetIspsByIds(targetFilter.IspIds...),
			Provinces: provinceService.GetProvincesByIds(targetFilter.ProvinceIds...),
			Cities: cityService.GetCity2sByIds(targetFilter.CityIds...),

			NameTags: nameTagService.GetNameTagsByIds(targetFilter.NameTagIds...),
			GroupTags: groupTagService.GetGroupTagsByIds(targetFilter.GroupTagIds...),
		},
		Output: &model.OutputDetail {
			Agent: q.Grouping.Agent,
			Target: q.Grouping.Target,
			Metrics: q.Output.Metrics,
		},
	}
}
