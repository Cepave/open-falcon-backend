package nqm

import (
	"fmt"
	cache "github.com/Cepave/open-falcon-backend/common/ccache"
	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
)

type AgentService struct {
	cache       *cache.CacheCtrl
	cacheConfig *cache.DataCacheConfig
}

func NewAgentService(cacheConfig cache.DataCacheConfig) *AgentService {
	return &AgentService{
		cacheConfig: &cacheConfig,
		cache:       cache.NewCacheCtrl(cache.NewDataCache(cacheConfig)),
	}
}

func (s *AgentService) GetSimpleAgent1ById(agentId int32) *nqmModel.SimpleAgent1 {
	v := s.cache.MustFetchNativeAndDoNotCacheEmpty(
		getKeyByAgentId(agentId),
		s.cacheConfig.Duration,
		func() interface{} {
			return nqmDb.GetSimpleAgent1ById(agentId)
		},
	)

	if v == nil {
		return nil
	}

	return v.(*nqmModel.SimpleAgent1)
}

// This function would cache agents by their id, the hit of database is always happened
func (s *AgentService) GetSimpleAgent1sByFilter(filter *nqmModel.AgentFilter) []*nqmModel.SimpleAgent1 {
	result := nqmDb.LoadSimpleAgent1sByFilter(filter)

	duration := s.cacheConfig.Duration
	for _, agent := range result {
		s.cache.Set(getKeyByAgentId(agent.Id), agent, duration)
	}

	return result
}

func getKeyByAgentId(agentId int32) string {
	return fmt.Sprintf("!aid!%d", agentId)
}
