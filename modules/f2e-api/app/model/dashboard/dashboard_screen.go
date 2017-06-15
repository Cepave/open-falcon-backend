package dashboard

import (
	con "github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

// +-------+------------------+------+-----+-------------------+-----------------------------+
// | Field | Type             | Null | Key | Default           | Extra                       |
// +-------+------------------+------+-----+-------------------+-----------------------------+
// | id    | int(11) unsigned | NO   | PRI | NULL              | auto_increment              |
// | pid   | int(11) unsigned | NO   | MUL | 0                 |                             |
// | name  | char(128)        | NO   |     | NULL              |                             |
// | time  | timestamp        | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +-------+------------------+------+-----+-------------------+-----------------------------+

type DashboardScreen struct {
	ID      int64  `json:"id" gorm:"column:id"`
	PID     int64  `json:"pid" gorm:"column:pid"`
	Name    string `json:"name" gorm:"column:name"`
	Creator string `json:"creator" gorm:"column:creator"`
}

func (this DashboardScreen) TableName() string {
	return "dashboard_screen"
}

func (this DashboardScreen) Graphs() []DashboardGraph {
	db := con.Con()
	graphs := []DashboardGraph{}
	db.Dashboard.Model(&graphs).Where("screen_id = ?", this.ID).Scan(&graphs)
	return graphs
}

func (mine DashboardScreen) Exist() bool {
	db := con.Con()
	rcount := 0
	db.Dashboard.Model(&mine).Where("id = ?", mine.ID).Count(&rcount)
	if rcount != 0 {
		return true
	} else {
		return false
	}
}

func (mine DashboardScreen) ExistName() bool {
	db := con.Con()
	rcount := 0
	if mine.ID == 0 {
		db.Dashboard.Model(&mine).Where("name = ?", mine.Name).Count(&rcount)
	} else {
		db.Dashboard.Table(mine.TableName()).Where("name = ? AND id != ?", mine.Name, mine.ID).Count(&rcount)
	}
	if rcount != 0 {
		return true
	} else {
		return false
	}
}
