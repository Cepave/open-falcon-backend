package imdb

type ResourceObject struct {
	ID         uint     `gorm:"primary_key" json:"id"`
	ObjectType int      `gorm:"not null" json:"object_type"`
	CreatedAt  JSONTime `json:"created_at"`
	UpdatedAt  JSONTime `json:"updated_at"`
}

func (ResourceObject) TableName() string {
	return "resource_objects"
}
