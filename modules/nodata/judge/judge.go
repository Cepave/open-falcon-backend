package judge

import (
	log "github.com/sirupsen/logrus"
)

func Start() {
	StartJudgeCron()
	log.Println("judge.Start ok")
}
