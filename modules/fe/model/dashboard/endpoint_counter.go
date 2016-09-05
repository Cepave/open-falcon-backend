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

func QueryEndpointIDByCounters(counters []string) (endpointID []string, err error) {
	q := orm.NewOrm()
	q.Using("graph")
	q.QueryTable("endpoint_counter")
	qb, _ := orm.NewQueryBuilder("mysql")
	qt := qb.Select("distinct(endpoint_id)").From("endpoint_counter").Where("counter").In(counters...)
	_, err = q.Raw(qt.String()).QueryRows(&endpointID)

	return
}

func QueryEndpointNamesByID(endpointIDs []string, limit int, filter string) (endpointNames []string, err error) {
	q := orm.NewOrm()
	q.Using("graph")
	q.QueryTable("endpoint")
	qb, _ := orm.NewQueryBuilder("mysql")
	var qt orm.QueryBuilder
	if filter == "" {
		qt = qb.Select("endpoint").From("endpoint").Where("id").In(endpointIDs...).Limit(limit)
	} else {
		regStr := fmt.Sprintf("endpoint regexp '%s'", filter)
		qt = qb.Select("endpoint").From("endpoint").Where(regStr).And("id").In(endpointIDs...).Limit(limit)
	}
	_, err = q.Raw(qt.String()).QueryRows(&endpointNames)

	return
}

func QueryEndpointsByCounter(counters []string, limit int, filter string) (endpoints []string, err error) {
	config := g.Config()
	if limit == 0 || limit > config.GraphDB.Limit {
		limit = config.GraphDB.Limit
	}

	enpIDs, IDErr := QueryEndpointIDByCounters(counters)
	if IDErr != nil {
		err = IDErr
		return
	}

	endpoints, NameErr := QueryEndpointNamesByID(enpIDs, limit, filter)
	if NameErr != nil {
		err = NameErr
		return
	}

	err = nil
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
