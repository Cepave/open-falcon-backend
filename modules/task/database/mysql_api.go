package database

import (
	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	graphSrv "github.com/Cepave/open-falcon-backend/common/service/graph"
	owlSrv "github.com/Cepave/open-falcon-backend/common/service/owl"
)

var QueryObjectService owlSrv.QueryService
var GraphService graphSrv.GraphService

func InitMySqlApi(restConfig *oHttp.RestfulClientConfig) {
	QueryObjectService = owlSrv.NewQueryService(
		owlSrv.QueryServiceConfig{restConfig},
	)

	GraphService = graphSrv.NewGraphService(
		&graphSrv.GraphServiceConfig{restConfig},
	)
}
