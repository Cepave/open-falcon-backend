package model

import (
	"github.com/gin-gonic/gin"

	oGin "github.com/Cepave/open-falcon-backend/common/gin"
)

type SyncHost struct {
	Activate int    `json:"activate"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
}

type SyncHostGroup struct {
	Creator string `json:"creator"`
	Name    string `json:"name"`
}

type SyncForAdding struct {
	Hosts      []*SyncHost         `json:"hosts"`
	Hostgroups []*SyncHostGroup    `json:"hostgroups"`
	Relations  map[string][]string `json:"relations"`
}

func (p *SyncForAdding) Bind(c *gin.Context) {
	oGin.BindJson(c, p)
}
