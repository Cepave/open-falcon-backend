package sender

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Cepave/transfer/g"
	"github.com/Cepave/transfer/proc"
	cpool "github.com/Cepave/transfer/sender/conn_pool"
	cmodel "github.com/open-falcon/common/model"
	nlist "github.com/toolkits/container/list"
)

const (
	DefaultSendQueueMaxSize = 102400 //10.24w
)

// 默认参数
var (
	MinStep int //最小上报周期,单位sec
)

// 服务节点的一致性哈希环
// pk -> node
var (
	JudgeNodeRing *ConsistentHashNodeRing
	GraphNodeRing *ConsistentHashNodeRing
)

// 发送缓存队列
// node -> queue_of_data
var (
	TsdbQueue      *nlist.SafeListLimited
	JudgeQueues    = make(map[string]*nlist.SafeListLimited)
	GraphQueues    = make(map[string]*nlist.SafeListLimited)
	InfluxdbQueues = make(map[string]*nlist.SafeListLimited)
	NqmRpcQueue    *nlist.SafeListLimited
)

// 连接池
// node_address -> connection_pool
var (
	JudgeConnPools       *cpool.SafeRpcConnPools
	TsdbConnPoolHelper   *cpool.TsdbConnPoolHelper
	GraphConnPools       *cpool.SafeRpcConnPools
	InfluxdbConnPools    *cpool.InfluxdbConnPools
	NqmRpcConnPoolHelper *cpool.NqmRpcConnPoolHelper
)

// 初始化数据发送服务, 在main函数中调用
func Start() {
	// 初始化默认参数
	MinStep = g.Config().MinStep
	if MinStep < 1 {
		MinStep = 30 //默认30s
	}
	//
	initConnPools()
	initSendQueues()
	initNodeRings()
	// SendTasks依赖基础组件的初始化,要最后启动
	startSendTasks()
	startSenderCron()
	log.Println("send.Start, ok")
}

// 将数据 打入 某个Judge的发送缓存队列, 具体是哪一个Judge 由一致性哈希 决定
func Push2JudgeSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		pk := item.PK()
		node, err := JudgeNodeRing.GetNode(pk)
		if err != nil {
			log.Println("E:", err)
			continue
		}

		// align ts
		step := int(item.Step)
		if step < MinStep {
			step = MinStep
		}
		ts := alignTs(item.Timestamp, int64(step))

		judgeItem := &cmodel.JudgeItem{
			Endpoint:  item.Endpoint,
			Metric:    item.Metric,
			Value:     item.Value,
			Timestamp: ts,
			JudgeType: item.CounterType,
			Tags:      item.Tags,
		}
		Q := JudgeQueues[node]
		isSuccess := Q.PushFront(judgeItem)

		// statistics
		if !isSuccess {
			proc.SendToJudgeDropCnt.Incr()
		}
	}
}

// 将数据 打入 某个Graph的发送缓存队列, 具体是哪一个Graph 由一致性哈希 决定
func Push2GraphSendQueue(items []*cmodel.MetaData) {
	cfg := g.Config().Graph

	for _, item := range items {
		graphItem, err := convert2GraphItem(item)
		if err != nil {
			log.Println("E:", err)
			continue
		}
		pk := item.PK()

		// statistics. 为了效率,放到了这里,因此只有graph是enbale时才能trace
		proc.RecvDataTrace.Trace(pk, item)
		proc.RecvDataFilter.Filter(pk, item.Value, item)

		node, err := GraphNodeRing.GetNode(pk)
		if err != nil {
			log.Println("E:", err)
			continue
		}

		cnode := cfg.ClusterList[node]
		errCnt := 0
		for _, addr := range cnode.Addrs {
			Q := GraphQueues[node+addr]
			if !Q.PushFront(graphItem) {
				errCnt += 1
			}
		}

		// statistics
		if errCnt > 0 {
			proc.SendToGraphDropCnt.Incr()
		}
	}
}

// 打到Graph的数据,要根据rrdtool的特定 来限制 step、counterType、timestamp
func convert2GraphItem(d *cmodel.MetaData) (*cmodel.GraphItem, error) {
	item := &cmodel.GraphItem{}

	item.Endpoint = d.Endpoint
	item.Metric = d.Metric
	item.Tags = d.Tags
	item.Timestamp = d.Timestamp
	item.Value = d.Value
	item.Step = int(d.Step)
	if item.Step < MinStep {
		item.Step = MinStep
	}
	item.Heartbeat = item.Step * 2

	if d.CounterType == g.GAUGE {
		item.DsType = d.CounterType
		item.Min = "U"
		item.Max = "U"
	} else if d.CounterType == g.COUNTER {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else if d.CounterType == g.DERIVE {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else {
		return item, fmt.Errorf("not_supported_counter_type")
	}

	item.Timestamp = alignTs(item.Timestamp, int64(item.Step)) //item.Timestamp - item.Timestamp%int64(item.Step)

	return item, nil
}

// 将原始数据入到tsdb发送缓存队列
func Push2TsdbSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		tsdbItem := convert2TsdbItem(item)
		isSuccess := TsdbQueue.PushFront(tsdbItem)

		if !isSuccess {
			proc.SendToTsdbDropCnt.Incr()
		}
	}
}

// 转化为tsdb格式
func convert2TsdbItem(d *cmodel.MetaData) *cmodel.TsdbItem {
	t := cmodel.TsdbItem{Tags: make(map[string]string)}

	for k, v := range d.Tags {
		t.Tags[k] = v
	}
	t.Tags["endpoint"] = d.Endpoint
	t.Metric = d.Metric
	t.Timestamp = d.Timestamp
	t.Value = d.Value
	return &t
}

func alignTs(ts int64, period int64) int64 {
	return ts - ts%period
}

// Push data to 3rd-party database
func Push2InfluxdbSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		// align ts
		step := int(item.Step)
		if step < MinStep {
			step = MinStep
		}
		ts := alignTs(item.Timestamp, int64(step))

		influxdbItem := &cmodel.JudgeItem{
			Endpoint:  item.Endpoint,
			Metric:    item.Metric,
			Value:     item.Value,
			Timestamp: ts,
			JudgeType: item.CounterType,
			Tags:      item.Tags,
		}
		Q := InfluxdbQueues["default"]
		isSuccess := Q.PushFront(influxdbItem)

		// statistics
		if !isSuccess {
			proc.SendToInfluxdbDropCnt.Incr()
		}
	}
}

// Push network quality metric pkt to the queue for RPC
func Push2NqmRpcSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		nqmitem, err := convert2NqmRpcItem(item)
		if err != nil {
			log.Println("NqmRpc converting error:", err)
			continue
		}
		isSuccess := NqmRpcQueue.PushFront(nqmitem)

		if !isSuccess {
			proc.SendToNqmRpcDropCnt.Incr()
		}
	}
}

func Demultiplex(items []*cmodel.MetaData) ([]*cmodel.MetaData, []*cmodel.MetaData) {
	nqms := []*cmodel.MetaData{}
	generics := []*cmodel.MetaData{}

	for _, item := range items {
		if strings.HasPrefix(item.Metric, "nqm-") {
			nqms = append(nqms, item)
		} else {
			generics = append(generics, item)
		}
	}

	return nqms, generics
}

func convert2NqmRpcItem(d *cmodel.MetaData) (*nqmRpcItem, error) {
	var t nqmRpcItem
	agent, err := convert2NqmEndpoint(d, "agent")
	if err != nil {
		return &t, err
	}
	target, err := convert2NqmEndpoint(d, "target")
	if err != nil {
		return &t, err
	}
	metrics, err := convert2NqmMetrics(d)
	if err != nil {
		return &t, err
	}

	t = nqmRpcItem{
		Timestamp: d.Timestamp,
		Agent:     *agent,
		Target:    *target,
		Metrics:   *metrics,
	}

	return &t, nil
}

func strToFloat32(out *float32, index string, dict map[string]string) error {
	var err error
	var ff float64
	if v, ok := dict[index]; ok {
		ff, err = strconv.ParseFloat(v, 32)
		if err != nil {
			return err
		}
		*out = float32(ff)
	}
	return nil
}

func strToInt32(out *int32, index string, dict map[string]string) error {
	var err error
	var ii int64
	if v, ok := dict[index]; ok {
		ii, err = strconv.ParseInt(v, 10, 32)
		if err != nil {
			return err
		}
		*out = int32(ii)
	}
	return nil
}

func strToInt16(out *int16, index string, dict map[string]string) error {
	var err error
	var ii int64
	if v, ok := dict[index]; ok {
		ii, err = strconv.ParseInt(v, 10, 16)
		if err != nil {
			return err
		}
		*out = int16(ii)
	}
	return nil
}

func convert2NqmEndpoint(d *cmodel.MetaData, endType string) (*nqmEndpoint, error) {
	t := nqmEndpoint{
		Id:         -1,
		IspId:      -1,
		ProvinceId: -1,
		CityId:     -1,
		NameTagId:  -1,
	}

	if err := strToInt32(&t.Id, endType+"-id", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt16(&t.IspId, endType+"-isp-id", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt16(&t.ProvinceId, endType+"-province-id", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt16(&t.CityId, endType+"-city-id", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt32(&t.NameTagId, endType+"-name-tag-id", d.Tags); err != nil {
		return nil, err
	}

	return &t, nil
}

// 轉化成 nqmMetrc 格式
func convert2NqmMetrics(d *cmodel.MetaData) (*nqmMetrics, error) {
	t := nqmMetrics{
		Rttmin:      -1,
		Rttavg:      -1,
		Rttmax:      -1,
		Rttmdev:     -1,
		Rttmedian:   -1,
		Pkttransmit: -1,
		Pktreceive:  -1,
	}
	var ff float32
	if err := strToFloat32(&ff, "rttmin", d.Tags); err != nil {
		return nil, err
	}
	t.Rttmin = int32(ff)
	if err := strToFloat32(&ff, "rttmax", d.Tags); err != nil {
		return nil, err
	}
	t.Rttmax = int32(ff)

	if err := strToFloat32(&t.Rttavg, "rttavg", d.Tags); err != nil {
		return nil, err
	}
	if err := strToFloat32(&t.Rttmdev, "rttmdev", d.Tags); err != nil {
		return nil, err
	}
	if err := strToFloat32(&t.Rttmedian, "rttmedian", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt32(&t.Pkttransmit, "pkttransmit", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt32(&t.Pktreceive, "pktreceive", d.Tags); err != nil {
		return nil, err
	}

	return &t, nil
}
