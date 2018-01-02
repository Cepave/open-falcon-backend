package rpc

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	cutils "github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	"github.com/Cepave/open-falcon-backend/modules/transfer/sender"
)

type Transfer int

type TransferResp struct {
	Msg        string
	Total      int
	ErrInvalid int
	Latency    int64
}

func (t *TransferResp) String() string {
	s := fmt.Sprintf("TransferResp total=%d, err_invalid=%d, latency=%dms",
		t.Total, t.ErrInvalid, t.Latency)
	if t.Msg != "" {
		s = fmt.Sprintf("%s, msg=%s", s, t.Msg)
	}
	return s
}

func (this *Transfer) Ping(req cmodel.NullRpcRequest, resp *cmodel.SimpleRpcResponse) error {
	return nil
}

func (t *Transfer) Update(args []*cmodel.MetricValue, reply *cmodel.TransferResponse) error {
	return RecvMetricValues(args, reply, "rpc")
}

// process new metric values
func RecvMetricValues(args []*cmodel.MetricValue, reply *cmodel.TransferResponse, from string) error {
	start := time.Now()
	reply.Invalid = 0

	cfg := g.Config()
	filters := cfg.Staging.Filters
	items := []*cmodel.MetaData{}
	stagingItems := []*cmodel.MetricValue{}
	for _, v := range args {
		if v == nil {
			reply.Invalid += 1
			continue
		}

		// 历史遗留问题.
		// 老版本agent上报的metric=kernel.hostname的数据,其取值为string类型,现在已经不支持了;所以,这里硬编码过滤掉
		if v.Metric == "kernel.hostname" {
			reply.Invalid += 1
			continue
		}

		if v.Metric == "" || v.Endpoint == "" {
			reply.Invalid += 1
			continue
		}

		if v.Type != g.COUNTER && v.Type != g.GAUGE && v.Type != g.DERIVE {
			reply.Invalid += 1
			continue
		}

		if v.Value == "" {
			reply.Invalid += 1
			continue
		}

		if v.Step <= 0 {
			reply.Invalid += 1
			continue
		}

		if len(v.Metric)+len(v.Tags) > 510 {
			reply.Invalid += 1
			continue
		}

		// TODO 呵呵,这里需要再优雅一点
		now := start.Unix()
		if v.Timestamp <= 0 || v.Timestamp > now*1+7200 {
			v.Timestamp = now
		}

		fv := &cmodel.MetaData{
			Metric:      v.Metric,
			Endpoint:    v.Endpoint,
			Timestamp:   v.Timestamp,
			Step:        v.Step,
			CounterType: v.Type,
			Tags:        cutils.DictedTagstring(v.Tags), //TODO tags键值对的个数,要做一下限制
		}

		valid := true
		var vv float64
		var err error

		switch cv := v.Value.(type) {
		case string:
			vv, err = strconv.ParseFloat(cv, 64)
			if err != nil {
				valid = false
			}
		case float64:
			vv = cv
		case int64:
			vv = float64(cv)
		default:
			valid = false
		}

		if !valid {
			reply.Invalid += 1
			continue
		}

		fv.Value = vv
		items = append(items, fv)

		// Filter Staging items through endpoint
		if cfg.Staging.Enabled {
			for _, filter := range filters {
				if strings.HasPrefix(v.Endpoint, filter) {
					sv := &cmodel.MetricValue{
						Endpoint:  v.Endpoint,
						Metric:    v.Metric,
						Value:     v.Value,
						Step:      v.Step,
						Type:      v.Type,
						Tags:      v.Tags,
						Timestamp: v.Timestamp,
					}
					stagingItems = append(stagingItems, sv)
					break
				}
			}
		}
	}

	// statistics
	cnt := int64(len(items))
	proc.RecvCnt.IncrBy(cnt)
	if from == "rpc" {
		proc.RpcRecvCnt.IncrBy(cnt)
	} else if from == "http" {
		proc.HttpRecvCnt.IncrBy(cnt)
	}

	// demultiplexing
	nqmFpingItems, nqmTcppingItems, nqmTcpconnItems, genericItems := sender.Demultiplex(items)

	if cfg.Staging.Enabled {
		sender.Push2StagingSendQueue(stagingItems)
	}

	if cfg.Graph.Enabled {
		sender.Push2GraphSendQueue(genericItems)
	}

	if cfg.Judge.Enabled {
		sender.Push2JudgeSendQueue(genericItems)
	}

	if cfg.Tsdb.Enabled {
		sender.Push2TsdbSendQueue(genericItems)
	}

	if cfg.Influxdb.Enabled {
		sender.Push2InfluxdbSendQueue(genericItems)
	}

	if cfg.NqmRest.Enabled {
		sender.Push2NqmIcmpSendQueue(nqmFpingItems)
		sender.Push2NqmTcpSendQueue(nqmTcppingItems)
		sender.Push2NqmTcpconnSendQueue(nqmTcpconnItems)
	}

	reply.Message = "ok"
	reply.Total = len(args)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}
