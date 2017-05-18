package rdb

import (
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
)

func ListHosts(paging commonModel.Paging) ([]*model.HostsResult, *commonModel.Paging) {
	var result []*model.HostsResult
	hostsSql := `
	SELECT h.hostname, h.id, GROUP_CONCAT(g.id ORDER BY g.id ASC SEPARATOR ',') AS gid, GROUP_CONCAT(g.grp_name ORDER BY g.id ASC SEPARATOR '\0') AS gname
	FROM host h
	  LEFT JOIN grp_host gh
	  ON h.id = gh.host_id
	  LEFT JOIN grp g
	  ON gh.grp_id = g.id
	GROUP BY h.id, h.hostname
	ORDER BY h.id ASC
	LIMIT 2, 15
	`

	hostgroupsSql := `
	SELECT g.id, g.grp_name, GROUP_CONCAT(pd.id ORDER BY pd.id ASC SEPARATOR ',') AS pdid, GROUP_CONCAT(pd.dir ORDER BY pd.id ASC SEPARATOR '\0') AS pddir
	FROM grp g
	  LEFT JOIN plugin_dir pd
	  ON g.id = pd.grp_id
	GROUP BY g.id, g.grp_name
	ORDER BY g.id ASC
	LIMIT 2, 15
	`

	return result, &paging
}
