package queue

import (
	"sync/atomic"

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
	flush   chan struct{}
	done    chan struct{}
}

func New(c *commonQueue.Config) *Queue {
	return &Queue{
		q:     commonQueue.New(),
		c:     c,
		flush: make(chan struct{}),
		done:  make(chan struct{}),
	}
}

func (q *Queue) Count() uint64 {
	return q.cnt
}

func (q *Queue) Start() {
	q.running = true
	go q.drain()
}

func (q *Queue) drain() {
	for {
		select {
		default:
			reqs := q.q.DrainNWithDurationByType(q.c, new(model.NqmAgentHeartbeatRequest)).([]*model.NqmAgentHeartbeatRequest)
			updateTx := &rdb.UpdateNqmAgentHeartbeatTx{
				Reqs: reqs,
			}
			rdb.DbFacade.SqlxDbCtrl.InTx(updateTx)

			d := uint64(len(reqs))
			atomic.AddUint64(&q.cnt, d)
			logger.Debugf("drained %d elements\n", d)
		case <-q.flush:
			c := &commonQueue.Config{Num: q.c.Num, Dur: 0}
			for q.q.Len() != 0 {
				reqs := q.q.DrainNWithDurationByType(c, new(model.NqmAgentHeartbeatRequest)).([]*model.NqmAgentHeartbeatRequest)
				updateTx := &rdb.UpdateNqmAgentHeartbeatTx{
					Reqs: reqs,
				}
				rdb.DbFacade.SqlxDbCtrl.InTx(updateTx)

				d := uint64(len(reqs))

				atomic.AddUint64(&q.cnt, d)
				logger.Debugf("flushed %d elements\n", d)
			}
			close(q.done)
			return
		}
	}
}

func (q *Queue) Stop() {
	q.running = false
	close(q.flush)
	<-q.done
	return
}

func (q *Queue) Put(v interface{}) {
	if q.running {
		q.q.Enqueue(v)
	}
	return
}

func (q *Queue) Len() int {
	return q.q.Len()
}

var NqmQueue *Queue

func InitNqmHeartbeat(c *commonQueue.Config) {
	NqmQueue = New(c)
	NqmQueue.Start()
}
