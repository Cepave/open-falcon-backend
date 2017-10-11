package imdb

type ValueModel struct {
	ID        int      `gorm:"primary_key" json:"id"`
	TagId     int      `gorm:"not null" json:"tag_id"`
	Value     string   `gorm:"not null;type:varchar(30)" json:"value"`
	CreatedAt JSONTime `json:"created_at"`
	UpdatedAt JSONTime `json:"updated_at"`
}

func (ValueModel) TableName() string {
	return "value_models"
}
