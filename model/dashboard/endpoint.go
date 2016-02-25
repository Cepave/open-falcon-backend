package dashboard

import (
	"github.com/astaxie/beego/orm"
)

func QueryEndpintByNameRegx(queryStr string) (enp []Endpoint, err error) {
	q := orm.NewOrm()
	q.Using("graph")
	_, err = q.Raw("select * from `endpoint` where endpoint regexp ?", queryStr).QueryRows(&enp)
	return
}
