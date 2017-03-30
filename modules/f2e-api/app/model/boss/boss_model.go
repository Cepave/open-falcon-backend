package boss

import (
	con "github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

type BossHost struct {
	Platform string `json:"platform" gorm:"column:platform"`
	Province string `json:"province" gorm:"column:province"`
	Isp      string `json:"isp"  gorm:"column:isp"`
	Idc      string `json:"idc" gorm:"column:idc"`
	Ip       string `json:"ip" gorm:"column:ip"`
	Hostname string `json:"hostname" gorm:"column:hostname"`
}

func (this BossHost) TableName() string {
	return "hosts"
}

func GetBossObjs() (res []BossHost) {
	db := con.Con()
	res = []BossHost{}
	db.Boss.Select("platform, province, isp, idc, ip, hostname").Table("hosts").Where("exist = 1").Scan(&res)
	return res
}
