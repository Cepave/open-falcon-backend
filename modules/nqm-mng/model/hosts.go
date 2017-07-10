package model

type HostsResult struct {
	Hostname string        `gorm:"column:hostname" json:"hostname" conform:"trim"`
	ID       int           `gorm:"primary_key:true;column:id" json:"id"`
	Groups   []*GroupField `json:"groups"`
}

func (HostsResult) TableName() string {
	return "host"
}

type GroupField struct {
	ID   int16  `gorm:"primary_key:true;column:id" json:"id"`
	Name string `gorm:"column:name" json:"name" conform:"trim"`
}
