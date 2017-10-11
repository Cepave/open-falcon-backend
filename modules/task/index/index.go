package index

import (
	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/task/g"
)

// 初始化索引功能模块
func Start() {
	cfg := g.Config()
	if !cfg.Index.Enable {
		log.Println("index.Start warning, not enable")
		return
	}

	InitDB()
	StartIndexUpdateAllTask()
	log.Println("index.Start ok")
}
