package ginHttp

import (
	"time"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"

	grahttp "github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/gin_http/grafana"
	"github.com/gin-gonic/gin"
)

type QueryInput struct {
	StartTs       time.Time
	EndTs         time.Time
	ComputeMethod string
	Endpoint      string
	Counter       string
}

//this function will generate query string obj for QueryRRDtool
func getq(q QueryInput) cmodel.GraphQueryParam {
	request := cmodel.GraphQueryParam{
		Start:     q.StartTs.Unix(),
		End:       q.EndTs.Unix(),
		ConsolFun: q.ComputeMethod,
		Endpoint:  q.Endpoint,
		Counter:   q.Counter,
	}
	return request
}

func StartLBWeb(serveEng *gin.Engine) *gin.Engine {
	handler := serveEng

	// deprecated
	// compute := handler.Group("/api/v1/func")
	// compute.GET("/compute", computeFunc.Compute)
	// compute.GET("/funcations", computeFunc.GetAvaibleFun)
	// compute.GET("/smapledata", computeFunc.GetTestData)

	// deprecated
	//openfalcon := handler.Group("/api/v1/owl")
	// openfalcon.GET("/endpoints", openFalcon.GetEndpoints)
	// openfalcon.GET("/queryrrd", openFalcon.QueryData)

	// will deprecated, only used for owl-portal top10 & lambda grafana plugin [only support old version] //
	grafanaV0 := handler.Group("/grafana")
	grafanaV0.GET("/", grahttp.GrafanaMain)
	grafanaV0.GET("/metrics/find", grahttp.GrafanaMain)
	grafanaV0.POST("/render", grahttp.GetQueryTargets)
	// will deprecated, only used for owl-portal top10 & lambda grafana plugin [only support old version] //

	// new api
	lambdaWeb := handler.Group("/api/v1/lambdaq")
	lambdaWeb.Use(utils.AuthSessionMidd)
	lambdaWeb.POST("/q", grahttp.LambdaQueryQ)

	return handler
}
