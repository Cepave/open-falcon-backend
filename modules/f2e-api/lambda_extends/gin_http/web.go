package ginHttp

import (
	"time"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/gin_http/computeFunc"
	grahttp "github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/gin_http/grafana"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/gin_http/openFalcon"
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

	compute := handler.Group("/func")
	compute.GET("/compute", computeFunc.Compute)
	compute.GET("/funcations", computeFunc.GetAvaibleFun)
	compute.GET("/smapledata", computeFunc.GetTestData)

	openfalcon := handler.Group("/owl")
	openfalcon.GET("/endpoints", openFalcon.GetEndpoints)
	openfalcon.GET("/queryrrd", openFalcon.QueryData)

	grafana := handler.Group("/api/grafana")
	grafana.GET("/", grahttp.GrafanaMain)
	grafana.GET("/metrics/find", grahttp.GrafanaMain)
	grafana.POST("/render", grahttp.GetQueryTargets)

	return handler
}
