package owl

// Isp represents data of ISP in RDB
type Isp struct {
	Id int16 `gorm:"primary_key:true;column:isp_id"`
	Name string `gorm:"column:isp_name"`
}

func (Isp) TableName() string {
	return "owl_isp"
}
