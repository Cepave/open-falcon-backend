package main

import (
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	commonOs "github.com/Cepave/open-falcon-backend/common/os"
	commonQueue "github.com/Cepave/open-falcon-backend/common/queue"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/restful"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/service"
	owlSrv "github.com/Cepave/open-falcon-backend/modules/mysqlapi/service/owl"
)

var logger = log.NewDefaultLogger("INFO")

func main() {
	/**
	 * Initialize loader of configurations
	 */
	confLoader := vipercfg.NewOwlConfigLoader()
	confLoader.FlagDefiner = pflagDefine

	confLoader.ProcessTrueValueCallbacks()
	// :~)

	config := confLoader.MustLoadConfigFile()

	rdb.InitRdb(toRdbConfig(config))
	restful.InitGin(toGinConfig(config))
	restful.InitCache(toCacheConfig(config))

	service.InitNqmHeartbeat(toNqmHeartbeatConfig(config))
	service.InitCachedTargetList(toTargetListConfig(config))
	owlSrv.InitQueryObjectService(*toQueryObjectServiceConfig(config))

	commonOs.HoldingAndWaitSignal(exitApp, syscall.SIGINT, syscall.SIGTERM)
}

func exitApp(signal os.Signal) {
	service.CloseNqmHeartbeat()
	service.CloseCachedTargetList()
	rdb.ReleaseRdb()
}

func toGinConfig(config *viper.Viper) *commonGin.GinConfig {
	return &commonGin.GinConfig{
		Mode: gin.ReleaseMode,
		Host: config.GetString("restful.listen.host"),
		Port: uint16(config.GetInt("restful.listen.port")),
	}
}

func toRdbConfig(config *viper.Viper) *commonDb.DbConfig {
	return &commonDb.DbConfig{
		Dsn:     config.GetString("rdb.dsn"),
		MaxIdle: config.GetInt("rdb.maxIdle"),
	}
}

func toQueryObjectServiceConfig(config *viper.Viper) *owlSrv.QueryObjectServiceConfig {
	return &owlSrv.QueryObjectServiceConfig {
		CacheSize: config.GetInt64("queryObject.cache.size"),
		CacheDuration: time.Duration(config.GetInt64("queryObject.cache.hourDuration")) * time.Hour,
	}
}

func toCacheConfig(config *viper.Viper) *restful.CacheConfig {
	return &restful.CacheConfig{
		Size:     config.GetInt("nqm.pingList.cache.size"),
		Lifetime: config.GetInt("nqm.pingList.cache.lifetime"),
	}
}

func toNqmHeartbeatConfig(config *viper.Viper) *commonQueue.Config {
	return &commonQueue.Config{
		Num: config.GetInt("heartbeat.nqm.batchSize"),
		Dur: time.Duration(config.GetInt("heartbeat.nqm.duration")) * time.Second,
	}
}

func toTargetListConfig(config *viper.Viper) *service.NqmCachedTargetListConfig {
	return &service.NqmCachedTargetListConfig{
		Size: config.GetInt("heartbeat.nqm.targetList.size"),
		Dur:  time.Duration(config.GetInt("heartbeat.nqm.targetList.duration")) * time.Minute,
	}
}

func pflagDefine() {
	pflag.StringP("config", "c", "cfg.json", "configuration file")
	pflag.BoolP("help", "h", false, "usage")
}
