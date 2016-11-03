package main

import (
	"os"
	"syscall"

	"github.com/chyeh/viper"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	commonOs "github.com/Cepave/open-falcon-backend/common/os"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/nqm-conf-mng/rdb"
	"github.com/Cepave/open-falcon-backend/modules/nqm-conf-mng/restful"
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

	commonOs.HoldingAndWaitSignal(exitApp, syscall.SIGINT, syscall.SIGTERM)
}

func exitApp(signal os.Signal) {
	rdb.ReleaseRdb()
}

func toGinConfig(config *viper.Viper) *commonGin.GinConfig {
	return &commonGin.GinConfig {
		Mode: gin.ReleaseMode,
		Host: config.GetString("restful.listen.host"),
		Port: uint16(config.GetInt("restful.listen.port")),
	}
}
func toRdbConfig(config *viper.Viper) *commonDb.DbConfig {
	return &commonDb.DbConfig {
		Dsn: config.GetString("rdb.dsn"),
		MaxIdle: config.GetInt("rdb.maxIdle"),
	}
}

func pflagDefine() {
	pflag.StringP("config", "c", "cfg.json", "configuration file")
	pflag.BoolP("help", "h", false, "usage")
}
