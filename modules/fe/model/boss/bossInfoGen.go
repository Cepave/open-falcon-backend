package boss

import (
	"strings"

	"github.com/astaxie/beego/orm"
)

func getOrmObj() (q orm.Ormer) {
	q = orm.NewOrm()
	q.Using("boss")
	return
}

func GetIPMap() (hostmap map[string]Hosts) {
	hosts := Gethosts()
	hostmap = map[string]Hosts{}
	for _, h := range hosts {
		hostmap[h.Hostname] = h
	}
	return
}

func Gethosts() (hosts []Hosts) {
	q := getOrmObj()
	hosts = []Hosts{}
	q.Raw("select * from `hosts` where exist = 1 and activate = 1").QueryRows(&hosts)
	return
}

func GetContactIfo() (contmap map[string]Contactor) {
	q := getOrmObj()
	contacts := []Contacts{}
	q.Raw("select name, phone, email from `contacts`").QueryRows(&contacts)
	contmap = map[string]Contactor{}
	for _, c := range contacts {
		contmap[c.Name] = Contactor{
			c.Phone,
			c.Email,
			c.Name,
		}
	}
	return
}

func GenContactMap() (platmap map[string][]Contactor) {
	q := getOrmObj()
	plat := []Platforms{}
	q.Raw("select * from `platforms`").QueryRows(&plat)
	contmap := GetContactIfo()
	platmap = map[string][]Contactor{}
	for _, p := range plat {
		names := strings.Split(p.Contacts, ",")
		contacts := make([]Contactor, len(names))
		for ind, n := range names {
			contacts[ind] = contmap[n]
		}
		platmap[p.Platform] = contacts
	}
	return
}
