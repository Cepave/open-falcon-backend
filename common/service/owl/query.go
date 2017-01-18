package owl

import (
	"github.com/satori/go.uuid"
	c "github.com/karlseguin/ccache"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	db "github.com/Cepave/open-falcon-backend/common/db"
	model "github.com/Cepave/open-falcon-backend/common/model/owl"
	"time"
)

type QueryServiceConfig struct {
	Name string
	CacheSize int64
	CacheDuration time.Duration
}
type QueryService struct {
	config *QueryServiceConfig
	cache *c.Cache
}

func KeyByUuid(uuid uuid.UUID) string {
	return uuid.String()
}
func KeyByDbUuid(dbUuid db.DbUuid) string {
	return KeyByUuid(uuid.UUID(dbUuid))
}

// The max size of cache and size of bucket
//
// 	<= 100 - 2
// 	<= 1000 - 4
// 	<= 5000 - 8
// 	> 10000 - 16
//
// 	The prune size is size of cache / 8
func NewQueryService(config QueryServiceConfig) *QueryService {
	var bucketSize uint32 = 16

	switch {
	case config.CacheSize <= 100:
		bucketSize = 2
	case config.CacheSize <= 1000:
		bucketSize = 4
	case config.CacheSize <= 5000:
		bucketSize = 8
	}

	cacheObject := c.New(
		c.Configure().
			MaxSize(config.CacheSize).
			Buckets(bucketSize).
			ItemsToPrune(uint32(config.CacheSize / 8)),
	)

	return &QueryService {
		&config, cacheObject,
	}
}

func (s *QueryService) LoadQueryByUuid(uuid uuid.UUID) *model.Query {
	now := time.Now()
	key := KeyByUuid(uuid)

	/**
	 * Updates access time to database if cache has this query
	 */
	itemInCache := s.cache.Get(key)
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
		s.config.Name, uuid, now,
	)
	if queryFromDatabase == nil {
		return nil
	}
	// :~)

	s.cache.Set(key, queryFromDatabase, s.config.CacheDuration)
	return s.cache.Get(key).Value().(*model.Query)
}

// Loads object of query or creating one
func (s *QueryService) CreateOrLoadQuery(query *model.Query) {
	now := time.Now()

	query.NamedId = s.config.Name
	owlDb.AddOrRefreshQuery(query, now)

	s.cache.Set(
		KeyByDbUuid(query.Uuid),
		query,
		s.config.CacheDuration,
	)
}
