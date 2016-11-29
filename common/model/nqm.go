package model

import (
	"fmt"
	"time"
)

const (
	// Value of undefined id
	UNDEFINED_ID = -1

	UNDEFINED_ISP_ID      = int16(UNDEFINED_ID)
	UNDEFINED_PROVINCE_ID = int16(UNDEFINED_ID)
	UNDEFINED_CITY_ID     = int16(UNDEFINED_ID)
	UNDEFINED_NAME_TAG_ID     = int16(UNDEFINED_ID)

	// Value of undefined string
	UNDEFINED_STRING = "<UNDEFINED>"
)

// NqmTaskRequest represents the request for ping task by NQM agent
type NqmTaskRequest struct {
	// The connection id of agent(used to identify task configruation)
	ConnectionId string `valid:"required"`
	// The hostname of agent
	Hostname string `valid:"required"`
	// The IP address of agent
	// Could be IPv4 or IPv6 format
	IpAddress string `valid:"required"`
}

type MeasurementsProperty struct {
	Enabled  bool
	Command  []string
	Interval time.Duration
}

// NqmTaskResponse represents the response for ping task requested from NQM agent
//
// If NeedPing is false, Targets and Command would be empty array
type NqmTaskResponse struct {
	// Whether or not the task should be performed
	NeedPing bool

	// The data of agent
	// nil if there is no need for ping
	Agent *NqmAgent

	// The list of target hosts to be probed(ping)
	// nil if there is no need for ping
	Targets []NqmTarget

	// The command/arguments of command to be executed
	// nil if there is no need for ping
	Measurements map[string]MeasurementsProperty
}

// Represents the data of agent
type NqmAgent struct {
	// The id of agent
	Id int

	// The name of agent
	Name string

	// The id of ISP, UNDEFINED_ID means there is not such data for this target
	IspId int16
	// The name of ISP
	IspName string

	// The id of province, UNDEFINED_ID means there is not such data for this target
	ProvinceId int16
	// The name of province
	ProvinceName string

	// The id of city, UNDEFINED_ID means there is not such data for this target
	CityId int16
	// The name of city
	CityName string

	NameTagId int16

	// The id of group tags
	GroupTagIds []int32
}

// Represents the data of target used by NQM agent
type NqmTarget struct {
	// The id of target
	Id int

	// The IP address or FQDN used by ping command
	Host string

	// The id of ISP, UNDEFINED_ID means there is not such data for this target
	IspId int16
	// The name of ISP
	IspName string

	// The id of province, UNDEFINED_ID means there is not such data for this target
	ProvinceId int16
	// The name of province
	ProvinceName string

	// The id of city, UNDEFINED_ID means there is not such data for this target
	CityId int16
	// The name of city
	CityName string

	// The id of name tag
	NameTagId int16
	// The name tag of the target
	NameTag string

	// The id of group tags
	GroupTagIds []int32
}

func (target *NqmTarget) String() string {
	return fmt.Sprintf(
		"Id: [%d] Host: [%s] Isp: \"%s\"(%d) Province: \"%s\"(%d), City: \"%s\"[%d], Name tag: [%s](%d), Group Tags: %v",
		target.Id, target.Host,
		target.IspName, target.IspId,
		target.ProvinceName, target.ProvinceId,
		target.CityName, target.CityId,
		target.NameTag, target.NameTagId,
		target.GroupTagIds,
	)
}

func (agent *NqmAgent) String() string {
	return fmt.Sprintf(
		"Id: [%d] Isp: \"%s\"(%d) Province: \"%s\"(%d), City: \"%s\"[%d], Name tag: (%d), Group Tag: %v",
		agent.Id,
		agent.IspName, agent.IspId,
		agent.ProvinceName, agent.ProvinceId,
		agent.CityName, agent.CityId,
		agent.NameTagId, agent.GroupTagIds,
	)
}
