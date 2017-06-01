package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/alarm"
	"github.com/spf13/viper"
	// "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/dashboardGraphOwl"
	// "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/dashboardScreenOWl"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/dashboard_graph"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/dashboard_screen"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/expression"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/graph"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/host"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/mockcfg"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/strategy"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/template"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/uic"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

var headers = []string{
	"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Cache-Control", "X-Requested-With",
	"accept", "origin", "Apitoken",
	"page-size", "page-pos", "order-by", "page-ptr", "total-count", "page-more", "previous-page", "next-page",
}

var corsConfig cors.Config

func StartGin(port string, r *gin.Engine, testMode bool) *gin.Engine {
	corsConfig = cors.Config{
		Methods:        strings.Join([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"}, ","),
		Origins:        "*",
		RequestHeaders: strings.Join(headers, ","),
		ExposedHeaders: strings.Join(headers, ","),
		Credentials:    true,
		MaxAge:         12 * time.Hour,
	}
	r.Use(cors.Middleware(corsConfig))
	r.Use(utils.CORS())
	// document routes
	r.StaticFS("doc", http.Dir("doc"))
	// markdown web page doc site routes
	if viper.GetBool("web_doc_site.enable") {
		docSitePath := viper.GetString("web_doc_site.folder")
		if docSitePath == "" {
			//default folder
			docSitePath = "./docs/_site"
		}
		log.Infof("markdown web started with path: %s\n", docSitePath)
		r.Static("docs", docSitePath)
	}
	r.GET("health", func(c *gin.Context) {
		db := config.Con()
		status, errorTable := db.HealthCheck()
		message := "everything is works!"
		if len(errorTable) > 0 {
			message = fmt.Sprintf("%s is down, please check it.", strings.Join(errorTable, ","))
		}
		c.JSON(200, map[string]interface{}{
			"rdb": map[string]interface{}{
				"ping_result":  status,
				"ping_message": message,
			},
			"http": map[string]interface{}{
				"listening": port,
			},
		})
		return
	})
	// frontend files
	if viper.GetBool("frontend.enable") {
		frontendFolder := viper.GetString("frontend.folder")
		frontendFolderStatic := viper.GetString("frontend.static_path")
		if frontendFolder == "" {
			frontendFolder = "./owlight"
		}
		if frontendFolderStatic == "" {
			frontendFolderStatic = "./owlight/static"
		}
		log.Infof("frontend started with path: %s, static: %s\n", frontendFolder, frontendFolderStatic)
		r.Static("web", frontendFolder)
		r.Static("static", frontendFolderStatic)
	}
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, I'm OWL (｡A｡)")
		return
	})

	graph.Routes(r)
	uic.Routes(r)
	template.Routes(r)
	strategy.Routes(r)
	host.Routes(r)
	expression.Routes(r)
	mockcfg.Routes(r)
	// dashboardScreenOWl.Routes(r)
	// dashboardGraphOwl.Routes(r)
	dashboard_graph.Routes(r)
	dashboard_screen.Routes(r)
	alarm.Routes(r)
	if !testMode {
		r.Run(port)
	}
	return r
}
