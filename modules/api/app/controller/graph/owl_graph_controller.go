package graph

import (
	"fmt"

	"github.com/emirpasic/gods/lists/arraylist"

	"regexp"
	"strconv"

	h "github.com/Cepave/open-falcon-backend/modules/api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/api/app/helper/filter"
	"github.com/Cepave/open-falcon-backend/modules/api/app/model/boss"
	"gopkg.in/gin-gonic/gin.v1"
)

func HostsSearching(c *gin.Context) {
	AccpectTypes := arraylist.New()
	AccpectTypes.Add("platform", "idc", "isp", "province", "hostname", "hostgroup")
	q := c.DefaultQuery("q", "--")
	limitStr := c.DefaultQuery("limit", "100")
	limit := 100
	if ok, _ := regexp.MatchString(`\d+`, limitStr); ok {
		limit, _ = strconv.Atoi(limitStr)
	}
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
	case "platform":
		res = map[string]interface{}{
			"platform": filter.PlatformFilter(bossList, q, limit),
		}
	case "idc":
		res = map[string]interface{}{
			"idc": filter.IdcFilter(bossList, q, limit),
		}
	case "isp":
		res = map[string]interface{}{
			"isp": filter.IspFilter(bossList, q, limit),
		}
	case "province":
		res = map[string]interface{}{
			"province": filter.ProvinceFilter(bossList, q, limit),
		}
	case "hostname":
		res = map[string]interface{}{
			"hostname": filter.HostNameFilter(bossList, q, limit),
		}
	case "hostgroup":
		res = map[string]interface{}{
			"hostgroup": filter.HostGroupFilter(q, limit),
		}
	case "all":
		res = map[string]interface{}{
			"platform":  filter.PlatformFilter(bossList, q, limit),
			"idc":       filter.IdcFilter(bossList, q, limit),
			"isp":       filter.IspFilter(bossList, q, limit),
			"province":  filter.ProvinceFilter(bossList, q, limit),
			"hostname":  filter.HostNameFilter(bossList, q, limit),
			"hostgroup": filter.HostGroupFilter(q, limit),
		}
	}
	h.JSONR(c, res)
	return
}
