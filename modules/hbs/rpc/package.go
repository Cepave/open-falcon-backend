package rpc

import (
	viper "github.com/spf13/viper"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	nqmService "github.com/Cepave/open-falcon-backend/common/service/nqm"
	hbsService "github.com/Cepave/open-falcon-backend/modules/hbs/service"
)

var logger = log.NewDefaultLogger("INFO")

func InitPackage(config *viper.Viper) {
	initNqmConfig(config)
	initFalconConfig(config)
}

func initNqmConfig(config *viper.Viper) {
	config.SetDefault("nqm.queue_size.refresh_agent_ping_list", 8)
	config.SetDefault("nqm.cache_minutes.agent_ping_list", 20)

	nqmConfig := nqmService.AgentHbsServiceConfig{
		QueueSizeOfRefreshCacheOfPingList: config.GetInt("nqm.queue_size.refresh_agent_ping_list"),
		CacheTimeoutMinutes:               config.GetInt("nqm.cache_minutes.agent_ping_list"),
	}

	/**
	 * If the mode is not in debug, the least timeout is 10 minutes
	 */
	if !config.GetBool("debug") && config.GetInt("nqm.cache_minutes.agent_ping_list") < 5 {
		nqmConfig.CacheTimeoutMinutes = 5
	}
	// :~)

	nqmAgentHbsService = nqmService.NewAgentHbsService(nqmConfig)

	logger.Infof("[NQM] Ping list of agent. Timeout: %d minutes. Queue Size: %d",
		nqmConfig.CacheTimeoutMinutes, nqmConfig.QueueSizeOfRefreshCacheOfPingList,
	)
}

func initFalconConfig(config *viper.Viper) {
	heartbeatConfig := &commonQueue.Config{
		Num: config.GetInt("heartbeat.falcon.batchSize"),
		Dur: config.GetDuration("heartbeat.falcon.duration"),
	}
	agentHeartbeatService = hbsService.NewAgentHeartbeatService(heartbeatConfig)
}
