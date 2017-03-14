package owl

// Isp represents data of ISP in RDB
type Isp struct {
	Id      int16  `gorm:"primary_key:true;column:isp_id" json:"id"`
	Name    string `gorm:"column:isp_name" json:"name"`
	Acronym string `gorm:"column:isp_acronym" json:"acronym"`
}

func (Isp) TableName() string {
	return "owl_isp"
}

// Isp represents data of ISP in RDB
type IspOfPingtaskView struct {
	Id   int    `gorm:"primary_key:true;column:isp_id" json:"id"`
	Name string `gorm:"column:isp_name" json:"name"`
}

func (IspOfPingtaskView) TableName() string {
	return "owl_isp"
}
