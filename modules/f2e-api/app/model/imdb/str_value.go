package imdb

type StrValue struct {
	ID               int    `gorm:"primary_key" json:"id"`
	TagId            int    `gorm:"not null" json:"tag_id"`
	ResourceObjectId int    `gorm:"not null" json:"resource_object_id"`
	ObjectTagId      int    `gorm:"not null" json:"object_tag_id"`
	Value            string `gorm:"not null;type:varchar(30)" json:"value"`

	ObjectTag ObjectTag `gorm:"ForeignKey:ObjectTagId"`
}

func (StrValue) TableName() string {
	return "str_values"
}
