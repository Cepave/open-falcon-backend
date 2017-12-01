package restful

import (
	"github.com/gin-gonic/gin"

	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	CmdbModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	CmdbDb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/cmdb"
)

func addNewCmdbSync(c *gin.Context) {
	/**
	 * Builds data from body of request
	 */
	targetForAdding := &CmdbModel.SyncForAdding{}

	commonGin.BindJson(c, targetForAdding)
	logger.Infof("Get data from API %s", targetForAdding)
	// :~)
	CmdbDb.SyncForHosts(targetForAdding)
}
