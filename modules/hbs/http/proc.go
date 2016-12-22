package http

import (
	"github.com/Cepave/open-falcon-backend/modules/hbs/cache"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func configProcRoutes(router *gin.Engine) {
	router.GET("/expressions", gin.WrapF(expressions))
	router.GET("/plugins/", gin.WrapF(plugins))
}

func expressions(w http.ResponseWriter, r *http.Request) {
	RenderDataJson(w, cache.ExpressionCache.Get())
}

func plugins(w http.ResponseWriter, r *http.Request) {
	hostname := r.URL.Path[len("/plugins/"):]
	RenderDataJson(w, cache.GetPlugins(hostname))
}
