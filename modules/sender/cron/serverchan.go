package cron

import (
	"github.com/Cepave/open-falcon-backend/modules/sender/g"
	"github.com/Cepave/open-falcon-backend/modules/sender/model"
	"github.com/Cepave/open-falcon-backend/modules/sender/proc"
	"github.com/Cepave/open-falcon-backend/modules/sender/redis"
	"github.com/toolkits/net/httplib"
	"log"
	"time"
)

func ConsumeServerchan() {
	queue := g.Config().Queue.Serverchan
	for {
		L := redis.PopAllServerchan(queue)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendServerchanList(L)
	}
}

func SendServerchanList(L []*model.Serverchan) {
	for _, serverchan := range L {
		ServerchanWorkerChan <- 1
		go SendServerchan(serverchan)
	}
}

func SendServerchan(serverchan *model.Serverchan) {
	defer func() {
		<-ServerchanWorkerChan
	}()

	sckey := serverchan.Tos
	if len(sckey) > 5 {
		url := g.Config().Api.Serverchan
		url += "/" + sckey + ".send"
		r := httplib.Post(url).SetTimeout(5*time.Second, 2*time.Minute)
		r.Param("text", serverchan.Subject)
		r.Param("desp", serverchan.Content)
		resp, err := r.String()
		if err != nil {
			log.Println(err)
		}

		if g.Config().Debug {
			log.Println("==serverchan==>>>>", serverchan)
			log.Println("<<<<==serverchan==", resp)
		}
	}
	proc.IncreServerchanCount()
}
