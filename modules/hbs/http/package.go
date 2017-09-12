package http

import (
	"time"

	"github.com/Cepave/open-falcon-backend/modules/hbs/service"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	oSrv "github.com/Cepave/open-falcon-backend/common/service"
)

var (
	MysqlApiService oSrv.MysqlApiService
	logger          = log.NewDefaultLogger("INFO")
	CLIENT_TIMEOUT  = 3 * time.Second
)

func NewMysqlApiCli() *gentleman.Client {
	return gentleman.New().Use(timeout.Request(CLIENT_TIMEOUT)).URL(service.GetMysqlApiUrl())
}

func InitMySqlApi(config *oHttp.RestfulClientConfig) {
	MysqlApiService = oSrv.NewMysqlApiService(
		oSrv.MysqlApiServiceConfig{
			config,
		},
	)
}
