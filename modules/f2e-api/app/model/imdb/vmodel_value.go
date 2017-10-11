package imdb

type VmodelValue struct {
	ID               int `gorm:"primary_key" json:"id"`
	TagId            int `gorm:"not null" json:"tag_id"`
	ObjectTagId      int `gorm:"not null" json:"object_tag_id"`
	ResourceObjectId int `gorm:"not null" json:"resource_object_id"`
	ValueModelId     int `gorm:"not null" json:"value_model_id"`

	ObjectTag  ObjectTag  `gorm:"ForeignKey:ObjectTagId"`
	ValueModel ValueModel `gorm:"ForeignKey:ValueModelId"`
}

func (VmodelValue) TableName() string {
	return "vmodel_values"
}
