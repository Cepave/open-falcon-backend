package dashboard

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/astaxie/beego/orm"
)

func QueryHostGroupByNameRegx(queryStr string) (hostgroup []HostGroup, err error) {
	hostgroup, err = getHostGroupRegex(queryStr)
	return
}

func GetHostsByHostGroupName(hostGroupName string) (hostgroup []Hosts, err error) {
	hostgroupid, aerr := getHostGroupIdByName(hostGroupName)
	if aerr != nil {
		err = aerr
		return
	}
	hostids, err := getHostIdsByHostGroupId(hostgroupid)
	if aerr != nil {
		err = aerr
		return
	}
	hostgroup, err = getHostsByHostIds(hostids)
	return
}

func getOrmObj() (q orm.Ormer) {
	q = orm.NewOrm()
	q.Using("falcon_portal")
	return
}

func getHostGroupRegex(queryStr string) (hostgroup []HostGroup, err error) {
	q := getOrmObj()
	_, err = q.Raw("select * from `grp` where grp_name regexp ?", queryStr).QueryRows(&hostgroup)
	return
}

func getHostGroupIdByName(hostgroup string) (id int64, err error) {
	q := getOrmObj()
	var hostgtmp []HostGroup
	_, err = q.Raw("select * from `grp` where grp_name = ?", hostgroup).QueryRows(&hostgtmp)
	if len(hostgtmp) != 0 {
		id = hostgtmp[0].Id
	}
	return
}

func getHostIdsByHostGroupId(hostgroupid int64) (host_ids []int64, err error) {
	q := getOrmObj()
	var mapping []HostGroupMapping
	_, err = q.Raw("select * from `grp_host` where grp_id = ?", hostgroupid).QueryRows(&mapping)
	if len(mapping) != 0 {
		for _, v := range mapping {
			host_ids = append(host_ids, v.HostId)
		}
	}
	return
}

func getHostsByHostIds(hostId []int64) (hosts []Hosts, err error) {
	if len(hostId) == 0 {
		err = errors.New("host is nil")
		return
	}
	q := getOrmObj()
	hostIdQueryTmp := ""
	for _, v := range hostId {
		hostIdQueryTmp += fmt.Sprintf("%d,", v)
	}
	pattn, _ := regexp.Compile("\\s*,\\s*$")
	queryperfix := fmt.Sprintf("select * from host where id IN(%s)", pattn.ReplaceAllString(hostIdQueryTmp, ""))
	_, err = q.Raw(queryperfix).QueryRows(&hosts)
	return
}
