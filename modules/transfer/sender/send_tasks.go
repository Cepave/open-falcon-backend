package sender

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	nsema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"
	nproc "github.com/toolkits/proc"
)

// send
const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

// TODO 添加对发送任务的控制,比如stop等
func startSendTasks() {
	cfg := g.Config()
	// init semaphore
	judgeConcurrent := cfg.Judge.MaxConns
	graphConcurrent := cfg.Graph.MaxConns
	tsdbConcurrent := cfg.Tsdb.MaxConns
	influxdbConcurrent := cfg.Influxdb.MaxIdle

	if tsdbConcurrent < 1 {
		tsdbConcurrent = 1
	}

	if judgeConcurrent < 1 {
		judgeConcurrent = 1
	}

	if graphConcurrent < 1 {
		graphConcurrent = 1
	}
	if influxdbConcurrent < 1 {
		influxdbConcurrent = 1
	}

	// init send go-routines
	for node := range cfg.Judge.Cluster {
		queue := JudgeQueues[node]
		go forward2JudgeTask(queue, node, judgeConcurrent)
	}

	for node, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			queue := GraphQueues[node+addr]
			go forward2GraphTask(queue, node, addr, graphConcurrent)
		}
	}

	if cfg.Tsdb.Enabled {
		go forward2TsdbTask(tsdbConcurrent)
	}

	go forward2InfluxdbTask(InfluxdbQueues["default"], influxdbConcurrent)

	if cfg.NqmRest.Enabled {
		go forward2NqmTask(NqmIcmpQueue, g.Config().NqmRest.Fping, proc.SendToNqmIcmpCnt, proc.SendToNqmIcmpFailCnt)
		go forward2NqmTask(NqmTcpQueue, g.Config().NqmRest.Tcpping, proc.SendToNqmTcpCnt, proc.SendToNqmTcpFailCnt)
		go forward2NqmTask(NqmTcpconnQueue, g.Config().NqmRest.Tcpconn, proc.SendToNqmTcpconnCnt, proc.SendToNqmTcpconnFailCnt)
	}

	if cfg.Staging.Enabled {
		go forward2StagingTask()
	}
}

// Judge定时任务, 将 Judge发送缓存中的数据 通过rpc连接池 发送到Judge
func forward2JudgeTask(Q *list.SafeListLimited, node string, concurrent int) {
	batch := g.Config().Judge.Batch // 一次发送,最多batch条数据
	addr := g.Config().Judge.Cluster[node]
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		judgeItems := make([]*cmodel.JudgeItem, count)
		for i := 0; i < count; i++ {
			judgeItems[i] = items[i].(*cmodel.JudgeItem)
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(addr string, judgeItems []*cmodel.JudgeItem, count int) {
			defer sema.Release()

			resp := &cmodel.SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = JudgeConnPools.Call(addr, "Judge.Send", judgeItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Errorf("send judge %s:%s fail: %v", node, addr, err)
				proc.SendToJudgeFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToJudgeCnt.IncrBy(int64(count))
			}
		}(addr, judgeItems, count)
	}
}

// Graph定时任务, 将 Graph发送缓存中的数据 通过rpc连接池 发送到Graph
func forward2GraphTask(Q *list.SafeListLimited, node string, addr string, concurrent int) {
	batch := g.Config().Graph.Batch // 一次发送,最多batch条数据
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		graphItems := make([]*cmodel.GraphItem, count)
		for i := 0; i < count; i++ {
			graphItems[i] = items[i].(*cmodel.GraphItem)
		}

		sema.Acquire()
		go func(addr string, graphItems []*cmodel.GraphItem, count int) {
			defer sema.Release()

			resp := &cmodel.SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = GraphConnPools.Call(addr, "Graph.Send", graphItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Errorf("send to graph %s:%s fail: %v", node, addr, err)
				proc.SendToGraphFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToGraphCnt.IncrBy(int64(count))
			}
		}(addr, graphItems, count)
	}
}

// Tsdb定时任务, 将数据通过api发送到tsdb
func forward2TsdbTask(concurrent int) {
	batch := g.Config().Tsdb.Batch // 一次发送,最多batch条数据
	retry := g.Config().Tsdb.MaxRetry
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := TsdbQueue.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		//  同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(itemList []interface{}) {
			defer sema.Release()

			var tsdbBuffer bytes.Buffer
			for i := 0; i < len(itemList); i++ {
				tsdbItem := itemList[i].(*cmodel.TsdbItem)
				tsdbBuffer.WriteString(tsdbItem.TsdbString())
				tsdbBuffer.WriteString("\n")
			}

			var err error
			for i := 0; i < retry; i++ {
				err = TsdbConnPoolHelper.Send(tsdbBuffer.Bytes())
				if err == nil {
					proc.SendToTsdbCnt.IncrBy(int64(len(itemList)))
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			if err != nil {
				proc.SendToTsdbFailCnt.IncrBy(int64(len(itemList)))
				log.Errorf("%v", err)
				return
			}
		}(items)
	}
}

// Influxdb schedule
func forward2InfluxdbTask(Q *list.SafeListLimited, concurrent int) {
	cfg := g.Config().Influxdb
	batch := cfg.Batch // 一次发送,最多batch条数据
	conn, err := parseDSN(cfg.Address)
	if err != nil {
		log.Errorf("syntax of influxdb address is wrong: %v", err)
		return
	}
	addr := conn.Address

	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		influxdbItems := make([]*cmodel.JudgeItem, count)
		for i := 0; i < count; i++ {
			influxdbItems[i] = items[i].(*cmodel.JudgeItem)
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(addr string, influxdbItems []*cmodel.JudgeItem, count int) {
			defer sema.Release()

			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = InfluxdbConnPools.Call(addr, influxdbItems)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Errorf("send influxdb %s fail: %v", addr, err)
				proc.SendToInfluxdbFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToInfluxdbCnt.IncrBy(int64(count))
			}
		}(addr, influxdbItems, count)
	}
}

func forwardNqmItems(nqmItem interface{}, nqmUrl string, s *nsema.Semaphore, cnt *nproc.SCounterQps, failCnt *nproc.SCounterQps) {
	defer s.Release()

	jsonItem, jsonErr := json.Marshal(nqmItem)
	if jsonErr != nil {
		log.Errorf("Error on serialization for nqm item(ICMP, TCP, or TCPCONN): %v", jsonErr)
		failCnt.IncrBy(1)
		return
	}

	log.Debugf("[ Cassandra ] JSON data to %s: %s", nqmUrl, string(jsonItem))
	postReq, err := http.NewRequest("POST", nqmUrl, bytes.NewBuffer(jsonItem))

	postReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	postReq.Header.Set("Connection", "close")
	httpClient := &http.Client{}
	postResp, err := httpClient.Do(postReq)
	if err != nil {
		log.Errorf("[ Cassandra ] Error on push: %v", err)
		failCnt.IncrBy(1)
		return
	}
	defer postResp.Body.Close()
	cnt.IncrBy(1)
}

func forward2NqmTask(Q *list.SafeListLimited, apiUrl string, cnt *nproc.SCounterQps, failCnt *nproc.SCounterQps) {
	batch := g.Config().NqmRest.Batch // 一次发送,最多batch条数据
	concurrent := g.Config().NqmRest.MaxConns
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		for _, v := range items {
			sema.Acquire()
			go forwardNqmItems(v, apiUrl, sema, cnt, failCnt)
		}
	}
}

func forward2StagingTask() {
	batch := g.Config().Staging.Batch
	retry := g.Config().Staging.MaxRetry
	concurrent := g.Config().Staging.MaxConns
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := StagingQueue.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		stagingItems := make([]*cmodel.MetricValue, count)
		for i := 0; i < count; i++ {
			stagingItems[i] = items[i].(*cmodel.MetricValue)
		}

		//	A synchronous call with limited concurrence
		sema.Acquire()
		go func(stagingItems []*cmodel.MetricValue, count int) {
			defer sema.Release()

			resp := &cmodel.SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < retry; i++ {
				err = StagingConnPoolHelper.Call("Transfer.Update", stagingItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Errorf("send staging fail: %v", err)
				proc.SendToStagingFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToStagingCnt.IncrBy(int64(count))
			}
		}(stagingItems, count)
	}
}
