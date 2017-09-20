package http

import (
	"testing"

	tFlag "github.com/Cepave/open-falcon-backend/common/testing/flag"
	"github.com/astaxie/beego/orm"

	_ "github.com/go-sql-driver/mysql" // import your used driver
)

// this test can only test when boss.hosts is empty.
func TestUpdateHostsTable(t *testing.T) {
	if !testFlags.HasMySqlOfOwlDb(tFlag.OWL_DB_BOSS) {
		t.Log("Skip test because the property of \"mysql.owl_boss=\" is not set")
		return
	}

	orm.RegisterDataBase("default", "mysql", testFlags.GetMysqlOfOwlDb(tFlag.OWL_DB_BOSS), 30)

	updateHostsTable(
		[]string{
			"host-01", "host-02", "host-03",
		},
		map[string]map[string]string{
			"host-01": {
				"hostname":  "host-01",
				"isp":       "isp-1",
				"province":  "province-1",
				"city":      "city-1",
				"platforms": "p1,p2",
				"platform":  "p1",
				"activate":  "1",
			},
			"host-02": {
				"hostname":  "host-02",
				"isp":       "isp-2",
				"province":  "province-2",
				"city":      "city-2",
				"platforms": "p1,p2",
				"platform":  "p1",
				"activate":  "1",
			},
			"host-03": {
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
