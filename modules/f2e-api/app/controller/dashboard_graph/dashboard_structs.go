package dashboard_graph

import (
	"fmt"
	"regexp"

	m "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/dashboard"
)

type GraphObj struct {
	Title        string
	Endpoints    []string
	Counters     []string
	TimeSpan     int64
	GraphType    string
	Method       string
	Position     int64
	FalconTags   string
	TimeRange    string
	YScale       string
	SortBy       string
	SampleMethod string
}

func (mine GraphObj) CustomCheck() (err error) {

	if mine.Method != "" && mine.Method != "NULL" && mine.Method != "sum" {
		err = fmt.Errorf("method not vaild: %v, only accept 'sum' or 'NULL'", mine.Method)
		return
	}

	if mine.GraphType != "h" && mine.GraphType != "k" && mine.GraphType != "a" {
		err = fmt.Errorf("value of graph_type only accept 'k' or 'h' or 'a', you typed: %v", mine.GraphType)
		return
	}

	if mine.SortBy != "" && mine.SortBy != "a-z" && mine.SortBy != "z-a" {
		err = fmt.Errorf("sort_by only accept 'a-z' or 'z-a', you typed: %v", mine.GraphType)
		return
	}

	// only accpet
	// ex. 1495508389,1496372791
	// ex. 3d
	if mine.TimeRange != "" {
		if match, _ := regexp.MatchString("^((\\d{8}\\d{1,5},\\d{8}\\d{1,5})|(\\d+(h|d|w|m|y)))$", mine.TimeRange); !match {
			err = fmt.Errorf(`you typed: %v, time_range only accept ex. "1495508389,1496372791" or "3d" [h,d,w,m,y]`, mine.TimeRange)
			return
		}
	}

	if mine.YScale != "" {
		// [\d\.]+
		// k m g t
		if match, _ := regexp.MatchString("^(\\d+(k|m|g|t)|[\\d\\.]+)(,(\\d+(k|m|g|t)|[\\d\\.]+))?$", mine.YScale); !match {
			err = fmt.Errorf(`y_scale value: '%v' not vaild, please check api document`, mine.YScale)
			return
		}
	}

	if mine.SampleMethod != "" {
		flag := false
		slist := []string{"AVERAGE", "MAX", "MIN"}
		for _, s := range slist {
			flag = mine.SampleMethod == s
			if flag {
				break
			}
		}
		if !flag {
			err = fmt.Errorf(`sample_method value: '%v' not vaild, only accept: "AVERAGE" or "MAX" or "MIN"`, mine.SampleMethod)
			return
		}
	}
	return
}

// Inputs struct for GraphCreateReqDataWithNewScreen
type APIGraphCreateReqDataWithNewScreenInputs struct {
	ScreenName   string   `json:"screen_name" form:"screen_name" binding:"required"`
	Title        string   `json:"title" form:"title" binding:"required"`
	Endpoints    []string `json:"endpoints" form:"endpoints" binding:"required"`
	Counters     []string `json:"counters" form:"counters" binding:"required"`
	TimeSpan     int64    `json:"timespan" form:"timespan"`
	GraphType    string   `json:"graph_type" form:"graph_type" binding:"required"`
	Method       string   `json:"method" form:"method"`
	Position     int64    `json:"position" form:"position"`
	FalconTags   string   `json:"falcon_tags" form:"falcon_tags"`
	TimeRange    string   `json:"time_range" form:"time_range"`
	YScale       string   `json:"y_scale" form:"y_scale"`
	SortBy       string   `json:"sort_by" form:"sort_by"`
	SampleMethod string   `json:"sample_method" form:"sample_method"`
}

func (mine APIGraphCreateReqDataWithNewScreenInputs) Check() (err error) {
	sc := m.DashboardScreen{Name: mine.ScreenName}
	// check screen_id
	if sc.ExistName() {
		err = fmt.Errorf("screen name:%v already existing", mine.ScreenName)
		return
	}

	graphobj := GraphObj{
		GraphType:    mine.GraphType,
		SortBy:       mine.SortBy,
		TimeRange:    mine.TimeRange,
		YScale:       mine.YScale,
		SampleMethod: mine.SampleMethod,
		Method:       mine.Method,
	}
	err = graphobj.CustomCheck()
	return
}

// Inputs struct for DashboardGraphCreate
type APIGraphCreateReqData struct {
	ScreenId     int64    `json:"screen_id" form:"screen_id" binding:"required"`
	Title        string   `json:"title" form:"title" binding:"required"`
	Endpoints    []string `json:"endpoints" form:"endpoints" binding:"required"`
	Counters     []string `json:"counters" form:"counters" binding:"required"`
	TimeSpan     int64    `json:"timespan" form:"timespan"`
	GraphType    string   `json:"graph_type" form:"graph_type" binding:"required"`
	Method       string   `json:"method" form:"method"`
	Position     int64    `json:"position" form:"position"`
	FalconTags   string   `json:"falcon_tags" form:"falcon_tags"`
	TimeRange    string   `json:"time_range" form:"time_range" binding:"required"`
	YScale       string   `json:"y_scale" form:"y_scale"`
	SortBy       string   `json:"sort_by" form:"sort_by"`
	SampleMethod string   `json:"sample_method" form:"sample_method"`
}

func (mine APIGraphCreateReqData) Check() (err error) {
	sc := m.DashboardScreen{ID: mine.ScreenId}
	// check screen_id
	if !sc.Exist() {
		err = fmt.Errorf("screen id:%v is not existing", mine.ScreenId)
		return
	}

	graphobj := GraphObj{
		TimeSpan:     mine.TimeSpan,
		GraphType:    mine.GraphType,
		SortBy:       mine.SortBy,
		TimeRange:    mine.TimeRange,
		YScale:       mine.YScale,
		SampleMethod: mine.SampleMethod,
	}
	err = graphobj.CustomCheck()
	return
}

// Inputs struct for DashboardGraphUpdate
type APIGraphUpdateReqData struct {
	ID           int64    `json:"id" form:"id" binding:"required"`
	ScreenId     int64    `json:"screen_id" form:"screen_id"`
	Title        string   `json:"title" form:"title"`
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

func (mine APIGraphUpdateReqData) Check() (err error) {
	sc := m.DashboardScreen{ID: mine.ScreenId}
	// check screen_id
	if mine.ScreenId != 0 && !sc.Exist() {
		err = fmt.Errorf("screen id:%v is not existing", mine.ScreenId)
		return
	}

	graphobj := GraphObj{
		TimeSpan:     mine.TimeSpan,
		GraphType:    mine.GraphType,
		SortBy:       mine.SortBy,
		TimeRange:    mine.TimeRange,
		YScale:       mine.YScale,
		SampleMethod: mine.SampleMethod,
	}
	err = graphobj.CustomCheck()
	return
}
