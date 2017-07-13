package boss

import "github.com/astaxie/beego/orm"

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
	q.Raw("select * from `hosts` where exist = 1").QueryRows(&hosts)
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

func GenContactMap() (contactMap map[string]Contactor) {
	contactMap = map[string]Contactor{}
	contmap := GetContactIfo()
	for _, ct := range contmap {
		contactMap[ct.Name] = ct
	}
	return
}
