package queue

import (
	"sync"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
)

var logger = log.NewDefaultLogger("INFO")

type mode byte

const (
	DRAIN mode = 1
	FLUSH mode = 2
)

type Queue struct {
	q       *commonQueue.Queue
	c       *commonQueue.Config
	cnt     uint64 // counter for the dequeued elements
	running bool
	flush   chan struct{}
	done    chan struct{}
	mutex   sync.Mutex
}

func New(c *commonQueue.Config) *Queue {
	return &Queue{
		q:     commonQueue.New(),
		c:     c,
		done:  make(chan struct{}),
		flush: make(chan struct{}),
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
		select {
		default:
			reqs := q.drainByMode(DRAIN, *q.c)
			if n := q.numToAccum(len(reqs)); n != 0 {
				update(reqs)
				logger.Infof("drained %d NQM agent heartbeat requests from queue", n)
			}
		case <-q.flush:
			for {
				reqs := q.drainByMode(FLUSH, *q.c)
				if n := q.numToAccum(len(reqs)); n != 0 {
					update(reqs)
					logger.Infof("flushed %d NQM agent heartbeat requests from queue", n)
				} else {
					close(q.done)
					return
				}
			}
		}
	}
}

func (q *Queue) drainByMode(m mode, c commonQueue.Config) []*model.NqmAgentHeartbeatRequest {
	if m == FLUSH {
		c.Dur = 0
	}
	return q.q.DrainNWithDurationByType(&c, new(model.NqmAgentHeartbeatRequest)).([]*model.NqmAgentHeartbeatRequest)
}

func (q *Queue) numToAccum(n int) int {
	if n == 0 {
		return 0
	}
	q.cnt += uint64(n)
	return n
}

func update(reqs []*model.NqmAgentHeartbeatRequest) {
	go utils.BuildPanicCapture(
		func() {
			rdb.UpdateNqmAgentHeartbeat(reqs)
		},
		func(p interface{}) {
			logger.Errorf("[PANIC] NQM agent's heartbeat requests(%v)", reqs)
		},
	)()
}

func (q *Queue) Stop() {
	q.mutex.Lock()
	if !q.running {
		q.mutex.Unlock()
		return
	}
	q.running = false
	q.mutex.Unlock()
	close(q.flush)
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

func CloseNqmHeartbeat() {
	logger.Info("Closing NQM heartbeat queue service...")
	NqmQueue.Stop()
	logger.Info("Closed NQM heartbeat queue service.")
}
