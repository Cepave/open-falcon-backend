package funcs

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"
)

func SocketStatSummaryMetrics() (L []*model.MetricValue) {
	ssMap, err := nux.SocketStatSummary()
	if err != nil {
		log.Println(err)
		return
	}

	for k, v := range ssMap {
		L = append(L, GaugeValue("ss."+k, v))
	}

	return
}
