package queue

import (
	"container/list"
	"sync"
	"time"
)

type Queue struct {
	l     *list.List // not thead safe
	mutex sync.Mutex
}

func (q *Queue) Enqueue(v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.l.PushBack(v)
}

func (q *Queue) Dequeue() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if e := q.l.Front(); e != nil {
		defer q.l.Remove(e)
		return e.Value
	}
	return nil
}

func (q *Queue) ThreadUnsafeLen() int { // Not thread-safe, not for accessing the data structure afterwards
	return q.l.Len()
}

func (q *Queue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.l.Len()
}

func (q *Queue) PeekFirst() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if f := q.l.Front(); f != nil {
		return f.Value
	}
	return nil
}

func (q *Queue) PeekLast() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if b := q.l.Back(); b != nil {
		return b.Value
	}
	return nil
}

func New() *Queue {
	return &Queue{
		l: list.New(),
	}
}

func (q *Queue) Poll(timeout time.Duration) interface{} {
	t := time.After(timeout)
	for {
		select {
		case <-t:
			return nil
		default:
			if e := q.Dequeue(); e != nil {
				return e
			}
		}
	}
}

func (q *Queue) DrainWithDuration(num int, dur time.Duration) []interface{} {
	es := make([]interface{}, 0, num)
	for {
		if e := q.Poll(dur); e != nil {
			es = append(es, e)
			if len(es) >= num {
				return es
			}
		} else {
			return es
		}
	}
}
