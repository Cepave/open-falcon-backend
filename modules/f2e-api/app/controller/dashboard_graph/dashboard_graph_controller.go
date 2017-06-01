package dashboard_graph

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	m "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/dashboard"
	"github.com/gin-gonic/gin"
)

type APIGraphCreateReqData struct {
	ScreenId   int64    `json:"screen_id" form:"screen_id" binding:"required"`
	Title      string   `json:"title" form:"title" binding:"required"`
	Endpoints  []string `json:"endpoints" form:"endpoints" binding:"required"`
	Counters   []string `json:"counters" form:"counters" binding:"required"`
	TimeSpan   int64    `json:"timespan" form:"timespan"`
	GraphType  string   `json:"graph_type" form:"graph_type" binding:"required"`
	Method     string   `json:"method" form:"method"`
	Position   int64    `json:"position" form:"position"`
	FalconTags string   `json:"falcon_tags" form:"falcon_tags"`
}

func (mine APIGraphCreateReqData) Check() (err error) {
	sc := m.DashboardScreen{ID: mine.ScreenId}
	// check screen_id
	if !sc.Exist() {
		err = fmt.Errorf("screen id:%v is not existing", mine.ScreenId)
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

func DashboardGraphCreate(c *gin.Context) {
	inputs := APIGraphCreateReqData{TimeSpan: 3600, GraphType: "h"}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	if err := inputs.Check(); err != nil {
		h.JSONR(c, badstatus, err)
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
		ScreenId:  int64(inputs.ScreenId),
		TimeSpan:  inputs.TimeSpan,
		GraphType: inputs.GraphType,
		Method:    inputs.Method,
		Position:  inputs.Position,
		Creator:   user.Name,
	}
	dt := db.Dashboard.Begin()
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

	h.JSONR(c, map[string]interface{}{"id": aid, "screen_id": d.ScreenId})

}

type APIGraphUpdateReqData struct {
	ID int64 `json:"id" form:"id" binding:"required"`
	APIGraphCreateReqData
}

func (mine APIGraphUpdateReqData) Check() (err error) {
	sc := m.DashboardScreen{ID: mine.ScreenId}
	// check screen_id
	if mine.ScreenId != 0 && !sc.Exist() {
		err = fmt.Errorf("screen id:%v is not existing", mine.ScreenId)
		return
	}

	if mine.ScreenId != 0 && mine.TimeSpan%60 != 0 {
		err = fmt.Errorf("value of timespan is not vaild: %v", mine.TimeSpan)
		return
	}

	if mine.GraphType != "" && mine.GraphType != "h" && mine.GraphType != "k" {
		err = fmt.Errorf("value of graph_type only accpet 'k' or 'h', you typed: %v", mine.GraphType)
		return
	}
	return
}

func DashboardGraphUpdate(c *gin.Context) {
	inputs := APIGraphUpdateReqData{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	d := m.DashboardGraph{ID: inputs.ID}

	if len(inputs.Endpoints) != 0 {
		es := inputs.Endpoints
		sort.Strings(es)
		es_string := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
		d.Hosts = es_string
	}
	if len(inputs.Counters) != 0 {
		cs := inputs.Counters
		sort.Strings(cs)
		cs_string := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
		d.Counters = cs_string
	}
	if inputs.Title != "" {
		d.Title = inputs.Title
	}
	if inputs.ScreenId != 0 {
		d.ScreenId = int64(inputs.ScreenId)
	}
	if inputs.TimeSpan != 0 {
		d.TimeSpan = inputs.TimeSpan
	}
	if inputs.GraphType != "" {
		d.GraphType = inputs.GraphType
	}
	if inputs.Method != "" {
		d.Method = inputs.Method
	}
	if inputs.Position != 0 {
		d.Position = inputs.Position
	}
	if inputs.FalconTags != "" {
		d.FalconTags = inputs.FalconTags
	}

	graph := m.DashboardGraph{}
	dt := db.Dashboard.Table(graph.TableName()).Model(&graph).Where("id = ?", inputs.ID).Updates(d)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, map[string]int64{"id": inputs.ID})

}

func DashboardGraphGet(c *gin.Context) {
	id := c.Param("id")
	gid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid graph id")
		return
	}

	graph := m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Where("id = ?", gid).First(&graph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
	cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)

	h.JSONR(c, map[string]interface{}{
		"graph_id":    graph.ID,
		"title":       graph.Title,
		"endpoints":   es,
		"counters":    cs,
		"screen_id":   graph.ScreenId,
		"graph_type":  graph.GraphType,
		"timespan":    graph.TimeSpan,
		"method":      graph.Method,
		"position":    graph.Position,
		"falcon_tags": graph.FalconTags,
	})

}

func DashboardGraphDelete(c *gin.Context) {
	id := c.Param("id")
	gid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid graph id")
		return
	}

	graph := m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Where("id = ?", gid).Delete(&graph)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, map[string]int{"id": gid})

}

func DashboardGraphGetsByScreenID(c *gin.Context) {
	id := c.Param("screen_id")
	sid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen id")
		return
	}
	limit := c.DefaultQuery("limit", "500")

	graphs := []m.DashboardGraph{}
	dt := db.Dashboard.Table("dashboard_graph").Where("screen_id = ?", sid).Limit(limit).Find(&graphs)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	ret := []map[string]interface{}{}
	for _, graph := range graphs {
		es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
		cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)

		r := map[string]interface{}{
			"graph_id":    graph.ID,
			"title":       graph.Title,
			"endpoints":   es,
			"counters":    cs,
			"screen_id":   graph.ScreenId,
			"graph_type":  graph.GraphType,
			"timespan":    graph.TimeSpan,
			"method":      graph.Method,
			"position":    graph.Position,
			"falcon_tags": graph.FalconTags,
		}
		ret = append(ret, r)
	}

	h.JSONR(c, ret)
}
