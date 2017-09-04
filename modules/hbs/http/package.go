package http

import (
	"time"

	"github.com/Cepave/open-falcon-backend/modules/hbs/service"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
)

var ClientTimeout = 3 * time.Second
var logger = log.NewDefaultLogger("INFO")

func NewMysqlApiCli() *gentleman.Client {
	return gentleman.New().Use(timeout.Request(ClientTimeout)).BaseURL(service.GetMysqlApiUrl())
}
