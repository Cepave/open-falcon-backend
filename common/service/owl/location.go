package owl

import (
	"fmt"
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

type ProvinceService struct {
	cache       *cache.CacheCtrl
	cacheConfig *cache.DataCacheConfig
}

func NewProvinceService(cacheConfig cache.DataCacheConfig) *ProvinceService {
	return &ProvinceService{
		cacheConfig: &cacheConfig,
		cache:       cache.NewCacheCtrl(cache.NewDataCache(cacheConfig)),
	}
}

// Loads province by id(or from cache)
//
// The id, which cannot be found, would not be put into cache.
func (s *ProvinceService) GetProvinceById(provinceId int16) *owlModel.Province {
	v := s.cache.MustFetchNativeAndDoNotCacheEmpty(
		provinceKeyById(provinceId),
		s.cacheConfig.Duration,
		func() interface{} {
			return owlDb.GetProvinceById(provinceId)
		},
	)

	if v == nil {
		return nil
	}

	return v.(*owlModel.Province)
}

func (s *ProvinceService) GetProvincesByIds(provinceIds ...int16) []*owlModel.Province {
	result := make([]*owlModel.Province, 0)

	for _, id := range provinceIds {
		result = append(result, s.GetProvinceById(id))
	}

	return result
}

// Loads province by name(or from cache)
//
// The name, which cannot find any data, would not be put into cache.
func (s *ProvinceService) GetProvincesByName(name string) []*owlModel.Province {
	v := s.cache.MustFetchNativeAndDoNotCacheEmpty(
		provinceKeyByName(name),
		s.cacheConfig.Duration,
		func() interface{} {
			return owlDb.GetProvincesByName(name)
		},
	)

	if v == nil {
		return nil
	}

	provinces := v.([]*owlModel.Province)

	/**
	 * Refresh cache indexed by id of province
	 */
	for _, province := range provinces {
		s.cache.Set(
			provinceKeyById(province.Id), province,
			s.cacheConfig.Duration,
		)
	}
	// :~)

	return provinces
}

func provinceKeyByName(name string) string {
	return fmt.Sprintf("!pname!%s", name)
}
func provinceKeyById(id int16) string {
	return fmt.Sprintf("!pid!%d", id)
}

type CityService struct {
	cache       *cache.CacheCtrl
	cacheConfig *cache.DataCacheConfig
}

func NewCityService(cacheConfig cache.DataCacheConfig) *CityService {
	return &CityService{
		cacheConfig: &cacheConfig,
		cache:       cache.NewCacheCtrl(cache.NewDataCache(cacheConfig)),
	}
}

// Loads city by id(or from cache)
//
// The id, which cannot be found, would not be put into cache.
func (s *CityService) GetCity2ById(cityId int16) *owlModel.City2 {
	v := s.cache.MustFetchNativeAndDoNotCacheEmpty(
		cityKeyById(cityId),
		s.cacheConfig.Duration,
		func() interface{} {
			return owlDb.GetCity2ById(cityId)
		},
	)

	if v == nil {
		return nil
	}

	return v.(*owlModel.City2)
}

func (s *CityService) GetCity2sByIds(cityIds ...int16) []*owlModel.City2 {
	result := make([]*owlModel.City2, 0)

	for _, id := range cityIds {
		result = append(result, s.GetCity2ById(id))
	}

	return result
}

// Loads city by name(or from cache)
//
// The name, which cannot find any data, would not be put into cache.
func (s *CityService) GetCity2sByName(name string) []*owlModel.City2 {
	v := s.cache.MustFetchNativeAndDoNotCacheEmpty(
		cityKeyByName(name),
		s.cacheConfig.Duration,
		func() interface{} {
			return owlDb.GetCity2sByName(name)
		},
	)

	if v == nil {
		return nil
	}

	cities := v.([]*owlModel.City2)

	/**
	 * Refresh cache indexed by id of city
	 */
	for _, city := range cities {
		s.cache.Set(
			cityKeyById(city.Id), city,
			s.cacheConfig.Duration,
		)
	}
	// :~)

	return cities
}

func cityKeyByName(name string) string {
	return fmt.Sprintf("!cname!%s", name)
}
func cityKeyById(id int16) string {
	return fmt.Sprintf("!cid!%d", id)
}
