package rpc

import (
	"fmt"
	"net"
	"strings"
	"time"

	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	dbNqm "github.com/Cepave/open-falcon-backend/common/db/nqm"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	"github.com/Cepave/open-falcon-backend/common/rpc"
	"github.com/asaskevich/govalidator"
)

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

	now := time.Now()

	/**
	 * Refresh the information of agent
	 */
	var currentAgent = nqmModel.NewNqmAgent(&request)
	var agentDetail = dbNqm.RefreshAgentInfo(currentAgent, now)
	if agentDetail == nil {
		return
	}
	// :~)

	targets := dbNqm.GetPingListFromCache(currentAgent, now)
	if len(targets) == 0 {
		return
	}

	response.NeedPing = true
	response.Agent = agentDetail
	response.Targets = targets
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
		if ipAddress := net.ParseIP(request.IpAddress)
			ipAddress == nil {
			err = fmt.Errorf("Cannot parse IP address: [%v]", request.IpAddress)
		}
	}
	// :~)

	return
}
