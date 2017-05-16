package queue

import (
	"container/list"
	"reflect"
	"sync"
	"time"

	"github.com/Cepave/open-falcon-backend/common/utils"
)

type Queue struct {
	l     *list.List // not thead safe
	mutex sync.Mutex
}

type Config struct {
	Num int
	Dur time.Duration
}

func (q *Queue) Enqueue(v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.l.PushBack(v)
}

// DequeueN dequeues UP TO N elements.
func (q *Queue) dequeueN(num int) []interface{} {
	elems := make([]interface{}, 0, num)
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for i := 0; i < num; i++ {
		if e := q.l.Front(); e != nil {
			elems = append(elems, e.Value)
			q.l.Remove(e)
		} else {
			break
		}
	}
	return elems
}

func (q *Queue) Len() int {
	return q.l.Len()
}

func New() *Queue {
	return &Queue{
		l: list.New(),
	}
}

func (q *Queue) pollN(c *Config) []interface{} {
	t := time.After(c.Dur)
	var elems []interface{}
	batchSize := c.Num
	for {
		select {
		case <-t:
			return elems
		default:
			b := q.dequeueN(batchSize)
			elems = append(elems, b...)
			if batchSize -= len(b); batchSize == 0 {
				return elems
			}
		}

	}
}

func (q *Queue) drainNWithDuration(c *Config) []interface{} {
	elems := make([]interface{}, 0, c.Num)
	for {
		if e := q.pollN(c); len(e) != 0 {
			elems = append(elems, e...)
			if len(elems) >= c.Num {
				return elems
			}
		} else {
			return elems
		}
	}
}

func (q *Queue) DrainNWithDurationByType(c *Config, eleValue interface{}) interface{} {
	return q.drainNWithDurationByReflectType(c, reflect.TypeOf(eleValue))
}

func (q *Queue) drainNWithDurationByReflectType(c *Config, eleType reflect.Type) interface{} {
	var result []interface{}
	if c.Num < 1 {
		result = []interface{}{}
	} else if c.Dur <= 0 {
		result = q.dequeueN(c.Num)
	} else {
		result = q.drainNWithDuration(c)
	}

	return utils.MakeAbstractArray(result).
		MapTo(
			func(v interface{}) interface{} {
				return reflect.ValueOf(v).Convert(eleType).Interface()
			},
			eleType,
		).GetArray()
}
