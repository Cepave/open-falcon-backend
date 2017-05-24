package queue

import (
	"sync"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
)

var logger = log.NewDefaultLogger("INFO")

type Queue struct {
	q       *commonQueue.Queue
	c       *commonQueue.Config
	cnt     uint64 // counter for the dequeued elements
	running bool
	done    chan struct{}
	mutex   sync.Mutex
}

func New(c *commonQueue.Config) *Queue {
	return &Queue{
		q:    commonQueue.New(),
		c:    c,
		done: make(chan struct{}),
	}
}

func (q *Queue) Count() uint64 {
	return q.cnt
}

func (q *Queue) Start() {
	q.mutex.Lock()
	if q.running {
		q.mutex.Unlock()
		return
	}
	q.running = true
	q.mutex.Unlock()
	go q.drain()
}

func (q *Queue) drain() {
	for {
		switch q.running {
		case true:
			reqs := q.q.DrainNWithDurationByType(q.c, new(model.NqmAgentHeartbeatRequest)).([]*model.NqmAgentHeartbeatRequest)
			d := uint64(len(reqs))
			q.cnt += d
			logger.Debugf("drained %d NQM agent heartbeat requests from queue\n", d)

			rdb.UpdateNqmAgentHeartbeat(reqs)
		case false:
			c := &commonQueue.Config{Num: q.c.Num, Dur: 0}
			for {
				reqs := q.q.DrainNWithDurationByType(c, new(model.NqmAgentHeartbeatRequest)).([]*model.NqmAgentHeartbeatRequest)
				if len(reqs) == 0 {
					close(q.done)
					return
				}
				d := uint64(len(reqs))
				q.cnt += d
				logger.Debugf("flushed %d NQM agent heartbeat requests from queue\n", d)

				rdb.UpdateNqmAgentHeartbeat(reqs)
			}
		}
	}
}

func (q *Queue) Stop() {
	q.mutex.Lock()
	if !q.running {
		q.mutex.Unlock()
		return
	}
	q.running = false
	q.mutex.Unlock()
	<-q.done
}

func (q *Queue) Put(v interface{}) {
	q.mutex.Lock()
	if !q.running {
		q.mutex.Unlock()
		return
	}
	q.mutex.Unlock()
	q.q.Enqueue(v)
}

func (q *Queue) Len() int {
	return q.q.Len()
}

var NqmQueue *Queue

func InitNqmHeartbeat(c *commonQueue.Config) {
	NqmQueue = New(c)
	NqmQueue.Start()
}
