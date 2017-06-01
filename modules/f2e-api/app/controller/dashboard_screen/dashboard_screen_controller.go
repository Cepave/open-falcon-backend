package dashboard_screen

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	m "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/dashboard"
	"github.com/gin-gonic/gin"
)

type APIScreenCreateInput struct {
	Pid  int64  `json:"pid" form:"pid"`
	Name string `json:"name" form:"name" binding:"required"`
}

func ScreenCreate(c *gin.Context) {
	inputs := APIScreenCreateInput{Pid: 0}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	newDS := m.DashboardScreen{
		PID:     inputs.Pid,
		Name:    inputs.Name,
		Creator: user.Name,
	}
	rcount := 0
	dt := db.Dashboard.Begin()
	dt = dt.Table(newDS.TableName()).Where("name = ?", newDS.Name).Count(&rcount)
	if rcount != 0 {
		h.JSONR(c, badstatus, fmt.Errorf("screen name '%s' alreay exist", newDS.Name))
		dt.Rollback()
		return
	}
	if inputs.Pid != 0 {
		dt = dt.Table(newDS.TableName()).Where("pid = ?", newDS.PID).Count(&rcount)
		if rcount == 0 {
			dt.Rollback()
			h.JSONR(c, badstatus, fmt.Errorf("parent id: %v, record not found", newDS.PID))
			return
		}
	}
	dt = dt.Save(&newDS)
	if dt.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	dt.Commit()
	h.JSONR(c, newDS)
}

func ScreenGet(c *gin.Context) {
	id := c.Param("screen_id")

	sid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen id")
		return
	}

	screen := m.DashboardScreen{}
	dt := db.Dashboard.Table("dashboard_screen").Where("id = ?", sid).First(&screen)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	graphsTmp := []m.DashboardGraph{}
	db.Dashboard.Model(&graphsTmp).Where("screen_id = ?", screen.ID).Scan(&graphsTmp)
	graphs := []map[string]interface{}{}
	for _, graph := range graphsTmp {
		es := strings.Split(graph.Hosts, TMP_GRAPH_FILED_DELIMITER)
		cs := strings.Split(graph.Counters, TMP_GRAPH_FILED_DELIMITER)
		graphs = append(graphs, map[string]interface{}{
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
	h.JSONR(c, map[string]interface{}{
		"scren":  screen,
		"graphs": graphs,
	})
}

func ScreenGetsByPid(c *gin.Context) {
	id := c.Param("pid")

	pid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen pid")
		return
	}

	screens := []m.DashboardScreen{}
	dt := db.Dashboard.Table("dashboard_screen").Where("pid = ?", pid).Find(&screens)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, screens)
}

type APIScreenGetsAllInputs struct {
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	KeyWord string `json:"key_word" form:"key_word"`
	Order   bool   `json:"desc" form:"desc"`
}

type APIScreenCreateOutput struct {
	GraphNames []string `json:"graph_names"`
	m.DashboardScreen
}

func ScreenGetsAll(c *gin.Context) {
	inputs := APIScreenGetsAllInputs{Limit: 500, Page: -1, Order: false}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	screens := []m.DashboardScreen{}
	totallCount := 0
	dt := db.Dashboard.Model(&screens)
	if inputs.KeyWord != "" {
		dt = dt.Where("name like ?", "%"+inputs.KeyWord+"%")
	}
	if inputs.Page <= 0 {
		if inputs.Order {
			dt = dt.Order("id desc")
		} else {
			dt = dt.Order("id asc")
		}
		dt = dt.Limit(inputs.Limit).Find(&screens)
		if dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			return
		} else {
			outputs := builtAllScreenOuput(screens)
			h.JSONR(c, outputs)
			return
		}
	} else {
		dt.Count(&totallCount)
		if inputs.Order {
			dt = dt.Order("id desc")
		} else {
			dt = dt.Order("id asc")
		}
		dt.Offset(inputs.Limit * (inputs.Page - 1)).Limit(inputs.Limit).Find(&screens)
		if dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			return
		} else {
			outputs := builtAllScreenOuput(screens)
			h.JSONR(c,
				map[string]interface{}{
					"current_page": inputs.Page,
					"totall_count": totallCount,
					"totall_page":  math.Ceil(float64(totallCount) / float64(inputs.Limit)),
					"data":         outputs,
					"order_by":     "id",
					"desc_order":   inputs.Order,
					"key_word":     inputs.KeyWord,
				})
			return
		}
	}
}

func ScreenDelete(c *gin.Context) {
	id := c.Param("screen_id")

	sid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen id")
		return
	}

	tx := db.Dashboard.Begin()
	screen := m.DashboardScreen{ID: int64(sid)}
	if dt := tx.Model(&screen).Where(&screen).Delete(&screen); dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	graph := m.DashboardGraph{ScreenId: int64(sid)}
	graphsID := []int64{}
	tx.Model(&graph).Where(&graph).Pluck("id", &graphsID)
	dt := tx.Model(&graph).Where(&graph).Delete(&graph)
	if dt := tx.Model(&graph).Where(&graph).Delete(&graph); dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	tx.Commit()
	h.JSONR(c, map[string]interface{}{
		"message":           "ok",
		"deleted_rows":      dt.RowsAffected,
		"deleted_graph_ids": graphsID,
	})
}

type APIScreenUpdateInputs struct {
	ID   int64  `json:"id" form:"id" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
}

// For rename usage
func ScreenUpdate(c *gin.Context) {
	inputs := APIScreenUpdateInputs{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if inputs.Name == "" {
		h.JSONR(c, badstatus, "name can not be empty")
		return
	}

	newData := m.DashboardScreen{
		ID:   inputs.ID,
		Name: inputs.Name,
	}

	if newData.ExistName() {
		h.JSONR(c, badstatus, fmt.Errorf("screen name '%s' alreay exist", newData.Name))
		return
	}
	dt := db.Dashboard.Table(newData.TableName()).Where("id = ?", inputs.ID).Update(newData)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, "ok")
}

// For clone screen by id
type APIScreenCloneInputs struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}

func ScreenClone(c *gin.Context) {
	inputs := APIScreenCloneInputs{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	originalScreen := m.DashboardScreen{ID: inputs.ID}
	if dt := db.Dashboard.Model(&originalScreen).Where(&originalScreen).Scan(&originalScreen); dt.Error != nil {
		h.JSONR(c, badstatus, fmt.Errorf("find screen by id:%d, got error:%s", inputs.ID, dt.Error.Error()))
		return
	}
	tx := db.Dashboard.Begin()
	newScreen := m.DashboardScreen{
		Name:    fmt.Sprintf("%s_copy", originalScreen.Name),
		Creator: user.Name,
	}
	if dt := tx.Model(&newScreen).Save(&newScreen); dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, fmt.Errorf("create new screen got error:%s", dt.Error.Error()))
		return
	}
	originalGList := originalScreen.Graphs()
	graphNames := make([]string, len(originalGList))
	for indx, gh := range originalGList {
		newGraph := m.DashboardGraph{
			Title:      gh.Title,
			Hosts:      gh.Hosts,
			Counters:   gh.Counters,
			ScreenId:   newScreen.ID,
			TimeSpan:   gh.TimeSpan,
			GraphType:  gh.GraphType,
			Method:     gh.Method,
			Position:   gh.Position,
			FalconTags: gh.FalconTags,
			Creator:    user.Name,
		}
		if dt := tx.Model(&newGraph).Save(&newGraph); dt.Error != nil {
			tx.Rollback()
			h.JSONR(c, badstatus, fmt.Errorf("create new graph with graph id: %d ,got error:%s", gh.ID, dt.Error.Error()))
			return
		}
		graphNames[indx] = gh.Title
	}
	tx.Commit()
	h.JSONR(c, APIScreenCreateOutput{
		DashboardScreen: newScreen,
		GraphNames:      graphNames,
	})
}

func builtAllScreenOuput(screens []m.DashboardScreen) []APIScreenCreateOutput {
	outputs := make([]APIScreenCreateOutput, len(screens))
	for indx, s := range screens {
		outputs[indx].DashboardScreen = s
		graphs := s.Graphs()
		if len(graphs) != 0 {
			gnames := make([]string, len(graphs))
			for indx, graphs := range graphs {
				gnames[indx] = graphs.Title
			}
			outputs[indx].GraphNames = gnames
		} else {
			outputs[indx].GraphNames = []string{}
		}
	}
	return outputs
}
