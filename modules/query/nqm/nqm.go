package nqm

import (
	"fmt"
	dsl "github.com/Cepave/open-falcon-backend/modules/query/dsl/nqm_parser" // As NQM intermediate representation
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	"reflect"
	"time"

	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	owlService "github.com/Cepave/open-falcon-backend/common/service/owl"
	nqmService "github.com/Cepave/open-falcon-backend/common/service/nqm"
	"github.com/Cepave/open-falcon-backend/common/utils"
)

/**
 * 1. Main controller for NQM reporting
 * 2. Makes the unit test more easily(by replacing lambda used in method)
 */
type ServiceController struct {
	GetStatisticsOfIcmpByDsl func(*NqmDsl) ([]IcmpResult, error)
	GetProvinceById func(int16) *owlModel.Province
	GetIspById func(int16) *owlModel.Isp
	GetCityById func(int16) *owlModel.City2
	GetTargetById func(int32) *nqmModel.SimpleTarget1
}

func GetDefaultServiceController() *ServiceController {
	return &ServiceController{}
}
// :~)

var queryService *owlService.QueryService

var ispService *owlService.IspService
var provinceService *owlService.ProvinceService
var cityService *owlService.CityService
var groupTagService *owlService.GroupTagService
var nameTagService *owlService.NameTagService

var agentService *nqmService.AgentService
var targetService *nqmService.TargetService

// Initilaize the service
func (srv *ServiceController) Init() {
	initIcmp()
	initServices()

	srv.GetStatisticsOfIcmpByDsl = getStatisticsOfIcmpByDsl
	srv.GetProvinceById = provinceService.GetProvinceById
	srv.GetCityById = cityService.GetCity2ById
	srv.GetIspById = ispService.GetIspById
	srv.GetTargetById = targetService.GetSimpleTarget1ById
}

const queryNamedId = "nqm.compound.report"

func initServices() {
	queryService = owlService.NewQueryService(
		owlService.QueryServiceConfig {
			queryNamedId,
			8,
			time.Hour * 8,
		},
	)

	// Cache for ISPs: Maximum 16 entities with 6 hours live time
	ispService = owlService.NewIspService(cache.DataCacheConfig{
		MaxSize: 16, Duration: time.Hour * 8,
	})

	// Cache for Provinces: Maximum 16 entities with 6 hours live time
	provinceService = owlService.NewProvinceService(cache.DataCacheConfig{
		MaxSize: 16, Duration: time.Hour * 16,
	})

	// Cache for Cities: Maximum 32 entities with 6 hours live time
	cityService = owlService.NewCityService(cache.DataCacheConfig{
		MaxSize: 32, Duration: time.Hour * 16,
	})

	// Cache for Targets: Maximum 150 entities with 2 hours live time
	nameTagService = owlService.NewNameTagService(cache.DataCacheConfig{
		MaxSize: 32, Duration: time.Hour * 8,
	})

	groupTagService = owlService.NewGroupTagService(cache.DataCacheConfig{
		MaxSize: 8, Duration: time.Hour * 8,
	})

	// Cache for Targets: Maximum 150 entities with 2 hours live time
	agentService = nqmService.NewAgentService(cache.DataCacheConfig{
		MaxSize: 150, Duration: time.Hour * 2,
	})

	// Cache for Targets: Maximum 150 entities with 2 hours live time
	targetService = nqmService.NewTargetService(cache.DataCacheConfig{
		MaxSize: 150, Duration: time.Hour * 2,
	})
}

// Query data for provinces
func (srv *ServiceController) ListByProvinces(dslParams *dsl.QueryParams) []ProvinceMetric {
	/**
	 * 1. Set-up the grouping column
	 * 2. Only for inter-province
	 */
	nqmDsl := toNqmDsl(dslParams)
	nqmDsl.GroupingColumns = []string { "ag_pv_id" }
	nqmDsl.ProvinceRelation = model.SameValue
	// :~)

	/**
	 * Loads statistics of ICMP
	 */
	rawIcmpData, err := srv.GetStatisticsOfIcmpByDsl(nqmDsl)
	if err != nil {
		panic(err)
	}
	// :~)

	/**
	 * Join data of provinces
	 */
	result := make([]ProvinceMetric, len(rawIcmpData))
	for i, v := range rawIcmpData {
		currentMetric := ProvinceMetric{}

		currentMetric.Province = srv.GetProvinceById(int16(v.grouping[0]))
		currentMetric.Metrics = v.metrics

		result[i] = currentMetric
	}
	// :~)

	return result
}

// Query data for detail of city
func (srv *ServiceController) ListTargetsWithCityDetail(dslParams *dsl.QueryParams) []CityMetric {
	/**
	 * Loads data with grouping by id of cities
	 */
	dslGroupByCity := toNqmDsl(dslParams)
	dslGroupByCity.GroupingColumns = []string { "tg_ct_id" }
	dslGroupByCity.ProvinceRelation = model.SameValue
	rawIcmpGroupByCity, errForCityReport := srv.GetStatisticsOfIcmpByDsl(dslGroupByCity)
	if errForCityReport != nil {
		panic(errForCityReport)
	}
	// :~)

	result := make([]CityMetric, 0, len(rawIcmpGroupByCity))

	/**
	 * Initialize map for city metrics
	 */
	idToCityMetrics := make(map[int16]*CityMetric)

	for i, rowByCity := range rawIcmpGroupByCity {
		cityId := int16(rowByCity.grouping[0])

		result = append(
			result,
			CityMetric{
				City: srv.GetCityById(cityId),
				Metrics: rowByCity.metrics,
				Targets: make([]TargetMetric, 0),
			},
		)
		idToCityMetrics[cityId] = &result[i]
	}
	// :~)

	/**
	 * Loads data with grouping by id of targets
	 */
	dslGroupByTarget := toNqmDsl(dslParams)
	dslGroupByTarget.GroupingColumns = []string { "tg_id", "tg_ct_id", "tg_isp_id" }
	rawIcmpGroupByTarget, errForTargetReport := srv.GetStatisticsOfIcmpByDsl(dslGroupByTarget)
	if errForTargetReport != nil {
		panic(errForTargetReport)
	}
	// :~)

	/**
	 * Collects the list of targets into matched city record
	 */
	for _, rowByTarget := range rawIcmpGroupByTarget {
		targetId := rowByTarget.grouping[0]
		cityId := int16(rowByTarget.grouping[1])
		ispId := int16(rowByTarget.grouping[2])

		cityRow, assertExisting := idToCityMetrics[cityId]
		if !assertExisting {
			panic(fmt.Errorf("Cannot find city[id: %v] for target row[target id: %v]", cityId, targetId))
		}

		// Loads data of target
		targetNode := srv.GetTargetById(targetId)

		/**
		 * Appends the found target
		 */
		cityRow.Targets = append(
			cityRow.Targets,
			TargetMetric {
				Id: targetNode.Id,
				Host: targetNode.Host,
				Isp: srv.GetIspById(ispId),
				Metrics: rowByTarget.metrics,
			},
		)
		// :~)
	}
	// :~)

	return result
}

// Converts the IR of DSL to specific data for query on Cassandra
func toNqmDsl(queryParams *dsl.QueryParams) *NqmDsl {
	nqmDsl := &NqmDsl{
		StartTime: EpochTime(queryParams.StartTime.Unix()),
		EndTime: EpochTime(queryParams.EndTime.Unix()),

		IdsOfAgents: safeIds(queryParams.AgentFilterById.MatchIds),
		IdsOfTargets: safeIds(queryParams.TargetFilterById.MatchIds),

		IspRelation: queryParams.IspRelation,
		ProvinceRelation: queryParams.ProvinceRelation,
		CityRelation: queryParams.CityRelation,
	}

	/**
	 * Loads filters of properties on agent
	 */
	nqmDsl.IdsOfAgentIsps = loadIds(queryParams.AgentFilter.MatchIsps, getIdOfIspByName, queryParams.AgentFilterById.MatchIsps)
	nqmDsl.IdsOfAgentProvinces = loadIds(queryParams.AgentFilter.MatchProvinces, getIdOfProvinceByName, queryParams.AgentFilterById.MatchProvinces)
	nqmDsl.IdsOfAgentCities = loadIds(queryParams.AgentFilter.MatchCities, getIdOfCityByName, queryParams.AgentFilterById.MatchCities)
	// :~)

	/**
	 * Loads filters of properties on target
	 */
	nqmDsl.IdsOfTargetIsps = loadIds(queryParams.TargetFilter.MatchIsps, getIdOfIspByName, queryParams.TargetFilterById.MatchIsps)
	nqmDsl.IdsOfTargetProvinces = loadIds(queryParams.TargetFilter.MatchProvinces, getIdOfProvinceByName, queryParams.TargetFilterById.MatchProvinces)
	nqmDsl.IdsOfTargetCities = loadIds(queryParams.TargetFilter.MatchCities, getIdOfCityByName, queryParams.TargetFilterById.MatchCities)
	// :~)

	return nqmDsl
}

func safeIds(ids []int32) []int32 {
	if ids == nil {
		return make([]int32, 0)
	}

	return ids
}

func loadIds(
	queryNames []string, loadIdFunc getIdFunc,
	additionalIds []int16,
) []int16 {
	var allIds []int16

	/**
	 * Loads ids from text of search condition
	 */
	for _, name := range queryNames {
		allIds = append(allIds, loadIdFunc(name)...)
	}
	// :~)

	// Loads ids from explicit value of id
	allIds = append(allIds, additionalIds...)

	return utils.MakeAbstractArray(allIds).
		FilterWith(utils.NewUniqueFilter(utils.TypeOfInt16)).
		GetArray().([]int16)
}

/**
 * Quick method for loader
 */
type getIdFunc func(string) []int16

func getIdOfProvinceByName(name string) []int16 {
	return loadIdsOrUnknown(provinceService.GetProvincesByName(name))
}
func getIdOfIspByName(name string) []int16 {
	return loadIdsOrUnknown(ispService.GetIspsByName(name))
}
func getIdOfCityByName(name string) []int16 {
	return loadIdsOrUnknown(cityService.GetCity2sByName(name))
}

func loadIdsOrUnknown(idObjects interface{}) []int16 {
	valueOfObjects := reflect.ValueOf(idObjects)

	if valueOfObjects.Len() == 0 {
		return []int16{ UNKNOWN_ID_FOR_QUERY }
	}

	ids := make([]int16, valueOfObjects.Len())
	for i := 0; i < valueOfObjects.Len(); i++ {
		ids[i] = valueOfObjects.Index(i).Elem().
			FieldByName("Id").Interface().(int16)
	}

	return ids
}
