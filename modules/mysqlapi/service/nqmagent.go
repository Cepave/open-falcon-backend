package service

import (
	"fmt"
	"reflect"
	"time"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
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

var typeOfNqmAgentHeartbeat = reflect.TypeOf(new(nqmModel.HeartbeatRequest))

type nqmAgentUpdateService struct {
	q                *commonQueue.Queue
	c                *commonQueue.Config
	cnt              uint64 // counter for the dequeued elements
	running          bool
	flush            chan struct{}
	done             chan struct{}
	updateToDatabase func([]*nqmModel.HeartbeatRequest)
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

func (q *nqmAgentUpdateService) Put(req *nqmModel.HeartbeatRequest) {
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

	var reqs []*nqmModel.HeartbeatRequest

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

func (q *nqmAgentUpdateService) drainFromQueue(config *commonQueue.Config) []*nqmModel.HeartbeatRequest {
	return q.q.DrainNWithDurationByType(
		config, typeOfNqmAgentHeartbeat,
	).([]*nqmModel.HeartbeatRequest)
}

func updateNqmAgentHeartbeatImpl(reqs []*nqmModel.HeartbeatRequest) {
	utils.BuildPanicCapture(
		func() {
			rdb.UpdateNqmAgentHeartbeat(reqs)
		},
		func(p interface{}) {
			logger.Errorf("[PANIC] Update heartbeats of NQM agent[#%d]: %v", len(reqs), p)
		},
	)()
}

var NqmCachedTargetList *nqmCachedTargetListService

func InitCachedTargetList(c *NqmCachedTargetListConfig) {
	NqmCachedTargetList = newNqmCachedTargetListService(c)
	logger.Infof("Target list service for agent. Timeout: %v. Queue Size: %v", c.Dur, c.Size)
	NqmCachedTargetList.Start()
}

func CloseCachedTargetList() {
	logger.Info("Closing NQM target list service...")
	NqmCachedTargetList.Stop()
	logger.Info("Finish.")
}

type NqmCachedTargetListConfig struct {
	Size int
	Dur  time.Duration
}

func newNqmCachedTargetListService(c *NqmCachedTargetListConfig) *nqmCachedTargetListService {
	return &nqmCachedTargetListService{
		agentIDQueueForRefreshCache: make(chan int32, c.Size),
		cacheTimeout:                c.Dur,
	}
}

type nqmCachedTargetListService struct {
	cacheTimeout                time.Duration
	agentIDQueueForRefreshCache chan int32
}

// Load get the current cached target list for an agent
func (s *nqmCachedTargetListService) Load(agentID int32) []*nqmModel.HeartbeatTarget {
	result, cacheLog := rdb.GetPingListFromCache(agentID, time.Now())

	go utils.BuildPanicCapture(
		func() {
			s.addRefreshCache(agentID, cacheLog)
		},
		func(p interface{}) {
			logger.Errorf("Cannot add agent id [%d] to queue: %v", agentID, p)
		},
	)()

	return result
}

// Start the service for refreshing cache of ping list
func (s *nqmCachedTargetListService) Start() {
	go func() {
		for agentID := range s.agentIDQueueForRefreshCache {
			logger.Debugf("Refresh cache for agent: [%d]", agentID)

			err := s.buildCacheOfPingList(agentID)
			if err != nil {
				logger.Errorf("Agent[%d]. Refresh has error: %v", agentID, err)
			}
		}
	}()
}

// Release resources of this service
func (s *nqmCachedTargetListService) Stop() {
	queueSize := len(s.agentIDQueueForRefreshCache)
	logger.Infof("Stopping nqmCachedTargetListService. Size of queue(refreshing cache): [%d]", queueSize)

	close(s.agentIDQueueForRefreshCache)

	/**
	 * Waiting for queue to be processed
	 */
	maxTimes := 30
	for queueSize > 0 && maxTimes > 0 {
		logger.Infof("Sleep for 2 seconds to wait queue to be processed... Current size: [%d]", queueSize)
		time.Sleep(2 * time.Second)
		maxTimes--
		queueSize = len(s.agentIDQueueForRefreshCache)
	}
	// :~)
}

func (s *nqmCachedTargetListService) addRefreshCache(agentID int32, cacheLog *model.PingListLog) {
	now := time.Now()

	/**
	 * If the timeout has reached, adds the id of agent into queue for refreshing cache
	 */
	diffDuration := now.Sub(cacheLog.RefreshTime)

	logger.Debugf(
		"Queue Size(refreshing cache of ping list): [%d]. Minutes: [%d]",
		len(s.agentIDQueueForRefreshCache),
		diffDuration/time.Minute,
	)

	if now.Sub(cacheLog.RefreshTime) >= s.cacheTimeout {
		s.agentIDQueueForRefreshCache <- agentID
	}
	// :~)
}

func (s *nqmCachedTargetListService) buildCacheOfPingList(agentID int32) (err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()

	rdb.BuildCacheOfPingList(agentID, time.Now())

	return nil
}
