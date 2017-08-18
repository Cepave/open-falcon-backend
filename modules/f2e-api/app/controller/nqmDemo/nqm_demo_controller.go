package nqmDemo

import (
	"encoding/json"
	"strings"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	f "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/falcon_portal"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
)

// http://localhost:10080/api/v1/nqm_demo/agents?status=1
func Agents(c *gin.Context) {
	var dat []map[string]interface{}
	json.Unmarshal(agentData, &dat)
	h.JSONR(c, dat)
}

func Isps(c *gin.Context) {
	var dat []map[string]interface{}
	json.Unmarshal(IspData, &dat)
	h.JSONR(c, dat)
}

func Provinces(c *gin.Context) {
	var dat []map[string]interface{}
	json.Unmarshal(provincesData, &dat)
	h.JSONR(c, dat)
}

func Targets(c *gin.Context) {
	var dat []map[string]interface{}
	json.Unmarshal(targetsData, &dat)
	h.JSONR(c, dat)
}

func PingTasks(c *gin.Context) {
	var dat []map[string]interface{}
	json.Unmarshal(pingTaskData, &dat)
	h.JSONR(c, dat)
}

func Cities(c *gin.Context) {
	var dat []map[string]interface{}
	json.Unmarshal(CityData, &dat)
	h.JSONR(c, dat)
}

func NameTags(c *gin.Context) {
	var dat []map[string]interface{}
	json.Unmarshal(nameTagsData, &dat)
	h.JSONR(c, dat)
}

func GroupTags(c *gin.Context) {
	var dat []map[string]interface{}
	json.Unmarshal(groupTagsData, &dat)
	h.JSONR(c, dat)
}

type APIEmailDemoInputs struct {
	ToS     string `json:"tos" form:"tos"`
	Subject string `json:"subject" form:"subject" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
	TplID   int64  `json:"tpl_id" form:"tpl_id" binding:"required"`
}

func EmailDemo(c *gin.Context) {
	inputs := APIEmailDemoInputs{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	tpl := f.Template{}
	dt := db.Falcon.Table(tpl.TableName()).Where("id = ?", inputs.TplID).Scan(&tpl)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	tplAction := tpl.ActionID
	if tplAction == 0 {
		h.JSONR(c, badstatus, "no action record found")
		return
	}
	action := f.Action{}
	dt = db.Falcon.Table(action.TableName()).Where("id = ?", tplAction).Scan(&action)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	teams := action.FindUics()
	tset := hashset.New()
	for _, ts := range teams {
		users, _ := ts.Members()
		for _, user := range users {
			if strings.Contains(user.Email, "@") {
				tset.Add(user.Email)
			}
		}
	}
	emails := make([]string, tset.Size())
	for indx, mem := range tset.Values() {
		emails[indx] = mem.(string)
		sendMail(mem.(string), inputs.Subject, inputs.Subject)
	}
	h.JSONR(c, map[string]interface{}{
		"emails": emails,
	})
	return
}
