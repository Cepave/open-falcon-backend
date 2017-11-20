package cmdb

type SyncHost struct {
	Activate int    `json:"activate"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
}

type SyncHostGroup struct {
	Creater string `json:"creater"`
	Name    string `json:"name"`
}

type SyncForAdding struct {
	Hosts      []SyncHost          `json:"hosts"`
	Hostgroups []SyncHostGroup     `json:"hostgroups"`
	Relations  map[string][]string `json:"relations"`
}

type SyncItem struct {
	Uuid      string
	StartTime int
}
