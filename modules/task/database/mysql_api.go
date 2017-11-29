package database

import (
	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	graphSrv "github.com/Cepave/open-falcon-backend/common/service/graph"
	owlSrv "github.com/Cepave/open-falcon-backend/common/service/owl"
)

var QueryObjectService owlSrv.QueryService
var GraphService graphSrv.GraphService
var ClearTaskLogEntryService owlSrv.ClearLogService

func InitMySqlApi(restConfig *oHttp.RestfulClientConfig) {
	QueryObjectService = owlSrv.NewQueryService(
		owlSrv.QueryServiceConfig{restConfig},
	)

	GraphService = graphSrv.NewGraphService(
		&graphSrv.GraphServiceConfig{restConfig},
	)

	ClearTaskLogEntryService = owlSrv.NewClearLogService(
		owlSrv.ClearLogServiceConfig{restConfig},
	)
}
