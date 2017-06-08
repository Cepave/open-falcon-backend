package dashboard_graph

import (
	"sort"
	"strings"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	m "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/dashboard"
	"github.com/gin-gonic/gin"
)

func GraphCreateReqDataWithNewScreen(c *gin.Context) {
	// set default value
	inputs := APIGraphCreateReqDataWithNewScreenInputs{
		TimeSpan:     3600,
		GraphType:    "h",
		TimeRange:    "3h",
		SortBy:       "a-z",
		SampleMethod: "AVERAGE",
	}
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
	user, err = h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	d := m.DashboardGraph{
		Title:        inputs.Title,
		Hosts:        esString,
		Counters:     csString,
		ScreenId:     sc.ID,
		TimeSpan:     inputs.TimeSpan,
		GraphType:    inputs.GraphType,
		Method:       inputs.Method,
		Position:     inputs.Position,
		Creator:      user.Name,
		TimeRange:    inputs.TimeRange,
		SortBy:       inputs.SortBy,
		SampleMethod: inputs.SampleMethod,
		YScale:       inputs.YScale,
	}
	dt = dt.Save(&d)
	if dt.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	dt.Commit()

	h.JSONR(c, map[string]interface{}{"graph": buildGraphGetOutput(d), "screen_id": d.ScreenId, "screen_name": inputs.ScreenName})
}
