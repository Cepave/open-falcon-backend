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

// DequeueN dequeues UP TO N elements.
func (q *Queue) DequeueN(num int) []interface{} {
	if num < 1 || q.l.Len() == 0 { // No need to use thread-safe Len()
		return []interface{}{}
	}
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

func (q *Queue) ThreadUnsafeLen() int { // Not thread-safe, not for accessing the data structure afterwards.
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

func (q *Queue) PollN(num int, timeout time.Duration) []interface{} {
	t := time.After(timeout)
	var elems []interface{}
	batchSize := num
	for {
		select {
		case <-t:
			return elems
		default:
			b := q.DequeueN(batchSize)
			elems = append(elems, b...)
			if batchSize -= len(b); batchSize == 0 {
				return elems
			}
		}

	}
}

func (q *Queue) DrainNWithDuration(num int, dur time.Duration) []interface{} {
	elems := make([]interface{}, 0, num)
	for {
		if e := q.PollN(num, dur); len(e) != 0 {
			elems = append(elems, e...)
			if len(elems) >= num {
				return elems
			}
		} else {
			return elems
		}
	}
}
