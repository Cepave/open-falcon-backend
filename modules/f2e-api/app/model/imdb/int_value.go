package imdb

type IntValue struct {
	ID               int `gorm:"primary_key" json:"id"`
	TagId            int `gorm:"not null" json:"tag_id"`
	ResourceObjectId int `gorm:"not null" json:"resource_object_id"`
	ObjectTagId      int `gorm:"not null" json:"object_tag_id"`
	Value            int `gorm:"not null;type:int(6)" json:"value"`
}

func (IntValue) TableName() string {
	return "int_values"
}
