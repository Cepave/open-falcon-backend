package graph

import (
	"fmt"
	"strings"

	"net/http"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	m "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/graph"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	g "github.com/Cepave/open-falcon-backend/modules/f2e-api/graph"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type APIEndpointRegexpQueryInputs struct {
	Q     string `json:"q" form:"q"`
	Label string `json:"tags" form:"tags"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

func EndpointRegexpQuery(c *gin.Context) {
	inputs := APIEndpointRegexpQueryInputs{
		//set default is 500
		Limit: 500,
		Page:  1,
		Q:     ".+",
	}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	if inputs.Q == "" && inputs.Label == "" {
		h.JSONR(c, http.StatusBadRequest, "q and labels are all missing")
		return
	}

	labels := []string{}
	if inputs.Label != "" {
		labels = strings.Split(inputs.Label, ",")
	}
	qs := []string{}
	if inputs.Q != "" {
		qs = strings.Split(inputs.Q, " ")
	}

	offset := 0
	if inputs.Page != 1 {
		offset = (inputs.Page - 1) * inputs.Limit
	}

	var endpoint []m.Endpoint
	var endpoint_id []int
	var dt *gorm.DB
	// query by labels , this is for support falcon-plus dashboard ui page
	if len(labels) != 0 {
		dt = db.Graph.Table("endpoint_counter").Select("distinct endpoint_id")
		for _, trem := range labels {
			dt = dt.Where(" counter like ? ", "%"+strings.TrimSpace(trem)+"%")
		}
		if inputs.Page > 0 {
			dt = dt.Offset(offset)
		}
		dt = dt.Limit(inputs.Limit).Pluck("distinct endpoint_id", &endpoint_id)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}
	}
	// query by endpoint regexp match directly
	if len(qs) != 0 {
		dt = db.Graph.Table("endpoint").
			Select("endpoint, id")
		// combine query result from labels
		if len(endpoint_id) != 0 {
			dt = dt.Where("id in (?)", endpoint_id)
		}

		for _, trem := range qs {
			dt = dt.Where(" endpoint regexp ? ", strings.TrimSpace(trem))
		}
		if inputs.Page > 0 {
			dt = dt.Offset(offset)
		}
		dt.Limit(inputs.Limit).Scan(&endpoint)
	} else if len(endpoint_id) != 0 {
		dt = db.Graph.Table("endpoint").
			Select("endpoint, id").
			Where("id in (?)", endpoint_id)
		if inputs.Page > 0 {
			dt = dt.Offset(offset)
		}
		dt.Limit(inputs.Limit).Scan(&endpoint)
	}
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, dt.Error)
		return
	}

	endpoints := []map[string]interface{}{}
	for _, e := range endpoint {
		endpoints = append(endpoints, map[string]interface{}{"id": e.ID, "endpoint": e.Endpoint})
	}

	h.JSONR(c, endpoints)
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

func QueryGraphLastPoint(c *gin.Context) {
	var inputs []cmodel.GraphLastParam
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	respData := []*cmodel.GraphLastResp{}

	for _, param := range inputs {
		one_resp, err := g.Last(param)
		if err != nil {
			log.Warn("query last point from graph fail:", err)
		} else {
			respData = append(respData, one_resp)
		}
	}

	h.JSONR(c, respData)
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

type APIEndpointCounterRegexpQueryInputs struct {
	Q     string `json:"metricQuery" form:"metricQuery"`
	Eid   string `json:"eid" form:"eid"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

// fastweb only
func EndpointCounterRegexpQuery(c *gin.Context) {
	inputs := APIEndpointCounterRegexpQueryInputs{
		Limit: 500,
		Q:     ".+",
		Page:  -1,
	}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if inputs.Eid == "" {
		h.JSONR(c, http.StatusBadRequest, "eid is missing")
	} else {
		eids := utils.ConverIntStringToList(inputs.Eid)
		if eids == "" {
			h.JSONR(c, http.StatusBadRequest, "input error, please check your input info.")
			return
		} else {
			eids = fmt.Sprintf("(%s)", eids)
			dt := db.Graph.Table("endpoint_counter").Select("DISTINCT(counter) as counter").
				Where(fmt.Sprintf("endpoint_id IN %s AND counter regexp '%s' ", eids, inputs.Q)).
				Limit(inputs.Limit)
			var offset int = 0
			if inputs.Page > 1 {
				offset = (inputs.Page - 1) * inputs.Limit
				dt = dt.Offset(offset)
			}
			countersResp := []string{}
			dt.Pluck("distinct counter", &countersResp)
			h.JSONR(c, countersResp)
			return
		}
	}
}

type APIEndpointStrCounterRegexpQueryInputs struct {
	Q         string `json:"metricQuery" form:"metricQuery"`
	Endpoints string `json:"endpoints" form:"endpoints"`
	Limit     int    `json:"limit" form:"limit"`
	Page      int    `json:"page" form:"page"`
}

func EndpointStrCounterRegexpQuery(c *gin.Context) {
	inputs := APIEndpointStrCounterRegexpQueryInputs{
		Limit: 500,
		Q:     ".+",
		Page:  -1,
	}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if inputs.Endpoints == "" {
		h.JSONR(c, http.StatusBadRequest, "endpoints is missing")
	} else {
		enps := strings.Split(inputs.Endpoints, ",")
		enpids := []int64{}
		db.Graph.Table("endpoint").Select("id").Where("endpoint IN (?)", enps).Pluck("id", &enpids)
		eids, _ := utils.ArrInt64ToString(enpids)
		if eids == "" {
			h.JSONR(c, http.StatusBadRequest, "input error, please check your input info.")
			return
		} else {
			eids = fmt.Sprintf("(%s)", eids)
		}
		dt := db.Graph.Table("endpoint_counter").
			Select("counter").
			Where(fmt.Sprintf("endpoint_id IN %s AND counter regexp '%s' ", eids, inputs.Q)).
			Limit(inputs.Limit)

		var offset int = 0
		if inputs.Page > 1 {
			offset = (inputs.Page - 1) * inputs.Limit
			dt = dt.Offset(offset)
		}

		countersResp := []string{}
		dt.Pluck("distinct counter", &countersResp)
		h.JSONR(c, countersResp)
		return
	}
}
