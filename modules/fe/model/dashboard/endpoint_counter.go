package dashboard

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/astaxie/beego/orm"
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

func QueryEndpointsByCounter(counter string, limit int, metricQuery string, negatePattern bool) (endpoints []string, err error) {
	config := g.Config()
	if limit == 0 || limit > config.GraphDB.Limit {
		limit = config.GraphDB.Limit
	}

	q := orm.NewOrm()
	q.Using("graph")
	q.QueryTable("endpoint_counter")
	var queryprefix string
	if metricQuery == "" {
		queryprefix = fmt.Sprintf("select endpoint from endpoint as table1 inner join (select distinct endpoint_id from endpoint_counter where counter = '%s') as table2 where table1.id = table2.endpoint_id limit %d", counter, limit)
	} else {
		var negate string
		if negatePattern {
			negate = "not"
		}
		queryprefix = fmt.Sprintf("select endpoint from endpoint as table1 inner join (select distinct endpoint_id from endpoint_counter where counter %s regexp '%s') as table2 where table1.id = table2.endpoint_id limit %d", negate, metricQuery, limit)
	}
	var enp []Endpoint
	_, err = q.Raw(queryprefix).QueryRows(&enp)
	for _, v := range enp {
		endpoints = append(endpoints, v.Endpoint)
	}

	return
}

func QueryCounterByEndpoints(endpoints []string, limit int, metricQuery string) (counters []string, err error) {
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
	var queryperfix string
	if metricQuery == "" {
		queryperfix = fmt.Sprintf("select distinct(counter) from endpoint_counter where endpoint_id IN(%s) limit %d", pattn.ReplaceAllString(endpoint_ids, ""), limit)
	} else {
		queryperfix = fmt.Sprintf("select distinct(counter) from endpoint_counter where endpoint_id IN(%s) and counter regexp '%s' limit %d", pattn.ReplaceAllString(endpoint_ids, ""), metricQuery, limit)
	}
	var enpc []EndpointCounter
	_, err = q.Raw(queryperfix).QueryRows(&enpc)
	for _, v := range enpc {
		counters = append(counters, v.Counter)
	}

	return
}
