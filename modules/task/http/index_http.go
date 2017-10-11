package http

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/task/index"
	srv "github.com/Cepave/open-falcon-backend/modules/task/service"
)

func configIndexHttpRoutes() {
	http.HandleFunc("/index/delete", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			p := recover()
			if p != nil {
				logger.Warn("Delete index of graph failed: %v", p)
			}
		}()

		srv.VacuumGraphIndex(7)

		RenderDataJson(w, "ok")
	})
	http.HandleFunc("/index/updateAll", func(w http.ResponseWriter, r *http.Request) {
		index.UpdateAllIndex()
		RenderDataJson(w, "ok")
	})
}
