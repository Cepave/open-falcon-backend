package nqm

import (
	"fmt"
	qcache "github.com/Cepave/open-falcon-backend/modules/query/cache"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/patrickmn/go-cache"
	"strconv"
	"time"
)

/**
 * Initializes the caching properties
 */
var (
	provinceCache qcache.CachePool
	ispCache      qcache.CachePool
	cityCache     qcache.CachePool
	targetCache   qcache.CachePool
)

/**
 * 1. Set-up ORM(on beego framework)
 */
func init() {
	orm.RegisterModel(
		new(Isp), new(Province), new(City),
		new(Target),
	)

	/**
	 * Caches data of provinces for 2 hours(10 minutes of interval for checking)
	 */
	provinceCache.Cache = cache.New(2*time.Hour, 10*time.Minute)
	// :~)

	/**
	 * Caches data of isps for 2 hours(10 minutes of interval for checking)
	 */
	ispCache.Cache = cache.New(2*time.Hour, 10*time.Minute)
	// :~)

	/**
	 * Caches data of isps for 4 hours(20 minutes of interval for checking)
	 */
	cityCache.Cache = cache.New(4*time.Hour, 20*time.Minute)
	// :~)

	/**
	 * Caches data of targets for 4 hours(20 minutes of interval for checking)
	 */
	targetCache.Cache = cache.New(4*time.Hour, 20*time.Minute)
	// :~)
}

// Lists agents in city by id of province
func ListAgentsInCityByProvinceId(provinceId int32) []*AgentsInCity {
	var rawResult []orm.Params

	/**
	 * Query data from database
	 */
	_, err := getOrmDb().Raw(
		`
		SELECT ag.ag_id, ag.ag_name, ag.ag_hostname, INET_NTOA(CONV(HEX(ag.ag_ip_address), 16, 10)) AS ag_ip_address,
			ct.ct_id, ct.ct_name, ct.ct_post_code
		FROM nqm_agent AS ag
			/**
			 * Only lists the agents(enabled) has ping task(enabled)
			 */
			INNER JOIN
			nqm_agent_ping_task AS apt
			ON ag.ag_id = apt.apt_ag_id
				AND ag.ag_status = TRUE
			INNER JOIN
			nqm_ping_task AS pt
			ON pt.pt_id = apt.apt_pt_id
				AND pt.pt_enable = TRUE
			# //:~)
			INNER JOIN
			owl_city AS ct
			ON ag.ag_ct_id = ct.ct_id
		WHERE ag.ag_pv_id = ?
		`,
		provinceId,
	).Values(&rawResult)
	if err != nil {
		log.Panicf("Query NQM agents by province[%v] has error: [%v]", provinceId, err)
	}
	// :~)

	/**
	 * Collects and grouping agents by city
	 */
	cityGrouping := make(map[int16]*AgentsInCity)

	for _, row := range rawResult {
		cityIdAsInt, _ := strconv.Atoi(row["ct_id"].(string))
		cityId := int16(cityIdAsInt)

		if _, hasCityId := cityGrouping[cityId]; !hasCityId {
			cityGrouping[cityId] = &AgentsInCity{
				City: &City{
					Id:       cityId,
					Name:     row["ct_name"].(string),
					PostCode: row["ct_post_code"].(string),
				},
				Agents: make([]Agent, 0),
			}
		}

		currentAgentsInCity := cityGrouping[cityId]

		agentIdAsInt, _ := strconv.Atoi(row["ag_id"].(string))
		currentAgentsInCity.Agents = append(
			currentAgentsInCity.Agents,
			Agent{
				Id:        int32(agentIdAsInt),
				Name:      row["ag_name"].(string),
				Hostname:  row["ag_hostname"].(string),
				IpAddress: row["ag_ip_address"].(string),
			},
		)
	}
	// :~)

	result := make([]*AgentsInCity, 0)
	for _, agentsInCity := range cityGrouping {
		result = append(result, agentsInCity)
	}

	return result
}

/**
 * The cache worker, which could be implemented by anonymous code
 */
type delegateCacheWorker struct {
	loadSourceData func() (interface{}, error)
	setCache       func(*cache.Cache, interface{})
}

func (worker delegateCacheWorker) LoadSourceData() (interface{}, error) {
	return worker.loadSourceData()
}
func (worker delegateCacheWorker) SetCache(cache *cache.Cache, object interface{}) {
	worker.setCache(cache, object)
}

// :~)

// Gets province by id
func getProvinceById(provinceId int16) *Province {
	cacheKey := (&Province{Id: provinceId}).getCacheKeyWithId()

	province, err := provinceCache.Get(
		cacheKey,
		delegateCacheWorker{
			loadSourceData: func() (interface{}, error) {
				return loadProvinceFromDbById(provinceId)
			},
			setCache: func(cachePool *cache.Cache, provinceFromDb interface{}) {
				cachePool.Set(cacheKey, provinceFromDb, cache.DefaultExpiration)
			},
		},
	)

	if err != nil {
		log.Panicf("Cannot load province. Id: %v. Error: %v.", provinceId, err)
	}

	return province.(*Province)
}

// Gets province by name(prefix)
func getProvinceByName(provinceName string) *Province {
	cacheKey := fmt.Sprintf("!name!%s", provinceName)

	province, err := provinceCache.Get(
		cacheKey,
		delegateCacheWorker{
			loadSourceData: func() (interface{}, error) {
				return loadProvinceFromDbByName(provinceName)
			},
			setCache: func(cachePool *cache.Cache, provinceFromDb interface{}) {
				// Refresh cache by id whenever the data is re-loaded
				cachePool.Set(provinceFromDb.(*Province).getCacheKeyWithId(), provinceFromDb, cache.DefaultExpiration)
				cachePool.Set(cacheKey, provinceFromDb, cache.DefaultExpiration)
				// :~)
			},
		},
	)

	if err != nil {
		log.Panicf("Cannot load province. Key: %v. Error: %v.", provinceName, err)
	}

	return province.(*Province)
}

// Gets ISP by id
func getIspById(ispId int16) *Isp {
	cacheKey := (&Isp{Id: ispId}).getCacheKeyWithId()

	isp, err := ispCache.Get(
		cacheKey,
		delegateCacheWorker{
			loadSourceData: func() (interface{}, error) {
				return loadIspFromDbById(ispId)
			},
			setCache: func(cachePool *cache.Cache, ispFromDb interface{}) {
				// Refresh cache by id whenever the data is re-loaded
				cachePool.Set(cacheKey, ispFromDb, cache.DefaultExpiration)
				// :~)
			},
		},
	)
	if err != nil {
		log.Panicf("Cannot load ISP. Id: %v. Error: %v.", ispId, err)
	}

	return isp.(*Isp)
}

// Gets ISP by name(prefix)
func getIspByName(ispName string) *Isp {
	cacheKey := fmt.Sprintf("!name!$s", ispName)

	isp, err := ispCache.Get(
		cacheKey,
		delegateCacheWorker{
			loadSourceData: func() (interface{}, error) {
				return loadIspFromDbByName(ispName)
			},
			setCache: func(cachePool *cache.Cache, ispFromDb interface{}) {
				// Refresh cache by id whenever the data is re-loaded
				cachePool.Set(ispFromDb.(*Isp).getCacheKeyWithId(), ispFromDb, cache.DefaultExpiration)
				cachePool.Set(cacheKey, ispFromDb, cache.DefaultExpiration)
				// :~)
			},
		},
	)
	if err != nil {
		log.Panicf("Cannot load ISP. Key: %v. Error: %v.", ispName, err)
	}

	return isp.(*Isp)
}

// Gets city by id
func getCityById(cityId int16) *City {
	cacheKey := (&City{Id: cityId}).getCacheKeyWithId()

	city, err := cityCache.Get(
		cacheKey,
		delegateCacheWorker{
			loadSourceData: func() (interface{}, error) {
				return loadCityFromDbById(cityId)
			},
			setCache: func(cachePool *cache.Cache, cityFromDb interface{}) {
				// Refresh cache by id whenever the data is re-loaded
				cachePool.Set(cacheKey, cityFromDb, cache.DefaultExpiration)
				// :~)
			},
		},
	)
	if err != nil {
		log.Panicf("Cannot load city. Id: %v. Error: %v.", cityId, err)
	}

	return city.(*City)
}

// Gets city by name(prefix)
func getCityByName(cityName string) *City {
	cacheKey := fmt.Sprintf("!name!$s", cityName)

	city, err := cityCache.Get(
		cacheKey,
		delegateCacheWorker{
			loadSourceData: func() (interface{}, error) {
				return loadCityFromDbByName(cityName)
			},
			setCache: func(cachePool *cache.Cache, cityFromDb interface{}) {
				// Refresh cache by id whenever the data is re-loaded
				cachePool.Set(cityFromDb.(*City).getCacheKeyWithId(), cityFromDb, cache.DefaultExpiration)
				cachePool.Set(cacheKey, cityFromDb, cache.DefaultExpiration)
				// :~)
			},
		},
	)
	if err != nil {
		log.Panicf("Cannot load city. Key: %v. Error: %v.", cityName, err)
	}

	return city.(*City)
}

// Gets target by id
func getTargetById(targetId int32) *Target {
	cacheKey := (&Target{Id: targetId}).getCacheKeyWithId()

	target, err := targetCache.Get(
		cacheKey,
		delegateCacheWorker{
			loadSourceData: func() (interface{}, error) {
				return loadTargetFromDbById(targetId)
			},
			setCache: func(cachePool *cache.Cache, targetFromDb interface{}) {
				// Refresh cache by id whenever the data is re-loaded
				cachePool.Set(cacheKey, targetFromDb, cache.DefaultExpiration)
				// :~)
			},
		},
	)
	if err != nil {
		log.Panicf("Cannot load target. Id: %v. Error: %v.", targetId, err)
	}

	return target.(*Target)
}

// Gets target by host
func getTargetByHost(targetHost string) *Target {
	cacheKey := fmt.Sprintf("!host!$s", targetHost)

	target, err := targetCache.Get(
		cacheKey,
		delegateCacheWorker{
			loadSourceData: func() (interface{}, error) {
				return loadTargetFromDbByHost(targetHost)
			},
			setCache: func(cachePool *cache.Cache, targetFromDb interface{}) {
				// Refresh cache by id whenever the data is re-loaded
				cachePool.Set(targetFromDb.(*Target).getCacheKeyWithId(), targetFromDb, cache.DefaultExpiration)
				cachePool.Set(cacheKey, targetFromDb, cache.DefaultExpiration)
				// :~)
			},
		},
	)
	if err != nil {
		log.Panicf("Cannot load target. Key: %v. Error: %v.", targetHost, err)
	}

	return target.(*Target)
}

func loadProvinceFromDbById(provinceId int16) (*Province, error) {
	var province Province

	loadedProvince, err := queryOneOrGetDefault(
		getOrmDb().QueryTable("owl_province").
			Filter("pv_id__exact", provinceId).
			Limit(1),
		&province,
		func() interface{} {
			province.Id = provinceId
			province.Name = UNKNOWN_NAME_FOR_QUERY
			return &province
		},
		nilProvince,
	)

	return loadedProvince.(*Province), err
}
func loadProvinceFromDbByName(provinceName string) (*Province, error) {
	var province Province

	loadedProvince, err := queryOneOrGetDefault(
		getOrmDb().QueryTable("owl_province").
			Filter("pv_name__istartswith", provinceName).
			Limit(1),
		&province,
		func() interface{} {
			province.Id = UNKNOWN_ID_FOR_QUERY
			province.Name = provinceName
			return &province
		},
		nilProvince,
	)

	return loadedProvince.(*Province), err
}

func loadIspFromDbById(ispId int16) (*Isp, error) {
	var isp Isp

	loadedIsp, err := queryOneOrGetDefault(
		getOrmDb().QueryTable("owl_isp").
			Filter("isp_id__exact", ispId).
			Limit(1),
		&isp,
		func() interface{} {
			isp.Id = ispId
			isp.Name = UNKNOWN_NAME_FOR_QUERY
			return &isp
		},
		nilIsp,
	)

	return loadedIsp.(*Isp), err
}
func loadIspFromDbByName(ispName string) (*Isp, error) {
	var isp Isp

	loadedIsp, err := queryOneOrGetDefault(
		getOrmDb().QueryTable("owl_isp").
			Filter("isp_name__istartswith", ispName).
			Limit(1),
		&isp,
		func() interface{} {
			isp.Id = UNKNOWN_ID_FOR_QUERY
			isp.Name = ispName
			return &isp
		},
		nilIsp,
	)

	return loadedIsp.(*Isp), err
}

func loadCityFromDbById(cityId int16) (*City, error) {
	var city City

	loadedCity, err := queryOneOrGetDefault(
		getOrmDb().QueryTable("owl_city").
			Filter("ct_id__exact", cityId).
			Limit(1),
		&city,
		func() interface{} {
			city.Id = cityId
			city.Name = UNKNOWN_NAME_FOR_QUERY
			city.PostCode = UNKNOWN_NAME_FOR_QUERY
			return &city
		},
		nilCity,
	)

	return loadedCity.(*City), err
}
func loadCityFromDbByName(cityName string) (*City, error) {
	var city City

	loadedCity, err := queryOneOrGetDefault(
		getOrmDb().QueryTable("owl_city").
			Filter("ct_name__istartswith", cityName).
			Limit(1),
		&city,
		func() interface{} {
			city.Id = UNKNOWN_ID_FOR_QUERY
			city.Name = cityName
			city.PostCode = UNKNOWN_NAME_FOR_QUERY
			return &city
		},
		nilCity,
	)

	return loadedCity.(*City), err
}

func loadTargetFromDbById(targetId int32) (*Target, error) {
	var target Target

	loadedTarget, err := queryOneOrGetDefault(
		getOrmDb().QueryTable("nqm_target").
			Filter("tg_id__exact", targetId).
			Limit(1),
		&target,
		func() interface{} {
			target.Id = targetId
			target.Host = UNKNOWN_NAME_FOR_QUERY
			return &target
		},
		nilTarget,
	)

	return loadedTarget.(*Target), err
}
func loadTargetFromDbByHost(targetHost string) (*Target, error) {
	var target Target

	loadedTarget, err := queryOneOrGetDefault(
		getOrmDb().QueryTable("nqm_target").
			Filter("tg_host__exact", targetHost).
			Limit(1),
		&target,
		func() interface{} {
			target.Id = UNKNOWN_ID_FOR_QUERY
			target.Host = targetHost
			return &target
		},
		nilTarget,
	)

	return loadedTarget.(*Target), err
}

func queryOneOrGetDefault(
	query orm.QuerySeter,
	resultHolder interface{}, defaultFunc func() interface{},
	nilValue interface{},
) (interface{}, error) {
	err := query.One(resultHolder)

	switch err {
	case nil:
		break
	case orm.ErrNoRows:
		return defaultFunc(), nil
	default:
		return nilValue, err
	}

	return resultHolder, nil
}

var ormDb orm.Ormer = nil

func getOrmDb() orm.Ormer {
	if ormDb == nil {
		ormDb = orm.NewOrm()
		ormDb.Using("default")
	}

	return ormDb
}
