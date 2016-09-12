package boss

import "time"

type PlatformList struct {
	Status int        `json:"status"`
	Info   string     `json:"info"`
	Result []Platform `json:"result"`
}

type Platform struct {
	Platform string
	IPList   []IPInfo `json:"ip_list"`
}

type IPInfo struct {
	IP       string `json:"ip"`
	HostName string `json:"hostname"`
	IPStatus string `json:"ip_status"`
	POPID    string `json:"pop_id"`
	Platform string `json:"platform"`
}

type Contactor struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

//orm model
type IDC struct {
	Id       int       `orm:"id"`
	PopId    int       `orm:"pop_id"`
	Name     string    `orm:"name"`
	Count    int       `orm:"count"`
	Area     string    `orm:"area"`
	Province string    `orm:"province"`
	City     string    `orm:"city"`
	UpdatAt  time.Time `orm:"update_at"`
}

//orm model
type Hosts struct {
	Id       int       `orm:"id"`
	Hostname string    `orm:"hostnam"`
	Exist    int       `orm:"exist"`
	Activate int       `orm:"activate"`
	Platform string    `orm:"platform"`
	Idc      string    `orm:"idc"`
	Ip       string    `orm:"ip"`
	Isp      string    `orm:"isp"`
	Province string    `orm:"province"`
	City     string    `orm:"city"`
	Status   string    `orm:"status"`
	Updated  time.Time `orm:"updated"`
}

//orm model
type Platforms struct {
	Id       int       `orm:"id"`
	Platform string    `orm:"platform"`
	Contacts string    `orm:"contacts"`
	Count    int       `orm:"count"`
	Updated  time.Time `orm:"updated"`
}

//orm model
type Contacts struct {
	Id      int       `orm:"id"`
	Name    string    `orm:"name"`
	Phone   string    `orm:"phone"`
	Email   string    `orm:"email"`
	Updated time.Time `orm:"updated"`
}
