package rpc

import (
	"bytes"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/rpc"
	"github.com/Cepave/open-falcon-backend/common/utils"

	"github.com/Cepave/open-falcon-backend/modules/hbs/cache"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	hbsService "github.com/Cepave/open-falcon-backend/modules/hbs/service"
)

var agentHeartbeatService *hbsService.AgentHeartbeatService

func (t *Agent) MinePlugins(args model.AgentHeartbeatRequest, reply *model.AgentPluginsResponse) (err error) {
	defer rpc.HandleError(&err)()

	if args.Hostname == "" {
		return nil
	}

	reply.Plugins = cache.GetPlugins(args.Hostname)
	reply.Timestamp = time.Now().Unix()
	reply.GitRepo = cache.GitRepo.Get()
	// deprecate the attributes: reply.GitUpdate, reply.GitRepoUpdate
	// git repo updating will be invoked only by reply.GitRepo
	reply.GitUpdate = false
	reply.GitRepoUpdate = false
	log.Debugln("show reply of MinePlugins: ", reply)

	return nil
}

func (t *Agent) ReportStatus(args *model.AgentReportRequest, reply *model.SimpleRpcResponse) (err error) {
	defer rpc.HandleError(&err)()

	if args.Hostname == "" {
		reply.Code = 1
		return nil
	}

	log.Debugln("show request of ReportStatus: ", args)
	agentHeartbeatService.Put(args, time.Now().Unix())

	return nil
}

// 需要checksum一下来减少网络开销？其实白名单通常只会有一个或者没有，无需checksum
func (t *Agent) TrustableIps(args *model.NullRpcRequest, ips *string) (err error) {
	defer rpc.HandleError(&err)()

	*ips = strings.Join(g.Config().Trustable, ",")
	return nil
}

// agent按照server端的配置，按需采集的metric，比如net.port.listen port=22 或者 proc.num name=zabbix_agentd
func (t *Agent) BuiltinMetrics(args *model.AgentHeartbeatRequest, reply *model.BuiltinMetricResponse) (err error) {
	defer rpc.HandleError(&err)()

	if args.Hostname == "" {
		return nil
	}

	metrics, err := cache.GetBuiltinMetrics(args.Hostname)
	if err != nil {
		return nil
	}

	checksum := ""
	if len(metrics) > 0 {
		checksum = DigestBuiltinMetrics(metrics)
	}

	if args.Checksum == checksum {
		reply.Metrics = []*model.BuiltinMetric{}
	} else {
		reply.Metrics = metrics
	}
	reply.Checksum = checksum
	reply.Timestamp = time.Now().Unix()

	return nil
}

func DigestBuiltinMetrics(items []*model.BuiltinMetric) string {
	sort.Sort(model.BuiltinMetricSlice(items))

	var buf bytes.Buffer
	for _, m := range items {
		buf.WriteString(m.String())
	}

	return utils.Md5(buf.String())
}
