package alarm

import "time"

// +---------------+------------------+------+-----+-------------------+----------------+
// | Field         | Type             | Null | Key | Default           | Extra          |
// +---------------+------------------+------+-----+-------------------+----------------+
// | id            | int(10) unsigned | NO   | PRI | NULL              | auto_increment |
// | name          | varchar(64)      | NO   |     | NULL              |                |
// | internal_data | tinyint(4)       | YES  |     | 1                 |                |
// | description   | varchar(255)     | NO   |     |                   |                |
// | created       | timestamp        | NO   |     | CURRENT_TIMESTAMP |                |
// | color         | varchar(20)      | NO   |     | black             |                |
// +---------------+------------------+------+-----+-------------------+----------------+

type AlarmTypes struct {
	ID           int        `json:"id" gorm:"column:id"`
	Name         string     `json:"name" gorm:"column:name"`
	InternalData int        `json:"internal_data" gorm:"internal_data"`
	Description  string     `json:"description" gorm:"description"`
	Color        string     `json:"color" gorm:"color"`
	CreateAt     *time.Time `json:"created" gorm:"created"`
}

func (this AlarmTypes) TableName() string {
	return "alarm_types"
}
