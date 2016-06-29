package http

import (
	"github.com/Cepave/hbs/cache"
	"net/http"
)

func configProcRoutes() {
	http.HandleFunc("/expressions", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, cache.ExpressionCache.Get())
	})

	http.HandleFunc("/plugins/", func(w http.ResponseWriter, r *http.Request) {
		hostname := r.URL.Path[len("/plugins/"):]
		RenderDataJson(w, cache.GetPlugins(hostname))
	})

}
