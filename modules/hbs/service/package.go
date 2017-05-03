package service

import (
	"net/url"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/dghubble/sling"
)

var logger = log.NewDefaultLogger("INFO")

var updateOnlyFlag bool
var mysqlApiSling *sling.Sling

func InitPackage() {
	mysqlApiConfig := g.Config().MysqlApi
	mysqlApiUrl := resolveUrl(mysqlApiConfig.Host, mysqlApiConfig.Resource)
	mysqlApiSling = sling.New().Base(mysqlApiUrl)

	if g.Config().Hosts != "" {
		updateOnlyFlag = true
	}
}

func NewSlingBase() *sling.Sling {
	return mysqlApiSling.New()
}

func resolveUrl(host string, resource string) string {
	base, err := url.Parse(host)
	if err != nil {
		logger.Errorln(err)
	}
	ref, err := url.Parse(resource)
	if err != nil {
		logger.Errorln(err)
	}

	return base.ResolveReference(ref).String()
}
