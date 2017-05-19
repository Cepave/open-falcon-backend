package queue

import (
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
)

type Queue struct {
	q *commonQueue.Queue
	c *commonQueue.Config
}

func New(c *commonQueue.Config) *Queue {
	return &Queue{
		q: commonQueue.New(),
		c: c,
	}
}

func (q *Queue) Start() {
	go q.drain()
}

func (q *Queue) drain() {
	for {
		reqs := q.q.DrainNWithDurationByType(q.c, new(model.NqmAgentHeartbeatRequest)).([]*model.NqmAgentHeartbeatRequest)
		updateTx := &rdb.UpdateNqmAgentProcessor{
			Reqs: reqs,
		}
		rdb.DbFacade.SqlxDbCtrl.InTx(updateTx)
	}
}

func (q *Queue) Put(v interface{}) {
	q.q.Enqueue(v)
	return
}

var NqmQueue *Queue

func InitNqmHeartbeat(c *commonQueue.Config) {
	NqmQueue = New(c)
	NqmQueue.Start()
}
