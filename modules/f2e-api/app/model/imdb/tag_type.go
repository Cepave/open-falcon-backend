package imdb

type TagType struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	TypeName    string `gorm:"not null" json:"type_name"`
	DbTableName string `gorm:"not null" json:"db_table_name"`
}

func (TagType) TableName() string {
	return "tag_types"
}
