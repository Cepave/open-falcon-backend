package nqmDemo

import (
	"encoding/json"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
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
