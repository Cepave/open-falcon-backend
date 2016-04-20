package model

import (
	"net"
	"github.com/Cepave/common/model"
	"strings"
)

// Represents the model of NQM agent, which is only used in HBS
type NqmAgent struct {
	Id int
	IpAddress net.IP

	rpcNqmAgent *model.NqmPingTaskRequest
}

// Constructs a new instance of NQM agent
func NewNqmAgent(rpcNqmAgent *model.NqmPingTaskRequest) *NqmAgent {
	var ipAddress = net.ParseIP(rpcNqmAgent.IpAddress);

	if strings.IndexAny(rpcNqmAgent.IpAddress, ".") >= 0 {
		ipAddress = ipAddress.To4()
	} else {
		ipAddress = ipAddress.To16()
	}

	return &NqmAgent{
		rpcNqmAgent: rpcNqmAgent,
		IpAddress: ipAddress,
	}
}

// Gets the value of connection id
func (agent *NqmAgent) ConnectionId() string {
	return agent.rpcNqmAgent.ConnectionId
}
// Gets the value of hostname
func (agent *NqmAgent) Hostname() string {
	return agent.rpcNqmAgent.Hostname
}
