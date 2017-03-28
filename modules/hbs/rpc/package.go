package rpc

import (
	viper "github.com/spf13/viper"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	nqmService "github.com/Cepave/open-falcon-backend/common/service/nqm"
)

var logger = log.NewDefaultLogger("INFO")

func InitPackage(config *viper.Viper) {
	initNqmConfig(config)
}

func initNqmConfig(config *viper.Viper) {
	config.SetDefault("nqm.queue_size.refresh_agent_ping_list", 8)
	config.SetDefault("nqm.cache_minutes.agent_ping_list", 20)

	nqmConfig := nqmService.AgentHbsServiceConfig {
		QueueSizeOfRefreshCacheOfPingList: config.GetInt("nqm.queue_size.refresh_agent_ping_list"),
		CacheTimeoutMinutes: config.GetInt("nqm.cache_minutes.agent_ping_list"),
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
