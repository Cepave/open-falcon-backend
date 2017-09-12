package service

import (
	"net/url"
	"time"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/model/config"
	oSrv "github.com/Cepave/open-falcon-backend/common/service"
	"github.com/dghubble/sling"
	"github.com/h2non/gentleman/plugins/timeout"
	"gopkg.in/h2non/gentleman.v2"
)

var (
	updateOnlyFlag  bool
	mysqlApiUrl     string
	logger          = log.NewDefaultLogger("INFO")
	MysqlApiService oSrv.MysqlApiService
	CLIENT_TIMEOUT  = 3 * time.Second
)

func InitMySqlApi(config *oHttp.RestfulClientConfig) {
	MysqlApiService = oSrv.NewMysqlApiService(
		oSrv.MysqlApiServiceConfig{
			config,
		},
	)
}

func InitPackage(cfg *config.MysqlApiConfig, hosts string) {
	mysqlApiUrl = resolveUrl(cfg.Host, cfg.Resource)
	logger.Infoln("[Config] MySQL_API=", mysqlApiUrl)

	if hosts != "" {
		updateOnlyFlag = true
	}
}

func GetMysqlApiUrl() string {
	if mysqlApiUrl == "" {
		logger.Errorln("[Config] Return empty MySQL_API_URL.")
	}
	return mysqlApiUrl
}

func NewSlingBase() *sling.Sling {
	return sling.New().Base(mysqlApiUrl)
}

func NewMysqlApiCli() *gentleman.Client {
	return gentleman.New().Use(timeout.Request(CLIENT_TIMEOUT)).URL(GetMysqlApiUrl())
}

func resolveUrl(host string, resource string) string {
	if host == "" {
		logger.Panicln("Empty host url of mysql_api.")
	}

	base, err := url.Parse(host)
	if err != nil {
		logger.Panicf("Cannot parse host of mysql_api: %s. Error: %v\n", host, err)
	}
	ref, err := url.Parse(resource)
	if err != nil {
		logger.Panicf("Cannot parse resource of mysql_api: %s. Error: %v\n", resource, err)
	}

	return base.ResolveReference(ref).String()
}
