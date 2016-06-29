package dashboard

import (
	"errors"
	"fmt"
	"github.com/Cepave/fe/g"
	"github.com/astaxie/beego/orm"
	"regexp"
	"strings"
)

func QueryHostGroupByNameRegx(queryStr string, limit int) (hostgroup []HostGroup, err error) {
	config := g.Config()
	if limit == 0 || limit > config.GraphDB.LimitHostGroup {
		limit = config.GraphDB.LimitHostGroup
	}
	hostgroup, err = getHostGroupRegex(queryStr, limit)
	return
}

func GetHostsByHostGroupName(hostGroupName []string) (hosts []Hosts, err error) {
	if len(hostGroupName) == 0 {
		err = errors.New("query string is empty")
		return
	}

	host_tmp := map[string]Hosts{}
	for _, v := range hostGroupName {
		v = strings.Replace(v, "\"", "", -1)
		hostgroupid, _ := getHostGroupIdByName(v)
		hostids, _ := getHostIdsByHostGroupId(hostgroupid)
		ho, _ := getHostsByHostIds(hostids)
		for _, v2 := range ho {
			if _, ok := host_tmp[v2.Hostname]; ok != true {
				host_tmp[v2.Hostname] = v2
			}
		}
	}
	for _, v := range host_tmp {
		hosts = append(hosts, v)
	}
	return
}

func getOrmObj() (q orm.Ormer) {
	q = orm.NewOrm()
	q.Using("falcon_portal")
	return
}

func getHostGroupRegex(queryStr string, limit int) (hostgroup []HostGroup, err error) {
	q := getOrmObj()
	_, err = q.Raw("select * from `grp` where grp_name regexp ? limit ?", queryStr, limit).QueryRows(&hostgroup)
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

func CountNumOfHostGroup() (c int, err error) {
	var h []HostGroup
	q := getOrmObj()
	_, err = q.Raw("select * from `grp`").QueryRows(&h)
	c = len(h)
	return
}
