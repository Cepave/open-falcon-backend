package model

type HostsResult struct {
	Hostname string        `json:"hostname" conform:"trim"`
	ID       int           `json:"id"`
	Groups   []*GroupField `json:"groups"`
}

type GroupField struct {
	ID   int16  `json:"id"`
	Name string `json:"name" conform:"trim"`
}
