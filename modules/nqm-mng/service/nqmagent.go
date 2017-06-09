package service

import (
	"time"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
)

var logger = log.NewDefaultLogger("INFO")

type mode byte

const (
	_DRAIN mode = 1
	_FLUSH mode = 2
)

type nqmAgentUpdateService struct {
	q       *commonQueue.Queue
	c       *commonQueue.Config
	cnt     uint64 // counter for the dequeued elements
	running bool
	flush   chan struct{}
	done    chan struct{}
}

func New(c *commonQueue.Config) *nqmAgentUpdateService {
	return &nqmAgentUpdateService{
		q:     commonQueue.New(),
		c:     c,
		done:  make(chan struct{}),
		flush: make(chan struct{}),
	}
}

func (q *nqmAgentUpdateService) Count() uint64 {
	return q.cnt
}

func (q *nqmAgentUpdateService) Start() {
	if q.running {
		return
	}
	q.running = true
	go q.drain()
}

func (q *nqmAgentUpdateService) drain() {
	for {
		select {
		default:
			reqs := q.drainByMode(_DRAIN, *q.c)
			if n := q.numToAccum(len(reqs)); n != 0 {
				update(reqs)
			}
		case <-q.flush:
			for {
				reqs := q.drainByMode(_FLUSH, *q.c)
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

func (q *nqmAgentUpdateService) drainByMode(m mode, c commonQueue.Config) []*model.NqmAgentHeartbeatRequest {
	if m == _FLUSH {
		c.Dur = 0
	}
	return q.q.DrainNWithDurationByType(&c, new(model.NqmAgentHeartbeatRequest)).([]*model.NqmAgentHeartbeatRequest)
}

func (q *nqmAgentUpdateService) numToAccum(n int) int {
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

func (q *nqmAgentUpdateService) Stop() {
	if !q.running {
		return
	}
	q.running = false
	time.Sleep(q.c.Dur) // for all `q.q.Enqueue()`s to be done
	close(q.flush)
	<-q.done
}

func (q *nqmAgentUpdateService) Put(v interface{}) {
	if !q.running {
		return
	}
	q.q.Enqueue(v)
}

func (q *nqmAgentUpdateService) Len() int {
	return q.q.Len()
}

var NqmQueue *nqmAgentUpdateService

func InitNqmHeartbeat(c *commonQueue.Config) {
	NqmQueue = New(c)
	NqmQueue.Start()
}

func CloseNqmHeartbeat() {
	logger.Info("Closing NQM heartbeat queue service...")
	NqmQueue.Stop()
	logger.Info("Finish.")
}
