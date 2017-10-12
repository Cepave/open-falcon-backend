package owl

import (
	"time"

	owlc "github.com/Cepave/open-falcon-backend/common/ccache"
	model "github.com/Cepave/open-falcon-backend/common/model/owl"
	owlDb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/owl"
	c "github.com/karlseguin/ccache"
	"github.com/satori/go.uuid"
)

type QueryObjectServiceConfig struct {
	CacheSize     int64
	CacheDuration time.Duration
}

var QueryObjectService *queryObjectService

func InitQueryObjectService(config QueryObjectServiceConfig) {
	QueryObjectService = newQueryObjectService(&config)
}
func StopQueryObjectService() {
	if QueryObjectService != nil {
		QueryObjectService.cache.Stop()
	}

	QueryObjectService = nil
}

type queryObjectService struct {
	config *QueryObjectServiceConfig
	cache  *c.Cache
}

func newQueryObjectService(config *QueryObjectServiceConfig) *queryObjectService {
	cacheObject := owlc.NewDataCache(
		owlc.DataCacheConfig{
			MaxSize:  config.CacheSize,
			Duration: config.CacheDuration,
		},
	)

	return &queryObjectService{
		config, cacheObject,
	}
}

func (s *queryObjectService) LoadQueryByUuid(uuid uuid.UUID) *model.Query {
	now := time.Now()
	stringUuid := uuid.String()

	/**
	 * Updates access time to database if cache has this query
	 *
	 * If the query object is existing in cache, but it is not existing in database.
	 * This function would re-insert the data with same uuid.
	 */
	itemInCache := s.cache.Get(stringUuid)
	if itemInCache != nil {
		query := itemInCache.Value().(*model.Query)
		owlDb.UpdateAccessTimeOrAddNewOne(query, now)

		return query
	}
	// :~)

	/**
	 * Loads query from database
	 */
	queryFromDatabase := owlDb.LoadQueryByUuidAndUpdateAccessTime(
		uuid, now,
	)
	if queryFromDatabase == nil {
		return nil
	}
	// :~)

	s.cache.Set(stringUuid, queryFromDatabase, s.config.CacheDuration)
	return s.cache.Get(stringUuid).Value().(*model.Query)
}

// Loads object of query or creating one
func (s *queryObjectService) CreateOrLoadQuery(query *model.Query) {
	now := time.Now()

	owlDb.AddOrRefreshQuery(query, now)

	s.cache.Set(
		query.Uuid.ToUuid().String(),
		query,
		s.config.CacheDuration,
	)
}
