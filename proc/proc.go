package proc

import (
	nproc "github.com/toolkits/proc"
	"log"
)

// trace
var (
	RecvDataTrace = nproc.NewDataTrace("RecvDataTrace", 3)
)

// filter
var (
	RecvDataFilter = nproc.NewDataFilter("RecvDataFilter", 5)
)

// 统计指标的整体数据
var (
	// 计数统计,正确计数,错误计数, ...
	RecvCnt       = nproc.NewSCounterQps("RecvCnt")
	RpcRecvCnt    = nproc.NewSCounterQps("RpcRecvCnt")
	HttpRecvCnt   = nproc.NewSCounterQps("HttpRecvCnt")
	SocketRecvCnt = nproc.NewSCounterQps("SocketRecvCnt")

	SendToJudgeCnt    = nproc.NewSCounterQps("SendToJudgeCnt")
	SendToTsdbCnt     = nproc.NewSCounterQps("SendToTsdbCnt")
	SendToGraphCnt    = nproc.NewSCounterQps("SendToGraphCnt")
	SendToInfluxdbCnt = nproc.NewSCounterQps("SendToInfluxdbCnt")
	SendToNqmRpcCnt   = nproc.NewSCounterQps("SendToNqmRpcCnt")

	SendToJudgeDropCnt    = nproc.NewSCounterQps("SendToJudgeDropCnt")
	SendToTsdbDropCnt     = nproc.NewSCounterQps("SendToTsdbDropCnt")
	SendToGraphDropCnt    = nproc.NewSCounterQps("SendToGraphDropCnt")
	SendToInfluxdbDropCnt = nproc.NewSCounterQps("SendToInfluxdbDropCnt")
	SendToNqmRpcDropCnt   = nproc.NewSCounterQps("SendToNqmRpcDropCnt")

	SendToJudgeFailCnt    = nproc.NewSCounterQps("SendToJudgeFailCnt")
	SendToTsdbFailCnt     = nproc.NewSCounterQps("SendToTsdbFailCnt")
	SendToGraphFailCnt    = nproc.NewSCounterQps("SendToGraphFailCnt")
	SendToInfluxdbFailCnt = nproc.NewSCounterQps("SendToInfluxdbFailCnt")
	SendToNqmRpcFailCnt   = nproc.NewSCounterQps("SendToNqmRpcFailCnt")

	// 发送缓存大小
	JudgeQueuesCnt    = nproc.NewSCounterBase("JudgeSendCacheCnt")
	TsdbQueuesCnt     = nproc.NewSCounterBase("TsdbSendCacheCnt")
	GraphQueuesCnt    = nproc.NewSCounterBase("GraphSendCacheCnt")
	InfluxdbQueuesCnt = nproc.NewSCounterBase("InfluxdbSendCacheCnt")
	NqmRpcQueuesCnt   = nproc.NewSCounterBase("NqmRpcSendCacheCnt")
)

func Start() {
	log.Println("proc.Start, ok")
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// recv cnt
	ret = append(ret, RecvCnt.Get())
	ret = append(ret, RpcRecvCnt.Get())
	ret = append(ret, HttpRecvCnt.Get())
	ret = append(ret, SocketRecvCnt.Get())

	// send cnt
	ret = append(ret, SendToJudgeCnt.Get())
	ret = append(ret, SendToTsdbCnt.Get())
	ret = append(ret, SendToGraphCnt.Get())
	ret = append(ret, SendToInfluxdbCnt.Get())
	ret = append(ret, SendToNqmRpcCnt.Get())

	// drop cnt
	ret = append(ret, SendToJudgeDropCnt.Get())
	ret = append(ret, SendToTsdbDropCnt.Get())
	ret = append(ret, SendToGraphDropCnt.Get())
	ret = append(ret, SendToInfluxdbDropCnt.Get())
	ret = append(ret, SendToNqmRpcDropCnt.Get())

	// send fail cnt
	ret = append(ret, SendToJudgeFailCnt.Get())
	ret = append(ret, SendToTsdbFailCnt.Get())
	ret = append(ret, SendToGraphFailCnt.Get())
	ret = append(ret, SendToInfluxdbFailCnt.Get())
	ret = append(ret, SendToNqmRpcFailCnt.Get())

	// cache cnt
	ret = append(ret, JudgeQueuesCnt.Get())
	ret = append(ret, TsdbQueuesCnt.Get())
	ret = append(ret, GraphQueuesCnt.Get())
	ret = append(ret, InfluxdbQueuesCnt.Get())
	ret = append(ret, NqmRpcQueuesCnt.Get())

	return ret
}
