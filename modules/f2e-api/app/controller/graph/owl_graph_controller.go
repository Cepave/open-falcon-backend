package graph

import (
	"github.com/emirpasic/gods/lists/arraylist"

	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper/filter"
	m "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/graph"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/boss"
	"github.com/gin-gonic/gin"
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

type APIEndpointsQuerySubMetricInputs struct {
	Endpoints   []string `json:"endpoints" form:"endpoints" binding:"required"`
	MetricQuery string   `json:"metric_query" form:"metric_query"`
	RawKeys     bool     `json:"raw_keys" form:"raw_keys"`
}

func EndpointsQuerySubMetric(c *gin.Context) {
	inputs := APIEndpointsQuerySubMetricInputs{}
	inputs.MetricQuery = ".+"
	inputs.RawKeys = true
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, http.StatusBadRequest, err.Error())
		return
	}
	if inputs.MetricQuery == "" {
		inputs.MetricQuery = ".+"
	}
	enps := inputs.Endpoints
	enpids := []int64{}
	db.Graph.Table("endpoint").Select("id").Where("endpoint IN (?)", enps).Pluck("id", &enpids)
	metricSlist := map[string]int{}
	if len(enpids) != 0 {
		counters := []string{}
		db.Graph.Table("endpoint_counter").Where("endpoint_id IN (?) AND counter regexp ? ", enpids, inputs.MetricQuery).Pluck("DISTINCT(counter)", &counters)
		if inputs.RawKeys {
			h.JSONR(c, counters)
		} else {
			for _, ct := range counters {
				c2 := getFirstCutStringOfMetric(ct)
				if v, ok := metricSlist[c2]; ok {
					metricSlist[c2] = v + 1
				} else {
					metricSlist[c2] = 1
				}
			}
			h.JSONR(c, metricSlist)
		}
		return
	}
	return
}

type APIEndpointsGetMetricBySubStarInput struct {
	Endpoints      []string `json:"endpoints" form:"endpoints" binding:"required"`
	MetricQuery    string   `json:"metric_query" form:"metric_query" binding:"required"`
	StarStr        []string `json:"star_strings" form:"star_strings" binding:"required"`
	SelectedMetric []string `json:"selected_metric" form:"selected_metric"`
}

func EndpointsGetMetricBySubStar(c *gin.Context) {
	//the same with other one , so reuse it
	inputs := APIEndpointsGetMetricBySubStarInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, http.StatusBadRequest, err.Error())
		return
	}
	enps := inputs.Endpoints
	enpids := []int{}
	db.Graph.Table("endpoint").Select("id").Where("endpoint IN (?)", enps).Pluck("id", &enpids)
	var regexFilterTmp []string
	for _, n := range inputs.StarStr {
		n2 := fmt.Sprintf("counter regexp '^%s\\.'", n)
		regexFilterTmp = append(regexFilterTmp, n2)
	}
	regexFilter := strings.Join(regexFilterTmp, " OR ")
	metricSlist := inputs.SelectedMetric
	if len(enpids) != 0 {
		var counters []m.EndpointCounter
		db.Graph.Table("endpoint_counter").
			Select("DISTINCT(counter) as counter").
			Where(fmt.Sprintf(
				"endpoint_id IN (?) AND (%s) AND counter regexp ? ", regexFilter), enpids, inputs.MetricQuery).
			Scan(&counters)

		for _, ct := range counters {
			metricSlist = append(metricSlist, ct.Counter)
		}
	}
	h.JSONR(c, metricSlist)
	return
}

func getFirstCutStringOfMetric(c string) string {
	patt := regexp.MustCompile("^([^\\.]+)\\.")
	matchPatt := patt.FindStringSubmatch(c)
	if len(matchPatt) == 0 {
		return c
	} else {
		return matchPatt[1]
	}
}
