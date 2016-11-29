package nqm

import (
	"net"
	"strings"

	"github.com/Cepave/open-falcon-backend/common/model"
)

// Represents the model of NQM agent, which is only used in HBS
type NqmAgent struct {
	Id        int
	IpAddress net.IP

	rpcNqmAgentReq *model.NqmTaskRequest
}

// Constructs a new instance of NQM agent
func NewNqmAgent(rpcNqmAgentReq *model.NqmTaskRequest) *NqmAgent {
	var ipAddress = net.ParseIP(rpcNqmAgentReq.IpAddress)

	if strings.IndexAny(rpcNqmAgentReq.IpAddress, ".") >= 0 {
		ipAddress = ipAddress.To4()
	} else {
		ipAddress = ipAddress.To16()
	}

	return &NqmAgent{
		rpcNqmAgentReq: rpcNqmAgentReq,
		IpAddress:      ipAddress,
	}
}

// Gets the value of connection id
func (agent *NqmAgent) ConnectionId() string {
	return agent.rpcNqmAgentReq.ConnectionId
}

// Gets the value of hostname
func (agent *NqmAgent) Hostname() string {
	return agent.rpcNqmAgentReq.Hostname
}
