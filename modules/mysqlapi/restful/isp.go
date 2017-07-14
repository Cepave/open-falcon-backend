package restful

import (
	"net/http"

	commonOwlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	"github.com/gin-gonic/gin"
)

func listISPs(c *gin.Context) {
	var ispName string
	if v, ok := c.GetQuery("name"); ok {
		ispName = v
	}
	ISPs := commonOwlDb.GetISPsByName(ispName)
	c.JSON(http.StatusOK, ISPs)
}

func getISPByID(
	p *struct {
		ISPID int16 `mvc:"param[isp_id]"`
	},
) mvc.OutputBody {
	return mvc.JsonOutputOrNotFound(
		commonOwlDb.GetIspById(p.ISPID),
	)
}
