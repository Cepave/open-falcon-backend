package dashboard

import (
	"github.com/Cepave/fe/g"
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

func CountNumOfHost() (c int, err error) {
	var h []Endpoint
	q := getOrmObj()
	_, err = q.Raw("select id from `endpoint`").QueryRows(&h)
	c = len(h)
	return
}
