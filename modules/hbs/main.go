package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/viper"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/model/config"
	oos "github.com/Cepave/open-falcon-backend/common/os"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/Cepave/open-falcon-backend/modules/hbs/http"
	"github.com/Cepave/open-falcon-backend/modules/hbs/rpc"
	"github.com/Cepave/open-falcon-backend/modules/hbs/service"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()

	service.InitPackage(
		toMysqlApiConfig(vipercfg.Config()),
		vipercfg.Config().GetString("hosts"),
	)
	service.InitMySqlApi(loadMySqlApiConfig(vipercfg.Config()))
	rpc.InitPackage(vipercfg.Config())

	go http.Start()
	go rpc.Start()

	oos.HoldingAndWaitSignal(
		func(signal os.Signal) {
			rpc.Stop()
		},
		os.Interrupt, os.Kill,
		syscall.SIGTERM,
	)
}

func toMysqlApiConfig(cfg *viper.Viper) *config.MysqlApiConfig {
	return &config.MysqlApiConfig{
		Host:     cfg.GetString("mysql_api.host"),
		Resource: cfg.GetString("mysql_api.resource"),
	}
}

func loadMySqlApiConfig(viper *viper.Viper) *oHttp.RestfulClientConfig {
	httpConfig := &client.HttpClientConfig{
		Url:            viper.GetString("mysql_api.host"),
		RequestTimeout: service.CLIENT_TIMEOUT,
	}

	resource := viper.GetString("mysql_api.resource")
	if resource != "" {
		httpConfig.Url += "/" + resource
	}

	return &oHttp.RestfulClientConfig{
		HttpClientConfig: httpConfig,
		FromModule:       "hbs",
	}
}
