package owl

// Province represents the data of province in RDB
type Province struct {
	Id   int16  `gorm:"primary_key:true;column:pv_id"`
	Name string `gorm:"column:pv_name"`
}

func (Province) TableName() string {
	return "owl_province"
}

// City represents the data of city in RDB
type City struct {
	Id         int16  `gorm:"primary_key:true;column:ct_id"`
	ProvinceId int16  `gorm:"column:ct_pv_id"`
	Name       string `gorm:"column:ct_name"`
	PostCode   string `gorm:"column:ct_post_code"`
}

func (City) TableName() string {
	return "owl_city"
}
