package service

import (
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/dghubble/sling"
)

var logger = log.NewDefaultLogger("INFO")

func NewSlingBase() *sling.Sling {
	return sling.New().Base(g.Config().MysqlApi)
}
