package model

type AgentHeartbeat struct {
	Hostname      string `json:"hostname" conform:"trim"`
	IP            string `json:"ip" conform:"trim"`
	AgentVersion  string `json:"agent_version" conform:"trim"`
	PluginVersion string `json:"plugin_version" conform:"trim"`
	UpdateTime    int64  `json:"update_time" conform:"trim"`
}

type AgentHeartbeatResult struct {
	RowsAffected int64 `json:"rows_affected"`
}
