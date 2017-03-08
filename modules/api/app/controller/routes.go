package controller

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/api/app/controller/dashboardGraphOwl"
	"github.com/Cepave/open-falcon-backend/modules/api/app/controller/dashboardScreenOWl"
	"github.com/Cepave/open-falcon-backend/modules/api/app/controller/expression"
	"github.com/Cepave/open-falcon-backend/modules/api/app/controller/graph"
	"github.com/Cepave/open-falcon-backend/modules/api/app/controller/host"
	"github.com/Cepave/open-falcon-backend/modules/api/app/controller/mockcfg"
	"github.com/Cepave/open-falcon-backend/modules/api/app/controller/strategy"
	"github.com/Cepave/open-falcon-backend/modules/api/app/controller/template"
	"github.com/Cepave/open-falcon-backend/modules/api/app/controller/uic"
	"github.com/Cepave/open-falcon-backend/modules/api/app/utils"
	"gopkg.in/gin-gonic/gin.v1"
)

func StartGin(port string, r *gin.Engine) {
	r.Use(utils.CORS())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, I'm OWL (｡A｡)")
	})
	graph.Routes(r)
	uic.Routes(r)
	template.Routes(r)
	strategy.Routes(r)
	host.Routes(r)
	expression.Routes(r)
	mockcfg.Routes(r)
	dashboardScreenOWl.Routes(r)
	dashboardGraphOwl.Routes(r)
	r.Run(port)
}
