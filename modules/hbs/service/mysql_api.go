package service

import (
	"net/http"
	"strconv"

	"github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	commonSling "github.com/Cepave/open-falcon-backend/common/sling"
)

func agentHeartbeatCall(agents []*model.FalconAgentHeartbeat) (rowsAffectedCnt int64, agentsDroppedCnt int64) {
	param := struct {
		UpdateOnly bool `json:"update_only"`
	}{updateOnlyFlag}
	req := NewSlingBase().Post("api/v1/agent/heartbeat").BodyJSON(agents).QueryStruct(&param)

	res := model.FalconAgentHeartbeatResult{}
	err := commonSling.ToSlintExt(req).DoReceive(http.StatusOK, &res)
	if err != nil {
		logger.Errorln("[Service] AgentHeartbeat HTTP error:", err)
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
		return nil, err
	}
	return resp, nil
}

func NqmAgentHeartbeatTargetList(agentID int32) ([]*nqmModel.HeartbeatTarget, error) {
	var resp []*nqmModel.HeartbeatTarget
	err := commonSling.ToSlintExt(
		NewSlingBase().
			Get("api/v1/heartbeat/nqm/agent/"+strconv.Itoa(int(agentID))+"/targets"),
	).DoReceive(http.StatusOK, &resp)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
