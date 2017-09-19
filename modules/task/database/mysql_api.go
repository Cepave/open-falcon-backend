package database

import (
	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	owlSrv "github.com/Cepave/open-falcon-backend/common/service/owl"
)

var QueryObjectService owlSrv.QueryService

func InitMySqlApi(restConfig *oHttp.RestfulClientConfig) {
	QueryObjectService = owlSrv.NewQueryService(
		owlSrv.QueryServiceConfig{restConfig},
	)
}
