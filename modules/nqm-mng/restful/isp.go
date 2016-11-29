package restful

import (
	"net/http"

	commonOwlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	"gopkg.in/gin-gonic/gin.v1"
)

func listISPs(c *gin.Context) {
	var ispName string
	if v, ok := c.GetQuery("name"); ok {
		ispName = v
	}
	ISPs := commonOwlDb.GetISPsByName(ispName)
	c.JSON(http.StatusOK, ISPs)
}
