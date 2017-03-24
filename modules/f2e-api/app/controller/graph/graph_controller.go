package graph

import (
	"fmt"
	"strconv"
	"strings"

	"net/http"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	m "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/graph"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	g "github.com/Cepave/open-falcon-backend/modules/f2e-api/graph"
)

func EndpointRegexpQuery(c *gin.Context) {
	q := c.DefaultQuery("q", "")
	limitTmp := c.DefaultQuery("limit", "500")
	limit, err := strconv.Atoi(limitTmp)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}
	if q == "" {
		h.JSONR(c, http.StatusBadRequest, "q is missing")
	} else {
		var endpoint []m.Endpoint
		db.Graph.Table("endpoint").Select("endpoint, id").Where("endpoint regexp ?", q).Limit(limit).Scan(&endpoint)
		endpoints := []map[string]interface{}{}
		for _, e := range endpoint {
			endpoints = append(endpoints, map[string]interface{}{"id": e.ID, "endpoint": e.Endpoint})
		}
		h.JSONR(c, endpoints)
	}
	return
}

func EndpointCounterRegexpQuery(c *gin.Context) {
	eid := c.DefaultQuery("eid", "")
	metricQuery := c.DefaultQuery("metricQuery", ".+")
	limitTmp := c.DefaultQuery("limit", "500")
	limit, err := strconv.Atoi(limitTmp)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}
	if eid == "" {
		h.JSONR(c, http.StatusBadRequest, "eid is missing")
	} else {
		eids := utils.ConverIntStringToList(eid)
		if eids == "" {
			h.JSONR(c, http.StatusBadRequest, "input error, please check your input info.")
			return
		} else {
			eids = fmt.Sprintf("(%s)", eids)
		}
		var counters []m.EndpointCounter
		db.Graph.Table("endpoint_counter").Select("counter").Where(fmt.Sprintf("endpoint_id IN %s AND counter regexp '%s' ", eids, metricQuery)).Scan(&counters)
		countersResp := []interface{}{}
		for _, c := range counters {
			countersResp = append(countersResp, c.Counter)
		}
		result := utils.UniqSet(countersResp)
		result = utils.MapTake(result, limit)
		h.JSONR(c, result)
	}
	return
}

func EndpointStrCounterRegexpQuery(c *gin.Context) {
	endpoints := c.DefaultQuery("endpoints", "")
	metricQuery := c.DefaultQuery("metricQuery", ".+")
	limitTmp := c.DefaultQuery("limit", "500")
	limit, err := strconv.Atoi(limitTmp)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}
	if endpoints == "" {
		h.JSONR(c, http.StatusBadRequest, "endpoints is missing")
	} else {
		enps := strings.Split(endpoints, ",")
		enpids := []int{}
		rows, _ := db.Graph.Table("endpoint").Select("id").Where("endpoint IN (?)", enps).Rows()
		for rows.Next() {
			id := struct {
				Id int
			}{}
			db.Falcon.ScanRows(rows, &id)
			enpids = append(enpids, id.Id)
		}
		eids, _ := utils.ArrIntToString(enpids)
		if eids == "" {
			h.JSONR(c, http.StatusBadRequest, "input error, please check your input info.")
			return
		} else {
			eids = fmt.Sprintf("(%s)", eids)
		}
		var counters []m.EndpointCounter
		db.Graph.Table("endpoint_counter").Select("counter").Where(fmt.Sprintf("endpoint_id IN %s AND counter regexp '%s' ", eids, metricQuery)).Scan(&counters)
		countersResp := []interface{}{}
		for _, c := range counters {
			countersResp = append(countersResp, c.Counter)
		}
		result := utils.UniqSet(countersResp)
		result = utils.MapTake(result, limit)
		h.JSONR(c, result)
	}
	return
}

type APIQueryGraphDrawData struct {
	HostNames []string `json:"hostnames" binding:"required"`
	Counters  []string `json:"counters" binding:"required"`
	ConsolFun string   `json:"consol_fun" binding:"required"`
	StartTime int64    `json:"start_time" binding:"required"`
	EndTime   int64    `json:"end_time" binding:"required"`
	Step      int      `json:"step" binding:"required"`
}

func QueryGraphDrawData(c *gin.Context) {
	var inputs APIQueryGraphDrawData
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	respData := []*cmodel.GraphQueryResponse{}
	for _, host := range inputs.HostNames {
		for _, counter := range inputs.Counters {
			data, _ := fetchData(host, counter, inputs.ConsolFun, inputs.StartTime, inputs.EndTime, inputs.Step)
			respData = append(respData, data)
		}
	}
	h.JSONR(c, respData)
	return
}

func fetchData(hostname string, counter string, consolFun string, startTime int64, endTime int64, step int) (resp *cmodel.GraphQueryResponse, err error) {
	qparm := g.GenQParam(hostname, counter, consolFun, startTime, endTime, step)
	log.Debugf("qparm: %v", qparm)
	resp, err = g.QueryOne(qparm)
	if err != nil {
		log.Debugf("query graph got error: %s", err.Error())
	}
	return
}
