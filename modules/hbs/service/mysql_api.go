package service

import (
	"fmt"
	"net/http"

	"github.com/juju/errors"

	"github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	commonSling "github.com/Cepave/open-falcon-backend/common/sling"
)

var annotateErr = errors.Annotate

func agentHeartbeatCall(agents []*model.FalconAgentHeartbeat) (rowsAffectedCnt int64, agentsDroppedCnt int64) {
	param := struct {
		UpdateOnly bool `json:"update_only"`
	}{updateOnlyFlag}
	req := NewSlingBase().Post("api/v1/agent/heartbeat").BodyJSON(agents).QueryStruct(&param)

	res := model.FalconAgentHeartbeatResult{}
	err := annotateErr(
		commonSling.ToSlintExt(req).DoReceive(http.StatusOK, &res),
		"calling of [api/v1/agent/heartbeat] has error",
	)
	if err != nil {
		logger.Errorf("[Service] AgentHeartbeat HTTP error: %v", errors.Details(err))
		return 0, int64(len(agents))
	}

	return res.RowsAffected, 0
}

func NqmAgentHeartbeat(req *nqmModel.HeartbeatRequest) (*nqmModel.AgentView, error) {
	resp := &nqmModel.AgentView{}
	err := commonSling.ToSlintExt(
		NewSlingBase().
			Post("api/v1/heartbeat/nqm/agent").
			BodyJSON(req),
	).DoReceive(http.StatusOK, &resp)
	if err != nil {
		return nil, annotateErr(err, "calling of [api/v1/heartbeat/nqm/agent] has error")
	}
	return resp, nil
}

func NqmAgentHeartbeatTargetList(agentID int32) ([]*nqmModel.HeartbeatTarget, error) {
	var resp []*nqmModel.HeartbeatTarget
	err := commonSling.ToSlintExt(
		NewSlingBase().
			Get(fmt.Sprintf("api/v1/heartbeat/nqm/agent/%d/targets", agentID)),
	).DoReceive(http.StatusOK, &resp)
	if err != nil {
		return nil, annotateErr(err, "calling of [api/v1/heartbeat/nqm/agent/] has error")
	}
	return resp, nil
}

func MinePlugins(hostname string) (*model.NewAgentPluginsResponse, error) {
	var resp *model.NewAgentPluginsResponse
	err := commonSling.ToSlintExt(
		NewSlingBase().Get("api/v1/agent/mineplugins").QueryStruct(struct {
			Hostname string `url:"hostname"`
		}{hostname}),
	).DoReceive(http.StatusOK, &resp)
	if err != nil {
		return nil, annotateErr(err, "calling of [api/v1/agent/mineplugins] has error")
	}
	return resp, nil
}

func Plugins(hostname string) ([]string, error) {
	var resp []string
	err := commonSling.ToSlintExt(
		NewSlingBase().Get("api/v1/agent/plugins/"+hostname),
	).DoReceive(http.StatusOK, &resp)
	if err != nil {
		return nil, annotateErr(err, "calling of [api/v1/agent/plugins/] has error")
	}
	return resp, nil
}

func BuiltinMetrics(hostname string, checksum string) (*model.NewBuiltinMetricResponse, error) {
	var resp *model.NewBuiltinMetricResponse
	err := commonSling.ToSlintExt(
		NewSlingBase().
			Get("api/v1/metrics/builtin").QueryStruct(struct {
			Hostname string `url:"hostname,omitempty"`
			Checksum string `url:"checksum,omitempty"`
		}{hostname, checksum}),
	).DoReceive(http.StatusOK, &resp)
	if err != nil {
		return nil, annotateErr(err, "calling of [api/v1/metrics/builtin] has error")
	}
	return resp, nil
}

func Strategies() ([]*model.NewHostStrategy, error) {
	var resp []*model.NewHostStrategy
	err := commonSling.ToSlintExt(
		NewSlingBase().
			Get("api/v1/strategies"),
	).DoReceive(http.StatusOK, &resp)
	if err != nil {
		return nil, annotateErr(err, "calling of [api/v1/strategies] has error")
	}
	for _, hs := range resp {
		for _, s := range hs.Strategies {
			s.AfterLoad()
		}
	}
	return resp, nil
}

func Expressions() ([]*model.NewExpression, error) {
	var resp []*model.NewExpression
	err := commonSling.ToSlintExt(
		NewSlingBase().
			Get("api/v1/expressions"),
	).DoReceive(http.StatusOK, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
