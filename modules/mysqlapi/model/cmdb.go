package model

import (
	"github.com/gin-gonic/gin"

	oGin "github.com/Cepave/open-falcon-backend/common/gin"
)

type SyncHost struct {
	Activate int
	Name     string
	IP       string
}

type SyncHostGroup struct {
	Creator string
	Name    string
}

type SyncForAdding struct {
	Hosts      []*SyncHost
	Hostgroups []*SyncHostGroup
	Relations  map[string][]string
}

func (p *SyncForAdding) Bind(c *gin.Context) {
	oGin.BindJson(c, p)
}
