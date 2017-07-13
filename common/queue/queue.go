package queue

import (
	"container/list"
	"reflect"
	"sync"
	"time"

	or "github.com/Cepave/open-falcon-backend/common/reflect"
	"github.com/Cepave/open-falcon-backend/common/utils"
)

type Config struct {
	Num int
	Dur time.Duration
}

type Queue struct {
	l     *list.List // not thead safe
	mutex *sync.Mutex
}

func New() *Queue {
	return &Queue{
		l:     list.New(),
		mutex: &sync.Mutex{},
	}
}

func (q *Queue) Enqueue(v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.l.PushBack(v)
}
func (q *Queue) DrainNWithDuration(c *Config) []interface{} {
	result := make([]interface{}, 0)

	haveWaited := false
	config := *c
	for config.Num > 0 {
		batchResult := q.dequeueN(config.Num)
		config.Num -= len(batchResult)

		/**
		 * Waiting and fetch next batch of data.
		 *
		 * If the next batch is still nothing after waiting, stop fetching process.
		 */
		if len(batchResult) == 0 {
			if haveWaited {
				break
			}

			haveWaited = true
			if config.Dur > time.Duration(0) {
				time.Sleep(config.Dur)
			}
			continue
		}
		// :~)

		haveWaited = false
		result = append(result, batchResult...)
	}

	return result
}
func (q *Queue) DrainNWithDurationByType(c *Config, eleValue interface{}) interface{} {
	return q.DrainNWithDurationByReflectType(c, or.TypeOfValue(eleValue))
}
func (q *Queue) DrainNWithDurationByReflectType(c *Config, eleType reflect.Type) interface{} {
	return utils.MakeAbstractArray(
		q.DrainNWithDuration(c),
	).
		MapTo(
			func(v interface{}) interface{} {
				return reflect.ValueOf(v).Convert(eleType).Interface()
			},
			eleType,
		).GetArray()
}
func (q *Queue) Len() int {
	return q.l.Len()
}

// DequeueN dequeues UP TO N elements.
func (q *Queue) dequeueN(num int) []interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	elems := make([]interface{}, 0, num)
	for i := 0; i < num; i++ {
		e := q.l.Front()
		if e == nil {
			break
		}

		q.l.Remove(e)
		elems = append(elems, e.Value)
	}
	return elems
}
