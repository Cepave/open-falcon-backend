package dashboard

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/astaxie/beego/orm"
)

func QueryEndpintByNameRegx(queryStr string, limit int) (enp []Endpoint, err error) {
	config := g.Config()
	if limit == 0 || limit > config.GraphDB.Limit {
		limit = config.GraphDB.Limit
	}
	q := orm.NewOrm()
	q.Using("graph")
	_, err = q.Raw("select * from `endpoint` where endpoint regexp ? limit ?", queryStr, limit).QueryRows(&enp)
	return
}

func QueryCounterByNameRegx(queryStr string, limit int) (counters []string, err error) {
	config := g.Config()
	if limit == 0 || limit > config.GraphDB.Limit {
		limit = config.GraphDB.Limit
	}
	q := orm.NewOrm()
	q.Using("graph")
	_, err = q.Raw("select distinct(counter) from `endpoint_counter` where counter regexp ? limit ?", queryStr, limit).QueryRows(&counters)
	return
}

func QueryEndpintByNameRegxForOps(queryStr string) (enp []Hosts, err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	_, err = q.Raw("select * from `host` where hostname regexp ?", queryStr).QueryRows(&enp)
	return
}

func QueryConfig(key string) (c model.Config, err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	err = q.Raw("select t.key, t.value from common_config as t where t.key = ?", key).QueryRow(&c)
	return
}

func CountNumOfHost() (c int, err error) {
	var h []Endpoint
	q := getOrmObj()
	_, err = q.Raw("select id from `endpoint`").QueryRows(&h)
	c = len(h)
	return
}
