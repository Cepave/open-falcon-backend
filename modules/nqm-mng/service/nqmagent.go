package service

import (
	"reflect"
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

var NqmQueue *nqmAgentUpdateService

func InitNqmHeartbeat(c *commonQueue.Config) {
	NqmQueue = newNqmAgentUpdateService(c)
	NqmQueue.updateToDatabase = updateNqmAgentHeartbeatImpl
	NqmQueue.Start()
}

func CloseNqmHeartbeat() {
	logger.Info("Closing NQM heartbeat queue service...")
	NqmQueue.Stop()
	logger.Info("Finish.")
}

var typeOfNqmAgentHeartbeat = reflect.TypeOf(new(model.NqmAgentHeartbeatRequest))

type nqmAgentUpdateService struct {
	q                *commonQueue.Queue
	c                *commonQueue.Config
	cnt              uint64 // counter for the dequeued elements
	running          bool
	flush            chan struct{}
	done             chan struct{}
	updateToDatabase func([]*model.NqmAgentHeartbeatRequest)
}

func newNqmAgentUpdateService(c *commonQueue.Config) *nqmAgentUpdateService {
	return &nqmAgentUpdateService{
		q:     commonQueue.New(),
		c:     c,
		done:  make(chan struct{}),
		flush: make(chan struct{}),
	}
}

// Gets the number of consumed updating requests(not guarantee on database)
func (q *nqmAgentUpdateService) ConsumedCount() uint64 {
	return q.cnt
}

// Gets the number of pending request of heartbeats
func (q *nqmAgentUpdateService) PendingLen() int {
	return q.q.Len()
}

func (q *nqmAgentUpdateService) Start() {
	if q.running {
		return
	}
	q.running = true
	go q.draining()
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

func (q *nqmAgentUpdateService) Put(req *model.NqmAgentHeartbeatRequest) {
	if !q.running {
		return
	}
	q.q.Enqueue(req)
}

func (q *nqmAgentUpdateService) draining() {
	for {
		select {
		default:
			q.syncToDatabase(_DRAIN)
		case <-q.flush:
			q.syncToDatabase(_FLUSH)
			close(q.done)
			return
		}
	}
}

func (q *nqmAgentUpdateService) syncToDatabase(m mode) {
	var config commonQueue.Config = *q.c

	var reqs []*model.NqmAgentHeartbeatRequest

	switch m {
	case _FLUSH:
		config.Dur = 0
		reqs = q.drainFromQueue(&config)
		if len(reqs) > 0 {
			logger.Infof("Flushing [%d] heartbeats of NQM agent from queue", len(reqs))
		}
	default:
		reqs = q.drainFromQueue(&config)
	}

	q.updateToDatabase(reqs)
	q.cnt += uint64(len(reqs))

	if len(reqs) > 0 {
		logger.Debugf("[%d] heartbeats of NQM agent from queue", len(reqs))

		if m == _FLUSH {
			q.syncToDatabase(m)
		}
	}
}

func (q *nqmAgentUpdateService) drainFromQueue(config *commonQueue.Config) []*model.NqmAgentHeartbeatRequest {
	return q.q.DrainNWithDurationByType(
		config, typeOfNqmAgentHeartbeat,
	).([]*model.NqmAgentHeartbeatRequest)
}

func updateNqmAgentHeartbeatImpl(reqs []*model.NqmAgentHeartbeatRequest) {
	utils.BuildPanicCapture(
		func() {
			rdb.UpdateNqmAgentHeartbeat(reqs)
		},
		func(p interface{}) {
			logger.Errorf("[PANIC] Update heartbeats of NQM agent[#%d]: %v", len(reqs), p)
		},
	)()
}
