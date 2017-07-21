package ccache

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/karlseguin/ccache"
	"reflect"
	"time"
)

// Used to indicate the controller that it should not put the object into cache
type CacheCtrl struct {
	*ccache.Cache
}

func NewCacheCtrl(cache *ccache.Cache) *CacheCtrl {
	return &CacheCtrl{cache}
}

// Fetchs object with panic if there is error occurred
//
// nil value wouldn't be fetched
func (cc *CacheCtrl) MustFetch(
	key string, duration time.Duration,
	fetch func() interface{},
) *ccache.Item {
	item, err := cc.Cache.Fetch(
		key, duration,
		func() (v interface{}, err error) {
			defer func() {
				p := recover()
				if p != nil {
					err = fmt.Errorf("%v", p)
				}
			}()

			return fetch(), nil
		},
	)

	if err != nil {
		panic(err)
	}

	return item
}

func (cc *CacheCtrl) MustFetchAndDoNotCacheEmpty(
	key string, duration time.Duration,
	fetch func() interface{},
) *ccache.Item {
	item := cc.Cache.Get(key)

	if item == nil {
		fetchedValue := fetch()

		if utils.ValueExt(reflect.ValueOf(fetchedValue)).IsViable() {
			cc.Set(key, fetchedValue, duration)
		}

		item = cc.Cache.Get(key)
	}

	return item
}

func (cc *CacheCtrl) MustFetchNative(
	key string, duration time.Duration,
	fetch func() interface{},
) interface{} {
	return cc.MustFetch(
		key, duration, fetch,
	).Value()
}

func (cc *CacheCtrl) MustFetchNativeAndDoNotCacheEmpty(
	key string, duration time.Duration,
	fetch func() interface{},
) interface{} {
	item := cc.MustFetchAndDoNotCacheEmpty(
		key, duration, fetch,
	)

	if item == nil {
		return nil
	}

	return item.Value()
}
