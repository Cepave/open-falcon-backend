package cron

import (
	"github.com/Cepave/open-falcon-backend/modules/sender/g"
)

var (
	SmsWorkerChan  chan int
	MailWorkerChan chan int
	QQWorkerChan   chan int
	ServerchanWorkerChan chan int
)

func InitWorker() {
	workerConfig := g.Config().Worker
	SmsWorkerChan = make(chan int, workerConfig.Sms)
	MailWorkerChan = make(chan int, workerConfig.Mail)
	QQWorkerChan = make(chan int, workerConfig.QQ)
	ServerchanWorkerChan = make(chan int, workerConfig.Serverchan)
}
