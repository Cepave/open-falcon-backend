package model

const (
	// Value of undefined id
	UNDEFINED_ID = -1
	// Value of undefined string
	UNDEFINED_STRING = "<UNDEFINED>"
)

// Represents the request for ping task by NQM agent
type NqmPingTaskRequest struct {
	// The connection id of agent(used to identify task configruation)
	ConnectionId string
	// The hostname of agent
	Hostname string
	// The IP address of agent
	// Could be IPv4 or IPv6 format
	IpAddress string
}

// Represents the response for ping task requested from NQM agent
//
// If NeedPing is false, Targets and Command would be empty array
type NqmPingTaskResponse struct {
	// Whether or not the task should be performed
	NeedPing bool

	// The list of target hosts to be probed(ping)
	Targets []NqmTarget

	// The command/arguments of command to be executed
	Command []string
}

// Represents the data of agent
type NqmAgent struct {
	// The id of agent
	Id int

	// The id of ISP, UNDEFINED_ID means there is not such data for this target
	IspId int16
	// The id of province, UNDEFINED_ID means there is not such data for this target
	ProvinceId int16
	// The id of city, UNDEFINED_ID means there is not such data for this target
	CityId int16
}

// Represents the data of target used by NQM agent
type NqmTarget struct {
	// The id of target
	Id int

	// The IP address or FQDN used by ping command
	Host string

	// The id of ISP, UNDEFINED_ID means there is not such data for this target
	IspId int16
	// The id of province, UNDEFINED_ID means there is not such data for this target
	ProvinceId int16
	// The id of city, UNDEFINED_ID means there is not such data for this target
	CityId int16
	// The tag of the target, UNDEFINED_STRING means no such data for this target
	NameTag string
}
