package dashboard

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"regexp"
)

func QueryEndpointidbyNames(endpoints []string) (enp []Endpoint, err error) {
	q := orm.NewOrm()
	q.Using("graph")
	q.QueryTable("endpoint")
	sql_wrapper := make([]string, len(endpoints))
	for i, v := range endpoints {
		sql_wrapper[i] = v
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	qt := qb.Select("*").From("endpoint").Where("endpoint").In(endpoints...)
	_, err = q.Raw(qt.String()).QueryRows(&enp)
	return
}

func QueryCounterByEndpoints(endpoints []string) (counters []string, err error) {
	enp, aerr := QueryEndpointidbyNames(endpoints)
	if aerr != nil {
		err = aerr
		return
	}
	q := orm.NewOrm()
	q.Using("graph")
	q.QueryTable("endpoint_counter")
	var endpoint_ids = ""
	for _, v := range enp {
		endpoint_ids += fmt.Sprintf("%d,", v.Id)
	}
	pattn, _ := regexp.Compile("\\s*,\\s*$")
	queryperfix := fmt.Sprintf("select distinct(counter) from endpoint_counter where endpoint_id IN(%s)", pattn.ReplaceAllString(endpoint_ids, ""))
	var enpc []EndpointCounter
	_, err = q.Raw(queryperfix).QueryRows(&enpc)
	for _, v := range enpc {
		counters = append(counters, v.Counter)
	}
	return
}
