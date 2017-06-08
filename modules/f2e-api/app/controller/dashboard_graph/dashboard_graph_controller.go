package dashboard_graph

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	m "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/dashboard"
	"github.com/gin-gonic/gin"
)

func DashboardGraphCreate(c *gin.Context) {
	//set default values
	inputs := APIGraphCreateReqData{
		TimeSpan:     3600,
		TimeRange:    "3h",
		GraphType:    "h",
		SortBy:       "a-z",
		SampleMethod: "AVERAGE"}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	if err := inputs.Check(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	screen := m.DashboardScreen{ID: inputs.ScreenId}
	if !screen.Exist() {
		h.JSONR(c, badstatus, fmt.Sprintf("screen id: %d record not found", screen.ID))
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
		Title:        inputs.Title,
		Hosts:        esString,
		Counters:     csString,
		ScreenId:     int64(inputs.ScreenId),
		TimeSpan:     inputs.TimeSpan,
		GraphType:    inputs.GraphType,
		Method:       inputs.Method,
		Position:     inputs.Position,
		Creator:      user.Name,
		TimeRange:    inputs.TimeRange,
		SampleMethod: inputs.SampleMethod,
		SortBy:       inputs.SortBy,
		YScale:       inputs.YScale,
	}
	dt := db.Dashboard.Begin()
	dt = dt.Save(&d)
	if dt.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	dt.Commit()
	db.Dashboard.Model(&screen).Where("id = ?", screen.ID).Scan(&screen)
	h.JSONR(c, map[string]interface{}{"graph": buildGraphGetOutput(d), "screen_id": screen.ID, "screen_name": screen.Name})
}

func DashboardGraphUpdate(c *gin.Context) {
	inputs := APIGraphUpdateReqData{Method: "NULL", FalconTags: "NULL"}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.Check(); err != nil {
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
	if inputs.Method != "NULL" {
		graph.Method = inputs.Method
	} else {
		graph.Method = ""
	}
	if inputs.Position != 0 {
		graph.Position = inputs.Position
	}
	if inputs.FalconTags != "NULL" {
		graph.FalconTags = inputs.FalconTags
	} else {
		graph.FalconTags = ""
	}

	if inputs.TimeRange != "" {
		graph.TimeRange = inputs.TimeRange
	}
	if inputs.SampleMethod != "" {
		graph.SampleMethod = inputs.SampleMethod
	}
	if inputs.YScale != "" {
		graph.YScale = inputs.YScale
	}
	if inputs.SortBy != "" {
		graph.SortBy = inputs.SortBy
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
	GraphID      int64    `json:"graph_id" form:"graph_id"`
	Title        string   `json:"title" form:"title"`
	ScreenId     int64    `json:"screen_id" form:"screen_id"`
	Endpoints    []string `json:"endpoints" form:"endpoints"`
	Counters     []string `json:"counters" form:"counters"`
	TimeSpan     int64    `json:"timespan" form:"timespan"`
	GraphType    string   `json:"graph_type" form:"graph_type"`
	Method       string   `json:"method" form:"method"`
	Position     int64    `json:"position" form:"position"`
	FalconTags   string   `json:"falcon_tags" form:"falcon_tags"`
	TimeRange    string   `json:"time_range" form:"time_range"`
	YScale       string   `json:"y_scale" form:"y_scale"`
	SortBy       string   `json:"sort_by" form:"sort_by"`
	SampleMethod string   `json:"sample_method" form:"sample_method"`
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
		GraphID:      graph.ID,
		Title:        graph.Title,
		Endpoints:    es,
		Counters:     cs,
		ScreenId:     graph.ScreenId,
		GraphType:    graph.GraphType,
		TimeSpan:     graph.TimeSpan,
		Method:       graph.Method,
		Position:     graph.Position,
		FalconTags:   graph.FalconTags,
		TimeRange:    graph.TimeRange,
		YScale:       graph.YScale,
		SampleMethod: graph.SampleMethod,
		SortBy:       graph.SortBy,
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

	ret := make([]APIDashboardGraphGetOuput, len(graphs))
	for _, graph := range graphs {
		r := buildGraphGetOutput(graph)
		ret = append(ret, r)
	}

	h.JSONR(c, ret)
}

type APIDashboardGraphCloneInputs struct {
	ID   int64  `json:"id" form:"id" binding:"required"`
	Name string `json:"name" form:"name"`
}

func DashboardGraphClone(c *gin.Context) {
	inputs := APIDashboardGraphCloneInputs{Name: "NaN"}
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
	if inputs.Name == "NaN" {
		inputs.Name = fmt.Sprintf("%s_copy_%d", originalGraph.Title, time.Now().Unix())
	}
	tx := db.Dashboard.Begin()
	newGraph := m.DashboardGraph{
		Title:        inputs.Name,
		Hosts:        originalGraph.Hosts,
		Counters:     originalGraph.Counters,
		ScreenId:     originalGraph.ScreenId,
		TimeSpan:     originalGraph.TimeSpan,
		GraphType:    originalGraph.GraphType,
		Method:       originalGraph.Method,
		Position:     originalGraph.Position,
		FalconTags:   originalGraph.FalconTags,
		Creator:      user.Name,
		SampleMethod: originalGraph.SampleMethod,
		TimeRange:    originalGraph.TimeRange,
		SortBy:       originalGraph.SortBy,
		YScale:       originalGraph.YScale,
	}
	if dt := tx.Model(&newGraph).Save(&newGraph); dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, fmt.Errorf("save new graph record got error: %s", dt.Error.Error()))
		return
	}
	tx.Commit()
	h.JSONR(c, buildGraphGetOutput(newGraph))
}
