package service

import (
	"net/url"

	"github.com/Cepave/open-falcon-backend/common/model/config"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/dghubble/sling"
)

var logger = log.NewDefaultLogger("INFO")

var updateOnlyFlag bool
var MysqlApiUrl string

func InitPackage(cfg *config.MysqlApiConfig, hosts string) {
	MysqlApiUrl = resolveUrl(cfg.Host, cfg.Resource)
	logger.Infoln("[Config] MySQL_API=", MysqlApiUrl)

	if hosts != "" {
		updateOnlyFlag = true
	}
}

func NewSlingBase() *sling.Sling {
	return sling.New().Base(MysqlApiUrl)
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
