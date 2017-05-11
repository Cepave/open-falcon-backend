package queue

import (
	"container/list"
	"reflect"
	"sync"
	"time"

	"github.com/mikelue/cepave-owl/common/utils"
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
	var elems []interface{}
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

func (q *Queue) PollN(c *Config) []interface{} {
	if c.Num < 1 || c.Dur <= 0 {
		return []interface{}{}
	}
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

func (q *Queue) DrainNWithDuration(c *Config) []interface{} {
	if c.Num < 1 || c.Dur <= 0 {
		return []interface{}{}
	}
	elems := make([]interface{}, 0, c.Num)
	for {
		if e := q.PollN(c); len(e) != 0 {
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
	if c.Num < 1 || c.Dur <= 0 {
		return []interface{}{}
	}
	return q.DrainNWithDurationByReflectType(c, reflect.TypeOf(eleValue))
}

func (q *Queue) DrainNWithDurationByReflectType(c *Config, eleType reflect.Type) interface{} {
	if c.Num < 1 || c.Dur <= 0 {
		return []interface{}{}
	}
	popValue := q.DrainNWithDuration(c)

	return utils.MakeAbstractArray(popValue).
		MapTo(
			func(v interface{}) interface{} {
				return reflect.ValueOf(v).Convert(eleType).Interface()
			},
			eleType,
		).GetArray()
}
