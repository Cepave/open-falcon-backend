package grafana

import (
	"regexp"
	"strings"

	"encoding/json"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/query/gin_http/computeFunc"
	"github.com/Cepave/open-falcon-backend/modules/query/gin_http/openFalcon"
	"github.com/Cepave/open-falcon-backend/modules/query/model"
	"github.com/Jeffail/gabs"
	log "github.com/sirupsen/logrus"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

type GrafanaResp struct {
	Expandable bool   `json:"expandable"`
	Text       string `json:"text"`
}

func ResultGen(alist []string) (result []GrafanaResp) {
	result = []GrafanaResp{}
	for _, enp := range alist {
		result = append(result, GrafanaResp{true, enp})
	}
	return
}

func parseTager(query string) (endpoints []string, counter string, ldfunction string) {
	log.Debugf("got query string: %s ", query)
	parseRgex := "^([^#]+)#(.+)"
	if strings.Contains(query, "#{") {
		parseRgex = "^([^#]+)#(.+)#({.+)"
	}
	r := regexp.MustCompile(parseRgex)
	matched := r.FindStringSubmatch(query)
	endpointstr := matched[1]
	counter = matched[2]
	ldfunction = "null"
	if len(matched) > 3 {
		ldfunction = matched[3]
	}
	repx := regexp.MustCompile("(^{|}$)")
	endpointstr = repx.ReplaceAllString(endpointstr, "")
	endpoints = strings.Split(endpointstr, ",")
	log.Debugf("counter: %s, endpoints: %v, ldfunction: %v", counter, endpoints, ldfunction)
	return
}

func GrafanaMain(c *gin.Context) {
	query := c.DefaultQuery("query", ".+")
	log.Debugf("got query string: %s, find keyword? %v", query, strings.Contains(query, "#"))
	if strings.Contains(query, "#") {
		endpoints, counter, _ := parseTager(query)
		endpid := model.EndpointIdQuery(endpoints)
		if len(endpid) == 0 {
			c.JSON(300, gin.H{"error": "can not find any endpoint id, please check your database."})
		}
		catchAll := false
		if strings.Contains(counter, "%") {
			catchAll = true
		}
		counter = strings.Replace(counter, ".*", "%", 1)
		counter = strings.Replace(counter, "#", ".", -1)
		counters := model.FindMatchedCounters(endpid, counter)
		result := []GrafanaResp{}
		perfix := strings.Replace(counter, "%", "", 1)
		tset := hashset.New()
		if perfix == "" {
			for _, c := range counters {
				keys := strings.Split(c, ".")
				if !tset.Contains(keys[0]) {
					expandable := true
					if len(keys) == 0 {
						expandable = false
					}
					result = append(result, GrafanaResp{
						expandable,
						keys[0],
					})
					tset.Add(keys[0])
				}
			}
		} else if catchAll {
			result = []GrafanaResp{GrafanaResp{false, "renderAll"}}
		} else {
			for _, c := range counters {
				skey := strings.Replace(c, perfix, "", 1)
				if skey == "" {
					result = append(result, GrafanaResp{false, c})
				} else {
					keys := strings.Split(skey, ".")
					if !tset.Contains(keys[0]) {
						expandable := true
						if len(keys) == 0 {
							expandable = false
						}
						result = append(result, GrafanaResp{expandable, keys[0]})
						tset.Add(keys[0])
					}
				}
			}
		}
		c.JSON(200, result)
		return
	}
	endpoints := model.EndpointQuery(query)
	var resp []GrafanaResp
	resp = ResultGen(endpoints)
	c.JSON(200, resp)
	return
}

//{"function":"sumAll","aliasName":"sumAll"}
func parseFunc(funjs string) (string, map[string]string) {
	parsedjson, err := gabs.ParseJSON([]byte(funjs))
	if err != nil {
		log.Errorf("during parse UDF got error with -> %s", err.Error())
	}
	funName := parsedjson.Search("function").Data().(string)
	c := computeFunc.GetFuncSetup(funName)
	gotKey := map[string]string{}
	for _, pa := range c.Params {
		pramArr := strings.Split(pa, ":")
		pname := pramArr[0]
		// ptype := pramArr[1]
		if parsedjson.Exists(pname) {
			gotKey[pname] = parsedjson.Search(pname).Data().(string)
		}
	}
	log.Debugf("got http map: %v, funName: %s", gotKey, funName)
	return funName, gotKey
}

type GrafanaPostPrams struct {
	Format        string   `form:"format" json:"format"`
	From          int64    `form:"from" json:"from"`
	MaxDataPoints int      `form:"maxDataPoints" json:"maxDataPoints"`
	Targets       []string `form:"targets" json:"targets"`
	Until         int64    `form:"until" json:"until"`
}

func GetQueryTargets(c *gin.Context) {
	// x := []string{"target", "from", "until", "format", "maxDataPoints", "ldfunction"}
	// for _, s := range x {
	// 	log.Debugf("params: %v", c.DefaultPostForm(s, "null"))
	// }
	var params GrafanaPostPrams
	err := c.Bind(&params)
	if err != nil {
		log.Error(err.Error())
	}
	mtarget := params.Targets
	var resResp []*cmodel.GraphQueryResponse
	for _, target := range mtarget {
		endpoints, counter, ldfunction := parseTager(target)
		endpid := model.EndpointIdQuery(endpoints)
		log.Debugf("got endpoints : %d itmes", len(endpid))
		counter = strings.Replace(counter, ".*", "%", 1)
		counter = strings.Replace(counter, "#", ".", -1)
		counters := model.FindMatchedCounters(endpid, counter)
		startTs := params.From
		endTs := params.Until
		result := []*cmodel.GraphQueryResponse{}
		log.Debugf("got counter : %d itmes", len(counters))
		for _, c := range counters {
			rrds := openFalcon.QueryOnce(startTs, endTs, "AVERAGE", 60, c, endpoints)
			if len(rrds) != 0 {
				for _, s := range rrds {
					result = append(result, s)
				}
			}
		}
		// ldfunction := c.DefaultPostForm("ldfunction", "null")
		if ldfunction == "null" || len(result) == 0 {
			for _, rs := range result {
				resResp = append(resResp, rs)
			}
		} else {
			vm := otto.New()
			funcName, tmpparams := parseFunc(ldfunction)
			funcInstance := computeFunc.GetFuncSetup(funcName)
			err := vm.Set("input", result)
			if err != nil {
				log.Debug(err.Error)
			}
			vm = computeFunc.SetParamsToJSVM(tmpparams, funcInstance.Params, vm)
			vm.Run(funcInstance.Codes)
			output, err := vm.Get("output")
			if err != nil {
				log.Error(err.Error())
				c.JSON(400, gin.H{
					"msg": err.Error(),
				})
				return
			}
			var res []*cmodel.GraphQueryResponse
			json.Unmarshal([]byte(output.String()), &res)
			for _, rs := range res {
				resResp = append(resResp, rs)
			}
			log.Debugf("outputStr: %v", output.String())
		}
	}
	c.JSON(200, resResp)
	return
}
