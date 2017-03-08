package dashboardGraphOwl

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
	h "github.com/Cepave/open-falcon-backend/modules/api/app/helper"
	d "github.com/Cepave/open-falcon-backend/modules/api/app/model/dashboard"
)

type APICreateGraphInput struct {
	Title     string `json:"title" binding:"required"`
	Hosts     string `json:"hosts" binding:"required"`
	Counters  string `json:"counters" binding:"required"`
	ScreenId  int64  `json:"screen_id" binding:"required"`
	Timespan  int64  `json:"timespan"`
	GraphType string `json:"graph_type"`
	Method    string `json:"method"`
}

func fillDefault(current APICreateGraphInput) APICreateGraphInput {
	newOne := current
	if newOne.Timespan == 0 {
		newOne.Timespan = 3600
	}
	if newOne.GraphType == "" {
		newOne.GraphType = "h"
	}
	return newOne
}

func CreateGraph(c *gin.Context) {
	var inputs APICreateGraphInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	inputs = fillDefault(inputs)
	hosts := converInputFormating(inputs.Hosts)
	counters := converInputFormating(inputs.Counters)
	graph := d.DashboardGraph{
		Title:     inputs.Title,
		Hosts:     hosts,
		Counters:  counters,
		ScreenId:  inputs.ScreenId,
		TimeSpan:  inputs.Timespan,
		GraphType: inputs.GraphType,
		Method:    inputs.Method,
	}
	tx := db.Dashboard.Begin()
	if dt := tx.Save(&graph); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		tx.Rollback()
		return
	}
	graph.Position = graph.ID
	tx.Model(&graph).Updates(&graph)
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("graph is created, id: %v", graph.ID))
	return
}

type APIUpdateGraphInput struct {
	ID        int64  `json:"id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Hosts     string `json:"hosts" binding:"required"`
	Counters  string `json:"counters" binding:"required"`
	Timespan  int64  `json:"timespan"`
	GraphType string `json:"graph_type"`
	Method    string `json:"method"`
	ScreenId  int64  `json:"screen_id" binding:"required"`
}

func UpdateGraph(c *gin.Context) {
	var inputs APIUpdateGraphInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	tx := db.Dashboard.Begin()
	hosts := converInputFormating(inputs.Hosts)
	counters := converInputFormating(inputs.Counters)
	graph := map[string]interface{}{
		"Title":     inputs.Title,
		"Hosts":     hosts,
		"Counters":  counters,
		"ScreenId":  inputs.ScreenId,
		"Timespan":  inputs.Timespan,
		"GraphType": inputs.GraphType,
		"Method":    inputs.Method,
	}
	if dt := tx.Model(&d.DashboardGraph{}).Where("id = ?", inputs.ID).Update(graph); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("graph:%d has been updated.", inputs.ID))
	return
}

type APIGetGraphOutput struct {
	ID        string `json:"id" binding:"required"`
	Title     string `json:"titile" binding:"required"`
	Hosts     string `json:"hosts" binding:"required"`
	Counters  string `json:"counters" binding:"required"`
	Timespan  int64  `json:"timespan"`
	GraphType string `json:"graph_type"`
	Method    string `json:"method"`
	ScreenId  int64  `json:"screen_id" binding:"required"`
}

func GetGraph(c *gin.Context) {
	gid := c.Params.ByName("gid")
	if gid == "" {
		h.JSONR(c, badstatus, "graph id is missing")
		return
	}
	gidi, err := strconv.Atoi(gid)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	graph := d.DashboardGraph{}
	if dt := db.Dashboard.Model(&graph).Where("id = ?", gidi).Scan(&graph); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	graph.Hosts = converOutpuFormating(graph.Hosts)
	graph.Counters = converOutpuFormating(graph.Counters)
	h.JSONR(c, graph)
	return
}

func DeleteGraph(c *gin.Context) {
	gid := c.Params.ByName("gid")
	if gid == "" {
		h.JSONR(c, badstatus, "graph id is missing")
		return
	}
	gidi, err := strconv.Atoi(gid)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	tx := db.Dashboard.Begin()
	graph := d.DashboardGraph{ID: int64(gidi)}
	if dt := tx.Delete(&graph); dt.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("graph:%d has been deleted", gidi))
	return
}

func converOutpuFormating(input string) string {
	return strings.Join(strings.Split(input, "|"), ",")
}

func converInputFormating(input string) string {
	return strings.Join(strings.Split(input, ","), "|")
}
