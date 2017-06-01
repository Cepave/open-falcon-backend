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
	ID         int64    `json:"id" form:"id" binding:"required"`
	ScreenId   int64    `json:"screen_id" form:"screen_id"`
	Title      string   `json:"title" form:"title"`
	Endpoints  []string `json:"endpoints" form:"endpoints"`
	Counters   []string `json:"counters" form:"counters"`
	TimeSpan   int64    `json:"timespan" form:"timespan"`
	GraphType  string   `json:"graph_type" form:"graph_type"`
	Method     string   `json:"method" form:"method"`
	Position   int64    `json:"position" form:"position"`
	FalconTags string   `json:"falcon_tags" form:"falcon_tags"`
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
	inputs := APIGraphUpdateReqData{Method: "NaN", FalconTags: "NaN"}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	graph := m.DashboardGraph{ID: inputs.ID}
	if dt := db.Dashboard.Model(&graph).Where("id = ?", inputs.ID).Scan(&graph); dt.Error != nil {
		h.JSONR(c, badstatus, fmt.Errorf("get graph id:%d, got error:%s", inputs.ID, dt.Error.Error()))
		return
	}

	if len(inputs.Endpoints) != 0 {
		es := inputs.Endpoints
		sort.Strings(es)
		es_string := strings.Join(es, TMP_GRAPH_FILED_DELIMITER)
		graph.Hosts = es_string
	}
	if len(inputs.Counters) != 0 {
		cs := inputs.Counters
		sort.Strings(cs)
		cs_string := strings.Join(cs, TMP_GRAPH_FILED_DELIMITER)
		graph.Counters = cs_string
	}
	if inputs.Title != "" {
		graph.Title = inputs.Title
	}
	if inputs.ScreenId != 0 {
		graph.ScreenId = int64(inputs.ScreenId)
	}
	if inputs.TimeSpan != 0 {
		graph.TimeSpan = inputs.TimeSpan
	}
	if inputs.GraphType != "" {
		graph.GraphType = inputs.GraphType
	}
	//method accpect empty, this is means not SUM action
	if inputs.Method != "NaN" {
		graph.Method = inputs.Method
	} else {
		graph.Method = ""
	}
	if inputs.Position != 0 {
		graph.Position = inputs.Position
	}
	if inputs.FalconTags != "NaN" {
		graph.FalconTags = inputs.FalconTags
	} else {
		graph.FalconTags = ""
	}

	tx := db.Dashboard.Begin()
	if dt := tx.Model(&graph).Where("id = ?", inputs.ID).Save(&graph); dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, fmt.Errorf("update graph id:%d, got error:%s", inputs.ID, dt.Error.Error()))
		return
	}
	tx.Commit()
	h.JSONR(c, buildGraphGetOutput(graph))
}

type APIDashboardGraphGetOuput struct {
	GraphID    int64    `json:"graph_id" form:"graph_id"`
	Title      string   `json:"title" form:"title"`
	ScreenId   int64    `json:"screen_id" form:"screen_id"`
	Endpoints  []string `json:"endpoints" form:"endpoints"`
	Counters   []string `json:"counters" form:"counters"`
	TimeSpan   int64    `json:"timespan" form:"timespan"`
	GraphType  string   `json:"graph_type" form:"graph_type"`
	Method     string   `json:"method" form:"method"`
	Position   int64    `json:"position" form:"position"`
	FalconTags string   `json:"falcon_tags" form:"falcon_tags"`
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

	h.JSONR(c, buildGraphGetOutput(graph))

}

func buildGraphGetOutput(graph m.DashboardGraph) APIDashboardGraphGetOuput {
	es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
	cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)
	return APIDashboardGraphGetOuput{
		GraphID:    graph.ID,
		Title:      graph.Title,
		Endpoints:  es,
		Counters:   cs,
		ScreenId:   graph.ScreenId,
		GraphType:  graph.GraphType,
		TimeSpan:   graph.TimeSpan,
		Method:     graph.Method,
		Position:   graph.Position,
		FalconTags: graph.FalconTags,
	}
}

type APIDashboardGraphDeleteInputs struct {
	ID int `json:"id" form:"id"  binding:"required"`
}

func DashboardGraphDelete(c *gin.Context) {
	id := c.Param("id")
	gid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid graph id")
		return
	}

	tx := db.Dashboard.Begin()
	graph := m.DashboardGraph{}
	dt := tx.Model(&graph).Where("id = ?", gid).Delete(&graph)
	if dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	tx.Commit()
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

type APIDashboardGraphCloneInputs struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}

func DashboardGraphClone(c *gin.Context) {
	inputs := APIDashboardGraphCloneInputs{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, fmt.Errorf("binding inputs got error:%s", err.Error()))
		return
	}
	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	originalGraph := m.DashboardGraph{ID: inputs.ID}
	if dt := db.Dashboard.Model(&originalGraph).Where(&originalGraph).Scan(&originalGraph); dt.Error != nil {
		h.JSONR(c, badstatus, fmt.Errorf("find graph with id: %d, got error:%s", inputs.ID, dt.Error.Error()))
		return
	}
	tx := db.Dashboard.Begin()
	newGraph := m.DashboardGraph{
		Title:      fmt.Sprintf("%s_copy", originalGraph.Title),
		Hosts:      originalGraph.Hosts,
		Counters:   originalGraph.Counters,
		ScreenId:   originalGraph.ScreenId,
		TimeSpan:   originalGraph.TimeSpan,
		GraphType:  originalGraph.GraphType,
		Method:     originalGraph.Method,
		Position:   originalGraph.Position,
		FalconTags: originalGraph.FalconTags,
		Creator:    user.Name,
	}
	if dt := tx.Model(&newGraph).Save(&newGraph); dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, fmt.Errorf("save new graph record got error: %s", dt.Error.Error()))
		return
	}
	tx.Commit()
	h.JSONR(c, newGraph)
}
