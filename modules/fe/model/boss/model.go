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
	Id       int       `orm:"column(id)"`
	PopId    int       `orm:"column(pop_id)"`
	Name     string    `orm:"column(name)"`
	Count    int       `orm:"column(count)"`
	Area     string    `orm:"column(area)"`
	Province string    `orm:"column(province)"`
	City     string    `orm:"column(city)"`
	UpdatAt  time.Time `orm:"column(update_at)"`
}

//orm model
type Hosts struct {
	Id       int       `orm:"column(id)"`
	Hostname string    `orm:"column(hostname)"`
	Exist    int       `orm:"column(exist)"`
	Activate int       `orm:"column(activate)"`
	Platform string    `orm:"column(platform)"`
	Idc      string    `orm:"column(idc)"`
	Ip       string    `orm:"column(ip)"`
	Isp      string    `orm:"column(isp)"`
	Province string    `orm:"column(province)"`
	City     string    `orm:"column(city)"`
	Status   string    `orm:"column(status)"`
	Updated  time.Time `orm:"column(updated)"`
}

//orm model
type Platforms struct {
	Id       int       `orm:"column(id)"`
	Platform string    `orm:"column(platform)"`
	Contacts string    `orm:"column(contacts)"`
	Count    int       `orm:"column(count)"`
	Updated  time.Time `orm:"column(updated)"`
}

//orm model
type Contacts struct {
	Id      int       `orm:"column(id)"`
	Name    string    `orm:"column(name)"`
	Phone   string    `orm:"column(phone)"`
	Email   string    `orm:"column(email)"`
	Updated time.Time `orm:"column(updated)"`
}
