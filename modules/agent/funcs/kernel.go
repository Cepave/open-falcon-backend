package funcs

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/nux"
)

func KernelMetrics() (L []*model.MetricValue) {

	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Println(err)
		return
	}

	L = append(L, GaugeValue("kernel.maxfiles", maxFiles))

	maxProc, err := nux.KernelMaxProc()
	if err != nil {
		log.Println(err)
		return
	}

	L = append(L, GaugeValue("kernel.maxproc", maxProc))

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Println(err)
		return
	}

	L = append(L, GaugeValue("kernel.files.allocated", allocateFiles))
	L = append(L, GaugeValue("kernel.files.left", maxFiles-allocateFiles))
	return
}
