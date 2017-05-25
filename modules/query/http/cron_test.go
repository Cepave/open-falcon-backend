package http

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"testing"
)

// this test can only test when boss.hosts is empty.
func TestUpdateHostsTable(t *testing.T) {
	orm.RegisterDataBase("default", "mysql", "root:password@tcp(10.20.30.40:3306)/boss?charset=utf8&loc=Asia%2FTaipei", 30)

	updateHostsTable(
		[]string{
			"host-01", "host-02", "host-03",
		},
		map[string]map[string]string{
			"host-01": map[string]string{
				"hostname":  "host-01",
				"isp":       "isp-1",
				"province":  "province-1",
				"city":      "city-1",
				"platforms": "p1,p2",
				"platform":  "p1",
				"activate":  "1",
			},
			"host-02": map[string]string{
				"hostname":  "host-02",
				"isp":       "isp-2",
				"province":  "province-2",
				"city":      "city-2",
				"platforms": "p1,p2",
				"platform":  "p1",
				"activate":  "1",
			},
			"host-03": map[string]string{
				"hostname":  "host-03",
				"isp":       "isp-3",
				"province":  "province-3",
				"city":      "city-3",
				"platforms": "p1,p2",
				"platform":  "p1",
				"activate":  "1",
			},
		},
	)
}
