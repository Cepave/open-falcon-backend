package http

import (
	"encoding/json"
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	"net/http"
)

func configPushRoutes() {
	http.HandleFunc("/v1/push", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "body is blank", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var metrics []*model.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "connot decode body", http.StatusBadRequest)
			return
		}

		toTransfer, toMQ := g.DemultiplexMetrics(metrics)
		g.SendToTransfer(toTransfer)
		g.SendToMQ(toMQ)
		w.Write([]byte("success"))
	})
}
