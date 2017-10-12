package model

import (
	"fmt"
)

type AgentReportRequest struct {
	Hostname      string
	IP            string
	AgentVersion  string
	PluginVersion string
	GitRepo       string
}

func (this *AgentReportRequest) String() string {
	return fmt.Sprintf(
		"<Hostname:%s, IP:%s, AgentVersion:%s, PluginVersion:%s, GitRepo: %s>",
		this.Hostname,
		this.IP,
		this.AgentVersion,
		this.PluginVersion,
		this.GitRepo,
	)
}

type AgentUpdateInfo struct {
	LastUpdate    int64
	ReportRequest *AgentReportRequest
}

type AgentHeartbeatRequest struct {
	Hostname string
	Checksum string
}

func (this *AgentHeartbeatRequest) String() string {
	return fmt.Sprintf(
		"<Hostname: %s, Checksum: %s>",
		this.Hostname,
		this.Checksum,
	)
}

type AgentPluginsResponse struct {
	Plugins       []string
	Timestamp     int64
	GitRepo       string
	GitUpdate     bool
	GitRepoUpdate bool
}

func (this *AgentPluginsResponse) String() string {
	return fmt.Sprintf(
		"<Plugins:%v, Timestamp:%v, GitRepo:%v, GitUpdate:%v, GitRepoUpdate:%v>",
		this.Plugins,
		this.Timestamp,
		this.GitRepo,
		this.GitUpdate,
		this.GitRepoUpdate,
	)
}

// e.g. net.port.listen or proc.num
type BuiltinMetric struct {
	Metric string
	Tags   string
}

func (this *BuiltinMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		this.Metric,
		this.Tags,
	)
}

type BuiltinMetricResponse struct {
	Metrics   []*BuiltinMetric
	Checksum  string
	Timestamp int64
}

func (this *BuiltinMetricResponse) String() string {
	return fmt.Sprintf(
		"<Metrics:%v, Checksum:%s, Timestamp:%v>",
		this.Metrics,
		this.Checksum,
		this.Timestamp,
	)
}

type BuiltinMetricSlice []*BuiltinMetric

func (this BuiltinMetricSlice) Len() int {
	return len(this)
}
func (this BuiltinMetricSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this BuiltinMetricSlice) Less(i, j int) bool {
	return this[i].String() < this[j].String()
}

type FalconAgentHeartbeat struct {
	Hostname      string `json:"hostname" conform:"trim"`
	IP            string `json:"ip" conform:"trim"`
	AgentVersion  string `json:"agent_version" conform:"trim"`
	PluginVersion string `json:"plugin_version" conform:"trim"`
	UpdateTime    int64  `json:"update_time" conform:"trim"`
}

type FalconAgentHeartbeatResult struct {
	RowsAffected int64 `json:"rows_affected"`
}

type NewAgentPluginsResponse struct {
	Plugins   []string `json:"plugins"`
	Timestamp int64    `json:"timestamp"`
	GitRepo   string   `json:"git_repo" conform:"trim"`
}

func (this *NewAgentPluginsResponse) String() string {
	return fmt.Sprintf(
		"<Plugins:%v, Timestamp:%v, GitRepo:%v>",
		this.Plugins,
		this.Timestamp,
		this.GitRepo,
	)
}

// e.g. net.port.listen or proc.num
type NewBuiltinMetric struct {
	Metric string `json:"metric"`
	Tags   string `json:"tags"`
}

func (this *NewBuiltinMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		this.Metric,
		this.Tags,
	)
}

type NewBuiltinMetricResponse struct {
	Metrics   []*NewBuiltinMetric `json:"metrics"`
	Checksum  string              `json:"checksum"`
	Timestamp int64               `json:"timestamp"`
}

func (this *NewBuiltinMetricResponse) String() string {
	return fmt.Sprintf(
		"<Metrics:%v, Checksum:%s, Timestamp:%v>",
		this.Metrics,
		this.Checksum,
		this.Timestamp,
	)
}

type NewBuiltinMetricSlice []*NewBuiltinMetric

func (this NewBuiltinMetricSlice) Len() int {
	return len(this)
}
func (this NewBuiltinMetricSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this NewBuiltinMetricSlice) Less(i, j int) bool {
	return this[i].String() < this[j].String()
}
