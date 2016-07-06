package dashboard

import "time"

//graph
type Endpoint struct {
	Id       int64     `json:"id"`
	Endpoint string    `json:"endpoint"`
	Ts       int64     `json:"ts"`
	TCreate  time.Time `json:"-"`
	TModify  time.Time `json:"-"`
	Ipv4     string    `json:"-"`
}

type EndpointCounter struct {
	Id         int64     `json:"id"`
	EndpointID int64     `json:"endpoint_id"`
	Counter    string    `json:"counter"`
	Step       int64     `json:"step"`
	Type       string    `json:"type"`
	Ts         int64     `json:"ts"`
	TCreate    time.Time `json:"-"`
	TModify    time.Time `json:"-"`
}

//falcon_portal
type HostGroup struct {
	Id         int64     `json:"id"`
	GrpName    string    `json:"grp_name"`
	CreateUser string    `json:"create_user"`
	CreateAt   time.Time `json:"create_at"`
	ComeFrom   int       `json:"come_from"`
}

type Hosts struct {
	Id            int64     `json:"id" `
	Hostname      string    `json:"hostname"`
	Ip            string    `json:"ip"`
	AgentVersion  string    `json:"agent_version"`
	PluginVersion string    `json:"plugin_version"`
	MaintainBegin int       `json:"maintain_begin"`
	MaintainEnd   int       `json:"maintain_end"`
	UpdateAt      time.Time `json:"update_at"`
}

type HostGroupMapping struct {
	GrpId  int64 `json:"grp_id"`
	HostId int64 `json:"host_id"`
}

type GitInfo struct {
	Hostname      string    `json:"hostname"`
	Ip            string    `json:"ip"`
	AgentVersion  string    `json:"agent_version"`
	PluginVersion string    `json:"plugin_version"`
	Title         string    `json:"title"`
	Date          time.Time `json:"commit_at"`
	Valid         bool      `json:"valid"`
}
