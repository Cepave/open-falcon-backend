package rpc

import (
	"strings"
	"time"

	commonModel "github.com/Cepave/common/model"
	"github.com/Cepave/hbs/db"
	"github.com/Cepave/hbs/model"
	"github.com/asaskevich/govalidator"
)

// Retrieve the configuration of ping task for certain client
//
// If the NqmPingTaskRequest.ConnectionId is not existing in database,
// this function would create one, but the NqmPingTiskResponse.NeedPing would be value of false.
func (t *NqmAgent) PingTask(request commonModel.NqmPingTaskRequest, response *commonModel.NqmPingTaskResponse) (err error) {

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
	response.Command = nil

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
	response.Command = []string{"fping", "-p", "20", "-i", "10", "-C", "4", "-q", "-a"}

	return
}

func validatePingTask(request *commonModel.NqmPingTaskRequest) (err error) {
	request.ConnectionId = strings.TrimSpace(request.ConnectionId)
	request.Hostname = strings.TrimSpace(request.Hostname)
	request.IpAddress = strings.TrimSpace(request.IpAddress)

	_, err = govalidator.ValidateStruct(request)
	return
}
