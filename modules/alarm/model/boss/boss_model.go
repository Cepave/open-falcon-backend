package boss

import (
	"strings"
)

type BossInfo struct {
	ID        int    `json:"id" orm:"column(id)"`
	Hostname  string `json:"hostname" orm:"column(hostname)"`
	Exist     int    `json:"exist" orm:"column(exist)"`
	Activate  int    `json:"activate" orm:"column(activate)"`
	Platform  string `json:"platform" orm:"column(platform)"`
	Platforms string `json:"platforms" orm:"column(platforms)"`
	Idc       string `json:"idc" orm:"column(idc)"`
	IP        string `json:"ip" orm:"column(ip)"`
	Isp       string `json:"isp"  orm:"column(isp)"`
	Province  string `json:"province" orm:"column(province)"`
	Contacts  string `json:"contacts" orm:"column(contacts)"`
}

func (mine BossInfo) Contact() string {
	contact := ""
	if mine.Contacts != "" {
		contact = strings.Split(mine.Contacts, ",")[0]
	}
	return contact
}

func (mine BossInfo) New() BossInfo {
	return BossInfo{
		ID:        0,
		Hostname:  "",
		Exist:     0,
		Activate:  0,
		Platform:  "",
		Platforms: "",
		Idc:       "",
		IP:        "",
		Isp:       "",
		Province:  "",
		Contacts:  "",
	}
}
