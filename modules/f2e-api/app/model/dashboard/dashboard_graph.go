package dashboard

// +---------------+------------------+------+-----+---------+----------------+
// | Field         | Type             | Null | Key | Default | Extra          |
// +---------------+------------------+------+-----+---------+----------------+
// | id            | int(11) unsigned | NO   | PRI | NULL    | auto_increment |
// | title         | char(128)        | NO   |     | NULL    |                |
// | hosts         | varchar(10240)   | NO   |     |         |                |
// | counters      | varchar(1024)    | NO   |     |         |                |
// | screen_id     | int(11) unsigned | NO   | MUL | NULL    |                |
// | timespan      | int(11) unsigned | NO   |     | 3600    |                |
// | graph_type    | char(2)          | NO   |     | h       |                |
// | method        | char(8)          | YES  |     |         |                |
// | position      | int(11) unsigned | NO   |     | 0       |                |
// | falcon_tags   | varchar(512)     | NO   |     |         |                |
// | creator       | varchar(50)      | YES  |     | root    |                |
// | time_range    | varchar(50)      | YES  |     | 3h      |                |
// | y_scale       | varchar(50)      | YES  |     | NULL    |                |
// | sort_by       | varchar(30)      | YES  |     | a-z     |                |
// | sample_method | varchar(20)      | YES  |     | AVERAGE |                |
// +---------------+------------------+------+-----+---------+----------------+

type DashboardGraph struct {
	ID           int64  `json:"id" gorm:"column:id"`
	Title        string `json:"title" gorm:"column:title"`
	Hosts        string `json:"hosts" gorm:"column:hosts"`
	Counters     string `json:"counters" gorm:"column:counters"`
	ScreenId     int64  `json:"screen_id" gorm:"column:screen_id"`
	TimeSpan     int64  `json:"timespan" gorm:"column:timespan"`
	GraphType    string `json:"graph_type" gorm:"column:graph_type"`
	Method       string `json:"method" gorm:"column:method"`
	Position     int64  `json:"position" gorm:"column:position"`
	FalconTags   string `json:"falcon_tags" gorm:"column:falcon_tags"`
	Creator      string `json:"creator" gorm:"column:creator"`
	TimeRange    string `json:"time_range" gorm:"column:time_range"`
	YScale       string `json:"y_scale" gorm:"column:y_scale"`
	SortBy       string `json:"sort_by" gorm:"column:sort_by"`
	SampleMethod string `json:"sample_method" gorm:"column:sample_method"`
}

func (this DashboardGraph) TableName() string {
	return "dashboard_graph"
}
