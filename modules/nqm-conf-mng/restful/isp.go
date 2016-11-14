package restful

import (
	"net/http"

	commonOwlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	"gopkg.in/gin-gonic/gin.v1"
)

func listISPs(c *gin.Context) {
	IspName := c.Param("name")
	ISPs := commonOwlDb.GetISPsByName(IspName)
	c.JSON(http.StatusOK, ISPs)
}
