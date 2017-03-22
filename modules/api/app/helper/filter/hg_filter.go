package filter

import (
	"fmt"

	con "github.com/Cepave/open-falcon-backend/modules/api/config"
)

type GrpHosts struct {
	GrpName  string `json:"grp_name";orm:"grp_name"`
	Hostname string `json:"hostname";orm:"hostname"`
}

func HostGroupFilter(filterTxt string, limit int) []GrpHosts {
	db := con.Con()
	sqlbuild := fmt.Sprintf(`select g2.id, g2.grp_name, h2.hostname from host h2 INNER JOIN (select g.id, g.grp_name, h.host_id from grp g INNER JOIN grp_host h
	on g.id = h.grp_id
	where g.grp_name regexp '%s') g2 on g2.host_id = h2.id limit %d`, filterTxt, limit)
	res := []GrpHosts{}
	db.Falcon.Raw(sqlbuild).Scan(&res)
	return res
}
