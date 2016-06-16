package rpc

import (
	"strings"
	"time"

	commonModel "github.com/Cepave/common/model"
	"github.com/Cepave/hbs/db"
	"github.com/Cepave/hbs/model"
	"github.com/asaskevich/govalidator"
)

// Task retrieves the configuration of measurement tasks for certain client
//
// If the NqmTaskRequest.ConnectionId is not existing in database,
// this function would create one, but the NqmTiskResponse.NeedPing would be value of false.
func (t *NqmAgent) Task(request commonModel.NqmTaskRequest, response *commonModel.NqmTaskResponse) (err error) {

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

	/**
	 * Refresh the information of agent
	 */
	var currentAgent = model.NewNqmAgent(&request)
	if err = db.RefreshAgentInfo(currentAgent); err != nil {
		return
	}
	// :~)

	/**
	 * Checks and loads agent which is needing performing ping task
	 */
	var nqmAgent *commonModel.NqmAgent
	if nqmAgent, err = db.GetAndRefreshNeedPingAgentForRpc(
		currentAgent.Id, time.Now(),
	); err != nil {
		return
	}

	if nqmAgent == nil {
		return
	}
	// :~)

	/**
	 * Loads matched targets
	 */
	var targets []commonModel.NqmTarget
	if targets, err = db.GetTargetsByAgentForRpc(currentAgent.Id); err != nil {

		return
	}
	// :~)

	response.NeedPing = true
	response.Agent = nqmAgent
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
	return
}
