package http

import (
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/toolkits/file"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"strings"
)

func configCommonRoutes(router *gin.Engine) {
	router.GET("/health", gin.WrapF(health))
	router.GET("/version", gin.WrapF(version))
	router.GET("/workdir", gin.WrapF(workDir))
	router.GET("/config/reload", gin.WrapF(reloadConfig))
}

// Checks health of web service
func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// Gets version of HBS
func version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(g.VERSION))
}

// Gets the working dir of current application
func workDir(w http.ResponseWriter, r *http.Request) {
	RenderDataJson(w, file.SelfDir())
}

// Reload the configuration
func reloadConfig(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
		g.ParseConfig(g.ConfigFile)
		RenderDataJson(w, g.Config())
	} else {
		w.Write([]byte("no privilege"))
	}
}
