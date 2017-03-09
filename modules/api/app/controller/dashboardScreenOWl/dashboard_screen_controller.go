package dashboardScreenOWl

import (
	"fmt"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
	h "github.com/Cepave/open-falcon-backend/modules/api/app/helper"
	d "github.com/Cepave/open-falcon-backend/modules/api/app/model/dashboard"
)

type APICreateScreenInput struct {
	Name string `json:"name" binding:"required"`
}

func CreateScreen(c *gin.Context) {
	var inputs APICreateScreenInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	screenName := inputs.Name
	if screenName == "" {
		h.JSONR(c, badstatus, "name is empty, please check it")
		return
	}
	screen := d.DashboardScreen{
		Name: screenName,
	}
	tx := db.Dashboard.Begin()
	if dt := tx.Save(&screen); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("screen is created, id: %v", screen.ID))
	return
}

type APIGetScreenListOuput struct {
	Screen      d.DashboardScreen `json:"screen"`
	GraphTitles []string          `json:"graph_names"`
	GraphIds    []int64           `json:"graph_ids"`
}

func GetScreenList(c *gin.Context) {
	screens := []d.DashboardScreen{}
	db.Dashboard.Model(&screens).Scan(&screens)
	result := []APIGetScreenListOuput{}
	for _, s := range screens {
		tmp := APIGetScreenListOuput{}
		tmp.Screen = s
		graphs := s.Graphs()
		graphTitles := []string{}
		graphIds := []int64{}
		for _, g := range graphs {
			graphTitles = append(graphTitles, g.Title)
			graphIds = append(graphIds, g.ID)
		}
		tmp.GraphTitles = graphTitles
		tmp.GraphIds = graphIds
		result = append(result, tmp)
	}
	h.JSONR(c, result)
	return
}

type APIUpdateScreenInput struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func UpdateScreen(c *gin.Context) {
	var inputs APIUpdateScreenInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	tx := db.Dashboard.Begin()
	screen := d.DashboardScreen{
		ID:   inputs.ID,
		Name: inputs.Name,
	}
	if dt := tx.Model(&screen).Where("id = ?", screen.ID).Update(&screen); dt.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("screen: %d has been updated.", inputs.ID))
	return
}

func GetScreen(c *gin.Context) {
	sid := c.Params.ByName("sid")
	if sid == "" {
		h.JSONR(c, badstatus, "screen id is missing")
		return
	}
	sidi, err := strconv.Atoi(sid)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	screen := d.DashboardScreen{}
	if dt := db.Dashboard.Model(&screen).Where("id = ?", sidi).Scan(&screen); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	var output APIGetScreenListOuput
	output.Screen = screen
	graphs := screen.Graphs()
	graphTitles := []string{}
	graphIds := []int64{}
	for _, g := range graphs {
		graphTitles = append(graphTitles, g.Title)
		graphIds = append(graphIds, g.ID)
	}
	output.GraphTitles = graphTitles
	output.GraphIds = graphIds
	h.JSONR(c, output)
	return
}

func DeleteScreen(c *gin.Context) {
	sid := c.Params.ByName("sid")
	if sid == "" {
		h.JSONR(c, badstatus, "screen id is missing")
		return
	}
	sidi, err := strconv.Atoi(sid)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	screen := d.DashboardScreen{ID: int64(sidi)}
	tx := db.Dashboard.Begin()
	if dt := tx.Delete(&screen); dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	graph := d.DashboardGraph{}
	if dt := tx.Where("screen_id = ?", sidi).Delete(&graph); dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("screen: %d and related graph is removed.", sidi))
	return
}
