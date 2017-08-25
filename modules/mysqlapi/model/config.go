package model

type AgentConfigResult struct {
	Key   string `gorm:"primary_key:true;column:key" json:"key" conform:"trim"`
	Value string `gorm:"column:value" json:"value" conform:"trim"`
}

func (AgentConfigResult) TableName() string {
	return "common_config"
}
