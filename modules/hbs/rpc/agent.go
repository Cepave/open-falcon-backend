package rpc

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/rpc"

	"github.com/Cepave/open-falcon-backend/modules/hbs/service"
)

var agentHeartbeatService *service.AgentHeartbeatService

func (t *Agent) MinePlugins(args model.AgentHeartbeatRequest, reply *model.AgentPluginsResponse) (err error) {
	if args.Hostname == "" {
		return nil
	}

	resp, err := service.MinePlugins(args.Hostname, args.Checksum)
	reply.Plugins = resp.Plugins
	reply.Timestamp = resp.Timestamp
	reply.GitRepo = resp.GitRepo
	reply.GitUpdate = resp.GitUpdate
	reply.GitRepoUpdate = resp.GitRepoUpdate
	log.Debugln("show reply of MinePlugins: ", reply)

	return err
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

// agent按照server端的配置，按需采集的metric，比如net.port.listen port=22 或者 proc.num name=zabbix_agentd
func (t *Agent) BuiltinMetrics(args *model.AgentHeartbeatRequest, reply *model.BuiltinMetricResponse) (err error) {
	defer rpc.HandleError(&err)()
	resp, err := service.BuiltinMetrics(args.Hostname, args.Checksum)

	reply.Checksum = resp.Checksum
	reply.Timestamp = resp.Timestamp
	for _, nm := range resp.Metrics {
		om := &model.BuiltinMetric{
			Metric: nm.Metric,
			Tags:   nm.Tags,
		}
		reply.Metrics = append(reply.Metrics, om)
	}

	return nil
}
