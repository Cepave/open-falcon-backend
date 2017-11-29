package restful

import (
	"github.com/gin-gonic/gin"
	"net/http"

	commonCmdbDb "github.com/Cepave/open-falcon-backend/common/db/cmdb"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	commonCmdbModel "github.com/Cepave/open-falcon-backend/common/model/cmdb"
)

func addNewCmdbSync(c *gin.Context) {
	/**
	 * Builds data from body of request
	 */
	targetForAdding := &commonCmdbModel.SyncForAdding{}

	commonGin.BindJson(c, targetForAdding)
	logger.Infof("Get data from API %s", targetForAdding)
	// :~)
	newSync, err := commonCmdbDb.SyncForHosts(targetForAdding)
	if err != nil {
		switch err.(type) {
		/*
			case commonCmdbDb.ErrDuplicatedNqmAgent:
				commonGin.JsonConflictHandler(
					c,
					commonGin.DataConflictError{
						ErrorCode:    1,
						ErrorMessage: err.Error(),
					},
				)
		*/
		default:
			panic(err)
		}

		return
	}

	c.JSON(http.StatusOK, newSync)
}
