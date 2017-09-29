package imdb

type Tag struct {
	ID          int      `gorm:"primary_key" json:"id"`
	Name        string   `gorm:"not null;type:varchar(30);unique_index" json:"name"`
	TagTypeId   uint     `gorm:"not null" json:"tag_type_id"`
	Description string   `gorm:"not null" json:"description"`
	Default     int      `gorm:"not null" json:"default"`
	CreatedAt   JSONTime `json:"created_at"`
	UpdatedAt   JSONTime `json:"updated_at"`
	TagType     TagType  `gorm:"ForeignKey:TagTypeId" json:"tag_type"`
}

func (Tag) TableName() string {
	return "tags"
}
