package model

type AgentHeartbeat struct {
	Hostname      string
	IP            string
	AgentVersion  string
	PluginVersion string
	UpdateTime    int64
}

type AgentHeartbeatResult struct {
	RowsAffected int64
}
