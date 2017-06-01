package dashboard_graph

import (
	"fmt"
	"sort"
	"strings"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	m "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/dashboard"
	"github.com/gin-gonic/gin"
)

type APIGraphCreateReqDataWithNewScreenInputs struct {
	ScreenName string   `json:"screen_name" form:"screen_name" binding:"required"`
	Title      string   `json:"title" form:"title" binding:"required"`
	Endpoints  []string `json:"endpoints" form:"endpoints" binding:"required"`
	Counters   []string `json:"counters" form:"counters" binding:"required"`
	TimeSpan   int64    `json:"timespan" form:"timespan"`
	GraphType  string   `json:"graph_type" form:"graph_type" binding:"required"`
	Method     string   `json:"method" form:"method"`
	Position   int64    `json:"position" form:"position"`
	FalconTags string   `json:"falcon_tags" form:"falcon_tags"`
}

func (mine APIGraphCreateReqDataWithNewScreenInputs) Check() (err error) {
	sc := m.DashboardScreen{Name: mine.ScreenName}
	// check screen_id
	if sc.ExistName() {
		err = fmt.Errorf("screen name:%v already existing", mine.ScreenName)
		return
	}

	if mine.TimeSpan%60 != 0 {
		err = fmt.Errorf("value of timespan is not vaild: %v", mine.TimeSpan)
		return
	}

	if mine.GraphType != "h" && mine.GraphType != "k" && mine.GraphType != "a" {
		err = fmt.Errorf("value of graph_type only accpet 'k' or 'h' or 'a', you typed: %v", mine.GraphType)
		return
	}
	return
}

func GraphCreateReqDataWithNewScreen(c *gin.Context) {
	inputs := APIGraphCreateReqDataWithNewScreenInputs{}
	// set default value
	inputs.TimeSpan = 3600
	inputs.GraphType = "h"
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	if err := inputs.Check(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	dt := db.Dashboard.Begin()
	sc := m.DashboardScreen{Name: inputs.ScreenName, Creator: user.Name}
	dt = dt.Save(&sc)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		dt.Rollback()
		return
	}

	es := inputs.Endpoints
	cs := inputs.Counters
	sort.Strings(es)
	sort.Strings(cs)
	esString := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
	csString := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	d := m.DashboardGraph{
		Title:     inputs.Title,
		Hosts:     esString,
		Counters:  csString,
		ScreenId:  sc.ID,
		TimeSpan:  inputs.TimeSpan,
		GraphType: inputs.GraphType,
		Method:    inputs.Method,
		Position:  inputs.Position,
		Creator:   user.Name,
	}
	dt = dt.Save(&d)
	if dt.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	var lid []int
	dt = dt.Table(d.TableName()).Raw("select LAST_INSERT_ID() as id").Pluck("id", &lid)
	if dt.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	dt.Commit()
	aid := lid[0]

	h.JSONR(c, map[string]interface{}{"id": aid, "screen_id": d.ScreenId, "screen_name": inputs.ScreenName})
}
