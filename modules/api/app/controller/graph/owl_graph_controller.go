package graph

import (
	"fmt"

	"github.com/emirpasic/gods/lists/arraylist"

	"github.com/gin-gonic/gin"
	h "github.com/Cepave/open-falcon-backend/modules/api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/api/app/helper/filter"
	"github.com/Cepave/open-falcon-backend/modules/api/app/model/boss"
)

func HostsSearching(c *gin.Context) {
	AccpectTypes := arraylist.New()
	AccpectTypes.Add("platfrom", "idc", "isp", "province", "hostname", "hostgroup")
	q := c.DefaultQuery("q", "--")
	ftype := c.DefaultQuery("filter_type", "all")
	if q == "--" {
		h.JSONR(c, badstatus, "q is empty, please check it!")
		return
	}
	if !(ftype == "all" || AccpectTypes.Contains(ftype)) {
		h.JSONR(c, badstatus, fmt.Sprintf("filter_type got error type, please check it!, only support: %v", AccpectTypes))
		return
	}
	bossList := boss.GetBossObjs()
	res := map[string]interface{}{}
	switch ftype {
	case "platfrom":
		res = map[string]interface{}{
			"platfrom": filter.PlatformFilter(bossList, q),
		}
	case "idc":
		res = map[string]interface{}{
			"idc": filter.IdcFilter(bossList, q),
		}
	case "isp":
		res = map[string]interface{}{
			"isp": filter.IspFilter(bossList, q),
		}
	case "province":
		res = map[string]interface{}{
			"province": filter.ProvinceFilter(bossList, q),
		}
	case "hostname":
		res = map[string]interface{}{
			"hostname": filter.HostNameFilter(bossList, q),
		}
	case "hostgroup":
		res = map[string]interface{}{
			"hostgroup": filter.HostGroupFilter(q),
		}
	case "all":
		res = map[string]interface{}{
			"platfrom":  filter.PlatformFilter(bossList, q),
			"idc":       filter.IdcFilter(bossList, q),
			"isp":       filter.IspFilter(bossList, q),
			"province":  filter.ProvinceFilter(bossList, q),
			"hostname":  filter.HostNameFilter(bossList, q),
			"hostgroup": filter.HostGroupFilter(q),
		}
	}
	h.JSONR(c, res)
	return
}
