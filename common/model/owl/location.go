package owl

// Province represents the data of province in RDB
type Province struct {
	Id   int16  `gorm:"primary_key:true;column:pv_id" json:"id"`
	Name string `gorm:"column:pv_name" json:"name"`
}

func (Province) TableName() string {
	return "owl_province"
}

// Province represents the data of province in RDB
type ProvinceOfPingtaskView struct {
	Id   int    `gorm:"primary_key:true;column:pv_id" json:"id"`
	Name string `gorm:"column:pv_name" json:"name"`
}

func (ProvinceOfPingtaskView) TableName() string {
	return "owl_province"
}

// City represents the data of city1 in RDB
type City1 struct {
	Id           int16  `gorm:"primary_key:true;column:ct_id" json:"id"`
	Name         string `gorm:"column:ct_name" json:"name"`
	PostCode     string `gorm:"column:ct_post_code" json:"post_code"`
	ProvinceId   int16  `gorm:"primary_key:true;column:pv_id" json:"province.id"`
	ProvinceName string `gorm:"column:pv_name" json:"province.name"`
}

func (City1) TableName() string {
	return "owl_city"
}

// City represents the data of city1 in RDB
type City2 struct {
	Id       int16  `gorm:"primary_key:true;column:ct_id" json:"id"`
	Name     string `gorm:"column:ct_name" json:"name"`
	PostCode string `gorm:"column:ct_post_code" json:"post_code"`
}

func (c *City2) GetKey() interface{} {
	return c.Id
}

func (City2) TableName() string {
	return "owl_city"
}

// CCityOfPingtaskViewity represents the data of city in the PingtaskView
type CityOfPingtaskView struct {
	Id   int    `gorm:"primary_key:true;column:ct_id" json:"id"`
	Name string `gorm:"column:ct_name" json:"name"`
}

func (CityOfPingtaskView) TableName() string {
	return "owl_city"
}
