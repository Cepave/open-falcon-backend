package nqm

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	utils "github.com/Cepave/open-falcon-backend/common/utils"

	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
)

// Configuration of HBS service for NQM agent
type AgentHbsServiceConfig struct {
	// The queue size of refresing cache of ping list
	QueueSizeOfRefreshCacheOfPingList int
	// The timeout of cache(minutes)
	CacheTimeoutMinutes int
}

func NewAgentHbsService(config AgentHbsServiceConfig) *AgentHbsService {
	return &AgentHbsService{
		agentIdQueueForRefreshCache: make(chan int32, config.QueueSizeOfRefreshCacheOfPingList),
		cacheTimeout:                time.Duration(config.CacheTimeoutMinutes) * time.Minute,
	}
}

// Main service of HBS for NQM agent
type AgentHbsService struct {
	cacheTimeout                time.Duration
	agentIdQueueForRefreshCache chan int32
}

// Loads ping list for agent at certain time
func (s *AgentHbsService) LoadPingList(agent *nqmModel.NqmAgent, checkedTime time.Time) []commonModel.NqmTarget {
	result, cacheLog := nqmDb.GetPingListFromCache(agent, checkedTime)

	go utils.BuildPanicCapture(
		func() {
			s.addRefreshCache(int32(agent.Id), cacheLog)
		},
		func(p interface{}) {
			logger.Errorf("Cannot add agent id [%d] to queue: %v", agent.Id, p)
		},
	)()

	return result
}

// Start the service for refreshing cache of ping list
func (s *AgentHbsService) Start() {
	go func() {
		for agentId := range s.agentIdQueueForRefreshCache {
			logger.Debugf("Refresh cache for agent: [%d]", agentId)

			err := s.buildCacheOfPingList(agentId)
			if err != nil {
				logger.Errorf("Agent[%d]. Refresh has error: %v", agentId, err)
			}
		}
	}()
}

// Release resources of this service
func (s *AgentHbsService) Stop() {
	queueSize := len(s.agentIdQueueForRefreshCache)
	logger.Infof("Stopping AgentHbsService. Size of queue(refreshing cache): [%d]", queueSize)

	close(s.agentIdQueueForRefreshCache)

	/**
	 * Waiting for queue to be processed
	 */
	maxTimes := 30
	for queueSize > 0 && maxTimes > 0 {
		logger.Infof("Sleep for 2 seconds to wait queue to be processed... Current size: [%d]", queueSize)
		time.Sleep(2 * time.Second)
		maxTimes--
		queueSize = len(s.agentIdQueueForRefreshCache)
	}
	// :~)
}

func (s *AgentHbsService) addRefreshCache(agentId int32, cacheLog *nqmModel.PingListLog) {
	now := time.Now()

	/**
	 * If the timeout has reached, adds the id of agent into queue for refreshing cache
	 */
	diffDuration := now.Sub(cacheLog.RefreshTime)
	if logger.Level == log.DebugLevel {
		logger.Debugf(
			"Queue Size(refreshing cache of ping list): [%d]. Minutes: [%d]",
			len(s.agentIdQueueForRefreshCache),
			diffDuration/time.Minute,
		)
	}
	if now.Sub(cacheLog.RefreshTime) >= s.cacheTimeout {
		s.agentIdQueueForRefreshCache <- agentId
	}
	// :~)
}

func (s *AgentHbsService) buildCacheOfPingList(agentId int32) (err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()

	nqmDb.BuildCacheOfPingList(agentId, time.Now())

	return nil
}
