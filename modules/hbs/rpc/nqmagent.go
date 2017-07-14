package rpc

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Cepave/open-falcon-backend/common/json"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	"github.com/Cepave/open-falcon-backend/common/rpc"
	nqmService "github.com/Cepave/open-falcon-backend/common/service/nqm"
	"github.com/Cepave/open-falcon-backend/modules/hbs/service"
	"github.com/asaskevich/govalidator"
	"github.com/dghubble/sling"
)

var nqmAgentHbsService *nqmService.AgentHbsService = nil

var mysqlApiSling *sling.Sling

// Task retrieves the configuration of measurement tasks for certain client
//
// If the NqmTaskRequest.ConnectionId is not existing in database,
// this function would create one, but the NqmTiskResponse.NeedPing would be value of false.
func (t *NqmAgent) Task(request commonModel.NqmTaskRequest, response *commonModel.NqmTaskResponse) (err error) {
	defer rpc.HandleError(&err)()

	/**
	 * Validates data
	 */
	if err = validatePingTask(&request); err != nil {
		return
	}
	// :~)

	response.NeedPing = false
	response.Agent = nil
	response.Targets = nil
	response.Measurements = nil

	agentHeartbeatReq := &nqmModel.HeartbeatRequest{
		ConnectionId: request.ConnectionId,
		Hostname:     request.Hostname,
		IpAddress:    json.NewIP(request.IpAddress),
		Timestamp:    json.JsonTime(time.Now()),
	}

	nqmAgentHeartbeatResp, err := service.NqmAgentHeartbeat(agentHeartbeatReq)
	if err != nil {
		return
	}

	if !nqmAgentHeartbeatResp.Status {
		return
	}

	nqmAgentHeartbeatTargetList, err := service.NqmAgentHeartbeatTargetList(nqmAgentHeartbeatResp.Id)
	if err != nil {
		return
	}

	response.NeedPing = true
	response.Agent = toOldAgent(nqmAgentHeartbeatResp)
	response.Targets = toOldTargets(nqmAgentHeartbeatTargetList)
	response.Measurements = map[string]commonModel.MeasurementsProperty{
		"fping":   {true, []string{"fping", "-p", "20", "-i", "10", "-C", "4", "-q", "-a"}, 300},
		"tcpping": {false, []string{"tcpping", "-i", "0.01", "-c", "4"}, 300},
		"tcpconn": {false, []string{"tcpconn"}, 300},
	}
	return
}

func validatePingTask(request *commonModel.NqmTaskRequest) (err error) {
	request.ConnectionId = strings.TrimSpace(request.ConnectionId)
	request.Hostname = strings.TrimSpace(request.Hostname)
	request.IpAddress = strings.TrimSpace(request.IpAddress)

	_, err = govalidator.ValidateStruct(request)

	/**
	 * Checks the validation of IP address
	 */
	if err == nil {
		if ipAddress := net.ParseIP(request.IpAddress); ipAddress == nil {
			err = fmt.Errorf("Cannot parse IP address: [%v]", request.IpAddress)
		}
	}
	// :~)

	return
}

func toOldAgent(a *nqmModel.AgentView) *commonModel.NqmAgent {
	return &commonModel.NqmAgent{
		Id:           int(a.Id),
		Name:         *a.Name,
		IspId:        a.ISP.ID,
		IspName:      a.ISP.Name,
		ProvinceId:   a.Province.ID,
		ProvinceName: a.Province.Name,
		CityId:       a.City.ID,
		CityName:     a.City.Name,
		NameTagId:    a.NameTag.ID,
		GroupTagIds:  a.GroupTags,
	}
}

func toOldTargets(l []*nqmModel.HeartbeatTarget) []commonModel.NqmTarget {
	var r []commonModel.NqmTarget
	for _, t := range l {
		var groupTagIDs []int32
		for _, id := range t.GroupTagIDs {
			groupTagIDs = append(groupTagIDs, int32(id))
		}
		r = append(r, commonModel.NqmTarget{
			Id:           int(t.ID),
			Host:         t.Host,
			IspId:        t.IspID,
			IspName:      t.IspName,
			ProvinceId:   t.ProvinceID,
			ProvinceName: t.ProvinceName,
			CityId:       t.CityID,
			CityName:     t.CityName,
			NameTagId:    t.NameTagID,
			NameTag:      t.NameTag,
			GroupTagIds:  groupTagIDs,
		})
	}
	return r
}
