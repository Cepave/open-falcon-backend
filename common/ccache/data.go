package ccache

import (
	"github.com/karlseguin/ccache"
	"time"
)

type DataCacheConfig struct {
	MaxSize  int64
	Duration time.Duration
}

type CacheSetup func(config *ccache.Configuration)

// Set-up the data cache
func NewDataCache(config DataCacheConfig) *ccache.Cache {
	return NewDataCacheWithSetup(config, nil)
}

// Sets-up the data cache with setup function, which you can override the configuration of Ccache
//
// By default, this function would set buckets with following rules:
//
// MaxSize <= 100 - 1 bucket
// MaxSize <= 500 - 2 bucket
// MaxSize <= 2000 - 4 bucket
// MaxSize <= 8000 - 8 bucket
// MaxSize > 8000 - 16 bucket
//
// The one eighths of MaxSize would be ItemsToPrune
func NewDataCacheWithSetup(config DataCacheConfig, setup CacheSetup) *ccache.Cache {
	cconfig := ccache.Configure()

	cconfig.MaxSize(config.MaxSize)

	var buckets uint32
	switch {
	case config.MaxSize <= 100:
		buckets = 1
	case config.MaxSize <= 500:
		buckets = 2
	case config.MaxSize <= 2000:
		buckets = 4
	case config.MaxSize <= 8000:
		buckets = 8
	default:
		buckets = 16
	}
	cconfig.Buckets(buckets)
	cconfig.ItemsToPrune(uint32(config.MaxSize / 8))

	if setup != nil {
		setup(cconfig)
	}

	return ccache.New(cconfig)
}
