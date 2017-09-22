package service

import (
	"net/url"
	"time"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	oConfig "github.com/Cepave/open-falcon-backend/common/model/config"
	oSrv "github.com/Cepave/open-falcon-backend/common/service"
	"github.com/dghubble/sling"
	"github.com/h2non/gentleman/plugins/timeout"
	"github.com/spf13/viper"
	"gopkg.in/h2non/gentleman.v2"
)

var (
	updateOnlyFlag  bool
	mysqlApiUrl     string
	logger          = log.NewDefaultLogger("INFO")
	MysqlApiService oSrv.MysqlApiService
	ClientTimeout   = 3 * time.Second
)

func InitPackage(vpConfig *viper.Viper) {
	apiConfig := ToMysqlApiConfig(vpConfig)
	SetMysqlApiUrl(apiConfig)

	// Set a flag of heartbeat call
	if vpConfig.GetString("hosts") != "" {
		updateOnlyFlag = true
	}

	InitMysqlApiService(buildRestfulConfig(apiConfig))
}

func InitMysqlApiService(config *oHttp.RestfulClientConfig) {
	MysqlApiService = oSrv.NewMysqlApiService(
		&oSrv.MysqlApiServiceConfig{
			config,
		},
	)
}

func ToMysqlApiConfig(config *viper.Viper) *oConfig.MysqlApiConfig {
	return &oConfig.MysqlApiConfig{
		Host:     config.GetString("mysql_api.host"),
		Resource: config.GetString("mysql_api.resource"),
	}
}

func buildRestfulConfig(config *oConfig.MysqlApiConfig) *oHttp.RestfulClientConfig {
	httpConfig := &client.HttpClientConfig{
		Url:            GetMysqlApiUrl(),
		RequestTimeout: ClientTimeout,
	}

	return &oHttp.RestfulClientConfig{
		HttpClientConfig: httpConfig,
		FromModule:       "hbs",
	}
}

func SetMysqlApiUrl(cfg *oConfig.MysqlApiConfig) {
	mysqlApiUrl = resolveUrl(cfg.Host, cfg.Resource)
	logger.Infoln("[Config] MySQL_API=", mysqlApiUrl)
}

func GetMysqlApiUrl() string {
	if mysqlApiUrl == "" {
		logger.Errorln("[Config] Return empty MySQL_API_URL.")
	}
	return mysqlApiUrl
}

func NewSlingBase() *sling.Sling {
	return sling.New().Base(GetMysqlApiUrl())
}

func NewMysqlApiCli() *gentleman.Client {
	return gentleman.New().Use(timeout.Request(ClientTimeout)).URL(GetMysqlApiUrl())
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
