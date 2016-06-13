package dashboard

import (
	"errors"
	"fmt"
	"github.com/Cepave/fe/g"
	"github.com/astaxie/beego/orm"
	"regexp"
)

func QueryEndpointidbyNames(endpoints []string, limit int) (enp []Endpoint, err error) {
	q := orm.NewOrm()
	q.Using("graph")
	q.QueryTable("endpoint")
	qb, _ := orm.NewQueryBuilder("mysql")
	qt := qb.Select("*").From("endpoint").Where("endpoint").In(endpoints...).Limit(limit)
	_, err = q.Raw(qt.String()).QueryRows(&enp)
	return
}

func QueryCounterByEndpoints(endpoints []string, limit int) (counters []string, err error) {
	config := g.Config()
	if limit == 0 || limit > config.GraphDB.Limit {
		limit = config.GraphDB.Limit
	}
	enp, aerr := QueryEndpointidbyNames(endpoints, limit)
	if aerr != nil {
		err = aerr
		return
	}
	if len(enp) == 0 {
		err = errors.New("The endpoints doesn't exist.")
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
	queryperfix := fmt.Sprintf("select distinct(counter) from endpoint_counter where endpoint_id IN(%s) limit %d", pattn.ReplaceAllString(endpoint_ids, ""), limit)
	var enpc []EndpointCounter
	_, err = q.Raw(queryperfix).QueryRows(&enpc)
	for _, v := range enpc {
		counters = append(counters, v.Counter)
	}

	return
}
