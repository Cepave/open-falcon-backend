package nqm

import (
	"fmt"
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
)

type TargetService struct {
	cache       *cache.CacheCtrl
	cacheConfig *cache.DataCacheConfig
}

func NewTargetService(cacheConfig cache.DataCacheConfig) *TargetService {
	return &TargetService{
		cacheConfig: &cacheConfig,
		cache:       cache.NewCacheCtrl(cache.NewDataCache(cacheConfig)),
	}
}

func (s *TargetService) GetSimpleTarget1ById(targetId int32) *nqmModel.SimpleTarget1 {
	v := s.cache.MustFetchNativeAndDoNotCacheEmpty(
		getKeyByTargetId(targetId),
		s.cacheConfig.Duration,
		func() interface{} {
			return nqmDb.GetSimpleTarget1ById(targetId)
		},
	)

	if v == nil {
		return nil
	}

	return v.(*nqmModel.SimpleTarget1)
}

// This function would cache targets by their id, the hit of database is always happened
func (s *TargetService) GetSimpleTarget1sByFilter(filter *nqmModel.TargetFilter) []*nqmModel.SimpleTarget1 {
	result := nqmDb.LoadSimpleTarget1sByFilter(filter)

	duration := s.cacheConfig.Duration
	for _, target := range result {
		s.cache.Set(getKeyByTargetId(target.Id), target, duration)
	}

	return result
}

func getKeyByTargetId(targetId int32) string {
	return fmt.Sprintf("!tid!%d", targetId)
}
