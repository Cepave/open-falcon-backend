package imdb

type DescriptionValue struct {
	ID               int    `gorm:"primary_key" json:"id"`
	TagId            int    `gorm:"not null" json:"tag_id"`
	ResourceObjectId int    `gorm:"not null" json:"resource_object_id"`
	ObjectTagId      int    `gorm:"not null" json:"object_tag_id"`
	Value            string `gorm:"not null;type:varchar(300)" json:"value"`

	ObjectTag ObjectTag `gorm:"ForeignKey:ObjectTagId"`
}

func (DescriptionValue) TableName() string {
	return "description_values"
}
