package owl

import (
	"fmt"
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

type IspService struct {
	cache       *cache.CacheCtrl
	cacheConfig *cache.DataCacheConfig
}

func NewIspService(cacheConfig cache.DataCacheConfig) *IspService {
	return &IspService{
		cacheConfig: &cacheConfig,
		cache:       cache.NewCacheCtrl(cache.NewDataCache(cacheConfig)),
	}
}

// Loads isp by id(or from cache)
//
// The id, which cannot be found, would not be put into cache.
func (s *IspService) GetIspById(ispId int16) *owlModel.Isp {
	v := s.cache.MustFetchNativeAndDoNotCacheEmpty(
		ispKeyById(ispId),
		s.cacheConfig.Duration,
		func() interface{} {
			return owlDb.GetIspById(ispId)
		},
	)

	if v == nil {
		return nil
	}

	return v.(*owlModel.Isp)
}

// Loads isp by name(or from cache)
//
// The name, which cannot find any data, would not be put into cache.
func (s *IspService) GetIspsByName(name string) []*owlModel.Isp {
	v := s.cache.MustFetchNativeAndDoNotCacheEmpty(
		ispKeyByName(name),
		s.cacheConfig.Duration,
		func() interface{} {
			return owlDb.GetISPsByName(name)
		},
	)

	if v == nil {
		return nil
	}

	isps := v.([]*owlModel.Isp)

	/**
	 * Refresh cache indexed by id of isp
	 */
	for _, isp := range isps {
		s.cache.Set(
			ispKeyById(isp.Id), isp,
			s.cacheConfig.Duration,
		)
	}
	// :~)

	return isps
}

func (s *IspService) GetIspsByIds(ids ...int16) []*owlModel.Isp {
	result := make([]*owlModel.Isp, 0)

	for _, id := range ids {
		result = append(result, s.GetIspById(id))
	}

	return result
}

func ispKeyByName(name string) string {
	return fmt.Sprintf("!pname!%s", name)
}
func ispKeyById(id int16) string {
	return fmt.Sprintf("!pid!%d", id)
}
