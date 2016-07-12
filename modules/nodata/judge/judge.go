package judge

import (
	log "github.com/Sirupsen/logrus"
)

func Start() {
	StartJudgeCron()
	log.Println("judge.Start ok")
}
